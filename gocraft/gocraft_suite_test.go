package gocraft_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGocraft(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gocraft Suite")
}
