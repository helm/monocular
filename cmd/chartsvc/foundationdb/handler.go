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
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/helm/monocular/cmd/chartsvc/foundationdb/datastore"
	"github.com/helm/monocular/cmd/chartsvc/models"
	"github.com/helm/monocular/cmd/chartsvc/utils"

	"github.com/gorilla/mux"
	"github.com/kubeapps/common/response"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// count is used to parse the result of a $count operation in the database
type count struct {
	Count int
}

var dbClient datastore.Client
var db datastore.Database
var dbCloser func()

//var db mongo.Database
var dbName string
var pathPrefix string

//SetPathPrefix sets the URL prefix for the ChartSVC API endpoint
func SetPathPrefix(prefix string) {
	pathPrefix = prefix
}

//InitDBConfig sets FDB Document-Layer client and DB config for the ChartSVC API handler
func InitDBConfig(client datastore.Client, name string) {
	dbClient = client
	db, dbCloser = dbClient.Database(name)
	dbName = name
}

// getPageNumberAndSize extracts the page number and size of a request. Default (1, 0) if not set
func getPageNumberAndSize(req *http.Request) (int, int) {
	page := req.FormValue("page")
	size := req.FormValue("size")
	pageInt, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		pageInt = 1
	}
	// ParseUint will return 0 if size is a not positive integer
	sizeInt, _ := strconv.ParseUint(size, 10, 64)
	return int(pageInt), int(sizeInt)
}

// showDuplicates returns if a request wants to retrieve charts. Default false
func showDuplicates(req *http.Request) bool {
	return len(req.FormValue("showDuplicates")) > 0
}

// min returns the minimum of two integers.
// We are not using math.Min since that compares float64
// and it's unnecessarily complex.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func uniqChartList(charts []*models.Chart) []*models.Chart {
	// We will keep track of unique digest:chart to avoid duplicates
	chartDigests := map[string]bool{}
	res := []*models.Chart{}
	for _, c := range charts {
		digest := c.ChartVersions[0].Digest
		// Filter out the chart if we've seen the same digest before
		if _, ok := chartDigests[digest]; !ok {
			chartDigests[digest] = true
			res = append(res, c)
		}
	}
	return res
}

func getPaginatedChartList(repo string, pageNumber, pageSize int, showDuplicates bool) (utils.ApiListResponse, interface{}, error) {
	log.Debugf("Request for paginated chart list..")

	//Find all charts for repo name and sort by chart name
	collection := db.Collection(chartCollection)
	filter := bson.M{}
	if repo != "" {
		filter = bson.M{"repo.name": repo}
	}
	var charts []*models.Chart
	err := collection.Find(context.Background(), filter, &charts, options.Find())
	if err != nil {
		log.WithError(err).Errorf(
			"Error fetching charts from DB for pagination %s",
			repo,
		)
		return newChartListResponse([]*models.Chart{}), utils.Meta{TotalPages: 0}, err
	}
	var tempChartMap map[string]*models.Chart = make(map[string]*models.Chart)

	chartsToSort := make([]*models.Chart, 0, len(tempChartMap))

	if !showDuplicates {
		for _, chart := range charts {
			log.Debugf("Chart digest: %v.", chart.ChartVersions[0].Digest)
			tempChartMap[chart.ChartVersions[0].Digest] = chart
		}
		log.Debugf("Charts in map: %v", len(tempChartMap))
		//Now just get all the values from our map
		for _, v := range tempChartMap {
			log.Debugf("Adding chart: %v to unique chart list.", *v)
			chartsToSort = append(chartsToSort, v)
		}
	} else {
		chartsToSort = charts
	}

	//Sort the list of paginated charts by name
	sort.Slice(chartsToSort, func(i, j int) bool {
		return chartsToSort[i].Name < chartsToSort[j].Name
	})

	sortedCharts := chartsToSort
	log.Debugf("Charts in sorted list: %v", len(sortedCharts))
	log.Debugf("Page size requested: %v", pageSize)
	var paginatedCharts = sortedCharts
	totalPages := 1
	if pageSize != 0 {
		// If a pageSize is given, returns only the the specified number of charts and
		// the number of pages
		cc := count{}
		cc.Count = len(sortedCharts)
		totalPages = int(math.Ceil(float64(cc.Count) / float64(pageSize)))

		// If the page number is out of range, return the last one
		if pageNumber > totalPages {
			pageNumber = totalPages
		}
		paginatedCharts = sortedCharts[(pageNumber-1)*pageSize : pageNumber*pageSize]
	}

	log.Debugf("Returning %v charts, Done.", len(paginatedCharts))
	return newChartListResponse(paginatedCharts), utils.Meta{TotalPages: totalPages}, nil
}

