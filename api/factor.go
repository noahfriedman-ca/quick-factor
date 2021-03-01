package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func init() {
	funcs = append(funcs, Factor)
}

// Struct defining the JSON response from the API function.
type FactorJSON struct {
	Result   string        `json:"result"` // A string representing the factoring result; either "full", "quadratic", "partial", or "not"
	Factored *FactoredJSON `json:"factored,omitempty"`
}

// Struct defining the JSON representation of a factored polynomial.
type FactoredJSON struct {
	Expression string   `json:"expression"`
	Intercepts []string `json:"intercepts,omitempty"`
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
	coefficients := make([]float64, degree+1)
	for i := range coefficients {
		// Extra whitespace is trimmed to avoid unintentional errors
		if s := strings.Trim(q.Get(fmt.Sprintf("x^%d", i)), " "); s == "" {
			// If an empty string is found, assume a value of 0
			coefficients[i] = 0
		} else if f, e := strconv.ParseFloat(s, 64); e != nil {
			http.Error(w, fmt.Sprintf("ERROR: could not parse value in query parameter 'x^%d'", i), http.StatusExpectationFailed)
			return
		} else if (i == 0 || i == int(degree)) && f == 0 {
			http.Error(w, fmt.Sprintf("ERROR: x^%d must not be 0", i), http.StatusExpectationFailed)
			return
		} else {
			coefficients[i] = f
		}
	}

	// Do the actual factoring
	var result *FactorJSON
	if degree == 2 {
		result = factorTrinomial(coefficients)
	} else {
		result = factorPolynomial(degree, coefficients)
	}

	// Write the result
	if b, e := json.Marshal(result); result == nil || e != nil {
		http.Error(w, "ERROR: failed to factor", http.StatusInternalServerError)
		if e == nil {
			log.Println("nil result")
		} else {
			log.Println(e)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}
}

// Implements specific factoring rules that can only be applied to a 3-term polynomial.
func factorTrinomial(coefficients []float64) *FactorJSON {
	// Ensure that there are the correct number of coefficients
	if len(coefficients) != 3 {
		log.Println("factorTrinomial called with not exactly 3 coefficients")
		return nil
	}

	// Check if any of the coefficients are not integers. If they are skip ahead to using the quadratic formula.
	for _, v := range coefficients {
		if v != math.Round(v) {
			goto formula
		}
	}

	// Attempt to factor by grouping

formula:
	// Calculate the components of the quadratic formula
	var (
		negativeB    = -1 * coefficients[1]
		discriminant = math.Pow(coefficients[1], 2) - (4 * coefficients[0] * coefficients[2])
		twoA         = 2 * coefficients[2]
	)

	// If the discriminant is negative, it has no square root, so do not factor further
	if !(discriminant > 0) {
		// Generate x-intercept strings
		intercepts := make([]string, 2)
		for i, v := range []byte{'+', '-'} {
			intercepts[i] = fmt.Sprintf("(%f %c √(%f)) / %f", negativeB, v, discriminant, twoA)
		}

		return &FactorJSON{
			Result: "quadratic",
			Factored: &FactoredJSON{
				Expression: fmt.Sprintf("(%s)(%s)", intercepts[0], intercepts[1]),
				Intercepts: intercepts,
			},
		}
	}

	return nil
}

// Implements general factorization rules that can be applied to any polynomial.
func factorPolynomial(degree uint, coefficients []float64) *FactorJSON {
	return nil
}

// Calculate factors of 'x'. 'x' must be positive.
func findFactorsOf(x uint) chan float64 {
	var (
		wg      sync.WaitGroup
		factors = make(chan float64, x) // Buffer the maximum amount of values that can be sent, though it will likely never be necessary to have that many values in the channel
	)

	// Check that x is not 0
	if x == 0 {
		log.Println("there are no factors of 0, and this shouldn't have happened")
		return nil
	}

	// Send '1' to the channel because everything can be factored by 1
	factors <- 1

	// Calculate if each value from 2 to x is a factor
	go func() {
		for i := uint(2); i < (x - 1); i++ {
			wg.Add(1)
			go func(i uint) {
				defer wg.Done()
				if f := float64(x) / float64(i); math.Round(f) == f {
					factors <- f
				}
			}(i)
		}
	}()

	// Ensure channel gets closed at the proper time
	go func() {
		wg.Wait()
		close(factors)
	}()

	// Return the channel because it can be iterated over, no reason to convert to array
	return factors
}
