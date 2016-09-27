package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
)

// GetAllChartsHandlerFunc turns a function with the right signature into a get all charts handler
type GetAllChartsHandlerFunc func(GetAllChartsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAllChartsHandlerFunc) Handle(params GetAllChartsParams) middleware.Responder {
	return fn(params)
}

// GetAllChartsHandler interface for that can handle valid get all charts params
type GetAllChartsHandler interface {
	Handle(GetAllChartsParams) middleware.Responder
}

// NewGetAllCharts creates a new http.Handler for the get all charts operation
func NewGetAllCharts(ctx *middleware.Context, handler GetAllChartsHandler) *GetAllCharts {
	return &GetAllCharts{Context: ctx, Handler: handler}
}

/*GetAllCharts swagger:route GET /v1/charts getAllCharts

get all charts from all repos

*/
type GetAllCharts struct {
	Context *middleware.Context
	Handler GetAllChartsHandler
}

func (o *GetAllCharts) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetAllChartsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

/*GetAllChartsOKBodyBody get all charts o k body body

swagger:model GetAllChartsOKBodyBody
*/
type GetAllChartsOKBodyBody struct {

	/* data

	Required: true
	*/
	Data []*models.Resource `json:"data"`
}

// Validate validates this get all charts o k body body
func (o *GetAllChartsOKBodyBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetAllChartsOKBodyBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("getAllChartsOK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	for i := 0; i < len(o.Data); i++ {

		if swag.IsZero(o.Data[i]) { // not required
			continue
		}

		if o.Data[i] != nil {

			if err := o.Data[i].Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}
