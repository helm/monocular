package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/jinzhu/copier"
	"github.com/kubeapps/common/datastore"
	"gopkg.in/mgo.v2/bson"
	helmrepo "k8s.io/helm/pkg/repo"
)

const (
	chartCollection = "charts"
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
	flag.Parse()

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

// Takes an entry from the index and constructs a database representation of the object
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
