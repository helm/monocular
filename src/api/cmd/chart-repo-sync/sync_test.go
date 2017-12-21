package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/arschles/assert"
	"github.com/disintegration/imaging"
	"github.com/kubeapps/common/datastore/mockstore"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

var validRepoIndexYAMLBytes, _ = ioutil.ReadFile("testdata/valid-index.yaml")
var validRepoIndexYAML = string(validRepoIndexYAMLBytes)

var invalidRepoIndexYAML = "invalid"

type badHTTPClient struct{}

func (h *badHTTPClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(500)
	return w.Result(), nil
}

type goodHTTPClient struct{}

func (h *goodHTTPClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	// Don't accept trailing slashes
	if strings.HasPrefix(req.URL.Path, "//") {
		w.WriteHeader(500)
	}
	// Ensure we're sending the right User-Agent
	if !strings.Contains(req.Header.Get("User-Agent"), "chart-repo-sync") {
		w.WriteHeader(500)
	}
	w.Write([]byte(validRepoIndexYAML))
	return w.Result(), nil
}

type badIconClient struct{}

func (h *badIconClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.Write([]byte("not-an-image"))
	return w.Result(), nil
}

type goodIconClient struct{}

func iconBytes() []byte {
	var b bytes.Buffer
	img := imaging.New(1, 1, color.White)
	imaging.Encode(&b, img, imaging.PNG)
	return b.Bytes()
}

func (h *goodIconClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.Write(iconBytes())
	return w.Result(), nil
}

type goodTarballClient struct {
	c          chart
	skipReadme bool
}

var testChartReadme = "# readme for chart\n\nBest chart in town"

func (h *goodTarballClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	gzw := gzip.NewWriter(w)
	files := []tarballFile{{h.c.Name + "/Chart.yaml", "should be a Chart.yaml here..."}}
	if !h.skipReadme {
		files = append(files, tarballFile{h.c.Name + "/README.md", testChartReadme})
	}
	createTestTarball(gzw, files)
	gzw.Flush()
	return w.Result(), nil
}

func Test_syncURLInvalidity(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
	}{
		{"invalid URL", "not-a-url"},
		{"invalid URL", "https//google.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := syncRepo("test", tt.repoURL)
			assert.ExistsErr(t, err, tt.name)
		})
	}
}

func Test_fetchRepoIndex(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
	}{
		{"valid HTTP URL", "http://my.examplerepo.com"},
		{"valid HTTPS URL", "https://my.examplerepo.com"},
		{"valid trailing URL", "https://my.examplerepo.com/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			netClient = &goodHTTPClient{}
			url, _ := url.ParseRequestURI(tt.repoURL)
			_, err := fetchRepoIndex(url)
			assert.NoErr(t, err)
		})
	}

	t.Run("failed request", func(t *testing.T) {
		netClient = &badHTTPClient{}
		url, _ := url.ParseRequestURI("https://my.examplerepo.com")
		_, err := fetchRepoIndex(url)
		assert.ExistsErr(t, err, "failed request")
	})
}

func Test_parseRepoIndex(t *testing.T) {
	tests := []struct {
		name     string
		repoYAML string
	}{
		{"invalid", "invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseRepoIndex([]byte(tt.repoYAML))
			assert.ExistsErr(t, err, tt.name)
		})
	}

	t.Run("valid", func(t *testing.T) {
		index, err := parseRepoIndex([]byte(validRepoIndexYAML))
		assert.NoErr(t, err)
		assert.Equal(t, len(index.Entries), 2, "number of charts")
		assert.Equal(t, index.Entries["acs-engine-autoscaler"][0].GetName(), "acs-engine-autoscaler", "chart version populated")
	})
}

func Test_chartsFromIndex(t *testing.T) {
	r := repo{Name: "test", URL: "http://testrepo.com"}
	index, _ := parseRepoIndex([]byte(validRepoIndexYAML))
	charts := chartsFromIndex(index, r)
	assert.Equal(t, len(charts), 2, "number of charts")

	indexWithDeprecated := validRepoIndexYAML + `
  deprecated-chart:
  - name: deprecated-chart
    deprecated: true`
	index2, err := parseRepoIndex([]byte(indexWithDeprecated))
	assert.NoErr(t, err)
	charts = chartsFromIndex(index2, r)
	assert.Equal(t, len(charts), 2, "number of charts")
}

