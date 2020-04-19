package subpackage

import (
	. "github.com/onsi/ginkgo_perf/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("nested sub packages", func() {
		GinkgoT().Fail(true)
	})
})
