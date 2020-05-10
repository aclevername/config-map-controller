package main_test

import (
	"os/exec"

	"github.com/onsi/gomega/gbytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

//go:generate counterfeiter -o fakes/fake_read_closer.go io.ReadCloser

var _ = Describe("Main", func() {

	var (
		session *gexec.Session
	)

	AfterEach(func() {
		if session != nil {
			Eventually(session.Terminate()).Should(gexec.Exit())
		}
	})

	When("no args are provided", func() {
		It("exits non-zero and gives a useful help message", func() {
			var err error
			cmd := exec.Command(binaryPath)
			stdErr := gbytes.NewBuffer()
			session, err = gexec.Start(cmd, GinkgoWriter, stdErr)
			Expect(err).NotTo(HaveOccurred())
			session.Wait()
			exitCode := session.ExitCode()
			Expect(exitCode).NotTo(Equal(0))
		})
	})

	When("no --kubeconfig value is provided", func() {
		It("exits non-zero", func() {
			var err error
			cmd := exec.Command(binaryPath, "--kubeconfig")
			stdErr := gbytes.NewBuffer()
			session, err = gexec.Start(cmd, GinkgoWriter, stdErr)
			Expect(err).NotTo(HaveOccurred())
			session.Wait()
			exitCode := session.ExitCode()
			Expect(exitCode).NotTo(Equal(0))
		})
	})

	When("an unrecognized flag is provided", func() {
		It("exits non-zero", func() {
			var err error
			cmd := exec.Command(binaryPath, "--incorrect-flag", "incorrect-value")
			stdErr := gbytes.NewBuffer()
			session, err = gexec.Start(cmd, GinkgoWriter, stdErr)
			Expect(err).NotTo(HaveOccurred())
			session.Wait()
			exitCode := session.ExitCode()
			Expect(exitCode).NotTo(Equal(0))
		})
	})

	When("the kubeconfig value isn't a real path", func() {
		It("exits non-zero", func() {
			var err error
			cmd := exec.Command(binaryPath, "--kubeconfig", "/path/to/nowhere")
			stdErr := gbytes.NewBuffer()
			session, err = gexec.Start(cmd, GinkgoWriter, stdErr)
			Expect(err).NotTo(HaveOccurred())
			session.Wait()
			exitCode := session.ExitCode()
			Expect(exitCode).NotTo(Equal(0))
			Expect(string(stdErr.Contents())).To(ContainSubstring("failed to build client config from: /path/to/nowhere"))
		})
	})
})
