package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/testutil"
)

func TestNotFound(t *testing.T) {
	const resource1 = "chart"
	const resource2 = "repo"
	w := httptest.NewRecorder()
	resp := notFound(resource1)
	assert.NotNil(t, resp, "notFound response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody1 models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody1))
	testutil.AssertErrBodyData(t, http.StatusNotFound, resource1, httpBody1)
	w = httptest.NewRecorder()
	var httpBody2 models.Error
	resp2 := notFound(resource2)
	assert.NotNil(t, resp2, "notFound response")
	resp2.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody2))
	testutil.AssertErrBodyData(t, http.StatusNotFound, resource2, httpBody2)
}

func AssertErrBodyData(t *testing.T, code int64, resource string, body models.Error) {
	assert.Equal(t, *body.Code, code, "response code in HTTP body data")
	assert.Equal(t, *body.Message, strconv.FormatInt(code, 10)+" "+resource+" not found", "error message in HTTP body data")
}

func AssertChartResourceBodyData(t *testing.T, chart *models.Resource, body *models.ResourceData) {
	attributes, err := testutil.ChartResourceAttributesFromHTTPResponse(body)
	assert.NoErr(t, err)
	links, err := testutil.ChartResourceLinksFromHTTPResponse(body)
	assert.NoErr(t, err)
	assert.Equal(t, chart.ID, body.Data.ID, "chart ID data in HTTP body data")
	assert.Equal(t, chart.Type, body.Data.Type, "chart type data in HTTP body data")
	assert.Equal(t, chart.Attributes.(*models.ChartResourceAttributes).Created, attributes.Created, "chart created data in HTTP body data")
	assert.Equal(t, chart.Attributes.(*models.ChartResourceAttributes).Description, attributes.Description, "chart descripion data in HTTP body data")
	assert.Equal(t, chart.Attributes.(*models.ChartResourceAttributes).Home, attributes.Home, "chart home data in HTTP body data")
	assert.Equal(t, chart.Attributes.(*models.ChartResourceAttributes).Name, attributes.Name, "chart name data in HTTP body data")
	assert.Equal(t, chart.Attributes.(*models.ChartResourceAttributes).Repo, attributes.Repo, "chart repo data in HTTP body data")
	assert.Equal(t, chart.Links.(*models.ChartResourceLinks).Latest, links.Latest, "chart link to latest data in HTTP body data")
}