// ListCharts returns a list of charts
func ListCharts(w http.ResponseWriter, req *http.Request) {
	log.Debug("Request for charts..")
	pageNumber, pageSize := getPageNumberAndSize(req)
	cl, meta, err := getPaginatedChartList("", pageNumber, pageSize, showDuplicates(req))
	if err != nil {
		log.WithError(err).Error("could not fetch charts")
		response.NewErrorResponse(http.StatusInternalServerError, "could not fetch all charts").Write(w)
		return
	}
	response.NewDataResponseWithMeta(cl, meta).Write(w)
	log.Debug("Done.")
}

// ListRepoCharts returns a list of charts in the given repo
func ListRepoCharts(w http.ResponseWriter, req *http.Request, params Params) {
	log.Debug("Request for charts..")
	pageNumber, pageSize := getPageNumberAndSize(req)
	cl, meta, err := getPaginatedChartList(params["repo"], pageNumber, pageSize, showDuplicates(req))
	if err != nil {
		log.WithError(err).Error("could not fetch charts")
		response.NewErrorResponse(http.StatusInternalServerError, "could not fetch all charts").Write(w)
		return
	}
	response.NewDataResponseWithMeta(cl, meta).Write(w)
	log.Debug("Done.")
}

// GetChart returns the chart from the given repo
func GetChart(w http.ResponseWriter, req *http.Request, params Params) {
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])

	chartCollection := db.Collection(chartCollection)
	filter := bson.M{"_id": chartID}
	findResult := chartCollection.FindOne(context.Background(), filter, &chart, options.FindOne())
	if findResult == mongo.ErrNoDocuments {
		log.WithError(findResult).Errorf("could not find chart with id %s", chartID)
		response.NewErrorResponse(http.StatusNotFound, "could not find chart").Write(w)
		return
	}

	cr := newChartResponse(&chart)
	response.NewDataResponse(cr).Write(w)
}

// ListChartVersions returns a list of chart versions for the given chart
func ListChartVersions(w http.ResponseWriter, req *http.Request, params Params) {
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])

	chartCollection := db.Collection(chartCollection)
	filter := bson.M{"_id": chartID}
	findResult := chartCollection.FindOne(context.Background(), filter, &chart, options.FindOne())
	if findResult == mongo.ErrNoDocuments {
		log.WithError(findResult).Errorf("could not find chart with id %s", chartID)
		response.NewErrorResponse(http.StatusNotFound, "could not find chart").Write(w)
		return
	}

	cvl := newChartVersionListResponse(&chart)
	response.NewDataResponse(cvl).Write(w)
}

// GetChartVersion returns the given chart version
func GetChartVersion(w http.ResponseWriter, req *http.Request, params Params) {
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])

	chartCollection := db.Collection(chartCollection)
	filter := bson.M{
		"_id":           chartID,
		"chartversions": bson.M{"$elemMatch": bson.M{"version": params["version"]}},
	}
	projection := bson.M{
		"name": 1, "repo": 1, "description": 1, "home": 1, "keywords": 1, "maintainers": 1, "sources": 1,
		"chartversions": 1,
	}
	findResult := chartCollection.FindOne(context.Background(), filter, &chart, options.FindOne().SetProjection(projection))
	if findResult == mongo.ErrNoDocuments {
		log.WithError(findResult).Errorf("could not find chart with id %s", chartID)
		response.NewErrorResponse(http.StatusNotFound, "could not find chart").Write(w)
		return
	}

	for i := range chart.ChartVersions {
		if chart.ChartVersions[i].Version == params["version"] {
			chart.ChartVersions = chart.ChartVersions[i : i+1]
			break
		}
	}
	// Cut the versions slice down to just one element
	cvr := newChartVersionResponse(&chart, chart.ChartVersions[0])
	response.NewDataResponse(cvr).Write(w)
}

