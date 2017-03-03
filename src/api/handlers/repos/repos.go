package repos

import (
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/config"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/handlers"
	reposapi "github.com/helm/monocular/src/api/swagger/restapi/operations/repositories"
)

// GetRepos returns all the enabled repositories
func GetRepos(params reposapi.GetAllReposParams) middleware.Responder {
	config, _ := config.GetConfig()
	resources := helpers.MakeRepoResources(config.Repos)

	payload := handlers.DataResourcesBody(resources)
	return reposapi.NewGetAllReposOK().WithPayload(payload)
}
