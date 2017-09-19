package repos

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/storage"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	reposapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/repositories"
)

// GetRepos returns all the enabled repositories
func GetRepos(params reposapi.GetAllReposParams) middleware.Responder {
	repos, err := storage.Driver.GetRepos()
	if err != nil {
		log.Error("unable to get Repos collection: ", err)
		return reposapi.NewGetAllReposDefault(http.StatusInternalServerError).WithPayload(internalServerErrorPayload())
	}
	resources := helpers.MakeRepoResources(repos)

	payload := handlers.DataResourcesBody(resources)
	return reposapi.NewGetAllReposOK().WithPayload(payload)
}

// GetRepo returns an enabled repo
func GetRepo(params reposapi.GetRepoParams) middleware.Responder {
	repo, found, err := storage.Driver.GetRepo(params.RepoName)
	if err != nil {
		log.Error("unable to get Repo: ", err)
		return reposapi.NewGetRepoDefault(http.StatusInternalServerError).WithPayload(internalServerErrorPayload())
	}
	if !found {
		log.Error("unable to find Repo: ", err)
		return reposapi.NewGetRepoDefault(http.StatusNotFound).WithPayload(notFoundPayload())
	}

	resource := helpers.MakeRepoResource(models.Repo(*repo))
	payload := handlers.DataResourceBody(resource)
	return reposapi.NewGetRepoOK().WithPayload(payload)
}

// CreateRepo adds a repo to the list of enabled repositories to index
func CreateRepo(params reposapi.CreateRepoParams, releasesEnabled bool) middleware.Responder {
	if !releasesEnabled {
		return errorResponse("Feature not enabled", http.StatusForbidden)
	}

	// Params validation
	format := strfmt.NewFormats()
	if err := params.Data.Validate(format); err != nil {
		return reposapi.NewCreateRepoDefault(http.StatusBadRequest).WithPayload(
			&models.Error{Code: pointerto.Int64(http.StatusBadRequest), Message: pointerto.String(err.Error())})
	}
	if _, err := url.ParseRequestURI(*params.Data.URL); err != nil {
		return reposapi.NewCreateRepoDefault(http.StatusBadRequest).WithPayload(
			&models.Error{Code: pointerto.Int64(http.StatusBadRequest), Message: pointerto.String("URL is invalid")})
	}

	repo := data.Repo(*params.Data)

	if err := storage.Driver.CreateRepo(&repo); err != nil {
		log.Error("unable to save Repo: ", err)
		return reposapi.NewCreateRepoDefault(http.StatusInternalServerError).WithPayload(internalServerErrorPayload())
	}

	resource := helpers.MakeRepoResource(models.Repo(repo))
	payload := handlers.DataResourceBody(resource)
	return reposapi.NewCreateRepoCreated().WithPayload(payload)
}

// DeleteRepo deletes a repo from the list of enabled repositories to index
func DeleteRepo(params reposapi.DeleteRepoParams, releasesEnabled bool) middleware.Responder {
	if !releasesEnabled {
		return errorResponse("Feature not enabled", http.StatusForbidden)
	}

	found, err := storage.Driver.DeleteRepo(params.RepoName)
	if err != nil {
		log.Error("unable to delete Repo: ", err)
		return reposapi.NewGetRepoDefault(http.StatusInternalServerError).WithPayload(internalServerErrorPayload())
	}

	if !found {
		return reposapi.NewGetRepoDefault(http.StatusNotFound).WithPayload(notFoundPayload())
	}

	repo := data.Repo{}
	resource := helpers.MakeRepoResource(models.Repo(repo))
	payload := handlers.DataResourceBody(resource)
	return reposapi.NewGetRepoOK().WithPayload(payload)
}

func notFoundPayload() *models.Error {
	return &models.Error{Code: pointerto.Int64(http.StatusNotFound), Message: pointerto.String("404 repository not found")}
}

func internalServerErrorPayload() *models.Error {
	return &models.Error{Code: pointerto.Int64(http.StatusInternalServerError), Message: pointerto.String("Internal server error")}
}

func errorResponse(message string, errorCode int64) middleware.Responder {
	return reposapi.NewGetAllReposDefault(int(errorCode)).WithPayload(
		&models.Error{Code: pointerto.Int64(errorCode), Message: &message},
	)
}