// GetChartIcon returns the icon for a given chart
func GetChartIcon(w http.ResponseWriter, req *http.Request, params Params) {
	var chart models.Chart
	chartID := fmt.Sprintf("%s/%s", params["repo"], params["chartName"])

	chartCollection := db.Collection(chartCollection)
	filter := bson.M{"_id": chartID}
	findResult := chartCollection.FindOne(context.Background(), filter, &chart, options.FindOne())
	if findResult == mongo.ErrNoDocuments {
		log.WithError(findResult).Errorf("could not find chart with id %s", chartID)
		http.NotFound(w, req)
		return
	}
	if chart.RawIcon == nil {
		http.NotFound(w, req)
		return
	}

	if chart.IconContentType != "" {
		// Force the Content-Type header because the autogenerated type does not work for
		// image/svg+xml. It is detected as plain text
		w.Header().Set("Content-Type", chart.IconContentType)
	}

	w.Write(chart.RawIcon)
}

// GetChartVersionReadme returns the README for a given chart
func GetChartVersionReadme(w http.ResponseWriter, req *http.Request, params Params) {

	var files models.ChartFiles
	fileID := fmt.Sprintf("%s/%s-%s", params["repo"], params["chartName"], params["version"])

	filesCollection := db.Collection(filesCollection)
	filter := bson.M{"_id": fileID}
	findResult := filesCollection.FindOne(context.Background(), filter, &files, options.FindOne())
	if findResult == mongo.ErrNoDocuments {
		log.WithError(findResult).Errorf("could not find files with id %s", fileID)
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

// GetChartVersionValues returns the values.yaml for a given chart
func GetChartVersionValues(w http.ResponseWriter, req *http.Request, params Params) {
	var files models.ChartFiles

	fileID := fmt.Sprintf("%s/%s-%s", params["repo"], params["chartName"], params["version"])
	filesCollection := db.Collection(filesCollection)
	filter := bson.M{"_id": fileID}
	findResult := filesCollection.FindOne(context.Background(), filter, &files, options.FindOne())
	if findResult == mongo.ErrNoDocuments {
		log.WithError(findResult).Errorf("could not find values.yaml with id %s", fileID)
		http.NotFound(w, req)
		return
	}

	w.Write([]byte(files.Values))
}

// GetChartVersionSchema returns the values.schema.json for a given chart
func GetChartVersionSchema(w http.ResponseWriter, req *http.Request, params Params) {

	var files models.ChartFiles

	fileID := fmt.Sprintf("%s/%s-%s", params["repo"], params["chartName"], params["version"])
	filter := bson.M{"_id": fileID}
	filesCollection := db.Collection(filesCollection)
	findResult := filesCollection.FindOne(context.Background(), filter, &files, options.FindOne())
	if findResult == mongo.ErrNoDocuments {
		log.WithError(findResult).Errorf("could not find values.schema.json with id %s", fileID)
		http.NotFound(w, req)
		return
	}

	w.Write([]byte(files.Schema))
}

// ListChartsWithFilters returns the list of repos that contains the given chart and the latest version found
func ListChartsWithFilters(w http.ResponseWriter, req *http.Request, params Params) {

	var charts []*models.Chart

	chartCollection := db.Collection(chartCollection)
	filter := bson.M{
		"name": params["chartName"],
		"chartversions": bson.M{
			"$elemMatch": bson.M{"version": req.FormValue("version"), "appversion": req.FormValue("appversion")},
		}}
	projection := bson.M{
		"name": 1, "repo": 1,
		"chartversions": bson.M{"$slice": 1},
	}
	err := chartCollection.Find(context.Background(), filter, &charts, options.Find().SetProjection(projection))
	if err != nil {
		log.WithError(err).Errorf(
			"Error finding charts with the given name %s, version %s and appversion %s",
			params["chartName"], req.FormValue("version"), req.FormValue("appversion"),
		)
		// continue to return empty list
	}

	chartResponse := charts
	if !showDuplicates(req) {
		chartResponse = uniqChartList(charts)
	}

	cl := newChartListResponse(chartResponse)
	response.NewDataResponse(cl).Write(w)
}

// SearchCharts returns the list of charts that matches the query param in any of these fields:
//  - name
//  - description
//  - repository name
//  - any keyword
//  - any source
//  - any maintainer name
func SearchCharts(w http.ResponseWriter, req *http.Request, params Params) {

	query := req.FormValue("q")
	var charts []*models.Chart

	chartCollection := db.Collection(chartCollection)
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query}},
			{"description": bson.M{"$regex": query}},
			{"repo.name": bson.M{"$regex": query}},
			{"keywords": bson.M{"$elemMatch": bson.M{"$regex": query}}},
			{"sources": bson.M{"$elemMatch": bson.M{"$regex": query}}},
			{"maintainers": bson.M{"$elemMatch": bson.M{"name": bson.M{"$regex": query}}}},
		},
	}
	if params["repo"] != "" {
		filter["repo.name"] = params["repo"]
	}
	err := chartCollection.Find(context.Background(), filter, &charts, options.Find())
	if err != nil {
		log.WithError(err).Errorf(
			"Error finding charts with the given name %s, version %s and appversion %s",
			params["chartName"], req.FormValue("version"), req.FormValue("appversion"),
		)
		// continue to return empty list
	}

	chartResponse := charts
	if !showDuplicates(req) {
		chartResponse = uniqChartList(charts)
	}

	cl := newChartListResponse(uniqChartList(chartResponse))
	response.NewDataResponse(cl).Write(w)
}

