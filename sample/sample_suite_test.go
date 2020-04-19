package sample_test

import (
	"testing"

	. "github.com/onsi/ginkgo_perf/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sample Suite")
}
