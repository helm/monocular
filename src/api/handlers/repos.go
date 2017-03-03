package handlers

import (
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/config"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/swagger/models"
	reposapi "github.com/helm/monocular/src/api/swagger/restapi/operations/repositories"
)

// GetRepos returns all the enabled repositories
func GetRepos(params reposapi.GetAllReposParams) middleware.Responder {
	config, _ := config.GetConfig()
	resources := helpers.MakeRepoResources(config.Repos)
	return reposHTTPBody(resources)
}

func reposHTTPBody(repos []*models.Resource) middleware.Responder {
	resourceArrayData := models.ResourceArrayData{
		Data: repos,
	}
	return reposapi.NewGetAllReposOK().WithPayload(&resourceArrayData)
}