func Test_newChart(t *testing.T) {
	r := repo{Name: "test", URL: "http://testrepo.com"}
	index, _ := parseRepoIndex([]byte(validRepoIndexYAML))
	c := newChart(index.Entries["wordpress"], r)
	assert.Equal(t, c.Name, "wordpress", "correctly built")
	assert.Equal(t, len(c.ChartVersions), 2, "correctly built")
	assert.Equal(t, c.Description, "new description!", "takes chart fields from latest entry")
	assert.Equal(t, c.Repo, r, "repo set")
	assert.Equal(t, c.ID, "test/wordpress", "id set")
}

func Test_importCharts(t *testing.T) {
	m := &mock.Mock{}
	// Ensure Upsert func is called with some arguments
	m.On("Upsert", mock.Anything)
	m.On("RemoveAll", mock.Anything)
	dbSession = mockstore.NewMockSession(m)
	index, _ := parseRepoIndex([]byte(validRepoIndexYAML))
	charts := chartsFromIndex(index, repo{Name: "test", URL: "http://testrepo.com"})
	importCharts(charts)

	m.AssertExpectations(t)
	// The Bulk Upsert method takes an array that consists of a selector followed by an interface to upsert.
	// So for x charts to upsert, there should be x*2 elements (each chart has it's own selector)
	// e.g. [selector1, chart1, selector2, chart2, ...]
	args := m.Calls[0].Arguments.Get(0).([]interface{})
	assert.Equal(t, len(args), len(charts)*2, "number of selector, chart pairs to upsert")
	for i := 0; i < len(args); i += 2 {
		c := args[i+1].(chart)
		assert.Equal(t, args[i], bson.M{"_id": "test/" + c.Name}, "selector")
	}
}

func Test_fetchAndImportIcon(t *testing.T) {
	t.Run("no icon", func(t *testing.T) {
		c := chart{ID: "test/acs-engine-autoscaler"}
		assert.NoErr(t, fetchAndImportIcon(c))
	})

	index, _ := parseRepoIndex([]byte(validRepoIndexYAML))
	charts := chartsFromIndex(index, repo{Name: "test", URL: "http://testrepo.com"})

	t.Run("failed download", func(t *testing.T) {
		netClient = &badHTTPClient{}
		c := charts[0]
		assert.Err(t, fmt.Errorf("500 %s", c.Icon), fetchAndImportIcon(c))
	})

	t.Run("bad icon", func(t *testing.T) {
		netClient = &badIconClient{}
		c := charts[0]
		assert.Err(t, image.ErrFormat, fetchAndImportIcon(c))
	})

	t.Run("valid icon", func(t *testing.T) {
		m := mock.Mock{}
		dbSession = mockstore.NewMockSession(&m)
		netClient = &goodIconClient{}
		c := charts[0]
		m.On("UpdateId", c.ID, bson.M{"$set": bson.M{"raw_icon": iconBytes()}}).Return(nil)
		assert.NoErr(t, fetchAndImportIcon(c))
		m.AssertExpectations(t)
	})
}

