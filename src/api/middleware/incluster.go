package middleware

import (
	"net/http"

	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/urfave/negroni"

	"github.com/kubernetes-helm/monocular/src/api/handlers/renderer"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// InClusterGate implements middleware to check if in cluster features are enabled before continuing
func InClusterGate(inCluster bool) negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		if !inCluster {
			renderer.Render.JSON(w, http.StatusForbidden,
				models.Error{Code: pointerto.Int64(http.StatusForbidden), Message: pointerto.String("feature not enabled")})
		} else {
			next(w, req)
		}
	}
}
