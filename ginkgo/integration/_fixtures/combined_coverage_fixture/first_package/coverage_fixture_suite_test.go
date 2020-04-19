package first_package_test

import (
	. "github.com/onsi/ginkgo_perf/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCoverageFixture(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CombinedFixture First Suite")
}
