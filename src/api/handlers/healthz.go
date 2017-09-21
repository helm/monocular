package handlers

import (
	"net/http"
)

// Healthz returns 200
func Healthz(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
