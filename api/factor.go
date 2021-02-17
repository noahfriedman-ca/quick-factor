package api

import "net/http"

func init() {
	funcs = append(funcs, Factor)
}

func Factor(w http.ResponseWriter, r *http.Request) {

}
