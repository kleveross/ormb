package oras

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
	"github.com/caicloud/ormb/pkg/oras/cache"
	cachemock "github.com/caicloud/ormb/pkg/oras/cache/mock"
	orasmock "github.com/caicloud/ormb/pkg/oras/mock"
	orasclientmock "github.com/caicloud/ormb/pkg/oras/orasclient/mock"
)

var _ = Describe("OCI Client", func() {
	Describe("with real use cases", func() {
		var c Interface
		var err error
		var rootPath string

		BeforeEach(func() {
			rootPath = "~/.ormb"
			c, err = NewClient(ClientOptRootPath(rootPath))
		})

		It("Should create the client successfully", func() {
			Expect(err).To(BeNil())
			Expect(c.(*Client).rootPath).To(Equal(rootPath))
		})
	})

	Describe("with all things mocked", func() {
		var c *Client

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			ctrl1 := gomock.NewController(GinkgoT())
			ctrl2 := gomock.NewController(GinkgoT())
			ctrl3 := gomock.NewController(GinkgoT())
			c = &Client{
				out:        os.Stdout,
				authorizer: &Authorizer{Client: orasmock.NewMockClient(ctrl)},
				resolver:   &Resolver{Resolver: orasmock.NewMockResolver(ctrl1)},
				cache:      cachemock.NewMockInterface(ctrl2),
				orasClient: orasclientmock.NewMockInterface(ctrl3),
			}
		})

		It("Should login successfully", func() {
			host := "test.harbor.com"
			user := "user"
			pwd := "pwd"
			insec := true
			c.authorizer.Client.(*orasmock.MockClient).EXPECT().Login(
				gomock.Any(),
				gomock.Eq(host),
				gomock.Eq(user),
				gomock.Eq(pwd),
				gomock.Eq(insec),
			).Return(nil).Times(1)
			Expect(c.Login(host, user, pwd, insec)).To(BeNil())
		})

		It("Should logout successfully", func() {
			host := "test.harbor.com"
			c.authorizer.Client.(*orasmock.MockClient).EXPECT().Logout(
				gomock.Any(),
				gomock.Eq(host),
			).Return(nil).Times(1)
			Expect(c.Logout(host)).To(BeNil())
		})

		It("Should save the model successfully", func() {
			refStr := "caicloud/resnet50:v1"
			ref, err := oci.ParseReference(refStr)
			Expect(err).To(BeNil())

			ch := &model.Model{
				Path: "/test",
			}
			returnedSummary := &cache.CacheRefSummary{
				Manifest: &ocispec.Descriptor{
					Digest: digest.Digest("sha256:123456"),
					Size:   int64(1),
				},
				Digest: digest.Digest("sha256:154saf"),
				Size:   int64(1),
				Name:   "test",
			}

			c.cache.(*cachemock.MockInterface).EXPECT().StoreReference(
				gomock.Eq(ref),
				gomock.Eq(ch),
			).Return(returnedSummary, nil).Times(1)
			c.cache.(*cachemock.MockInterface).EXPECT().AddManifest(
				gomock.Eq(ref),
				gomock.Eq(returnedSummary.Manifest),
			).Return(nil).Times(1)
			Expect(c.SaveModel(ch, ref)).To(BeNil())
		})

		It("Should push the model successfully", func() {
			refStr := "caicloud/resnet50:v1"
			ref, err := oci.ParseReference(refStr)
			Expect(err).To(BeNil())

			contentLayer := &ocispec.Descriptor{
				Digest: digest.Digest("sha256:gvdfgd"),
				Size:   int64(1),
			}
			layers := []ocispec.Descriptor{*contentLayer}
			returnedSummary := &cache.CacheRefSummary{
				Manifest: &ocispec.Descriptor{
					Digest: digest.Digest("sha256:123456"),
					Size:   int64(1),
				},
				Config: &ocispec.Descriptor{
					Digest: digest.Digest("sha256:kfgdv"),
					Size:   int64(1),
				},
				ContentLayer: contentLayer,
				Digest:       digest.Digest("sha256:123456"),
				Size:         int64(1),
				Name:         refStr,
				Exists:       true,
			}

			c.cache.(*cachemock.MockInterface).EXPECT().FetchReference(
				gomock.Eq(ref),
			).Return(returnedSummary, nil).Times(1)
			c.cache.(*cachemock.MockInterface).EXPECT().Provider().Return(nil).Times(1)
			c.orasClient.(*orasclientmock.MockInterface).EXPECT().Push(
				gomock.Any(),
				gomock.Eq(c.resolver),
				gomock.Eq(refStr),
				gomock.Nil(),
				gomock.Eq(layers),
				gomock.Any(),
				gomock.Any(),
			).Return(ocispec.Descriptor{}, nil).Times(1)
			Expect(c.PushModel(ref)).To(BeNil())
		})
	})
})
