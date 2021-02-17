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
		for _, v := range []string{"funcOne", "funcTwo", "func_three"} {
			r, e := getResponse(v)
			Expect(e).NotTo(HaveOccurred())

			s := string(r)
			Expect(s).NotTo(ContainSubstring("404 page not found"))
			Expect(s).NotTo(ContainSubstring("\"available\":")) // This is done to check if the "help" text is displayed
		}
	})
	It("should lowercase the first character of the function name", func() {
		f1, e := getResponse("funcOne")
		Expect(e).NotTo(HaveOccurred())

		F1, e := getResponse("FuncOne")
		Expect(e).NotTo(HaveOccurred())

		sf1 := string(f1)
		sF1 := string(F1)
		Expect(sf1).NotTo(ContainSubstring("404 page not found"))
		Expect(sf1).NotTo(ContainSubstring("\"available\":"))
		Expect(sF1).NotTo(ContainSubstring("404 page not found"))
		Expect(sF1).To(ContainSubstring("\"available\":"))
	})
	DescribeTable("when the list of available functions should be displayed",
		func(path string) {
			var resp struct {
				Available []string `json:"available"`
			}

			r, e := getResponse(path)
			Expect(e).NotTo(HaveOccurred())
			Expect(string(r)).NotTo(ContainSubstring("404 page not found"))

			Expect(json.Unmarshal(r, &resp)).To(Succeed())

			Expect(resp.Available).To(ContainElements("funcOne", "funcTwo", "func_three"))
		},
		Entry("the root path", ""),
		Entry("any unrecognized path", "/test/path"),
	)
})
