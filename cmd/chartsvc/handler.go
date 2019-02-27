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

package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/helm/monocular/cmd/chartsvc/models"
	"github.com/kubeapps/common/datastore"
	"github.com/kubeapps/common/response"
	log "github.com/sirupsen/logrus"
)

// Params a key-value map of path params
type Params map[string]string

// WithParams can be used to wrap handlers to take an extra arg for path params
type WithParams func(http.ResponseWriter, *http.Request, Params)

func (h WithParams) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	h(w, req, vars)
}

const chartCollection = "charts"
const filesCollection = "files"

type apiResponse struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Attributes    interface{} `json:"attributes"`
	Links         interface{} `json:"links"`
	Relationships relMap      `json:"relationships"`
}

type apiListResponse []*apiResponse

type selfLink struct {
	Self string `json:"self"`
}

type relMap map[string]rel
type rel struct {
	Data  interface{} `json:"data"`
	Links selfLink    `json:"links"`
}

type meta struct {
	TotalPages int `json:"totalPages"`
}

// getPageNumberAndSize extracts the page number and size of a request. Default (1, 0) if not set
func getPageNumberAndSize(req *http.Request) (int, int) {
	pageNumber := req.FormValue("page")
	pageSize := req.FormValue("size")
	numberInt, err := strconv.Atoi(pageNumber)
	if err != nil || numberInt < 1 {
		numberInt = 1
	}
	sizeInt, _ := strconv.Atoi(pageSize)
	// Safeguard for negative numbers
	if sizeInt < 0 {
		sizeInt = 0
	}
	return numberInt, sizeInt
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getChartList(repo string, pageNumber, pageSize int) (apiListResponse, interface{}, error) {
	db, closer := dbSession.DB()
	defer closer()
	var charts []*models.Chart

	var query datastore.Query
	c := db.C(chartCollection)
	if repo != "" {
		query = c.Find(bson.M{"repo.name": repo})
	} else {
		query = c.Find(nil)
	}
	query = query.Sort("name")

	// TODO(andresmgot): We should not query the charts just yet in case we paginate the response
	// but we need to do so since it's not possible to filter & count charts (countDocuments) until
	// https://github.com/kubeapps/common/issues/8 is addressed
	err := query.All(&charts)
	if err != nil {
		return apiListResponse{}, 0, err
	}

	totalPages := 1
	if pageSize != 0 {
		// If a pageSize is given, returns only the the specified number of charts and
		// the number of pages
		length := len(charts)
		totalPages = int(math.Ceil(float64(length) / float64(pageSize)))

		// If the page number is out of range, return the last one
		if pageNumber > totalPages {
			pageNumber = totalPages
		}

		first := pageSize * (pageNumber - 1)
		last := min(first+pageSize, length)
		charts = charts[first:last]
	}

	return newChartListResponse(charts), meta{totalPages}, nil
}

// listCharts returns a list of charts
func listCharts(w http.ResponseWriter, req *http.Request) {
	pageNumber, pageSize := getPageNumberAndSize(req)
	cl, meta, err := getChartList("", pageNumber, pageSize)
	if err != nil {
		log.WithError(err).Error("could not fetch charts")
		response.NewErrorResponse(http.StatusInternalServerError, "could not fetch all charts").Write(w)
		return
	}
	response.NewDataResponseWithMeta(cl, meta).Write(w)
}

// listRepoCharts returns a list of charts in the given repo
func listRepoCharts(w http.ResponseWriter, req *http.Request, params Params) {
	pageNumber, pageSize := getPageNumberAndSize(req)
	cl, meta, err := getChartList(params["repo"], pageNumber, pageSize)
	if err != nil {
		log.WithError(err).Error("could not fetch charts")
		response.NewErrorResponse(http.StatusInternalServerError, "could not fetch all charts").Write(w)
		return
	}
	response.NewDataResponseWithMeta(cl, meta).Write(w)
}

// getChart returns the chart from the given repo
func getChart(w http.ResponseWriter, req *http.Request, params Params) {
	db, closer := dbSession.DB()
	defer closer()
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])
	if err := db.C(chartCollection).FindId(chartID).One(&chart); err != nil {
		log.WithError(err).Errorf("could not find chart with id %s", chartID)
		response.NewErrorResponse(http.StatusNotFound, "could not find chart").Write(w)
		return
	}

	cr := newChartResponse(&chart)
	response.NewDataResponse(cr).Write(w)
}

// listChartVersions returns a list of chart versions for the given chart
func listChartVersions(w http.ResponseWriter, req *http.Request, params Params) {
	db, closer := dbSession.DB()
	defer closer()
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])
	if err := db.C(chartCollection).FindId(chartID).One(&chart); err != nil {
		log.WithError(err).Errorf("could not find chart with id %s", chartID)
		response.NewErrorResponse(http.StatusNotFound, "could not find chart").Write(w)
		return
	}

	cvl := newChartVersionListResponse(&chart)
	response.NewDataResponse(cvl).Write(w)
}

// getChartVersion returns the given chart version
func getChartVersion(w http.ResponseWriter, req *http.Request, params Params) {
	db, closer := dbSession.DB()
	defer closer()
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])
	if err := db.C(chartCollection).Find(bson.M{
		"_id":           chartID,
		"chartversions": bson.M{"$elemMatch": bson.M{"version": params["version"]}},
	}).Select(bson.M{
		"name": 1, "repo": 1, "description": 1, "home": 1, "keywords": 1, "maintainers": 1, "sources": 1,
		"chartversions.$": 1,
	}).One(&chart); err != nil {
		log.WithError(err).Errorf("could not find chart with id %s", chartID)
		response.NewErrorResponse(http.StatusNotFound, "could not find chart version").Write(w)
		return
	}

	cvr := newChartVersionResponse(&chart, chart.ChartVersions[0])
	response.NewDataResponse(cvr).Write(w)
}

