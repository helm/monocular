package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/strfmt"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
)

// GetChartsInRepoHandlerFunc turns a function with the right signature into a get charts in repo handler
type GetChartsInRepoHandlerFunc func(GetChartsInRepoParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetChartsInRepoHandlerFunc) Handle(params GetChartsInRepoParams) middleware.Responder {
	return fn(params)
}

// GetChartsInRepoHandler interface for that can handle valid get charts in repo params
type GetChartsInRepoHandler interface {
	Handle(GetChartsInRepoParams) middleware.Responder
}

// NewGetChartsInRepo creates a new http.Handler for the get charts in repo operation
func NewGetChartsInRepo(ctx *middleware.Context, handler GetChartsInRepoHandler) *GetChartsInRepo {
	return &GetChartsInRepo{Context: ctx, Handler: handler}
}

/*GetChartsInRepo swagger:route GET /v1/charts/{repo} getChartsInRepo

get all charts by repo

*/
type GetChartsInRepo struct {
	Context *middleware.Context
	Handler GetChartsInRepoHandler
}

func (o *GetChartsInRepo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetChartsInRepoParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

/*GetChartsInRepoOKBodyBody get charts in repo o k body body

swagger:model GetChartsInRepoOKBodyBody
*/
type GetChartsInRepoOKBodyBody struct {

	/* data

	Required: true
	*/
	Data *models.Resource `json:"data"`
}

// Validate validates this get charts in repo o k body body
func (o *GetChartsInRepoOKBodyBody) Validate(formats strfmt.Registry) error {
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

func (o *GetChartsInRepoOKBodyBody) validateData(formats strfmt.Registry) error {

	if o.Data != nil {

		if err := o.Data.Validate(formats); err != nil {
			return err
		}
	}

	return nil
}
