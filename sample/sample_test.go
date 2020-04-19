package sample_test

import (
	"fmt"

	. "github.com/onsi/ginkgo_perf/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo_perf/sample"
)

var _ = Describe("Sample", func() {
	for i := 0; i < 20; i += 1 {
		It(fmt.Sprintf("is true #%d", i), func() {
			Ω(Sample()).Should(BeTrue())
		})
	}
})
