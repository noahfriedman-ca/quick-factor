package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	funcs = append(funcs, Factor)
}

// Struct defining the JSON response from the API function.
type FactorJSON struct {
	Result   string `json:"result"` // A string representing the factoring result; either "full", "partial", or "not"
	Factored struct {
		Expression string    `json:"expression"`
		Intercepts []float64 `json:"intercepts,omitempty"`
	} `json:"factored,omitempty"`
}

// API function for factoring polynomials
func Factor(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var degree uint
	if s := strings.Trim(q.Get("degree"), " "); s == "" { // Extra whitespace is trimmed to avoid unintentional errors
		http.Error(w, "ERROR: Missing required query parameter 'degree'", http.StatusExpectationFailed)
		return
	} else if d, e := strconv.Atoi(s); e != nil || d < 2 {
		http.Error(w, "ERROR: Query parameter 'degree' must be an integer >= 2", http.StatusExpectationFailed)
		return
	} else {
		degree = uint(d) // Can be safely converted to uint because it must be a positive integer
	}

	// Extract each exponent from query parameters
	exps := make([]float64, degree+1)
	for i := range exps {
		// Extra whitespace is trimmed to avoid unintentional errors
		if s := strings.Trim(q.Get(fmt.Sprintf("x^%d", i)), " "); s == "" {
			// If an empty string is found, assume a value of 0
			exps[i] = 0
		} else if f, e := strconv.ParseFloat(s, 64); e != nil {
			http.Error(w, fmt.Sprintf("ERROR: could not parse value in query parameter 'x^%d'", i), http.StatusExpectationFailed)
			return
		} else if (i == 0 || i == int(degree)) && f == 0 {
			http.Error(w, fmt.Sprintf("ERROR: x^%d must not be 0", i), http.StatusExpectationFailed)
			return
		} else {
			exps[i] = f
		}
	}

}

// The actual mathematical logic for this API function.
func actualFactoring(degree uint, exponents ...float64) (uint, []float64) {
	if degree < 2 {
		// This shouldn't be a possible outcome, which is why panic is permitted
		log.Panicln("degree is smaller than 2")
		return 0, nil
	}

	return 0, nil
}
