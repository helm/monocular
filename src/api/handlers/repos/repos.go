package repos

import (
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	reposapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/repositories"
)

// GetRepos returns all the enabled repositories
func GetRepos(params reposapi.GetAllReposParams) middleware.Responder {
	repos := []*data.Repo{}
	cache.Repos.FindAll(&repos)
	resources := helpers.MakeRepoResources(repos)

	payload := handlers.DataResourcesBody(resources)
	return reposapi.NewGetAllReposOK().WithPayload(payload)
}
