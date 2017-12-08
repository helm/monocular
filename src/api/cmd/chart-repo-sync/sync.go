package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/disintegration/imaging"
	"github.com/ghodss/yaml"
	"github.com/jinzhu/copier"
	"github.com/kubeapps/common/datastore"
	"gopkg.in/mgo.v2/bson"
	helmrepo "k8s.io/helm/pkg/repo"
)

const (
	chartCollection       = "charts"
	chartReadmeCollection = "readmes"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var dbSession datastore.Session

var netClient httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [REPO NAME] [REPO URL]\n", os.Args[0])
		flag.PrintDefaults()
	}
	dbURL := flag.String("mongo-url", "localhost", "MongoDB URL (see https://godoc.org/labix.org/v2/mgo#Dial for format)")
	dbName := flag.String("mongo-database", "charts", "MongoDB database")
	dbUsername := flag.String("mongo-user", "", "MongoDB user")
	dbPassword := os.Getenv("MONGO_PASSWORD")
	debug := flag.Bool("debug", false, "verbose logging")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	mongoConfig := datastore.Config{URL: *dbURL, Database: *dbName, Username: *dbUsername, Password: dbPassword}
	var err error
	dbSession, err = datastore.NewSession(mongoConfig)
	if err != nil {
		log.WithFields(log.Fields{"host": *dbURL}).Fatal(err)
	}

	if err := sync(flag.Arg(0), flag.Arg(1)); err != nil {
		log.WithError(err).Error("sync failed")
		os.Exit(1)
	}
}

// Syncing is performed in the following sequential steps:
// 1. Update database to match chart metadata from index
// 3. Process icons for charts
// 4. Process READMEs for each chart version
//
// These steps are processed in this way to ensure relevant chart data is
// imported into the database as fast as possible. E.g. we want all icons for
// charts before fetching readmes for each chart and version pair.
func sync(repoName, repoURL string) error {
	url, err := url.ParseRequestURI(repoURL)
	if err != nil {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("failed to parse URL")
		return err
	}

	index, err := fetchRepoIndex(url)
	if err != nil {
		return err
	}

	charts := chartsFromIndex(index, repo{Name: repoName, URL: repoURL})
	err = importCharts(charts)
	if err != nil {
		return err
	}

	for _, c := range charts {
		if err := fetchAndImportIcon(c); err != nil {
			log.WithFields(log.Fields{"name": c.Name}).WithError(err).Error("failed to import icon")
		}
	}

	for _, c := range charts {
		fetchAndImportReadmes(c)
	}

	return nil
}

func fetchRepoIndex(repoURL *url.URL) (*helmrepo.IndexFile, error) {
	indexURL := *repoURL
	indexURL.Path = "/index.yaml"
	req, err := http.NewRequest("GET", indexURL.String(), nil)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("could not build repo index request")
		return nil, err
	}
	res, err := netClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("error requesting repo index")
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("error requesting repo index, are you sure this is a chart repository?")
		return nil, errors.New("repo index request failed")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseRepoIndex(body)
}

func parseRepoIndex(body []byte) (*helmrepo.IndexFile, error) {
	var index helmrepo.IndexFile
	err := yaml.Unmarshal(body, &index)
	if err != nil {
		return nil, err
	}
	index.SortEntries()
	return &index, nil
}

func chartsFromIndex(index *helmrepo.IndexFile, r repo) []chart {
	var charts []chart
	for _, entry := range index.Entries {
		if entry[0].GetDeprecated() {
			log.WithFields(log.Fields{"name": entry[0].GetName()}).Info("skipping deprecated chart")
			continue
		}
		charts = append(charts, newChart(entry, r))
	}
	return charts
}

// Takes an entry from the index and constructs a database representation of the
// object.
func newChart(entry helmrepo.ChartVersions, r repo) chart {
	var c chart
	copier.Copy(&c, entry[0])
	copier.Copy(&c.ChartVersions, entry)
	c.Repo = r
	c.ID = fmt.Sprintf("%s/%s", r.Name, c.Name)
	return c
}

