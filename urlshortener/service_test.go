package urlshortener

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	cg "github.com/woraphol-j/url-shortener/pkg/codegenerator"
	"github.com/woraphol-j/url-shortener/pkg/mongo"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	testResult := os.Getenv("TEST_RESULTS")
	if testResult != "" {
		junitReporter := reporters.NewJUnitReporter(testResult)
		RunSpecsWithDefaultAndCustomReporters(
			t,
			"Test Service Suite",
			[]Reporter{junitReporter},
		)
	} else {
		RunSpecs(t, "Test Service Suite")
	}
}

var _ = Describe("Service", func() {
	const (
		code        = "abcde"
		originalURL = "http://www.google.com"
	)

	var (
		mockCtrl          *gomock.Controller
		mockDAO           *mongo.MockDAO
		mockCodeGenerator *cg.MockCodeGenerator
		service           Service
		data              *mongo.ShortURL = &mongo.ShortURL{
			Code: code,
			URL:  originalURL,
		}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDAO = mongo.NewMockDAO(mockCtrl)
		mockCodeGenerator = cg.NewMockCodeGenerator(mockCtrl)
		service = NewService(mockDAO, mockCodeGenerator)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should save new short url entry correctly", func() {
		mockCodeGenerator.EXPECT().Generate().Return(code, nil).Times(1)
		mockDAO.EXPECT().Save(data).Return(nil).Times(1)
		shortURL, err := service.GenerateShortURL(originalURL)
		Expect(shortURL).To(Equal(code))
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should retrieve existing url entry correctly", func() {
		mockDAO.EXPECT().Get(code).Return(data, nil).Times(1)
		originalURL, err := service.GetOriginalURL(code)
		Expect(originalURL).To(Equal(originalURL))
		Expect(err).ShouldNot(HaveOccurred())
	})
})
