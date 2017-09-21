package repos

import (
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"

	"github.com/go-openapi/strfmt"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/handlers/renderer"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// GetRepos returns all the enabled repositories
func GetRepos(w http.ResponseWriter, req *http.Request) {
	reposCollection, err := data.GetRepos()
	if err != nil {
		log.WithError(err).Error("unable to get Repos collection")
		renderer.Render.JSON(w, http.StatusInternalServerError, internalServerErrorPayload())
		return
	}
	var repos []*data.Repo
	reposCollection.FindAll(&repos)
	resources := helpers.MakeRepoResources(repos)

	payload := handlers.DataResourcesBody(resources)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// GetRepo returns an enabled repo
func GetRepo(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	var repo data.Repo
	reposCollection, err := data.GetRepos()
	if err != nil {
		log.WithError(err).Error("unable to get Repos collection")
		renderer.Render.JSON(w, http.StatusInternalServerError, internalServerErrorPayload())
		return
	}
	err = reposCollection.Find(params["repo"], &repo)
	if err != nil {
		log.WithError(err).Error("unable to find Repo")
		renderer.Render.JSON(w, http.StatusNotFound, notFoundPayload())
		return
	}

	resource := helpers.MakeRepoResource(models.Repo(repo))
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// CreateRepo adds a repo to the list of enabled repositories to index
func CreateRepo(w http.ResponseWriter, req *http.Request) {
	reposCollection, err := data.GetRepos()
	if err != nil {
		log.WithError(err).Error("unable to get Repos collection")
		renderer.Render.JSON(w, http.StatusInternalServerError, internalServerErrorPayload())
		return
	}

	// Params validation
	format := strfmt.NewFormats()
	var params models.Repo
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "unable to parse request body")
		return
	}
	if err := params.Validate(format); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := url.ParseRequestURI(*params.URL); err != nil {
		errorResponse(w, http.StatusBadRequest, "URL is invalid")
		return
	}

	repo := data.Repo(params)
	if err := reposCollection.Save(&repo); err != nil {
		log.WithError(err).Error("unable to save Repo")
		errorResponse(w, http.StatusInternalServerError, err.Error())
	}

	resource := helpers.MakeRepoResource(models.Repo(repo))
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusCreated, payload)
}

// DeleteRepo deletes a repo from the list of enabled repositories to index
func DeleteRepo(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	reposCollection, err := data.GetRepos()
	if err != nil {
		log.WithError(err).Error("unable to get Repos collection")
		renderer.Render.JSON(w, http.StatusInternalServerError, internalServerErrorPayload())
		return
	}

	repo := data.Repo{}
	found, err := reposCollection.Delete(params["repo"])
	if err != nil {
		log.WithError(err).Error("unable to delete Repo")
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		renderer.Render.JSON(w, http.StatusNotFound, notFoundPayload())
		return
	}

	resource := helpers.MakeRepoResource(models.Repo(repo))
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

func notFoundPayload() *models.Error {
	return &models.Error{Code: pointerto.Int64(http.StatusNotFound), Message: pointerto.String("404 repository not found")}
}

func internalServerErrorPayload() *models.Error {
	return &models.Error{Code: pointerto.Int64(http.StatusInternalServerError), Message: pointerto.String("Internal server error")}
}

func errorResponse(w http.ResponseWriter, errorCode int64, message string) {
	renderer.Render.JSON(w, int(errorCode), models.Error{Code: pointerto.Int64(errorCode), Message: &message})
}
