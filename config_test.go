package thruster_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/thruster"
)

var _ = Describe("Config", func() {

	Describe("NewConfig", func() {
		It("maps the contents of 'path' to a Config correctly", func() {
			config, err := thruster.NewConfig("fixtures/sample.config.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(config.Hostname).To(Equal("localhost"))
			Expect(config.Port).To(Equal(8888))
			Expect(config.TLS).To(Equal(true))
			Expect(config.Certificate).To(Equal("/etc/certificate1"))
			Expect(config.PublicKey).To(Equal("/etc/public_key"))
			Expect(config.HTTPAuth).To(Equal([]thruster.HTTPAuth{
				thruster.HTTPAuth{Username: "admin", Password: "12345"},
				thruster.HTTPAuth{Username: "user1", Password: "6666"},
			}))
		})
	})

})
