package handlers

import (
	"net/http"

	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"github.com/kubernetes-helm/monocular/src/api/models"
)

// UserHandlers defines handlers that serve user data
type UserHandlers struct {
	dbSession            datastore.Session
	chartsImplementation data.Charts
}

// NewUserHandlers takes a datastore.Session implementation and returns a UserHandlers struct
func NewUserHandlers(dbSession datastore.Session, chartsImplementation data.Charts) *UserHandlers {
	return &UserHandlers{dbSession, chartsImplementation}
}

// StarChart marks the chart as starred by the user
func (u *UserHandlers) StarChart(w http.ResponseWriter, r *http.Request, params Params) {
	user := currentUser(r)
	db, closer := u.dbSession.DB()
	defer closer()

	if _, err := u.chartsImplementation.ChartFromRepo(params["repo"], params["chart"]); err != nil {
		errorResponse(w, http.StatusNotFound, "404 chart not found")
		return
	}

	// Create chart in case it does not exist
	// TODO: once chart storage is moved to database, we shouldn't need to create it here
	chart := &models.Chart{Name: params["chart"], Repo: params["repo"]}
	if err := models.CreateChart(db, chart); err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not save Chart")
		return
	}
	chart, err := models.GetChartByName(db, params["repo"], params["chart"])
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not find Chart")
		return
	}

	if err := user.StarChart(db, chart); err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not star Chart")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnstarChart unmarks the chart as starred by the user
func (u *UserHandlers) UnstarChart(w http.ResponseWriter, r *http.Request, params Params) {
	user := currentUser(r)
	db, closer := u.dbSession.DB()
	defer closer()

	if _, err := u.chartsImplementation.ChartFromRepo(params["repo"], params["chart"]); err != nil {
		errorResponse(w, http.StatusNotFound, "404 chart not found")
		return
	}

	chart, err := models.GetChartByName(db, params["repo"], params["chart"])
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not find Chart")
		return
	}

	if err := user.UnstarChart(db, chart); err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not unstar Chart")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func currentUser(r *http.Request) *models.User {
	return r.Context().Value(models.UserKey).(models.UserClaims).User
}
