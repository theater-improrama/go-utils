package pbkdf2_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPbkdf2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pbkdf2 Suite")
}
