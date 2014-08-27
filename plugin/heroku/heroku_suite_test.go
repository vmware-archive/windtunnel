package heroku_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHeroku(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Heroku Suite")
}
