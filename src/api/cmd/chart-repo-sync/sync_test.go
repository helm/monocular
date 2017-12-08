package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/disintegration/imaging"
	"github.com/kubeapps/common/datastore/mockstore"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

var validRepoIndexYAML = `
entries:
  acs-engine-autoscaler:
  - apiVersion: v1
    appVersion: 2.1.1
    created: 2017-12-06T18:48:59.568323124Z
    description: Scales worker nodes within agent pools
    digest: 39e66eb53c310529bd9dd19776f8ba662e063a4ebd51fc5ec9f2267e2e073e3e
    icon: https://github.com/kubernetes/kubernetes/blob/master/logo/logo.png
    maintainers:
    - email: ritazh@microsoft.com
      name: ritazh
    - email: wibuch@microsoft.com
      name: wbuchwalter
    name: acs-engine-autoscaler
    sources:
    - https://github.com/wbuchwalter/Kubernetes-acs-engine-autoscaler
    urls:
    - https://kubernetes-charts.storage.googleapis.com/acs-engine-autoscaler-2.1.1.tgz
    version: 2.1.1
  wordpress:
  - appVersion: 4.9.1
    created: 2017-12-06T18:48:59.644981487Z
    description: new description!
    digest: 74889e60a35dcffa4686f88bb23de863fed2b6e63a69b1f4858dde37c301885c
    engine: gotpl
    home: http://www.wordpress.com/
    icon: https://bitnami.com/assets/stacks/wordpress/img/wordpress-stack-220x234.png
    keywords:
    - wordpress
    - cms
    - blog
    - http
    - web
    - application
    - php
    maintainers:
    - email: containers@bitnami.com
      name: bitnami-bot
    name: wordpress
    sources:
    - https://github.com/bitnami/bitnami-docker-wordpress
    urls:
    - https://kubernetes-charts.storage.googleapis.com/wordpress-0.7.5.tgz
    version: 0.7.5
  - appVersion: 4.9.0
    created: 2017-12-01T11:49:00.136950565Z
    description: Web publishing platform for building blogs and websites.
    digest: a69139ef3008eeb11ca60261ec2ded61e84ce7db32bb3626056e84bcff7ec270
    engine: gotpl
    home: http://www.wordpress.com/
    icon: https://bitnami.com/assets/stacks/wordpress/img/wordpress-stack-220x234.png
    keywords:
    - wordpress
    - cms
    - cms
    - blog
    - http
    - web
    - application
    - php
    maintainers:
    - email: containers@bitnami.com
      name: bitnami-bot
    name: wordpress
    sources:
    - https://github.com/bitnami/bitnami-docker-wordpress
    urls:
    - https://kubernetes-charts.storage.googleapis.com/wordpress-0.7.4.tgz
    version: 0.7.4
`

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
			err := sync("test", tt.repoURL)
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
	})
}
