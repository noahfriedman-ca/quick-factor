package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
)

var funcs []http.HandlerFunc
var funcsJSON = struct {
	Available []string `json:"available"`
}{}

// Create a router configured properly for this program.
func Router() *mux.Router {
	rtr := mux.NewRouter()
	r := rtr.PathPrefix("/projects/quick-factor/api").Subrouter()

	for _, v := range funcs {
		// Extract function names
		t := string(regexp.MustCompile(`\.([^.]*)$`).FindSubmatch([]byte(runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()))[1])

		// Map all functions
		funcsJSON.Available = append(funcsJSON.Available, t)
		r.Path("/" + t).HandlerFunc(v)
	}

	rtr.PathPrefix("/projects/quick-factor/api").HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if b, e := json.MarshalIndent(funcsJSON, "", "    "); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(b)
		}
	})

	return rtr
}
