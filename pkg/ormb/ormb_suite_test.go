package ormb

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestORMB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ormb Suite")
}
