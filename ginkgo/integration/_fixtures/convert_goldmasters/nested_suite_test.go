package nested_test

import (
	"testing"

	. "github.com/onsi/ginkgo_perf/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNested(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nested Suite")
}
