package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/noahfriedman-ca/quick-factor/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http/httptest"
	"strconv"
)

var _ = Describe("the Factor function", func() {
	rand.Seed(GinkgoRandomSeed())

	getResponse := func(queries string) string {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "http://example.com?"+queries, nil)

		api.Factor(w, r)

		b, e := ioutil.ReadAll(w.Result().Body)
		Expect(e).NotTo(HaveOccurred())

		return string(b)
	}

	DescribeTable("when an error should be thrown",
		func(degree string) {
			resp := getResponse("degree=" + degree)
			Expect(resp).To(ContainSubstring("ERROR:"))
		},
		Entry("should throw an error when the 'degree' query isn't present", ""),
		Entry("should throw an error when the 'degree' query isn't numeric", "notanumber"),
		Entry("should throw an error when the 'degree' query isn't an integer", "3.14"),
		Entry("should throw an error when the 'degree' query is < 2", "1"),
	)

	DescribeTable("the results when attempting to factor certain polynomials",
		func(genPolynomial func() ([]float64, []string), expected func(intercepts []string) *api.FactorJSON) {
			polynomial, intercepts := genPolynomial()
			queries := fmt.Sprintf("degree=%d", len(polynomial)-1)
			for i, v := range polynomial {
				queries += fmt.Sprintf("&x^%d=%f", i, v)
			}

			resp := getResponse(queries)
			Expect(resp).NotTo(ContainSubstring("ERROR:"))

			var respJSON api.FactorJSON
			Expect(json.Unmarshal([]byte(resp), &respJSON)).To(Succeed())

			Expect(respJSON).To(Equal(*expected(intercepts)))
		},
		Entry("should factor 'x^2 + 7x + 10' into '(x - 2)(x - 5)'",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{1, 7, 10}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "full",
					Factored: &api.FactoredJSON{
						Expression: "(x - 2)(x - 5)",
						Intercepts: []string{"2", "5"},
					},
				}
			},
		),
		Entry("should factor '3x^2 - 12' into '(x - 6)^2'",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{3, 0, -12}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "full",
					Factored: &api.FactoredJSON{
						Expression: "(x - 6)^2",
						Intercepts: []string{"6", "-6"},
					},
				}
			},
		),
		Entry("should factor 'x^3 - 2x^2 - 5x + 6' into '(x - 1)(x + 2)(x - 3)'",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{6, -5, -2, 1}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result:   "full",
					Factored: &api.FactoredJSON{Expression: "(x - 1)(x + 2)(x - 3)", Intercepts: []string{"1", "-2", "3"}},
				}
			}),
		Entry("should be unable to factor '2x^3 + 7x^2 + 4",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{4, 0, 7, 2}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{Result: "not"}
			}),
		Entry("should be able to factor a randomly generated factorable polynomial",
			func() ([]float64, []string) {
				var (
					a    = make([]float64, 3)
					ints = make([]string, 3)
				)
				for i := range a {
					r := rand.Intn(8) + 1
					b := rand.Intn(1) != 0

					if b {
						r *= -1
					}

					a[i] = float64(r)
					ints[i] = strconv.Itoa(r)
				}

				// Expand (x - a0)(x - a1)(x - a2) into x^3 - (a0 + a1 + a2)x^2 + (a0a1 + (a0 + a1)a2)x - a0a1a2
				expanded := []float64{a[0] * a[1] * a[2], (a[0] * a[1]) + ((a[0] + a[1]) * a[2]), a[0] + a[1] + a[2], 1}

				return expanded, ints
			},
			func(ints []string) *api.FactorJSON {
				// Extract proper operators based on intercept values
				var ops [3]string
				for i, v := range ints {
					if v[0] == '-' {
						ops[i] = "+"
					} else {
						ops[i] = "-"
					}
				}

				// Define expected result
				return &api.FactorJSON{
					Result: "full",
					Factored: &api.FactoredJSON{
						Expression: fmt.Sprintf("(x %s %s)(x %s %s)(x %s %s)", ops[0], ints[0], ops[1], ints[1], ops[2], ints[2]),
						Intercepts: ints,
					},
				}
			},
		),
		Entry("should be able to solve using the quadratic formula",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{-2, 10, 1}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "quadratic",
					Factored: &api.FactoredJSON{
						Expression: "(x - 0.19615)(x + 10.19615)",
						Intercepts: []string{"0.19615", "-10.19615"},
					},
				}
			},
		),
		Entry("should factor into (-b ± √(b²-4ac)) ÷ 2a) form when the discriminant is negative",
			func() ([]float64, []string) {
				// Randomly generate coefficients
				coefficients := make([]float64, 3)
			generate:
				for i := 0; i < 5; i++ {
					for i := range coefficients {
						coefficients[i] = float64(rand.Intn(1999)+1) / 100
					}

					// Check if b² is larger than or equal to 4ac. If it is, restart the loop.
					if math.Pow(coefficients[1], 2) >= (4 * coefficients[0] * coefficients[2]) {
						coefficients = make([]float64, 3)
						continue generate
					}

					// Make sure that a trinomial with these coefficients cannot be factored by grouping
					for _, v := range coefficients {
						// If any of the coefficients are not whole numbers, the whole thing can't be factored by grouping
						if v != math.Round(v) {
							break generate
						}
					}

					// If the loop is still running, all coefficients must be whole numbers
					// Now, check if the trinomial can be factored by grouping
					var (
						product = coefficients[0] * coefficients[2]
						half    = math.Ceil(coefficients[1] / 2)
					)
					for a, b := 1., coefficients[1]-1; a < half; a, b = a+1, b-1 {
						if a*b == product {
							// If the trinomial is factorable by grouping, reset the array and restart the loop
							coefficients = make([]float64, 3)
							continue generate
						}
					}
				}

				// If coefficients[0] is 0, then everything must be 0 (because 0 cannot be randomly generated)
				if coefficients[0] == 0 {
					Fail("took too long to generate coefficients")
					return nil, nil
				}

				ints := make([]string, 2)
				for i, v := range []byte{'+', '-'} {
					ints[i] = fmt.Sprintf("(%f %c √(%f)) / %f", coefficients[1]*-1, v, math.Pow(coefficients[1], 2)-(4*coefficients[2]*coefficients[0]), 2*coefficients[2])
				}

				return coefficients, ints
			},
			func(ints []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "quadratic",
					Factored: &api.FactoredJSON{
						Expression: fmt.Sprintf("(%s)(%s)", ints[0], ints[1]),
						Intercepts: ints,
					},
				}
			},
		),
	)
})
