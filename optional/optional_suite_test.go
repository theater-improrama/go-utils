package optional_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOptional(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Optional Suite")
}
