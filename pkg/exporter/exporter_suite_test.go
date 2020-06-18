package exporter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExporter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Exporter Suite")
}