func importCharts(charts []chart) error {
	var pairs []interface{}
	var chartIDs []string
	for _, c := range charts {
		chartIDs = append(chartIDs, c.ID)
		// charts to upsert - pair of selector, chart
		pairs = append(pairs, bson.M{"_id": c.ID}, c)
	}

	db, closer := dbSession.DB()
	defer closer()
	bulk := db.C(chartCollection).Bulk()

	// Upsert pairs of selectors, charts
	bulk.Upsert(pairs...)

	// Remove	charts no longer existing in index
	bulk.RemoveAll(bson.M{
		"_id": bson.M{
			"$nin": chartIDs,
		},
		"repo.name": charts[0].Repo.Name,
	})

	_, err := bulk.Run()
	return err
}

func fetchAndImportIcon(c chart) error {
	if c.Icon == "" {
		log.WithFields(log.Fields{"name": c.Name}).Info("icon not found")
		return nil
	}

	req, err := http.NewRequest("GET", c.Icon, nil)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		return err
	}

	res, err := netClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s", res.StatusCode, c.Icon)
	}

	orig, err := imaging.Decode(res.Body)
	if err != nil {
		return err
	}

	// TODO: make this configurable?
	icon := imaging.Fit(orig, 160, 160, imaging.Lanczos)

	var b bytes.Buffer
	imaging.Encode(&b, icon, imaging.PNG)

	db, closer := dbSession.DB()
	defer closer()
	return db.C(chartCollection).UpdateId(c.ID, bson.M{"$set": bson.M{"raw_icon": b.Bytes()}})
}

func fetchAndImportReadmes(c chart) {
	// TODO: This should be using a worker pool for concurrent processing later
	// TODO: we should prioritise the latest chartVersion for each chart
	for _, cv := range c.ChartVersions {
		if err := fetchAndImportReadme(c, cv); err != nil {
			log.WithFields(log.Fields{"name": c.Name, "version": cv.Version}).WithError(err).Error("failed to import readme")
		}
	}
}

func fetchAndImportReadme(c chart, cv chartVersion) error {
	chartReadmeID := fmt.Sprintf("%s/%s-%s", c.Repo.Name, c.Name, cv.Version)
	db, closer := dbSession.DB()
	defer closer()
	if err := db.C(chartReadmeCollection).FindId(chartReadmeID).One(&chartReadme{}); err == nil {
		log.WithFields(log.Fields{"name": c.Name, "version": cv.Version}).Debug("skipping existing readme")
		return nil
	}
	log.WithFields(log.Fields{"name": c.Name, "version": cv.Version}).Debug("fetching readme")

	url := chartTarballURL(c.Repo, cv)
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		return err
	}

	res, err := netClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// We read the whole chart into memory, this should be okay since the chart
	// tarball needs to be small enough to fit into a GRPC call (Tiller
	// requirement)
	gzf, err := gzip.NewReader(res.Body)
	if err != nil {
		return err
	}
	defer gzf.Close()

	tarf := tar.NewReader(gzf)

	readmeFileName := c.Name + "/README.md"
	readme, err := extractFileFromTarball(readmeFileName, tarf)
	if err != nil && !strings.Contains(err.Error(), "file not found") {
		return err
	}

	// Even if the readme doesn't exist, we create an empty entry to avoid
	// refetching for this chart version in the future
	if readme == "" {
		log.WithFields(log.Fields{"name": c.Name, "version": cv.Version}).Info("readme not found")
	}

	db.C(chartReadmeCollection).Insert(chartReadme{chartReadmeID, readme})

	return nil
}

func chartTarballURL(r repo, cv chartVersion) string {
	source := cv.URLs[0]
	if _, err := url.ParseRequestURI(source); err != nil {
		// If the chart URL is not absolute, join with repo URL. It's fine if the
		// URL we build here is invalid as we can catch this error when actually
		// making the request
		u, _ := url.Parse(r.URL)
		u.Path = path.Join(u.Path, source)
		return u.String()
	}
	return source
}

func extractFileFromTarball(filename string, tarf *tar.Reader) (string, error) {
	for {
		header, err := tarf.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if header.Name == filename {
			var b bytes.Buffer
			io.Copy(&b, tarf)
			return string(b.Bytes()), nil
		}
	}
	return "", fmt.Errorf("%s file not found", filename)
}
