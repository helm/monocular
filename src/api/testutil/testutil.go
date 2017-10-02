package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// constants
const (
	RepoName           = "stable"
	BogusRepo          = "bogon"
	ChartName          = "drupal"
	ChartVersionString = "0.3.4"
	UnparseableRepo    = "unparseable"
)

// AssertErrBodyData asserts expected HTTP error response body data
func AssertErrBodyData(t *testing.T, code int64, resource string, body models.Error) {
	assert.Equal(t, *body.Code, code, "response code in HTTP body data")
	assert.Equal(t, *body.Message, strconv.FormatInt(code, 10)+" "+resource+" not found", "error message in HTTP body data")
}

// AssertChartResourceBodyData asserts expected HTTP chart resource body data
func AssertChartResourceBodyData(t *testing.T, chart *models.Resource, body *models.ResourceData) {
	attributes, err := ChartResourceAttributesFromHTTPResponse(body)
	assert.NoErr(t, err)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.ID, *body.Data.ID, "chart ID data in HTTP body data")
	assert.Equal(t, *chart.Type, *body.Data.Type, "chart type data in HTTP body data")
	assert.Equal(t, *chart.Attributes.(*models.Chart).Description, *attributes.Description, "chart descripion data in HTTP body data")
	assert.Equal(t, *chart.Attributes.(*models.Chart).Home, *attributes.Home, "chart home data in HTTP body data")
	assert.Equal(t, *chart.Attributes.(*models.Chart).Name, *attributes.Name, "chart name data in HTTP body data")
	assert.Equal(t, *chart.Attributes.(*models.Chart).Repo, *attributes.Repo, "chart repo data in HTTP body data")
}

// AssertChartVersionResourceBodyData asserts expected HTTP "chart version" resource body data
func AssertChartVersionResourceBodyData(t *testing.T, chart *models.Resource, body *models.ResourceData) {
	attributes, err := ChartVersionResourceAttributesFromHTTPResponse(body)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.ID, *body.Data.ID, "chart ID data in HTTP body data")
	assert.Equal(t, *chart.Type, *body.Data.Type, "chart type data in HTTP body data")
	assert.Equal(t, *chart.Attributes.(*models.ChartVersion).Created, *attributes.Created, "chartVersion created data in HTTP body data")
	assert.Equal(t, *chart.Attributes.(*models.ChartVersion).Digest, *attributes.Digest, "chartVersion digest data in HTTP body data")
	assert.Equal(t, chart.Attributes.(*models.ChartVersion).Urls, attributes.Urls, "chartVersion URLs data in HTTP body data")
	assert.Equal(t, *chart.Attributes.(*models.ChartVersion).Version, *attributes.Version, "chartVersion version data in HTTP body data")
}

// ResourceArrayDataFromJSON is a convenience that converts JSON to a models.ResourceArrayData
func ResourceArrayDataFromJSON(r io.Reader, resource *models.ResourceArrayData) error {
	return json.NewDecoder(r).Decode(resource)
}

// ResourceDataFromJSON is a convenience that converts JSON to a models.ResourceData
func ResourceDataFromJSON(r io.Reader, resource *models.ResourceData) error {
	return json.NewDecoder(r).Decode(resource)
}

// ErrorModelFromJSON is a convenience that converts JSON to a models.Error
func ErrorModelFromJSON(r io.Reader, errorModel *models.Error) error {
	return json.NewDecoder(r).Decode(errorModel)
}

// ChartResourceAttributesFromHTTPResponse is a convenience that grabs the Attributes interface from
// a chart resource in generic models.ResourceData form and converts to a models.Chart
func ChartResourceAttributesFromHTTPResponse(body *models.ResourceData) (*models.Chart, error) {
	attributes := new(models.Chart)
	b, err := json.Marshal(body.Data.Attributes.(map[string]interface{}))
	if err != nil {
		return attributes, err
	}
	err = json.Unmarshal(b, attributes)
	return attributes, err
}

// ChartVersionResourceAttributesFromHTTPResponse is a convenience that grabs the Attributes interface from
// a chart resource in generic models.ResourceData form and converts to a models.ChartPackage
func ChartVersionResourceAttributesFromHTTPResponse(body *models.ResourceData) (*models.ChartVersion, error) {
	attributes := new(models.ChartVersion)
	b, err := json.Marshal(body.Data.Attributes.(map[string]interface{}))
	if err != nil {
		return attributes, err
	}
	err = json.Unmarshal(b, attributes)
	return attributes, err
}

// ToJSONBody returns a JSON representation of a struct
func ToJSONBody(t *testing.T, body interface{}) *bytes.Buffer {
	jsonParams, err := json.Marshal(body)
	assert.NoErr(t, err)
	return bytes.NewBuffer(jsonParams)
}
