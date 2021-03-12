package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"sort"
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
	if b, e := json.MarshalIndent(result, "", "  "); result == nil || e != nil {
		http.Error(w, "ERROR: failed to factor", http.StatusInternalServerError)
		if e != nil {
			log.Println(e)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}
}

// Implements specific factoring rules that can only be applied to a 3-term polynomial.
func factorTrinomial(coefficients []float64) *FactorJSON {
	// Definitions to avoid errors with labels
	var (
		result                        [2]float64
		resultType                    = "full"
		negativeB, discriminant, twoA float64
	)

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

	// Grouping is done differently depending on operators (+/-)
	if product := coefficients[0] * coefficients[2]; product < 0 {
		if coefficients[1] < 0 {
			for a, b := 1., coefficients[1]-1; a*b >= product; a, b = a+1, b-1 {
				if a*b == product {
					result = [2]float64{a * -1, b * -1}
				}
			}
		} else if coefficients[1] > 0 {
			for a, b := -1., coefficients[1]+1; a*b >= product; a, b = a-1, b+1 {
				if a*b == product {
					result = [2]float64{a, b}
				}
			}
		} else { // coefficients[1] == 0
			r := strconv.FormatFloat(math.Sqrt(math.Abs(product)), 'g', 5, 64)

			return &FactorJSON{
				Result: "full",
				Factored: &FactoredJSON{
					Expression: fmt.Sprintf("(x + %[1]s)(x - %[1]s)", r),
					Intercepts: []string{r, "-" + r},
				},
			}
		}
	} else { // coefficients[0] and coefficients[2] both cannot be 0, so neither can 'product'
		// If coefficients[1] == 0, it cannot be factored
		if coefficients[1] == 0 {
			return &FactorJSON{
				Result: "not",
			}
		}

		// Make sure to recall that the value should be negative
		var negative bool
		if coefficients[1] < 0 {
			coefficients[1] = math.Abs(coefficients[1])
			negative = true
		}

		half := math.Ceil(coefficients[1] / 2)
		for a, b := 1., coefficients[1]-1; a < half; a, b = a+1, b-1 {
			if a*b == product {
				if negative {
					result = [2]float64{a, b}
				} else {
					result = [2]float64{-1 * a, -1 * b}
				}

				break
			}
		}
	}

	if result == [2]float64{0, 0} {
		goto formula
	} else {
		for i := range result {
			result[i] /= coefficients[2]
		}
		goto ret
	}

formula:
	// Calculate the components of the quadratic formula
	negativeB = -1 * coefficients[1]
	discriminant = math.Pow(coefficients[1], 2) - (4 * coefficients[0] * coefficients[2])
	twoA = 2 * coefficients[2]

	// If the discriminant is negative, it has no square root, so do not factor further
	if !(discriminant > 0) {
		// Generate x-intercept strings
		intercepts := make([]string, 2)
		for i, v := range []byte{'+', '-'} {
			intercepts[i] = fmt.Sprintf("(%s %c √(%s)) / %s", formatFloat(negativeB), v, formatFloat(discriminant), formatFloat(twoA))
		}

		return &FactorJSON{
			Result: "quadratic",
			Factored: &FactoredJSON{
				Expression: fmt.Sprintf("(%s)(%s)", intercepts[0], intercepts[1]),
				Intercepts: intercepts,
			},
		}
	}

	// Factor fully using the quadratic formula and return
	discriminant = math.Sqrt(discriminant)
	result = [2]float64{(negativeB + discriminant) / twoA, (negativeB - discriminant) / twoA}
	resultType = "quadratic"
	goto ret

ret:
	sorted := result[:]
	sort.Float64s(sorted)
	res := make([]string, 2)
	for i, v := range sorted {
		res[i] = formatFloat(v)
	}

	return &FactorJSON{
		Result: resultType,
		Factored: &FactoredJSON{
			Expression: fmt.Sprintf("(x%s)(x%s)", getOp(sorted[0]*-1), getOp(sorted[1]*-1)),
			Intercepts: res,
		},
	}
}

