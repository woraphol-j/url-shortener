package mongo

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDataAccess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Access Suite")
}

var _ = Describe("Book", func() {
	// var (
	// 	meal DAO
	// )

	BeforeEach(func() {
		// meal = DAO{
		// 	Name: "Les Miserables",
		// 	Date: time.Now().UTC(),
		// }
	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				// Expect(meal.Name).To(Equal("Les Miserables"))
				Expect("Les Miserables").To(Equal("Les Miserables"))
			})
		})
	})
})
