package main_test

import (
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfigMapController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var (
	binaryPath string
)

var _ = BeforeSuite(func() {
	var err error
	binaryPath, err = gexec.Build("github.com/aclevername/config-map-controller")
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
