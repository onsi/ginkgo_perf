package writer_test

import (
	. "github.com/onsi/ginkgo_perf/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWriter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Writer Suite")
}
