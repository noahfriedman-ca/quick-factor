package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func init() {
	funcs = append(funcs, Factor)
}

func Factor(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if s := q.Get("degree"); s == "" {
		http.Error(w, "ERROR: Missing required query parameter 'degree'", http.StatusExpectationFailed)
	} else if d, e := strconv.Atoi(s); e != nil || d < 2 {
		http.Error(w, "ERROR: Query parameter 'degree' must be an integer >= 2", http.StatusExpectationFailed)
	} else {
		fmt.Println(d)
	}
}
