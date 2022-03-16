package weibo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWeibo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Weibo Suite")
}
