package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFormatGoImports(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FormatGoImports Suite")
}
