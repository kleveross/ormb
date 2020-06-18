package cache

import (
	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	Describe("with real use cases", func() {
		var c *Cache
		var err error
		var rootPath string

		BeforeEach(func() {
			var i Interface
			rootPath = ".cache"

			i, err = New(CacheOptRoot(rootPath))
			c = i.(*Cache)
		})

		It("Should create the cache successfully", func() {
			Expect(err).To(BeNil())
			Expect(c.rootDir).To(Equal(rootPath))
		})

		It("Should store the reference successfully", func() {
			m := &model.Model{
				Metadata: &model.Metadata{
					Format: "SavedModel",
				},
			}
			refStr := "caicloud/test:v1"
			ref, err := oci.ParseReference(refStr)
			Expect(err).To(BeNil())

			actual, err := c.StoreReference(ref, m)
			Expect(err).To(BeNil())
			Expect(actual.Model).To(Equal(m))
			Expect(actual.Name).To(Equal(refStr))
		})
	})
})
