package urlshortener

import (
	"net/http"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/woraphol-j/url-shortener/pkg/mongo"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEndpoints(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Endpoint Suite")
}

var _ = Describe("Book", func() {
	var httpExpect *httpexpect.Expect

	BeforeSuite(func() {
		urlDAO := mongo.NewDAO("", "", "")
		service := NewService(urlDAO)
		var logger log.Logger
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		httpLogger := log.With(logger, "component", "http")
		handler := MakeHandler(service, httpLogger)

		httpExpect = httpexpect.WithConfig(httpexpect.Config{
			BaseURL: "http://localhost",
			Client: &http.Client{
				Transport: httpexpect.NewBinder(handler),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
			Printers: []httpexpect.Printer{
				httpexpect.NewDebugPrinter(t, true),
			},
		})
	})

	BeforeEach(func() {
		// httptest.NewServer(handler)
		// meal = DAO{
		// 	Name: "Les Miserables",
		// 	Date: time.Now().UTC(),
		// }
	})

	It("should create and retrieve the short url code", func() {
		// Expect(meal.Name).To(Equal("Les Miserables"))
		Expect("Les Miserables").To(Equal("Les Miserables"))
	})
})
