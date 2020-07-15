package exporter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kleveross/ormb/pkg/exporter"
	"github.com/kleveross/ormb/pkg/saver"
)

var _ = Describe("Exporter", func() {
	var e exporter.Interface

	BeforeEach(func() {
		e = exporter.New()
	})

	It("Should export the model successfully", func() {
		s := saver.New()
		pwd := "../../examples/PMML-model"

		m, err := s.Save(pwd)
		Expect(err).To(BeNil())
		Expect(m.Metadata.Format).To(Equal("PMML"))

		_, err = e.Export(m, "/tmp")
		Expect(err).To(BeNil())
	})
})
