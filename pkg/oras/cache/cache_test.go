package cache

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kleveross/ormb/pkg/model"
	"github.com/kleveross/ormb/pkg/oci"
)

var _ = Describe("Cache", func() {
	Describe("with real use cases", func() {
		var c *Cache
		var err error
		var rootPath string

		BeforeEach(func() {
			var i Interface
			rootPath = ".cache"

			if err := os.RemoveAll(rootPath); err != nil {
				Expect(err).To(BeNil())
			}
			i, err = New(CacheOptRoot(rootPath), CacheOptDebug(true), CacheOptWriter(os.Stdout))
			Expect(err).To(BeNil())
			c = i.(*Cache)
		})

		It("Should create the cache successfully", func() {
			Expect(c.rootDir).To(Equal(rootPath))
		})

		Describe("with a cached artifact caicloud/test:v1", func() {
			var m *model.Model
			var ref *oci.Reference

			BeforeEach(func() {
				m = &model.Model{
					Metadata: &model.Metadata{
						Format: "SavedModel",
					},
					Content: []byte("test12345"),
				}
				refStr := "caicloud/test:v1"
				ref, err = oci.ParseReference(refStr)
				Expect(err).To(BeNil())

				actual, err := c.StoreReference(ref, m)
				Expect(err).To(BeNil())
				Expect(actual.Model).To(Equal(m))
				Expect(actual.Name).To(Equal(refStr))

				Expect(c.AddManifest(ref, actual.Manifest)).To(BeNil())
			})

			It("Should list the cached artifact successfully", func() {
				actual, err := c.ListReferences()
				Expect(err).To(BeNil())
				Expect(len(actual)).To(Equal(1))
			})

			It("Should delete the reference successfully", func() {
				actual, err := c.DeleteReference(ref)
				Expect(err).To(BeNil())
				Expect(actual.Name).To(Equal(ref.FullName()))
			})

			It("Should tag the reference successfully", func() {
				target, err := oci.ParseReference("test:1")
				Expect(err).To(BeNil())
				Expect(c.TagReference(ref, target)).To(BeNil())
			})
		})

	})
})
