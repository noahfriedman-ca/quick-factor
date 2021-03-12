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
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var _ = Describe("the Factor function", func() {
	rand.Seed(1615477392)

	getResponse := func(queries string) string {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://example.com?"+queries, nil)

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
				return []float64{10, 7, 1}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "full",
					Factored: &api.FactoredJSON{
						Expression: "(x + 5)(x + 2)",
						Intercepts: []string{"-5", "-2"},
					},
				}
			},
		),
		Entry("should factor '3x^2 - 12' into '(x + 6)(x - 6)'",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{3, 0, -12}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "full",
					Factored: &api.FactoredJSON{
						Expression: "(x + 6)(x - 6)",
						Intercepts: []string{"6", "-6"},
					},
				}
			},
		),
		Entry("should factor 'x^3 - 2x^2 - 5x + 6' into '(x + 2)(x - 1)(x - 3)'",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{6, -5, -2, 1}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result:   "full",
					Factored: &api.FactoredJSON{Expression: "(x + 2)(x - 1)(x - 3)", Intercepts: []string{"-2", "1", "3"}},
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
		Entry("should be able to solve using the quadratic formula",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{-2, 10, 1}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "quadratic",
					Factored: &api.FactoredJSON{
						Expression: "(x + 10.19615)(x - 0.19615)",
						Intercepts: []string{"-10.19615", "0.19615"},
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
					Skip("took too long to generate coefficients")
					return nil, nil
				}

				ints := make([]string, 2)
				for i, v := range []byte{'+', '-'} {
					var parts [3]string
					for i, v := range []float64{coefficients[1] * -1, math.Pow(coefficients[1], 2) - (4 * coefficients[2] * coefficients[0]), 2 * coefficients[2]} {
						parts[i] = strings.TrimRight(strings.TrimRight(strconv.FormatFloat(v, 'f', 5, 64), "0"), ".")
					}
					ints[i] = fmt.Sprintf("(%s %c √(%s)) / %s", parts[0], v, parts[1], parts[2])
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
		Entry("should be able to factor a randomly generated factorable polynomial",
			func() ([]float64, []string) {
				// Randomly define degree to be in range 2-6
				degree := rand.Intn(4) + 2

				// Randomly generate intercepts
				intercepts := make([]float64, degree)
				for i := 0; i < degree; i++ {
					var r float64
				noDupe:
					for i := 0; i < 5; i++ {
						r = float64(rand.Intn(8) + 1)
						if rand.Intn(2) == 1 {
							r *= -1
						}

						for _, v := range intercepts {
							if r == v {
								r = 0
								continue noDupe
							}
						}

						break noDupe
					}
					if r == 0 {
						Skip("generating intercepts took too long")
					}

					intercepts[i] = r
				}

				// Sort intercept array to make sure test doesn't fail incorrectly
				sort.Float64s(intercepts)

				// Convert intercepts to strings asynchronously
				var (
					interceptStrings = make([]string, len(intercepts))
					done             = make(chan bool, 1)
				)
				go func() {
					for i, v := range intercepts {
						interceptStrings[i] = strconv.Itoa(int(v))
					}
					done <- true
					close(done)
				}()

				// Define array for coefficients
				var coefficients []float64

				coefficients = []float64{-intercepts[0], 1}

				// Continue expanding until there aren't any intercepts left
				for _, v := range intercepts[1:] {
					// Multiply the polynomial by x
					xProducts := append([]float64{0}, coefficients...)

					// Multiply the polynomial by the intercept
					intProducts := make([]float64, len(coefficients))
					for i, a := range coefficients {
						intProducts[i] = a * -v
					}

					// Add those results together
					products := make([]float64, len(xProducts))
					for i, v := range xProducts {
						if i >= len(intProducts) {
							products[i] = v
						} else {
							products[i] = v + intProducts[i]
						}
					}

					coefficients = products
				}

				// Try to reduce all coefficients
				smallest := math.Abs(coefficients[0])
				for _, v := range coefficients[1:] {
					v = math.Abs(v)
					if v < smallest {
						smallest = v
					}
				}
			reduce:
				for i := 2.; i <= smallest; i++ {
					r := make([]float64, len(coefficients))
					for a, v := range coefficients {
						r[a] = v / i
					}
					for _, v := range r {
						if v != math.Round(v) {
							continue reduce
						}
					}

					// If the loop has gotten this far, the reduction was successful
					i = 1
					coefficients = r
				}

				// Use a By field to log the resulting polynomial
				By(fmt.Sprintf("using polynomial from array: %v\n\twith intercepts: %v", coefficients, intercepts))

				// Make sure all intercepts have been converted to strings, then return
				<-done
				return coefficients, interceptStrings
			},
			func(ints []string) *api.FactorJSON {
				// Remove duplicate intercepts
				intMap := make(map[string]int)
				for _, v := range ints {
					intMap[v] += 1
				}

				var (
					newInts []string
					done    = make(chan bool, 1)
				)
				go func() {
					for k := range intMap {
						newInts = append(newInts, k)
					}

					sort.Slice(newInts, func(i, j int) bool {
						one, _ := strconv.Atoi(newInts[i])
						two, _ := strconv.Atoi(newInts[j])

						return one < two
					})
					done <- true
					close(done)
				}()

				var exprs []string
				for k, v := range intMap {
					r := "(x "
					if k[0] == '-' {
						r += "+ " + strings.TrimLeft(k, "-")
					} else {
						r += "- " + k
					}
					r += ")"
					if v >= 2 {
						r += fmt.Sprintf("^%d", v)
					}

					exprs = append(exprs, r)
				}

				regex := regexp.MustCompile(`\(x ([+-]) ([0-9]*)\)(?:\^([0-9]*))?`)
				sort.Slice(exprs, func(i, j int) bool {
					var vals, exps [2]float64
					for i, v := range []int{i, j} {
						r := regex.FindStringSubmatch(exprs[v])

						vals[i], _ = strconv.ParseFloat(r[2], 64)
						if r[1] == "+" {
							vals[i] *= -1
						}

						if r[3] == "" {
							exps[i] = 1
						} else {
							exps[i], _ = strconv.ParseFloat(r[3], 64)
						}
					}

					if exps[0] != exps[1] {
						return exps[0] < exps[1]
					} else {
						return vals[0] < vals[1]
					}
				})

				var expr string
				for _, v := range exprs {
					expr += v
				}

				<-done
				return &api.FactorJSON{
					Result: "full",
					Factored: &api.FactoredJSON{
						Expression: expr,
						Intercepts: newInts,
					},
				}
			},
		),
		Entry("should partially factor x^4 + 2x^3 - 37x^2 + 14x - 20 into (x^3 + 7x^2 - 2x + 4)(x - 5)",
			func() ([]float64, []string) {
				// No intercept array is returned because its not needed
				return []float64{-20, 14, -37, 2, 1}, nil
			},
			func(_ []string) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "partial",
					Factored: &api.FactoredJSON{
						Expression: "(x^3 + 7x^2 - 2x + 4)(x - 5)",
						Intercepts: []string{"5"},
					},
				}
			},
		),
	)
})
