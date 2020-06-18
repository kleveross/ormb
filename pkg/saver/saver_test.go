package saver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/caicloud/ormb/pkg/saver"
)

var _ = Describe("Saver", func() {
	var s saver.Interface

	BeforeEach(func() {
		s = saver.New()
	})

	It("Should save the model successfully", func() {
		pwd := "../../examples/PMML-model"
		m, err := s.Save(pwd)
		Expect(err).To(BeNil())
		Expect(m.Metadata.Format).To(Equal("PMML"))
	})
})
