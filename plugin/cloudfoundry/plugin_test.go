package cloudfoundry_test

import (
	. "github.com/cf-platform-eng/windtunnel/plugin/cloudfoundry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plugin", func() {
	It("Should Return My CF Bearer Token", func() {
		cf := new(Plugin)

		token := cf.Authenticate()
		Ω(token).ShouldNot(BeZero())
	})

	It("Should Return App Status", func() {
		cf := new(Plugin)
		token := cf.Authenticate()

		Ω(cf.Status(token, "bfed49e3-376e-47ad-942a-12ec1d6f9fe9/")).Should(Equal([]int{1,1}))
	})
})
