/*
Copyright (c) 2018 The Helm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package foundationdb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/helm/monocular/cmd/chart-repo/common"
	"github.com/helm/monocular/cmd/chart-repo/utils"

	"github.com/arschles/assert"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var validRepoIndexYAMLBytes, _ = ioutil.ReadFile("../testdata/valid-index.yaml")
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
	// If subpath repo URL test, check that index.yaml is correctly added to the
	// subpath
	if req.URL.Host == "subpath.test" && req.URL.Path != "/subpath/index.yaml" {
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

type svgIconClient struct{}

func (h *svgIconClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.Write([]byte("foo"))
	res := w.Result()
	res.Header.Set("Content-Type", "image/svg")
	return res, nil
}

type goodTarballClient struct {
	c          common.Chart
	skipReadme bool
	skipValues bool
	skipSchema bool
}

var testChartReadme = "# readme for chart\n\nBest chart in town"
var testChartValues = "image: test"
var testChartSchema = `{"properties": {}}`

func (h *goodTarballClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	gzw := gzip.NewWriter(w)
	files := []tarballFile{{h.c.Name + "/Chart.yaml", "should be a Chart.yaml here..."}}
	if !h.skipValues {
		files = append(files, tarballFile{h.c.Name + "/values.yaml", testChartValues})
	}
	if !h.skipReadme {
		files = append(files, tarballFile{h.c.Name + "/README.md", testChartReadme})
	}
	if !h.skipSchema {
		files = append(files, tarballFile{h.c.Name + "/values.schema.json", testChartSchema})
	}
	createTestTarball(gzw, files)
	gzw.Flush()
	return w.Result(), nil
}

type authenticatedTarballClient struct {
	c common.Chart
}

func (h *authenticatedTarballClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()

	// Ensure we're sending the right Authorization header
	if !strings.Contains(req.Header.Get("Authorization"), "Bearer ThisSecretAccessTokenAuthenticatesTheClient") {
		w.WriteHeader(500)
	} else {
		gzw := gzip.NewWriter(w)
		files := []tarballFile{{h.c.Name + "/Chart.yaml", "should be a Chart.yaml here..."}}
		files = append(files, tarballFile{h.c.Name + "/values.yaml", testChartValues})
		files = append(files, tarballFile{h.c.Name + "/README.md", testChartReadme})
		files = append(files, tarballFile{h.c.Name + "/values.schema.json", testChartSchema})
		createTestTarball(gzw, files)
		gzw.Flush()
	}
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
	m := mock.Mock{}
	dbClient := NewMockClient(&m)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := syncRepo(dbClient, "test", "test", tt.repoURL, "")
			assert.ExistsErr(t, err, tt.name)
		})
	}
}

func Test_fetchRepoIndex(t *testing.T) {
	tests := []struct {
		name string
		r    common.Repo
	}{
		{"valid HTTP URL", common.Repo{URL: "http://my.examplerepo.com"}},
		{"valid HTTPS URL", common.Repo{URL: "https://my.examplerepo.com"}},
		{"valid trailing URL", common.Repo{URL: "https://my.examplerepo.com/"}},
		{"valid subpath URL", common.Repo{URL: "https://subpath.test/subpath/"}},
		{"valid URL with trailing spaces", common.Repo{URL: "https://subpath.test/subpath/  "}},
		{"valid URL with leading spaces", common.Repo{URL: "  https://subpath.test/subpath/"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			netClient = &goodHTTPClient{}
			_, err := common.FetchRepoIndex(tt.r, netClient)
			assert.NoErr(t, err)
		})
	}

	t.Run("failed request", func(t *testing.T) {
		netClient = &badHTTPClient{}
		_, err := common.FetchRepoIndex(common.Repo{URL: "https://my.examplerepo.com"}, netClient)
		assert.ExistsErr(t, err, "failed request")
	})
}

func Test_fetchRepoIndexUserAgent(t *testing.T) {
	tests := []struct {
		name              string
		version           string
		userAgentComment  string
		expectedUserAgent string
	}{
		{"default user agent", "", "", "chart-repo/devel"},
		{"custom version no app", "1.0", "", "chart-repo/1.0"},
		{"custom version and app", "1.0", "monocular/1.2", "chart-repo/1.0 (monocular/1.2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override global variables used to generate the userAgent
			if tt.version != "" {
				utils.Version = tt.version
			}

			if tt.userAgentComment != "" {
				utils.UserAgentComment = tt.userAgentComment
			}

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tt.expectedUserAgent, req.Header.Get("User-Agent"), "expected user agent")
				rw.Write([]byte(validRepoIndexYAML))
			}))
			// Close the server when test finishes
			defer server.Close()

			netClient = server.Client()

			_, err := common.FetchRepoIndex(common.Repo{URL: server.URL}, netClient)
			assert.NoErr(t, err)
		})
	}
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
			_, err := common.ParseRepoIndex([]byte(tt.repoYAML))
			assert.ExistsErr(t, err, tt.name)
		})
	}

	t.Run("valid", func(t *testing.T) {
		index, err := common.ParseRepoIndex([]byte(validRepoIndexYAML))
		assert.NoErr(t, err)
		assert.Equal(t, len(index.Entries), 2, "number of charts")
		assert.Equal(t, index.Entries["acs-engine-autoscaler"][0].GetName(), "acs-engine-autoscaler", "chart version populated")
	})
}

func Test_chartsFromIndex(t *testing.T) {
	r := common.Repo{Name: "test", URL: "http://testrepo.com"}
	index, _ := common.ParseRepoIndex([]byte(validRepoIndexYAML))
	charts := common.ChartsFromIndex(index, r)
	assert.Equal(t, len(charts), 2, "number of charts")

	indexWithDeprecated := validRepoIndexYAML + `
  deprecated-chart:
  - name: deprecated-chart
    deprecated: true`
	index2, err := common.ParseRepoIndex([]byte(indexWithDeprecated))
	assert.NoErr(t, err)
	charts = common.ChartsFromIndex(index2, r)
	assert.Equal(t, len(charts), 2, "number of charts")
}

func Test_newChart(t *testing.T) {
	r := common.Repo{Name: "test", URL: "http://testrepo.com"}
	index, _ := common.ParseRepoIndex([]byte(validRepoIndexYAML))
	c := common.NewChart(index.Entries["wordpress"], r)
	assert.Equal(t, c.Name, "wordpress", "correctly built")
	assert.Equal(t, len(c.ChartVersions), 2, "correctly built")
	assert.Equal(t, c.Description, "new description!", "takes chart fields from latest entry")
	assert.Equal(t, c.Repo, r, "repo set")
	assert.Equal(t, c.ID, "test/wordpress", "id set")
}

func Test_importCharts(t *testing.T) {
	m := &mock.Mock{}
	// Ensure Upsert func is called with some arguments
	m.On("BulkWrite", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.BulkWriteResult{}, nil)
	m.On("DeleteMany", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, nil)
	dbClient := NewMockClient(m)
	db, _ := dbClient.Database("test")
	index, _ := common.ParseRepoIndex([]byte(validRepoIndexYAML))
	charts := common.ChartsFromIndex(index, common.Repo{Name: "test", URL: "http://testrepo.com"})
	importCharts(db, "test", charts)

	m.AssertExpectations(t)
	// The BulkWrite method takes an array of WriteModels.
	// For x charts to upsert, there should be x elements.
	// Each element has a selector (filter) and an "update" - the update document to apply
	args := m.Calls[0].Arguments.Get(1).([]mongo.WriteModel)
	assert.Equal(t, len(args), len(charts), "number of charts to upsert")
	for i := 0; i < len(args); i++ {
		updateModel := args[i].(*mongo.UpdateOneModel)
		assert.Equal(t, updateModel.Filter, bson.M{"_id": "test/" + updateModel.Update.(bson.M)["name"].(string)}, "selector")
	}
}

func Test_DeleteRepo(t *testing.T) {
	m := &mock.Mock{}
	m.On("DeleteMany", mock.Anything, bson.M{
		"repo.name": "test",
	}, mock.Anything).Return(&mongo.DeleteResult{}, nil)
	dbClient := NewMockClient(m)

	err := deleteRepo(dbClient, "test", "test")
	if err != nil {
		t.Errorf("failed to delete chart repo test: %v", err)
	}
	m.AssertExpectations(t)
}

func Test_fetchAndImportIcon(t *testing.T) {
	t.Run("no icon", func(t *testing.T) {
		m := mock.Mock{}
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		c := common.Chart{ID: "test/acs-engine-autoscaler"}
		assert.NoErr(t, fetchAndImportIcon(db, c))
	})

	index, _ := common.ParseRepoIndex([]byte(validRepoIndexYAML))
	charts := common.ChartsFromIndex(index, common.Repo{Name: "test", URL: "http://testrepo.com"})

	t.Run("failed download", func(t *testing.T) {
		netClient = &badHTTPClient{}
		c := charts[0]
		m := mock.Mock{}
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		assert.Err(t, fmt.Errorf("500 %s", c.Icon), fetchAndImportIcon(db, c))
	})

	t.Run("bad icon", func(t *testing.T) {
		netClient = &badIconClient{}
		c := charts[0]
		m := mock.Mock{}
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		assert.Err(t, image.ErrFormat, fetchAndImportIcon(db, c))
	})

	t.Run("valid icon", func(t *testing.T) {
		netClient = &goodIconClient{}
		c := charts[0]
		m := mock.Mock{}
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		m.On("UpdateOne", mock.Anything, bson.M{"_id": c.ID}, bson.M{"$set": bson.M{"raw_icon": iconBytes(), "icon_content_type": "image/png"}}, mock.Anything).Return(&mongo.UpdateResult{}, nil)
		assert.NoErr(t, fetchAndImportIcon(db, c))
		m.AssertExpectations(t)
	})

	t.Run("valid SVG icon", func(t *testing.T) {
		netClient = &svgIconClient{}
		c := common.Chart{
			ID:   "foo",
			Icon: "https://foo/bar/logo.svg",
			Repo: common.Repo{},
		}
		m := mock.Mock{}
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		m.On("UpdateOne", mock.Anything, bson.M{"_id": c.ID}, bson.M{"$set": bson.M{"raw_icon": []byte("foo"), "icon_content_type": "image/svg"}}, mock.Anything).Return(&mongo.UpdateResult{}, nil)
		assert.NoErr(t, fetchAndImportIcon(db, c))
		m.AssertExpectations(t)
	})
}

func Test_fetchAndImportFiles(t *testing.T) {
	index, _ := common.ParseRepoIndex([]byte(validRepoIndexYAML))
	charts := common.ChartsFromIndex(index, common.Repo{Name: "test", URL: "http://testrepo.com", AuthorizationHeader: "Bearer ThisSecretAccessTokenAuthenticatesTheClient1s"})
	cv := charts[0].ChartVersions[0]

	t.Run("http error", func(t *testing.T) {
		m := mock.Mock{}
		m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments)
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		netClient = &badHTTPClient{}
		assert.Err(t, io.EOF, fetchAndImportFiles(db, charts[0].Name, charts[0].Repo, cv))
	})

	t.Run("file not found", func(t *testing.T) {
		netClient = &goodTarballClient{c: charts[0], skipValues: true, skipReadme: true, skipSchema: true}
		m := mock.Mock{}
		m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments)
		chartFilesID := fmt.Sprintf("%s/%s-%s", charts[0].Repo.Name, charts[0].Name, cv.Version)
		chartFiles := common.ChartFiles{ID: chartFilesID, Readme: "", Values: "", Schema: "", Repo: charts[0].Repo, Digest: cv.Digest}
		chartBSON, _ := bson.Marshal(&chartFiles)
		var doc bson.M
		bson.Unmarshal(chartBSON, &doc)
		delete(doc, "_id")
		update := bson.M{"$set": doc}
		m.On("UpdateOne", mock.Anything, bson.M{"_id": chartFilesID}, update, mock.Anything).Return(&mongo.UpdateResult{}, nil)
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		err := fetchAndImportFiles(db, charts[0].Name, charts[0].Repo, cv)
		assert.NoErr(t, err)
		m.AssertExpectations(t)
	})

	t.Run("authenticated request", func(t *testing.T) {
		netClient = &authenticatedTarballClient{c: charts[0]}
		m := mock.Mock{}
		m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments)
		chartFilesID := fmt.Sprintf("%s/%s-%s", charts[0].Repo.Name, charts[0].Name, cv.Version)
		chartFiles := common.ChartFiles{ID: chartFilesID, Readme: testChartReadme, Values: testChartValues, Schema: testChartSchema, Repo: charts[0].Repo, Digest: cv.Digest}
		chartBSON, _ := bson.Marshal(&chartFiles)
		var doc bson.M
		bson.Unmarshal(chartBSON, &doc)
		delete(doc, "_id")
		update := bson.M{"$set": doc}
		m.On("UpdateOne", mock.Anything, bson.M{"_id": chartFilesID}, update, mock.Anything).Return(&mongo.UpdateResult{}, nil)
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		err := fetchAndImportFiles(db, charts[0].Name, charts[0].Repo, cv)
		assert.NoErr(t, err)
		m.AssertExpectations(t)
	})

	t.Run("valid tarball", func(t *testing.T) {
		netClient = &goodTarballClient{c: charts[0]}
		m := mock.Mock{}
		m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments)
		chartFilesID := fmt.Sprintf("%s/%s-%s", charts[0].Repo.Name, charts[0].Name, cv.Version)
		chartFiles := common.ChartFiles{ID: chartFilesID, Readme: testChartReadme, Values: testChartValues, Schema: testChartSchema, Repo: charts[0].Repo, Digest: cv.Digest}
		chartBSON, _ := bson.Marshal(&chartFiles)
		var doc bson.M
		bson.Unmarshal(chartBSON, &doc)
		delete(doc, "_id")
		update := bson.M{"$set": doc}
		m.On("UpdateOne", mock.Anything, bson.M{"_id": chartFilesID}, update, mock.Anything).Return(&mongo.UpdateResult{}, nil)
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		err := fetchAndImportFiles(db, charts[0].Name, charts[0].Repo, cv)
		assert.NoErr(t, err)
		m.AssertExpectations(t)
	})

	t.Run("file exists", func(t *testing.T) {
		m := mock.Mock{}
		// don't return an error when checking if files already exists
		m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		dbClient := NewMockClient(&m)
		db, _ := dbClient.Database("test")
		err := fetchAndImportFiles(db, charts[0].Name, charts[0].Repo, cv)
		assert.NoErr(t, err)
		m.AssertNotCalled(t, "UpdateOne", mock.Anything, mock.Anything, mock.Anything)
	})
}

func Test_chartTarballURL(t *testing.T) {
	r := common.Repo{Name: "test", URL: "http://testrepo.com"}
	tests := []struct {
		name   string
		cv     common.ChartVersion
		wanted string
	}{
		{"absolute url", common.ChartVersion{URLs: []string{"http://testrepo.com/wordpress-0.1.0.tgz"}}, "http://testrepo.com/wordpress-0.1.0.tgz"},
		{"relative url", common.ChartVersion{URLs: []string{"wordpress-0.1.0.tgz"}}, "http://testrepo.com/wordpress-0.1.0.tgz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, common.ChartTarballURL(r, tt.cv), tt.wanted, "url")
		})
	}
}

func Test_extractFilesFromTarball(t *testing.T) {
	tests := []struct {
		name     string
		files    []tarballFile
		filename string
		want     string
	}{
		{"file", []tarballFile{{"file.txt", "best file ever"}}, "file.txt", "best file ever"},
		{"multiple file tarball", []tarballFile{{"file.txt", "best file ever"}, {"file2.txt", "worst file ever"}}, "file2.txt", "worst file ever"},
		{"file in dir", []tarballFile{{"file.txt", "best file ever"}, {"test/file2.txt", "worst file ever"}}, "test/file2.txt", "worst file ever"},
		{"filename ignore case", []tarballFile{{"Readme.md", "# readme for chart"}, {"values.yaml", "key: value"}}, "README.md", "# readme for chart"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			createTestTarball(&b, tt.files)
			r := bytes.NewReader(b.Bytes())
			tarf := tar.NewReader(r)
			files, err := common.ExtractFilesFromTarball([]string{tt.filename}, tarf)
			assert.NoErr(t, err)
			assert.Equal(t, files[tt.filename], tt.want, "file body")
		})
	}

	t.Run("extract multiple files", func(t *testing.T) {
		var b bytes.Buffer
		tFiles := []tarballFile{{"file.txt", "best file ever"}, {"file2.txt", "worst file ever"}}
		createTestTarball(&b, tFiles)
		r := bytes.NewReader(b.Bytes())
		tarf := tar.NewReader(r)
		files, err := common.ExtractFilesFromTarball([]string{tFiles[0].Name, tFiles[1].Name}, tarf)
		assert.NoErr(t, err)
		assert.Equal(t, len(files), 2, "matches")
		for _, f := range tFiles {
			assert.Equal(t, files[f.Name], f.Body, "file body")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		var b bytes.Buffer
		createTestTarball(&b, []tarballFile{{"file.txt", "best file ever"}})
		r := bytes.NewReader(b.Bytes())
		tarf := tar.NewReader(r)
		name := "file2.txt"
		files, err := common.ExtractFilesFromTarball([]string{name}, tarf)
		assert.NoErr(t, err)
		assert.Equal(t, files[name], "", "file body")
	})

	t.Run("not a tarball", func(t *testing.T) {
		b := make([]byte, 4)
		rand.Read(b)
		r := bytes.NewReader(b)
		tarf := tar.NewReader(r)
		files, err := common.ExtractFilesFromTarball([]string{"file2.txt"}, tarf)
		assert.Err(t, io.ErrUnexpectedEOF, err)
		assert.Equal(t, len(files), 0, "file body")
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

func Test_initNetClient(t *testing.T) {
	// Test env
	otherDir, err := ioutil.TempDir("", "ca-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(otherDir)

	// Create cert
	caCert := `-----BEGIN CERTIFICATE-----
MIIC6jCCAdKgAwIBAgIUKVfzA7lfBgSYP8enCVhlm0ql5YwwDQYJKoZIhvcNAQEL
BQAwDTELMAkGA1UEAxMCQ0EwHhcNMTgxMjEyMTQxNzAwWhcNMjMxMjExMTQxNzAw
WjANMQswCQYDVQQDEwJDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
ALZU3fsAgvoUuLSHr24apslaYyuX84wGoZQmtFtQ+A3DF9KL/2nn3yZ6qJPkH0TF
sbObEQRNi+P6vQ3nI/dSNMX5PzMBP2CB6L7zEXzZQEHtAK0Bzva5CKEBGX7OfIKl
aBvs+dzKVJBdb+Oh0maacMwa4QbcD6ejzF90jUbaO65lpQpcL7KQdppKOGNclRaA
hQTV2VsxrV4hH7K9btaTTxso+8W6p8v6X9vf40Ywx72p+SKnGh+FCrOp1gYLBLwo
4SM0OUQHRvqUlj0XhZk5pW0dMRwHcoz1S2GmE5bj4edr4j+zGzGxa2wRGKvM0OCn
Do84AVszTFPmUf+mCl4pJNECAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgEGMA8GA1Ud
EwEB/wQFMAMBAf8wHQYDVR0OBBYEFI5l5k+MEhrbOQ29dOW1qJhI0yKaMA0GCSqG
SIb3DQEBCwUAA4IBAQByDebUOKzn6jfmXlW62vm09V+ipqId01wm21G9XMtMEqhc
xtun6YwQeTuGPtdepWG+NXuSsiX/HNAHeaumJaaljHhdKDisnMQ0CTnNsu8NPkAl
9iMEB3iXLWkb7+HgfPJAHZVGcMqMxNEMZYHB1Fh0G2Ne376X94+GYJ08qR2C8rUP
BShhMSktB578h4GtPIWSjPhDUWg1fGe7sewR+GPyuL9859hOD0wGm9tUixBKloCu
b90fhqZZ3FqZD7W1qJGKvz/8geqi0noip+uq/dokK1jarRkOVEJP+EvXkHo0tIuc
h251U/Daz6NiQBM9AxyAw6EHm8XAZBvCuebfzyrT
-----END CERTIFICATE-----`
	otherCA := path.Join(otherDir, "ca.crt")
	err = ioutil.WriteFile(otherCA, []byte(caCert), 0644)
	if err != nil {
		t.Error(err)
	}

	_, err = common.InitNetClient(otherCA, defaultTimeoutSeconds)
	if err != nil {
		t.Error(err)
	}
}

var emptyRepoIndexYAMLBytes, _ = ioutil.ReadFile("testdata/empty-repo-index.yaml")
var emptyRepoIndexYAML = string(emptyRepoIndexYAMLBytes)

type emptyChartRepoHTTPClient struct{}

func (h *emptyChartRepoHTTPClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.Write([]byte(emptyRepoIndexYAML))
	return w.Result(), nil
}

func Test_emptyChartRepo(t *testing.T) {
	netClient = &emptyChartRepoHTTPClient{}
	m := mock.Mock{}
	dbClient := NewMockClient(&m)
	//Expect a call to test the DB readiness
	m.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{}, nil)
	m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	err := syncRepo(dbClient, "test", "testRepo", "https://my.examplerepo.com", "")
	assert.ExistsErr(t, err, "Failed Request")
}

func Test_getSha256(t *testing.T) {
	sha, err := common.GetSha256([]byte("this is a test"))
	assert.Equal(t, err, nil, "Unable to get sha")
	assert.Equal(t, sha, "2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c", "Unable to get sha")
}

func Test_repoAlreadyProcessed(t *testing.T) {
	tests := []struct {
		name            string
		checksum        string
		mockedLastCheck common.RepoCheck
		processed       bool
	}{
		{"not processed yet", "bar", common.RepoCheck{}, false},
		{"already processed", "bar", common.RepoCheck{Checksum: "bar"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mock.Mock{}
			repo := &common.RepoCheck{}
			m.On("FindOne", mock.Anything, mock.Anything, repo, mock.Anything).Run(func(args mock.Arguments) {
				*args.Get(2).(*common.RepoCheck) = tt.mockedLastCheck
			}).Return(nil)
			dbClient := NewMockClient(&m)
			db, _ := dbClient.Database("test")
			res := repoAlreadyProcessed(db, "", tt.checksum)
			if res != tt.processed {
				t.Errorf("Expected alreadyProcessed to be %v got %v", tt.processed, res)
			}
		})
	}
}

func Test_updateLastCheck(t *testing.T) {
	m := mock.Mock{}
	repoName := "foo"
	checksum := "bar"
	now := time.Now()
	selector := bson.M{"_id": repoName}
	m.On("UpdateOne", mock.Anything, selector, bson.M{"$set": bson.M{"last_update": now, "checksum": checksum}}, mock.Anything).Return(&mongo.UpdateResult{}, nil)
	dbClient := NewMockClient(&m)
	db, _ := dbClient.Database("test")
	err := updateLastCheck(db, repoName, checksum, now)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if len(m.Calls) != 1 {
		t.Errorf("Expected one call got %d", len(m.Calls))
	}
}