func Test_fetchAndImportReadme(t *testing.T) {
	index, _ := parseRepoIndex([]byte(validRepoIndexYAML))
	charts := chartsFromIndex(index, repo{Name: "test", URL: "http://testrepo.com"})
	cv := charts[0].ChartVersions[0]

	t.Run("http error", func(t *testing.T) {
		m := mock.Mock{}
		m.On("One", mock.Anything).Return(errors.New("return an error when checking if readme already exists to force fetching"))
		dbSession = mockstore.NewMockSession(&m)
		netClient = &badHTTPClient{}
		assert.Err(t, io.EOF, fetchAndImportReadme(charts[0].Name, charts[0].Repo, cv))
	})

	t.Run("readme not found", func(t *testing.T) {
		netClient = &goodTarballClient{c: charts[0], skipReadme: true}
		m := mock.Mock{}
		m.On("One", mock.Anything).Return(errors.New("return an error when checking if readme already exists to force fetching"))
		m.On("Insert", chartReadme{fmt.Sprintf("%s/%s-%s", charts[0].Repo.Name, charts[0].Name, cv.Version), ""})
		dbSession = mockstore.NewMockSession(&m)
		err := fetchAndImportReadme(charts[0].Name, charts[0].Repo, cv)
		assert.NoErr(t, err)
		m.AssertExpectations(t)
	})

	t.Run("valid tarball", func(t *testing.T) {
		netClient = &goodTarballClient{c: charts[0]}
		m := mock.Mock{}
		m.On("One", mock.Anything).Return(errors.New("return an error when checking if readme already exists to force fetching"))
		m.On("Insert", chartReadme{fmt.Sprintf("%s/%s-%s", charts[0].Repo.Name, charts[0].Name, cv.Version), testChartReadme})
		dbSession = mockstore.NewMockSession(&m)
		err := fetchAndImportReadme(charts[0].Name, charts[0].Repo, cv)
		assert.NoErr(t, err)
		m.AssertExpectations(t)
	})

	t.Run("readme exists", func(t *testing.T) {
		m := mock.Mock{}
		// don't return an error when checking if readme already exists
		m.On("One", mock.Anything).Return(nil)
		dbSession = mockstore.NewMockSession(&m)
		err := fetchAndImportReadme(charts[0].Name, charts[0].Repo, cv)
		assert.NoErr(t, err)
		m.AssertNotCalled(t, "Insert", mock.Anything)
	})
}

func Test_chartTarballURL(t *testing.T) {
	r := repo{Name: "test", URL: "http://testrepo.com"}
	tests := []struct {
		name   string
		cv     chartVersion
		wanted string
	}{
		{"absolute url", chartVersion{URLs: []string{"http://testrepo.com/wordpress-0.1.0.tgz"}}, "http://testrepo.com/wordpress-0.1.0.tgz"},
		{"relative url", chartVersion{URLs: []string{"wordpress-0.1.0.tgz"}}, "http://testrepo.com/wordpress-0.1.0.tgz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, chartTarballURL(r, tt.cv), tt.wanted, "url")
		})
	}
}

func Test_extractFileFromTarball(t *testing.T) {
	tests := []struct {
		name     string
		files    []tarballFile
		filename string
		want     string
	}{
		{"file", []tarballFile{{"file.txt", "best file ever"}}, "file.txt", "best file ever"},
		{"multiple files", []tarballFile{{"file.txt", "best file ever"}, {"file2.txt", "worst file ever"}}, "file2.txt", "worst file ever"},
		{"file in dir", []tarballFile{{"file.txt", "best file ever"}, {"test/file2.txt", "worst file ever"}}, "test/file2.txt", "worst file ever"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			createTestTarball(&b, tt.files)
			r := bytes.NewReader(b.Bytes())
			tarf := tar.NewReader(r)
			readme, err := extractFileFromTarball(tt.filename, tarf)
			assert.NoErr(t, err)
			assert.Equal(t, readme, tt.want, "file body")
		})
	}

	t.Run("file not found", func(t *testing.T) {
		var b bytes.Buffer
		createTestTarball(&b, []tarballFile{{"file.txt", "best file ever"}})
		r := bytes.NewReader(b.Bytes())
		tarf := tar.NewReader(r)
		readme, err := extractFileFromTarball("file2.txt", tarf)
		assert.Err(t, errors.New("file2.txt file not found"), err)
		assert.Equal(t, readme, "", "file body")
	})

	t.Run("not a tarball", func(t *testing.T) {
		b := make([]byte, 4)
		rand.Read(b)
		r := bytes.NewReader(b)
		tarf := tar.NewReader(r)
		readme, err := extractFileFromTarball("file2.txt", tarf)
		assert.Err(t, io.ErrUnexpectedEOF, err)
		assert.Equal(t, readme, "", "file body")
	})
}

type tarballFile struct {
	Name, Body string
}

func createTestTarball(w io.Writer, files []tarballFile) {
	// Create a new tar archive.
	tarw := tar.NewWriter(w)

	// Add files to the archive.
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tarw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tarw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}
	// Make sure to check the error on Close.
	if err := tarw.Close(); err != nil {
		log.Fatal(err)
	}
}
