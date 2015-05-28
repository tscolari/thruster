package thruster_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestThruster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Thuster Suite")
}
