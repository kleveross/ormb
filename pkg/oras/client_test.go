package oras

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cachemock "github.com/caicloud/ormb/pkg/oras/cache/mock"
	orasmock "github.com/caicloud/ormb/pkg/oras/mock"
	orasclientmock "github.com/caicloud/ormb/pkg/oras/orasclient/mock"
)

var _ = Describe("OCI Client", func() {
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
	})
})
