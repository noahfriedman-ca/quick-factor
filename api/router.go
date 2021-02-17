package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var funcs []http.HandlerFunc

// Create a router configured properly for this program.
func Router() *mux.Router {
	rtr := mux.NewRouter()
	r := rtr.PathPrefix("/projects/quick-factor/api").Subrouter()

	var funcsJSON = struct {
		Available []string `json:"available"`
	}{}

	for _, v := range funcs {
		// Extract function names
		t := string(regexp.MustCompile(`\.([^.]*)$`).FindSubmatch([]byte(runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()))[1])

		// Convert the first character to lowercase
		t = strings.ToLower(t[0:1]) + t[1:]

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
