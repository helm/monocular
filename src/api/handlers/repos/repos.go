package repos

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"

	log "github.com/Sirupsen/logrus"

	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/handlers/renderer"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// RepoHandlers defines handlers that serve chart data
type RepoHandlers struct {
	dbSession datastore.Session
}

// NewRepoHandlers takes a datastore.Session implementation and returns a RepoHandlers struct
func NewRepoHandlers(db datastore.Session) *RepoHandlers {
	return &RepoHandlers{db}
}

// ListRepos returns all repositories
func (r *RepoHandlers) ListRepos(w http.ResponseWriter, req *http.Request) {
	db, closer := r.dbSession.DB()
	defer closer()
	repos, err := models.ListRepos(db)
	if err != nil {
		log.WithError(err).Error("unable to get fetch repos")
		renderer.Render.JSON(w, http.StatusInternalServerError, internalServerErrorPayload())
		return
	}
	resources := helpers.MakeRepoResources(repos)

	payload := handlers.DataResourcesBody(resources)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// GetRepo returns a repo
func (r *RepoHandlers) GetRepo(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	db, closer := r.dbSession.DB()
	defer closer()
	repo, err := models.GetRepo(db, params["repo"])
	if err != nil {
		log.WithError(err).Error("unable to find Repo")
		renderer.Render.JSON(w, http.StatusNotFound, notFoundPayload())
		return
	}

	resource := helpers.MakeRepoResource(repo)
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// CreateRepo adds a repo to the list of enabled repositories to index
func (r *RepoHandlers) CreateRepo(w http.ResponseWriter, req *http.Request) {
	db, closer := r.dbSession.DB()
	defer closer()

	// Params validation
	var repo *models.Repo
	if err := json.NewDecoder(req.Body).Decode(&repo); err != nil {
		errorResponse(w, http.StatusBadRequest, "unable to parse request body: "+err.Error())
		return
	}

	if _, err := govalidator.ValidateStruct(repo); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := models.CreateRepo(db, repo); err != nil {
		log.WithError(err).Error("unable to save Repo")
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resource := helpers.MakeRepoResource(repo)
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusCreated, payload)
}

// DeleteRepo deletes a repo from the list of enabled repositories to index
func (r *RepoHandlers) DeleteRepo(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	db, closer := r.dbSession.DB()
	defer closer()

	repo, err := models.GetRepo(db, params["repo"])
	if err != nil {
		log.WithError(err).Error("unable to find Repo")
		renderer.Render.JSON(w, http.StatusNotFound, notFoundPayload())
		return
	}

	err = models.DeleteRepo(db, params["repo"])
	if err != nil {
		log.WithError(err).Error("unable to delete Repo")
		renderer.Render.JSON(w, http.StatusInternalServerError, internalServerErrorPayload())
		return
	}

	resource := helpers.MakeRepoResource(repo)
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

func notFoundPayload() *swaggermodels.Error {
	return &swaggermodels.Error{Code: pointerto.Int64(http.StatusNotFound), Message: pointerto.String("404 repository not found")}
}

func internalServerErrorPayload() *swaggermodels.Error {
	return &swaggermodels.Error{Code: pointerto.Int64(http.StatusInternalServerError), Message: pointerto.String("Internal server error")}
}

func errorResponse(w http.ResponseWriter, errorCode int64, message string) {
	renderer.Render.JSON(w, int(errorCode), swaggermodels.Error{Code: pointerto.Int64(errorCode), Message: &message})
}
