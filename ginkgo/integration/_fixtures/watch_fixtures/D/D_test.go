package D_test

import (
	. "github.com/onsi/ginkgo_perf/ginkgo/integration/_fixtures/watch_fixtures/C"

	. "github.com/onsi/ginkgo_perf/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("D", func() {
	It("should do it", func() {
		Ω(DoIt()).Should(Equal("done!"))
	})
})
