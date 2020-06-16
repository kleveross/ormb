package ormb

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	exportermock "github.com/caicloud/ormb/pkg/exporter/mock"
	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
	ocimock "github.com/caicloud/ormb/pkg/oci/mock"
	savermock "github.com/caicloud/ormb/pkg/saver/mock"
)

var _ = Describe("ormb golang library", func() {
	var ociORMB *ORMB

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		ctrl1 := gomock.NewController(GinkgoT())
		ctrl2 := gomock.NewController(GinkgoT())
		ociORMB = &ORMB{
			client:   ocimock.NewMockInterface(ctrl),
			saver:    savermock.NewMockInterface(ctrl1),
			exporter: exportermock.NewMockInterface(ctrl2),
		}

	})

	It("Should login successfully", func() {
		host := "test.harbor.com"
		user := "user"
		pwd := "pwd"
		insec := true
		ociORMB.client.(*ocimock.MockInterface).EXPECT().Login(
			gomock.Eq(host),
			gomock.Eq(user),
			gomock.Eq(pwd),
			gomock.Eq(insec)).Return(nil).Times(1)
		Expect(ociORMB.Login(host, user, pwd, insec)).To(BeNil())
	})

	It("Should push the model successfully", func() {
		refStr := "caicloud/resnet50:v1"
		ref, err := oci.ParseReference(refStr)
		Expect(err).To(BeNil())

		ociORMB.client.(*ocimock.MockInterface).EXPECT().PushModel(
			gomock.Eq(ref),
		).Return(nil).Times(1)
		Expect(ociORMB.Push(refStr)).To(BeNil())
	})

	It("Should pull the model successfully", func() {
		refStr := "caicloud/resnet50:v1"
		ref, err := oci.ParseReference(refStr)
		Expect(err).To(BeNil())

		ociORMB.client.(*ocimock.MockInterface).EXPECT().PullModel(
			gomock.Eq(ref),
		).Return(nil).Times(1)
		Expect(ociORMB.Pull(refStr)).To(BeNil())
	})

	It("Should save the model successfully", func() {
		src := "/test"
		refStr := "caicloud/resnet50:v1"
		ch := &model.Model{
			Path: src,
		}
		ref, err := oci.ParseReference(refStr)
		Expect(err).To(BeNil())

		ociORMB.saver.(*savermock.MockInterface).EXPECT().Save(
			gomock.Eq(src),
		).Return(ch, nil).Times(1)
		ociORMB.client.(*ocimock.MockInterface).EXPECT().SaveModel(
			gomock.Eq(ch),
			gomock.Eq(ref),
		).Return(nil).Times(1)
		Expect(ociORMB.Save(src, refStr)).To(BeNil())
	})

	It("Should export the model successfully", func() {
		dst := "/test"
		refStr := "caicloud/resnet50:v1"
		ch := &model.Model{
			Path: dst,
		}
		ref, err := oci.ParseReference(refStr)
		Expect(err).To(BeNil())

		ociORMB.exporter.(*exportermock.MockInterface).EXPECT().Export(
			gomock.Eq(ch),
			gomock.Eq(dst),
		).Return("", nil).Times(1)
		ociORMB.client.(*ocimock.MockInterface).EXPECT().LoadModel(
			gomock.Eq(ref),
		).Return(ch, nil).Times(1)
		Expect(ociORMB.Export(refStr, dst)).To(BeNil())
	})
})
