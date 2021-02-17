package api_test

import (
	"github.com/noahfriedman-ca/quick-factor/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http/httptest"
)

var _ = Describe("the Factor function", func() {
	var w *httptest.ResponseRecorder
	BeforeEach(func() {
		w = httptest.NewRecorder()
	})

	DescribeTable("when an error should be thrown",
		func(queries string) {
			r := httptest.NewRequest("", "http://example.com?"+queries, nil)

			api.Factor(w, r)

			b, e := ioutil.ReadAll(w.Result().Body)
			Expect(e).NotTo(HaveOccurred())

			Expect(string(b)).To(ContainSubstring("ERROR:"))
		},
		Entry("should throw an error when the 'degree' query isn't present", ""),
		Entry("should throw an error when the 'degree' query isn't numeric", "degree=notanumber"),
		Entry("should throw an error when the 'degree' query isn't an integer", "degree=3.14"),
		Entry("should throw an error when the 'degree' query is < 2", "degree=1"),
	)
})
