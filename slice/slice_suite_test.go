package slice_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSlice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Slice Suite")
}
