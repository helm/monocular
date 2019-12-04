/*
Copyright (c) 2019

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
	"bytes"
	"encoding/json"
	"image/color"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/helm/monocular/cmd/chartsvc/foundationdb/datastore"
	"github.com/helm/monocular/cmd/chartsvc/models"
	"github.com/helm/monocular/cmd/chartsvc/utils"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var cc count

const testChartReadme = "# Quickstart\n\n```bash\nhelm install my-repo/my-chart\n```"
const testChartValues = "image:\n  registry: docker.io\n  repository: my-repo/my-chart\n  tag: 0.1.0"

const fdbURL = "mongodb://fdb-service/27016"
const fDB = "charts"

func iconBytes() []byte {
	var b bytes.Buffer
	img := imaging.New(1, 1, color.White)
	imaging.Encode(&b, img, imaging.PNG)
	return b.Bytes()
}

func Test_chartAttributes(t *testing.T) {
	tests := []struct {
		name  string
		chart models.Chart
	}{
		{"chart has no icon", models.Chart{
			ID: "stable/wordpress",
		}},
		{"chart has a icon", models.Chart{
			ID: "repo/mychart", RawIcon: iconBytes(), IconContentType: "image/svg",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := chartAttributes(tt.chart)
			assert.Equal(t, tt.chart.ID, c.ID)
			assert.Equal(t, tt.chart.RawIcon, c.RawIcon)
			if len(tt.chart.RawIcon) == 0 {
				assert.Equal(t, len(c.Icon), 0, "icon url should be undefined")
			} else {
				assert.Equal(t, c.Icon, pathPrefix+"/assets/"+tt.chart.ID+"/logo", "the icon url should be the same")
				assert.Equal(t, c.IconContentType, tt.chart.IconContentType, "the icon content type should be the same")
			}
		})
	}
}

func Test_chartVersionAttributes(t *testing.T) {
	tests := []struct {
		name  string
		chart models.Chart
	}{
		{"my-chart", models.Chart{
			ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := chartVersionAttributes(tt.chart.ID, tt.chart.ChartVersions[0])
			assert.Equal(t, cv.Version, tt.chart.ChartVersions[0].Version, "version string should be the same")
			assert.Equal(t, cv.Readme, pathPrefix+"/assets/"+tt.chart.ID+"/versions/"+tt.chart.ChartVersions[0].Version+"/README.md", "README.md resource path should be the same")
			assert.Equal(t, cv.Values, pathPrefix+"/assets/"+tt.chart.ID+"/versions/"+tt.chart.ChartVersions[0].Version+"/values.yaml", "values.yaml resource path should be the same")
		})
	}
}

func Test_newChartResponse(t *testing.T) {
	tests := []struct {
		name  string
		chart models.Chart
	}{
		{"chart has only one version", models.Chart{
			ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "1.2.3"}}},
		},
		{"chart has many versions", models.Chart{
			ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.2"}, {Version: "0.1.0"}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cResponse := newChartResponse(&tt.chart)
			assert.Equal(t, cResponse.Type, "chart", "response type is chart")
			assert.Equal(t, cResponse.ID, tt.chart.ID, "chart ID should be the same")
			assert.Equal(t, cResponse.Relationships["latestChartVersion"].Data.(models.ChartVersion).Version, tt.chart.ChartVersions[0].Version, "latestChartVersion should match version at index 0")
			assert.Equal(t, cResponse.Links.(utils.SelfLink).Self, pathPrefix+"/charts/"+tt.chart.ID, "self link should be the same")
			assert.Equal(t, len(cResponse.Attributes.(models.Chart).ChartVersions), len(tt.chart.ChartVersions), "number of chart versions in the response should be the same")
		})
	}
}

func Test_newChartListResponse(t *testing.T) {
	tests := []struct {
		name   string
		input  []*models.Chart
		result []*models.Chart
	}{
		{"no charts", []*models.Chart{}, []*models.Chart{}},
		{"has one chart", []*models.Chart{
			{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}, []*models.Chart{
			{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}},
		{"has two charts", []*models.Chart{
			{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
			{ID: "stable/wordpress", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}, {Version: "1.2.2", Digest: "12345"}}},
		}, []*models.Chart{
			{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
			{ID: "stable/wordpress", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}, {Version: "1.2.2", Digest: "12345"}}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clResponse := newChartListResponse(tt.input)
			assert.Equal(t, len(clResponse), len(tt.result), "number of charts in response should be the same")
			for i := range tt.result {
				assert.Equal(t, clResponse[i].Type, "chart", "response type is chart")
				assert.Equal(t, clResponse[i].ID, tt.result[i].ID, "chart ID should be the same")
				assert.Equal(t, clResponse[i].Relationships["latestChartVersion"].Data.(models.ChartVersion).Version, tt.result[i].ChartVersions[0].Version, "latestChartVersion should match version at index 0")
				assert.Equal(t, clResponse[i].Links.(utils.SelfLink).Self, pathPrefix+"/charts/"+tt.result[i].ID, "self link should be the same")
				assert.Equal(t, len(clResponse[i].Attributes.(models.Chart).ChartVersions), len(tt.result[i].ChartVersions), "number of chart versions in the response should be the same")
			}
		})
	}
}

func Test_newChartVersionResponse(t *testing.T) {
	tests := []struct {
		name  string
		chart models.Chart
	}{
		{"my-chart", models.Chart{
			ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}, {Version: "0.2.3"}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.chart.ChartVersions {
				cvResponse := newChartVersionResponse(&tt.chart, tt.chart.ChartVersions[i])
				assert.Equal(t, cvResponse.Type, "chartVersion", "response type is chartVersion")
				assert.Equal(t, cvResponse.ID, tt.chart.ID+"-"+tt.chart.ChartVersions[i].Version, "reponse id should have chart version suffix")
				assert.Equal(t, cvResponse.Links.(interface{}).(utils.SelfLink).Self, pathPrefix+"/charts/"+tt.chart.ID+"/versions/"+tt.chart.ChartVersions[i].Version, "self link should be the same")
				assert.Equal(t, cvResponse.Attributes.(models.ChartVersion).Version, tt.chart.ChartVersions[i].Version, "chart version in the response should be the same")
				assert.Equal(t, cvResponse.Relationships["chart"].Data.(interface{}).(models.Chart), tt.chart, "chart in relatioship matches")
			}
		})
	}
}

func Test_newChartVersionListResponse(t *testing.T) {
	tests := []struct {
		name  string
		chart models.Chart
	}{
		{"chart has no versions", models.Chart{
			ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{},
		}},
		{"chart has one version", models.Chart{
			ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1"}},
		}},
		{"chart has many versions", models.Chart{
			ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1"}, {Version: "0.0.2"}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cvListResponse := newChartVersionListResponse(&tt.chart)
			assert.Equal(t, len(cvListResponse), len(tt.chart.ChartVersions), "number of chart versions in response should be the same")
			for i := range tt.chart.ChartVersions {
				assert.Equal(t, cvListResponse[i].Type, "chartVersion", "response type is chartVersion")
				assert.Equal(t, cvListResponse[i].ID, tt.chart.ID+"-"+tt.chart.ChartVersions[i].Version, "reponse id should have chart version suffix")
				assert.Equal(t, cvListResponse[i].Links.(interface{}).(utils.SelfLink).Self, pathPrefix+"/charts/"+tt.chart.ID+"/versions/"+tt.chart.ChartVersions[i].Version, "self link should be the same")
				assert.Equal(t, cvListResponse[i].Attributes.(models.ChartVersion).Version, tt.chart.ChartVersions[i].Version, "chart version in the response should be the same")
			}
		})
	}
}

func Test_listCharts(t *testing.T) {
	pageSize := 2
	tests := []struct {
		name            string
		query           string
		dbQueryResult   []*models.Chart
		chartListResult []*models.Chart
		meta            utils.Meta
	}{
		{"no charts", "", []*models.Chart{}, []*models.Chart{}, utils.Meta{TotalPages: 1}},
		{"one chart", "", []*models.Chart{
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}, []*models.Chart{
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}, utils.Meta{TotalPages: 1}},
		{"two charts", "", []*models.Chart{
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
			{ID: "stable/dokuwiki", Name: "dokuwiki", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}, {Version: "1.2.2", Digest: "12345"}}},
		}, []*models.Chart{
			{ID: "stable/dokuwiki", Name: "dokuwiki", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}, {Version: "1.2.2", Digest: "12345"}}},
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}, utils.Meta{TotalPages: 1}},
		// Pagination tests
		{"four charts with pagination", "?size=" + strconv.Itoa(pageSize), []*models.Chart{
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
			{ID: "stable/dokuwiki", Name: "dokuwiki", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}}},
			{ID: "stable/drupal", Name: "drupal", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "12345"}}},
			{ID: "stable/wordpress", Name: "wordpress", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "123456"}}},
		}, []*models.Chart{
			{ID: "stable/dokuwiki", Name: "dokuwiki", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}}},
			{ID: "stable/drupal", Name: "drupal", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "12345"}}},
		}, utils.Meta{TotalPages: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")
			var chartsList []*models.Chart

			m.On("Find", mock.Anything, mock.Anything, &chartsList, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				*args.Get(2).(*[]*models.Chart) = tt.dbQueryResult
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/charts"+tt.query, nil)
			ListCharts(w, req)

			m.AssertExpectations(t)
			assert.Equal(t, http.StatusOK, w.Code)

			var b utils.BodyAPIListResponse
			json.NewDecoder(w.Body).Decode(&b)
			if b.Data == nil {
				t.Fatal("chart list shouldn't be null")
			}
			data := *b.Data

			assert.Len(t, data, len(tt.chartListResult))

			for i, resp := range data {
				assert.Equal(t, resp.ID, tt.chartListResult[i].ID, "chart id in the response should be the same")
				assert.Equal(t, resp.Type, "chart", "response type is chart")
				assert.Equal(t, resp.Links.(map[string]interface{})["self"], pathPrefix+"/charts/"+tt.chartListResult[i].ID, "self link should be the same")
				assert.Equal(t, resp.Relationships["latestChartVersion"].Data.(map[string]interface{})["version"], tt.chartListResult[i].ChartVersions[0].Version, "latestChartVersion should match version at index 0")
			}
			assert.Equal(t, b.Meta, tt.meta, "response meta should be the same")
		})
	}
}

func Test_listRepoCharts(t *testing.T) {
	pageSize := 2
	tests := []struct {
		name            string
		repo            string
		query           string
		dbQueryResult   []*models.Chart
		chartListResult []*models.Chart
		meta            utils.Meta
	}{
		{"repo has no charts", "my-repo", "", []*models.Chart{}, []*models.Chart{}, utils.Meta{TotalPages: 1}},
		{"repo has one chart", "my-repo", "", []*models.Chart{
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}, []*models.Chart{
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}, utils.Meta{TotalPages: 1}},
		{"repo has many charts", "my-repo", "", []*models.Chart{
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
			{ID: "my-repo/dokuwiki", Name: "dokuwiki", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}, {Version: "1.2.2", Digest: "12345"}}},
		}, []*models.Chart{
			{ID: "my-repo/dokuwiki", Name: "dokuwiki", ChartVersions: []models.ChartVersion{{Version: "1.2.3", Digest: "1234"}, {Version: "1.2.2", Digest: "12345"}}},
			{ID: "my-repo/my-chart", Name: "my-chart", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
		}, utils.Meta{TotalPages: 1}},
		{"repo has many charts with pagination", "my-repo", "?size=" + strconv.Itoa(pageSize), []*models.Chart{
			{ID: "my-repo/my-chart3", Name: "my-chart3", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "123"}}},
			{ID: "my-repo/my-chart1", Name: "my-chart1", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "1234"}}},
			{ID: "my-repo/my-chart2", Name: "my-chart2", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "12345"}}},
		}, []*models.Chart{
			{ID: "my-repo/my-chart1", Name: "my-chart1", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "1234"}}},
			{ID: "my-repo/my-chart2", Name: "my-chart2", ChartVersions: []models.ChartVersion{{Version: "0.0.1", Digest: "12345"}}},
		}, utils.Meta{TotalPages: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")
			var chartsList []*models.Chart
			m.On("Find", mock.Anything, bson.M{"repo.name": "my-repo"}, &chartsList, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				*args.Get(2).(*[]*models.Chart) = tt.dbQueryResult
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/charts/"+tt.repo+tt.query, nil)
			params := Params{
				"repo": "my-repo",
			}

			ListRepoCharts(w, req, params)

			m.AssertExpectations(t)
			assert.Equal(t, http.StatusOK, w.Code)

			var b utils.BodyAPIListResponse
			json.NewDecoder(w.Body).Decode(&b)
			data := *b.Data
			assert.Len(t, data, len(tt.chartListResult))
			for i, resp := range data {
				assert.Equal(t, resp.ID, tt.chartListResult[i].ID, "chart id in the response should be the same")
				assert.Equal(t, resp.Type, "chart", "response type is chart")
				assert.Equal(t, resp.Relationships["latestChartVersion"].Data.(map[string]interface{})["version"], tt.chartListResult[i].ChartVersions[0].Version, "latestChartVersion should match version at index 0")
			}
			assert.Equal(t, b.Meta, tt.meta, "response meta should be the same")
		})
	}
}

func Test_getChart(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		chart    models.Chart
		wantCode int
	}{
		{
			"chart does not exist",
			mongo.ErrNoDocuments,
			models.Chart{ID: "my-repo/my-chart"},
			http.StatusNotFound,
		},
		{
			"chart exists",
			nil,
			models.Chart{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}}},
			http.StatusOK,
		},
		{
			"chart has multiple versions",
			nil,
			models.Chart{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}, {Version: "0.0.1"}}},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")

			if tt.err != nil {
				m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			} else {
				m.On("FindOne", mock.Anything, mock.Anything, &models.Chart{}, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					*args.Get(2).(*models.Chart) = tt.chart
				})
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/charts/"+tt.chart.ID, nil)
			parts := strings.Split(tt.chart.ID, "/")
			params := Params{
				"repo":      parts[0],
				"chartName": parts[1],
			}

			GetChart(w, req, params)

			m.AssertExpectations(t)
			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantCode == http.StatusOK {
				var b utils.BodyAPIResponse
				json.NewDecoder(w.Body).Decode(&b)
				assert.Equal(t, b.Data.ID, tt.chart.ID, "chart id in the response should be the same")
				assert.Equal(t, b.Data.Type, "chart", "response type is chart")
				assert.Equal(t, b.Data.Links.(map[string]interface{})["self"], pathPrefix+"/charts/"+tt.chart.ID, "self link should be the same")
				assert.Equal(t, b.Data.Relationships["latestChartVersion"].Data.(map[string]interface{})["version"], tt.chart.ChartVersions[0].Version, "latestChartVersion should match version at index 0")
			}
		})
	}
}

func Test_listChartVersions(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		chart    models.Chart
		wantCode int
	}{
		{
			"chart does not exist",
			mongo.ErrNoDocuments,
			models.Chart{ID: "my-repo/my-chart"},
			http.StatusNotFound,
		},
		{
			"chart exists",
			nil,
			models.Chart{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}}},
			http.StatusOK,
		},
		{
			"chart has multiple versions",
			nil,
			models.Chart{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}, {Version: "0.0.1"}}},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")

			if tt.err != nil {
				m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			} else {
				m.On("FindOne", mock.Anything, mock.Anything, &models.Chart{}, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					*args.Get(2).(*models.Chart) = tt.chart
				})
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/charts/"+tt.chart.ID+"/versions", nil)
			parts := strings.Split(tt.chart.ID, "/")
			params := Params{
				"repo":      parts[0],
				"chartName": parts[1],
			}

			ListChartVersions(w, req, params)

			m.AssertExpectations(t)
			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantCode == http.StatusOK {
				var b utils.BodyAPIListResponse
				json.NewDecoder(w.Body).Decode(&b)
				data := *b.Data
				for i, resp := range data {
					assert.Equal(t, resp.ID, tt.chart.ID+"-"+tt.chart.ChartVersions[i].Version, "chart id in the response should be the same")
					assert.Equal(t, resp.Type, "chartVersion", "response type is chartVersion")
					assert.Equal(t, resp.Attributes.(map[string]interface{})["version"], tt.chart.ChartVersions[i].Version, "chart version should match")
				}
			}
		})
	}
}

func Test_getChartVersion(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		chart    models.Chart
		wantCode int
	}{
		{
			"chart does not exist",
			mongo.ErrNoDocuments,
			models.Chart{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}}},
			http.StatusNotFound,
		},
		{
			"chart exists",
			nil,
			models.Chart{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}}},
			http.StatusOK,
		},
		{
			"chart has multiple versions",
			nil,
			models.Chart{ID: "my-repo/my-chart", ChartVersions: []models.ChartVersion{{Version: "0.1.0"}, {Version: "0.0.1"}}},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")

			if tt.err != nil {
				m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			} else {
				m.On("FindOne", mock.Anything, mock.Anything, &models.Chart{}, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					*args.Get(2).(*models.Chart) = tt.chart
				})
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/charts/"+tt.chart.ID+"/versions/"+tt.chart.ChartVersions[0].Version, nil)
			parts := strings.Split(tt.chart.ID, "/")
			params := Params{
				"repo":      parts[0],
				"chartName": parts[1],
				"version":   tt.chart.ChartVersions[0].Version,
			}

			GetChartVersion(w, req, params)

			m.AssertExpectations(t)
			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantCode == http.StatusOK {
				var b utils.BodyAPIResponse
				json.NewDecoder(w.Body).Decode(&b)
				assert.Equal(t, b.Data.ID, tt.chart.ID+"-"+tt.chart.ChartVersions[0].Version, "chart id in the response should be the same")
				assert.Equal(t, b.Data.Type, "chartVersion", "response type is chartVersion")
				assert.Equal(t, b.Data.Attributes.(map[string]interface{})["version"], tt.chart.ChartVersions[0].Version, "chart version should match")
			}
		})
	}
}

func Test_getChartIcon(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		chart    models.Chart
		wantCode int
	}{
		{
			"chart does not exist",
			mongo.ErrNoDocuments,
			models.Chart{ID: "my-repo/my-chart"},
			http.StatusNotFound,
		},
		{
			"chart has icon",
			nil,
			models.Chart{ID: "my-repo/my-chart", RawIcon: iconBytes(), IconContentType: "image/png"},
			http.StatusOK,
		},
		{
			"chart does not have a icon",
			nil,
			models.Chart{ID: "my-repo/my-chart"},
			http.StatusNotFound,
		},
		{
			"chart has icon with custom type",
			nil,
			models.Chart{ID: "my-repo/my-chart", RawIcon: iconBytes(), IconContentType: "image/svg"},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")

			if tt.err != nil {
				m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			} else {
				m.On("FindOne", mock.Anything, mock.Anything, &models.Chart{}, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					*args.Get(2).(*models.Chart) = tt.chart
				})
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/assets/"+tt.chart.ID+"/logo", nil)
			parts := strings.Split(tt.chart.ID, "/")
			params := Params{
				"repo":      parts[0],
				"chartName": parts[1],
			}

			GetChartIcon(w, req, params)

			m.AssertExpectations(t)
			assert.Equal(t, tt.wantCode, w.Code, "http status code should match")
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, w.Body.Bytes(), tt.chart.RawIcon, "raw icon data should match")
				assert.Equal(t, w.Header().Get("Content-Type"), tt.chart.IconContentType, "icon content type should match")
			}
		})
	}
}

func Test_getChartVersionReadme(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		err      error
		files    models.ChartFiles
		wantCode int
	}{
		{
			"chart does not exist",
			"0.1.0",
			mongo.ErrNoDocuments,
			models.ChartFiles{ID: "my-repo/my-chart"},
			http.StatusNotFound,
		},
		{
			"chart exists",
			"1.2.3",
			nil,
			models.ChartFiles{ID: "my-repo/my-chart", Readme: testChartReadme},
			http.StatusOK,
		},
		{
			"chart does not have a readme",
			"1.1.1",
			nil,
			models.ChartFiles{ID: "my-repo/my-chart"},
			http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")

			if tt.err != nil {
				m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			} else {
				m.On("FindOne", mock.Anything, mock.Anything, &models.ChartFiles{}, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					*args.Get(2).(*models.ChartFiles) = tt.files
				})
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/assets/"+tt.files.ID+"/versions/"+tt.version+"/README.md", nil)
			parts := strings.Split(tt.files.ID, "/")
			params := Params{
				"repo":      parts[0],
				"chartName": parts[1],
				"version":   "0.1.0",
			}

			GetChartVersionReadme(w, req, params)

			m.AssertExpectations(t)
			assert.Equal(t, tt.wantCode, w.Code, "http status code should match")
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, string(w.Body.Bytes()), tt.files.Readme, "content of the readme should match")
			}
		})
	}
}

func Test_getChartVersionValues(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		err      error
		files    models.ChartFiles
		wantCode int
	}{
		{
			"chart does not exist",
			"0.1.0",
			mongo.ErrNoDocuments,
			models.ChartFiles{ID: "my-repo/my-chart"},
			http.StatusNotFound,
		},
		{
			"chart exists",
			"3.2.1",
			nil,
			models.ChartFiles{ID: "my-repo/my-chart", Values: testChartValues},
			http.StatusOK,
		},
		{
			"chart does not have values.yaml",
			"2.2.2",
			nil,
			models.ChartFiles{ID: "my-repo/my-chart"},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mock.Mock
			dbClient = datastore.NewMockClient(&m)
			db, _ = dbClient.Database("test")

			if tt.err != nil {
				m.On("FindOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			} else {
				m.On("FindOne", mock.Anything, mock.Anything, &models.ChartFiles{}, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					*args.Get(2).(*models.ChartFiles) = tt.files
				})
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/assets/"+tt.files.ID+"/versions/"+tt.version+"/values.yaml", nil)
			parts := strings.Split(tt.files.ID, "/")
			params := Params{
				"repo":      parts[0],
				"chartName": parts[1],
				"version":   "0.1.0",
			}

			GetChartVersionValues(w, req, params)

			m.AssertExpectations(t)
			assert.Equal(t, tt.wantCode, w.Code, "http status code should match")
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, string(w.Body.Bytes()), tt.files.Values, "content of values.yaml should match")
			}
		})
	}
}

func Test_findLatestChart(t *testing.T) {
	t.Run("returns mocked chart", func(t *testing.T) {
		chart := &models.Chart{
			Name: "foo",
			ID:   "foo",
			Repo: models.Repo{Name: "bar"},
			ChartVersions: []models.ChartVersion{
				models.ChartVersion{Version: "1.0.0", AppVersion: "0.1.0"},
				models.ChartVersion{Version: "0.0.1", AppVersion: "0.1.0"},
			},
		}
		charts := []*models.Chart{chart}
		reqVersion := "1.0.0"
		reqAppVersion := "0.1.0"

		var m mock.Mock
		dbClient = datastore.NewMockClient(&m)
		db, _ = dbClient.Database("test")
		var chartsList []*models.Chart
		m.On("Find", mock.Anything, mock.Anything, &chartsList, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*[]*models.Chart) = charts
		}).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/charts?name="+chart.Name+"&version="+reqVersion+"&appversion="+reqAppVersion, nil)
		params := Params{
			"name":       chart.Name,
			"version":    reqVersion,
			"appversion": reqAppVersion,
		}

		ListChartsWithFilters(w, req, params)

		var b utils.BodyAPIListResponse
		json.NewDecoder(w.Body).Decode(&b)
		if b.Data == nil {
			t.Fatal("chart list shouldn't be null")
		}
		data := *b.Data

		if data[0].ID != chart.ID {
			t.Errorf("Expecting %v, received %v", chart, data[0].ID)
		}
	})
	t.Run("ignores duplicated chart", func(t *testing.T) {
		charts := []*models.Chart{
			{Name: "foo", ID: "stable/foo", Repo: models.Repo{Name: "bar"}, ChartVersions: []models.ChartVersion{models.ChartVersion{Version: "1.0.0", AppVersion: "0.1.0", Digest: "123"}}},
			{Name: "foo", ID: "bitnami/foo", Repo: models.Repo{Name: "bar"}, ChartVersions: []models.ChartVersion{models.ChartVersion{Version: "1.0.0", AppVersion: "0.1.0", Digest: "123"}}},
		}
		reqVersion := "1.0.0"
		reqAppVersion := "0.1.0"

		var m mock.Mock
		dbClient = datastore.NewMockClient(&m)
		db, _ = dbClient.Database("test")
		var chartsList []*models.Chart
		m.On("Find", mock.Anything, mock.Anything, &chartsList, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*[]*models.Chart) = charts
		}).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/charts?name="+charts[0].Name+"&version="+reqVersion+"&appversion="+reqAppVersion, nil)
		params := Params{
			"name":       charts[0].Name,
			"version":    reqVersion,
			"appversion": reqAppVersion,
		}

		ListChartsWithFilters(w, req, params)

		var b utils.BodyAPIListResponse
		json.NewDecoder(w.Body).Decode(&b)
		if b.Data == nil {
			t.Fatal("chart list shouldn't be null")
		}
		data := *b.Data

		assert.Equal(t, len(data), 1, "it should return a single chart")
		if data[0].ID != charts[0].ID {
			t.Errorf("Expecting %v, received %v", charts[0], data[0].ID)
		}
	})
	t.Run("includes duplicated charts when showDuplicates param set", func(t *testing.T) {
		charts := []*models.Chart{
			{Name: "foo", ID: "stable/foo", Repo: models.Repo{Name: "bar"}, ChartVersions: []models.ChartVersion{models.ChartVersion{Version: "1.0.0", AppVersion: "0.1.0", Digest: "123"}}},
			{Name: "foo", ID: "bitnami/foo", Repo: models.Repo{Name: "bar"}, ChartVersions: []models.ChartVersion{models.ChartVersion{Version: "1.0.0", AppVersion: "0.1.0", Digest: "123"}}},
		}
		reqVersion := "1.0.0"
		reqAppVersion := "0.1.0"

		var m mock.Mock
		dbClient = datastore.NewMockClient(&m)
		db, _ = dbClient.Database("test")
		var chartsList []*models.Chart
		m.On("Find", mock.Anything, mock.Anything, &chartsList, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*[]*models.Chart) = charts
		}).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/charts?showDuplicates=true&name="+charts[0].Name+"&version="+reqVersion+"&appversion="+reqAppVersion, nil)
		params := Params{
			"name":       charts[0].Name,
			"version":    reqVersion,
			"appversion": reqAppVersion,
		}

		ListChartsWithFilters(w, req, params)

		var b utils.BodyAPIListResponse
		json.NewDecoder(w.Body).Decode(&b)
		if b.Data == nil {
			t.Fatal("chart list shouldn't be null")
		}
		data := *b.Data

		assert.Equal(t, 2, len(data), "it should return both charts")
	})
}