func newChartResponse(c *models.Chart) *utils.ApiResponse {
	latestCV := c.ChartVersions[0]
	return &utils.ApiResponse{
		Type:       "chart",
		ID:         c.ID,
		Attributes: chartAttributes(*c),
		Links:      utils.SelfLink{Self: pathPrefix + "/charts/" + c.ID},
		Relationships: utils.RelMap{
			"latestChartVersion": utils.Rel{
				Data:  chartVersionAttributes(c.ID, latestCV),
				Links: utils.SelfLink{Self: pathPrefix + "/charts/" + c.ID + "/versions/" + latestCV.Version},
			},
		},
	}
}

func newChartListResponse(charts []*models.Chart) utils.ApiListResponse {
	cl := utils.ApiListResponse{}
	for _, c := range charts {
		cl = append(cl, newChartResponse(c))
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
		c.Icon = pathPrefix + "/assets/" + c.ID + "/logo"
	} else {
		// If the icon wasn't processed, it is either not set or invalid
		c.Icon = ""
	}
	return c
}

func newChartVersionResponse(c *models.Chart, cv models.ChartVersion) *utils.ApiResponse {
	return &utils.ApiResponse{
		Type:       "chartVersion",
		ID:         fmt.Sprintf("%s-%s", c.ID, cv.Version),
		Attributes: chartVersionAttributes(c.ID, cv),
		Links:      utils.SelfLink{Self: pathPrefix + "/charts/" + c.ID + "/versions/" + cv.Version},
		Relationships: utils.RelMap{
			"chart": utils.Rel{
				Data:  chartAttributes(*c),
				Links: utils.SelfLink{Self: pathPrefix + "/charts/" + c.ID},
			},
		},
	}
}

func newChartVersionListResponse(c *models.Chart) utils.ApiListResponse {
	var cvl utils.ApiListResponse
	for _, cv := range c.ChartVersions {
		cvl = append(cvl, newChartVersionResponse(c, cv))
	}

	return cvl
}
