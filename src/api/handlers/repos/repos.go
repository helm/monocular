package repos

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	reposapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/repositories"
)

// GetRepos returns all the enabled repositories
func GetRepos(params reposapi.GetAllReposParams) middleware.Responder {
	reposCollection, err := data.GetRepos()
	if err != nil {
		return reposapi.NewGetAllReposDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Code: pointerto.Int64(http.StatusInternalServerError), Message: pointerto.String("Internal server error")},
		)
	}
	repos := []*data.Repo{}
	reposCollection.FindAll(&repos)
	resources := helpers.MakeRepoResources(repos)

	payload := handlers.DataResourcesBody(resources)
	return reposapi.NewGetAllReposOK().WithPayload(payload)
}
