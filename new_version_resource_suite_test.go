package new_version_resource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHtmlCsspathVersionResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HtmlCsspathVersionResource Suite")
}
