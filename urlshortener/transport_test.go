package urlshortener

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	cg "github.com/woraphol-j/url-shortener/pkg/codegenerator"
	"github.com/woraphol-j/url-shortener/pkg/mongo"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

type GinkgoTestReporter struct{}

func (g GinkgoTestReporter) Errorf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g GinkgoTestReporter) Fatalf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func TestTransport(t *testing.T) {
	RegisterFailHandler(Fail)
	testResult := os.Getenv("TEST_RESULTS")
	if testResult != "" {
		junitReporter := reporters.NewJUnitReporter(testResult)
		RunSpecsWithDefaultAndCustomReporters(
			t,
			"Test Transport Suite",
			[]Reporter{junitReporter},
		)
	} else {
		RunSpecs(t, "Test Transport")
	}
}

var _ = Describe("Transport", func() {
	const (
		baseURL         = "http://localhost:8080"
		code            = "abcde"
		nonExistingCode = "wxyz"
		originalURL     = "http://www.google.com"
	)

	var (
		tt         GinkgoTestReporter
		mockCtrl   *gomock.Controller
		httpExpect *httpexpect.Expect
		urlDAO     mongo.Repository
	)

	BeforeSuite(func() {
		os.Setenv("BASE_URL", baseURL)
		mongoURL := os.Getenv("MONGO_URL")
		mongoDb := os.Getenv("MONGO_DATABASE")
		mongoColl := os.Getenv("MONGO_COLLECTION")
		urlDAO = mongo.NewMongoRepository(mongoURL, mongoDb, mongoColl)
		mockCtrl = gomock.NewController(GinkgoT())
		mockCodeGenerator := cg.NewMockCodeGenerator(mockCtrl)
		mockCodeGenerator.EXPECT().Generate().Return(code, nil)

		service := NewService(urlDAO, mockCodeGenerator)
		var logger log.Logger
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		httpLogger := log.With(logger, "component", "http")
		handler := MakeHandler(service, httpLogger)

		httpExpect = httpexpect.WithConfig(httpexpect.Config{
			BaseURL: baseURL,
			Client: &http.Client{
				Transport: httpexpect.NewBinder(handler),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: httpexpect.NewAssertReporter(tt),
		})
	})

	BeforeEach(func() {
		urlDAO.Truncate()
	})

	It("should create a new entry correctly", func() {
		shortURLRequest := map[string]interface{}{
			"url": originalURL,
		}
		httpExpect.POST("/shorturls").
			WithJSON(shortURLRequest).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ValueEqual("shortUrl", baseURL+"/"+code)
	})

	It("should return 404 if the entry does not exist", func() {
		httpExpect.GET("/" + nonExistingCode).
			Expect().
			Status(http.StatusNotFound)
	})
})
