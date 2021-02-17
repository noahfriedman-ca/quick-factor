package api

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func FuncOne(_ http.ResponseWriter, _ *http.Request)    {}
func funcTwo(_ http.ResponseWriter, _ *http.Request)    {}
func func_three(_ http.ResponseWriter, _ *http.Request) {}

var _ = Describe("the router", func() {
	BeforeSuite(func() {
		funcs = []http.HandlerFunc{
			FuncOne,
			funcTwo,
			func_three,
		}
	})

	var ms *httptest.Server
	BeforeEach(func() {
		ms = httptest.NewServer(Router())
	})

	getResponse := func(path string) ([]byte, error) {
		if r, e := http.Get(ms.URL + "/projects/quick-factor/api/" + path); e != nil {
			return nil, e
		} else if b, e := ioutil.ReadAll(r.Body); e != nil {
			return nil, e
		} else {
			return b, nil
		}
	}

	It("should route each function in the 'funcs' array", func() {
		for _, v := range []string{"FuncOne", "funcTwo", "func_three"} {
			r, e := getResponse(v)
			Expect(e).NotTo(HaveOccurred())
			Expect(string(r)).NotTo(ContainSubstring("404 page not found"))
		}
	})

	DescribeTable("when the list of available functions should be displayed",
		func(path string) {
			var resp struct {
				Available []string `json:"available"`
			}

			r, e := getResponse(path)
			Expect(e).NotTo(HaveOccurred())

			Expect(json.Unmarshal(r, &resp)).To(Succeed())

			Expect(resp.Available).To(ContainElements("FuncOne", "funcTwo", "func_three"))
		},
		Entry("the root path", ""),
		Entry("any unrecognized path", "/test/path"),
	)
})