// getChartIcon returns the icon for a given chart
func getChartIcon(w http.ResponseWriter, req *http.Request, params Params) {
	db, closer := dbSession.DB()
	defer closer()
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])
	if err := db.C(chartCollection).FindId(chartID).One(&chart); err != nil {
		log.WithError(err).Errorf("could not find chart with id %s", chartID)
		http.NotFound(w, req)
		return
	}

	if chart.RawIcon == nil {
		http.NotFound(w, req)
		return
	}

	w.Write(chart.RawIcon)
}

// getChartVersionReadme returns the README for a given chart
func getChartVersionReadme(w http.ResponseWriter, req *http.Request, params Params) {
	db, closer := dbSession.DB()
	defer closer()
	var files models.ChartFiles
	fileID := fmt.Sprintf("%s/%s-%s", params["repo"], params["chartName"], params["version"])
	if err := db.C(filesCollection).FindId(fileID).One(&files); err != nil {
		log.WithError(err).Errorf("could not find files with id %s", fileID)
		http.NotFound(w, req)
		return
	}
	readme := []byte(files.Readme)
	if len(readme) == 0 {
		log.Errorf("could not find a README for id %s", fileID)
		http.NotFound(w, req)
		return
	}
	w.Write(readme)
}

// getChartVersionValues returns the values.yaml for a given chart
func getChartVersionValues(w http.ResponseWriter, req *http.Request, params Params) {
	db, closer := dbSession.DB()
	defer closer()
	var files models.ChartFiles
	fileID := fmt.Sprintf("%s/%s-%s", params["repo"], params["chartName"], params["version"])
	if err := db.C(filesCollection).FindId(fileID).One(&files); err != nil {
		log.WithError(err).Errorf("could not find values.yaml with id %s", fileID)
		http.NotFound(w, req)
		return
	}

	w.Write([]byte(files.Values))
}

// listChartsWithFilters returns the list of repos that contains the given chart and the latest version found
func listChartsWithFilters(w http.ResponseWriter, req *http.Request, params Params) {
	db, closer := dbSession.DB()
	defer closer()

	var charts []*models.Chart
	if err := db.C(chartCollection).Find(bson.M{
		"name": params["chartName"],
		"chartversions": bson.M{
			"$elemMatch": bson.M{"version": req.FormValue("version"), "appversion": req.FormValue("appversion")},
		}}).Select(bson.M{
		"name": 1, "repo": 1,
		"chartversions": bson.M{"$slice": 1},
	}).All(&charts); err != nil {
		log.WithError(err).Errorf(
			"could not find charts with the given name %s, version %s and appversion %s",
			params["chartName"], req.FormValue("version"), req.FormValue("appversion"),
		)
		// continue to return empty list
	}

	cl := newChartListResponse(charts)
	response.NewDataResponse(cl).Write(w)
}

func newChartResponse(c *models.Chart) *apiResponse {
	latestCV := c.ChartVersions[0]
	return &apiResponse{
		Type:       "chart",
		ID:         c.ID,
		Attributes: chartAttributes(*c),
		Links:      selfLink{pathPrefix + "/charts/" + c.ID},
		Relationships: relMap{
			"latestChartVersion": rel{
				Data:  chartVersionAttributes(c.ID, latestCV),
				Links: selfLink{pathPrefix + "/charts/" + c.ID + "/versions/" + latestCV.Version},
			},
		},
	}
}

func newChartListResponse(charts []*models.Chart) apiListResponse {
	// We will keep track of unique digest:chart to avoid duplicates
	chartDigests := map[string]bool{}
	cl := apiListResponse{}
	for _, c := range charts {
		digest := c.ChartVersions[0].Digest
		// Filter out the chart if we've seen the same digest before
		if _, ok := chartDigests[digest]; !ok {
			chartDigests[digest] = true
			cl = append(cl, newChartResponse(c))
		}
	}
	return cl
}

func chartVersionAttributes(cid string, cv models.ChartVersion) models.ChartVersion {
	cv.Readme = pathPrefix + "/assets/" + cid + "/versions/" + cv.Version + "/README.md"
	cv.Values = pathPrefix + "/assets/" + cid + "/versions/" + cv.Version + "/values.yaml"
	return cv
}

func chartAttributes(c models.Chart) models.Chart {
	if c.RawIcon != nil {
		c.Icon = pathPrefix + "/assets/" + c.ID + "/logo-160x160-fit.png"
	} else {
		// If the icon wasn't processed, it is either not set or invalid
		c.Icon = ""
	}
	return c
}

func newChartVersionResponse(c *models.Chart, cv models.ChartVersion) *apiResponse {
	return &apiResponse{
		Type:       "chartVersion",
		ID:         fmt.Sprintf("%s-%s", c.ID, cv.Version),
		Attributes: chartVersionAttributes(c.ID, cv),
		Links:      selfLink{pathPrefix + "/charts/" + c.ID + "/versions/" + cv.Version},
		Relationships: relMap{
			"chart": rel{
				Data:  chartAttributes(*c),
				Links: selfLink{pathPrefix + "/charts/" + c.ID},
			},
		},
	}
}

func newChartVersionListResponse(c *models.Chart) apiListResponse {
	var cvl apiListResponse
	for _, cv := range c.ChartVersions {
		cvl = append(cvl, newChartVersionResponse(c, cv))
	}

	return cvl
}
