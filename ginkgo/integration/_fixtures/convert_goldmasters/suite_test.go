package tmp_test

import (
	"testing"

	. "github.com/onsi/ginkgo_perf/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConvertFixtures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ConvertFixtures Suite")
}
