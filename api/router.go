package api

import (
	"github.com/gorilla/mux"
)

// Create a router configured properly for this program.
func Router() *mux.Router {
	r := mux.NewRouter()
	return r
}