// Implements general factorization rules that can be applied to any polynomial.
func factorPolynomial(degree uint, coefficients []float64) *FactorJSON {
	// Validate degree value
	if degree < 2 {
		log.Println("degree was smaller than 2, this shouldn't have happened")
		return nil
	} else if degree == 2 {
		log.Println("factorPolynomial was called when factorTrinomial should have been")
	}
	if int(degree) != len(coefficients)-1 {
		log.Println("degree does not match up with number of coefficients")
		return nil
	}

	// If any coefficient is not a whole number the polynomial cannot be factored.
	for _, v := range coefficients {
		if v != math.Round(v) {
			return &FactorJSON{Result: "not"}
		}
	}

	// Attempt to implement the rational root theorem
	var (
		wg            sync.WaitGroup
		interceptChan = make(chan float64)
	)
	for _, num := range findFactorsOf(uint(math.Abs(coefficients[0]))) {
		for _, den := range findFactorsOf(uint(math.Abs(coefficients[len(coefficients)-1]))) {
			wg.Add(1)
			go func(num, den float64) {
				defer wg.Done()

				x := num / den

				var resultPos, resultNeg float64
				for i := len(coefficients) - 1; i >= 0; i-- {
					resultPos += coefficients[i] * math.Pow(x, float64(i))
					resultNeg += coefficients[i] * math.Pow(x*-1, float64(i))
				}

				if resultPos == 0 {
					interceptChan <- x
				} else if resultNeg == 0 {
					interceptChan <- x * -1
				}
			}(num, den)
		}
	}

	// Create a channel to notify when the WaitGroup is empty
	doneWaiting := make(chan bool, 1)
	go func() {
		wg.Wait()
		doneWaiting <- true
		close(doneWaiting)
	}()

	var intercept float64
	select {
	case intercept = <-interceptChan:
	case <-doneWaiting: // If this happens there are no valid factors
		return &FactorJSON{Result: "not"}
	}

	// Divide the polynomial by the discovered intercept
	newCoefficients := make([]float64, len(coefficients)-1)
	newCoefficients[len(newCoefficients)-1] = coefficients[len(coefficients)-1]
	for i := len(coefficients) - 2; i > 0; i-- {
		newCoefficients[i-1] = coefficients[i] + (intercept * newCoefficients[i])
	}
	if coefficients[0]+(intercept*newCoefficients[0]) != 0 {
		log.Println("an intercept marked as valid was not")
		return nil
	}

	// Recursion
	var d *FactorJSON
	if degree > 3 {
		d = factorPolynomial(degree-1, newCoefficients)
	} else {
		d = factorTrinomial(newCoefficients)
	}

	if d.Result == "not" {
		firstCoefficient := formatFloat(newCoefficients[len(newCoefficients)-1])
		if firstCoefficient == "1" {
			firstCoefficient = ""
		} else if firstCoefficient == "-1" {
			firstCoefficient = "-"
		}

		expr := fmt.Sprintf("(%sx^%d", firstCoefficient, len(newCoefficients)-1)

		for i := len(newCoefficients) - 2; i > 1; i-- {
			expr += fmt.Sprintf("%sx^%d", getOp(newCoefficients[i]), i)
		}
		expr += fmt.Sprintf("%sx%s)(x%s)", getOp(newCoefficients[1]), getOp(newCoefficients[0]), getOp(intercept*-1))

		return &FactorJSON{
			Result: "partial",
			Factored: &FactoredJSON{
				Expression: expr,
				Intercepts: []string{formatFloat(intercept)},
			},
		}
	} else {
		d.Factored.Expression += fmt.Sprintf("(x%s)", getOp(intercept*-1))
		d.Factored.Intercepts = append(d.Factored.Intercepts, formatFloat(intercept))
		sort.Slice(d.Factored.Intercepts, func(i, j int) bool {
			var strs [2]string
			strs[0] = d.Factored.Intercepts[i]
			strs[1] = d.Factored.Intercepts[j]

			var isQuad, isNeg [2]bool
			for i, v := range strs {
				if strings.Contains(v, "√") {
					isQuad[i] = true

					if strings.Contains(v, "- √") {
						isNeg[i] = true
					}
				}
			}

			if isQuad[0] && !isQuad[1] {
				return false
			} else if !isQuad[0] && isQuad[1] {
				return true
			} else if isQuad[0] && isQuad[1] {
				return isQuad[1]
			} else {
				var flts [2]float64
				for i, v := range strs {
					flts[i], _ = strconv.ParseFloat(v, 64)
				}
				return flts[0] < flts[1]
			}
		})
		for i := 0; i < (len(d.Factored.Intercepts) - 1); i++ {
			if d.Factored.Intercepts[i] == d.Factored.Intercepts[i+1] {
				if i+2 < len(d.Factored.Intercepts) {
					d.Factored.Intercepts = append(d.Factored.Intercepts[:i+1], d.Factored.Intercepts[i+2:]...)
				} else {
					d.Factored.Intercepts = d.Factored.Intercepts[:len(d.Factored.Intercepts)-1]
				}
			}
		}

		var (
			sortedExpr string
			regex      = regexp.MustCompile(`\(x ([+-]) ([0-9])*\)(?:\^([0-9]*))?`)
			exprs      = regex.FindAllStringSubmatch(d.Factored.Expression, -1)
		)
		if s := regexp.MustCompile(`^[^)]*\)`).FindString(d.Factored.Expression); !regex.MatchString(s) {
			sortedExpr += s
		}

		sort.Slice(exprs, func(i, j int) bool {
			var vals, exps [2]float64
			for i, v := range []int{i, j} {
				matches := regex.FindStringSubmatch(exprs[v][0])

				r, _ := strconv.ParseFloat(matches[2], 64)
				if matches[1][0] == '+' {
					r *= -1
				}
				vals[i] = r

				if matches[3] == "" {
					exps[i] = 1
				} else {
					exps[i], _ = strconv.ParseFloat(matches[3], 64)
				}
			}

			if exps[0] != exps[1] {
				return exps[0] < exps[1]
			} else {
				return vals[0] < vals[1]
			}
		})

		for _, v := range exprs {
			sortedExpr += v[0]
		}

		d.Factored.Expression = sortedExpr

		return d
	}
}

