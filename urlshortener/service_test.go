package urlshortener

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/woraphol-j/url-shortener/mocks"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Service Suite")
}

var _ = Describe("Service", func() {
	const (
		code        = "abcde"
		originalURL = "www.goggle.com"
	)

	var (
		mockCtrl *gomock.Controller
		mockDAO  *mocks.MockDAO
		service  Service
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDAO = mocks.NewMockDAO(mockCtrl)
		service = NewService(mockDAO)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should save new short url entry correctly", func() {
		mockDAO.EXPECT().Save("").Return(nil).Times(1)
		_, err := service.GenerateShortURL(originalURL)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should retrieve existing short url entry correctly", func() {
		mockDAO.EXPECT().Get(code).Return(originalURL).Times(1)
		_, err := service.GetOriginalURL(code)
		Expect(err).ShouldNot(HaveOccurred())
	})
})
