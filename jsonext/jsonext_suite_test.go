package jsonext_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJsonext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jsonext Suite")
}
