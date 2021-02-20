package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/noahfriedman-ca/quick-factor/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
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
		func(genPolynomial func() ([]float64, []float64), expected func(intercepts []float64) *api.FactorJSON) {
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
		Entry("should factor 'x^3 - 2x^2 - 5x + 6' into '(x - 1)(x + 2)(x - 3)'",
			func() ([]float64, []float64) {
				// No intercept array is returned because its not needed
				return []float64{6, -5, -2, 1}, nil
			},
			func(_ []float64) *api.FactorJSON {
				return &api.FactorJSON{
					Result: "full",
					Factored: struct {
						Expression string    `json:"expression"`
						Intercepts []float64 `json:"intercepts,omitempty"`
					}{Expression: "(x - 1)(x + 2)(x - 3)", Intercepts: []float64{}},
				}
			}),
		Entry("should be unable to factor '2x^3 + 7x^2 + 4",
			func() ([]float64, []float64) {
				// No intercept array is returned because its not needed
				return []float64{4, 0, 7, 2}, nil
			},
			func(_ []float64) *api.FactorJSON {
				return &api.FactorJSON{Result: "not"}
			}),
		Entry("should be able to factor a randomly generated factorable polynomial",
			func() ([]float64, []float64) {
				a := make([]float64, 3)
				for i := range a {
					r := rand.Intn(8) + 1
					b := rand.Intn(1) != 0

					if b {
						r *= -1
					}

					a[i] = float64(r)
				}

				// Expand (x - a0)(x - a1)(x - a2) into x^3 - (a0 + a1 + a2)x^2 + (a0a1 + (a0 + a1)a2)x - a0a1a2
				expanded := []float64{a[0] * a[1] * a[2], (a[0] * a[1]) + ((a[0] + a[1]) * a[2]), a[0] + a[1] + a[2], 1}

				return expanded, a
			},
			func(ints []float64) *api.FactorJSON {
				// Extract proper operators based on intercept values
				var ops [3]string
				for i, v := range ints {
					if v < 0 {
						ops[i] = "+"
					} else {
						ops[i] = "-"
					}
				}

				// Define expected result
				return &api.FactorJSON{
					Result: "full",
					Factored: struct {
						Expression string    `json:"expression"`
						Intercepts []float64 `json:"intercepts,omitempty"`
					}{
						Expression: fmt.Sprintf("(x %s %f)(x %s %f)(x %s %f)", ops[0], ints[0], ops[1], ints[1], ops[2], ints[2]),
						Intercepts: ints,
					},
				}
			},
		),
	)
})
