package acceptance_test

import (
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var (
	kubeconfigPath string
	session        *gexec.Session
)

var _ = BeforeSuite(func() {
	kubeconfigPath = mustGetEnv("KUBECONFIG")

	binaryPath, err := gexec.Build("github.com/aclevername/config-map-controller")
	Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command(binaryPath, "--kubeconfig", kubeconfigPath)
	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {
	Eventually(session.Terminate()).Should(gexec.Exit())
	gexec.CleanupBuildArtifacts()
})

func mustGetEnv(keyName string) string {
	val := os.Getenv(keyName)
	if val == "" {
		Fail("Need " + keyName + " for the test")
	}
	return val
}
