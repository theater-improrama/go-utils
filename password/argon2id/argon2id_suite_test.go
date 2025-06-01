package argon2id_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestArgon2id(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Argon2id Suite")
}
