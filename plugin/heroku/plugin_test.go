package heroku_test

import (
	. "github.com/cf-platform-eng/windtunnel/plugin/heroku"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plugin", func() {
	It("Should Return My Heroku Oauth Token", func() {
		heroku := new(Plugin)

		token := heroku.Authenticate()
		Ω(token).ShouldNot(BeZero())
	})

	It("Should Return App Status", func() {
		heroku := new(Plugin)
		token := heroku.Authenticate()

		Ω(heroku.Status(token, "lit-wave-5074")).Should(Equal([]int{1, 1}))
	})
})