// Calculate factors of 'x'. 'x' must be positive.
func findFactorsOf(x uint) []float64 {
	var (
		wg      sync.WaitGroup
		factors = make(chan float64, x) // Buffer the maximum amount of values that can be sent, though it will likely never be necessary to have that many values in the channel
	)

	// Check that x is not 0
	if x == 0 {
		log.Println("there are no factors of 0, and this shouldn't have happened")
		return nil
	}

	// Calculate if each value from 2 to x is a factor
	for i := uint(2); i < (x - 1); i++ {
		wg.Add(1)
		go func(i uint) {
			defer wg.Done()
			if f := float64(x) / float64(i); math.Round(f) == f {
				factors <- f
			}
		}(i)
	}

	// Ensure channel gets closed at the proper time
	go func() {
		wg.Wait()
		close(factors)
	}()

	// Move channel values into array and return
	r := []float64{1} // 1 is a factor of everything
	for v := range factors {
		r = append(r, v)
	}

	// Sort the array and return
	sort.Float64s(r)
	return r
}

// Return v formatted with the correct operator in front of it
//  getOp(-45) -> " - 45"
func getOp(v float64) string {
	var r string
	if v < 0 {
		r = " - "
	} else {
		r = " + "
	}

	return r + formatFloat(math.Abs(v))
}

// Take 0's off the end of floats
func formatFloat(v float64) string {
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(v, 'f', 5, 64), "0"), ".")
}
