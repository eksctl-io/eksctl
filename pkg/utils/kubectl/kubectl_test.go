package kubectl_test

import (
	"fmt"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/utils/kubectl"
)

var _ = Describe("Kubectl", func() {
	var client kubectl.KubernetesClient
	var genericError = "genericError"

	Context("FmtCmd", func() {
		BeforeEach(func() {
			client = kubectl.NewClient()
		})

		It("should properly format the command", func() {
			args := []string{"arg1", "arg2"}
			Expect(client.FmtCmd(args)).To(Equal("kubectl arg1 arg2"))
		})
	})

	Context("GetClientVersion", func() {

		BeforeEach(func() {
			client = kubectl.NewClient()
		})
		AfterEach(func() {
			kubectl.SetExecCommand(exec.Command)
		})

		It("should return an error if kubectl call fails", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `fail`)
			})
			_, err := client.GetClientVersion()
			Expect(err).To(MatchError(ContainSubstring("error running `kubectl version`: exit status 1")))
		})

		It("should return an error if parsing the version fails", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), invalidCommandOutput)
			})
			_, err := client.GetClientVersion()
			Expect(err).To(MatchError(ContainSubstring("error parsing `kubectl version` output")))
		})

		It("should return the version successfully", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), commandOutput)
			})
			clientVersion, err := client.GetClientVersion()
			Expect(err).NotTo(HaveOccurred())
			Expect(clientVersion).To(Equal("v1.28.1"))
		})
	})

	Context("GetServerVersion", func() {
		BeforeEach(func() {
			client = kubectl.NewClient()
		})
		AfterEach(func() {
			kubectl.SetExecCommand(exec.Command)
		})

		It("should return an error if env is not set", func() {
			_, err := client.GetServerVersion()
			Expect(err).To(MatchError(ContainSubstring("client env should be set before trying to fetch server version")))
		})

		Context("env is set appropriately", func() {
			BeforeEach(func() {
				client.SetEnv([]string{"env"})
			})

			It("should return an error if kubectl call fails", func() {
				kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
					return exec.Command(filepath.Join("testdata", "fake-version"), `fail`)
				})
				_, err := client.GetServerVersion()
				Expect(err).To(MatchError(ContainSubstring("error running `kubectl version`: exit status 1")))
			})

			It("should return an error if parsing the version fails", func() {
				kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
					return exec.Command(filepath.Join("testdata", "fake-version"), invalidCommandOutput)
				})
				_, err := client.GetServerVersion()
				Expect(err).To(MatchError(ContainSubstring("error parsing `kubectl version` output")))
			})

			It("should return the version successfully", func() {
				kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
					return exec.Command(filepath.Join("testdata", "fake-version"), commandOutput)
				})
				serverVersion, err := client.GetServerVersion()
				Expect(err).NotTo(HaveOccurred())
				Expect(serverVersion).To(Equal("v1.24.16-eks-2d98532"))
			})
		})
	})

	Context("CheckKubectlVersion", func() {
		BeforeEach(func() {
			client = kubectl.NewClient()
		})
		It("should return an error if kubectl is not present", func() {
			kubectl.SetExecLookPath(func(file string) (string, error) {
				return "", fmt.Errorf(genericError)
			})
			Expect(client.CheckKubectlVersion()).To(MatchError(ContainSubstring("kubectl not found, v1.10.0 or newer is required")))
		})

		Context("kubectl is present", func() {
			BeforeEach(func() {
				kubectl.SetExecLookPath(func(file string) (string, error) {
					return "path", nil
				})
			})

			It("should return an error if it fails to parse kubectl version", func() {
				kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
					return exec.Command(filepath.Join("testdata", "fake-version"), invalidClientVersionCommandOutput)
				})
				Expect(client.CheckKubectlVersion()).To(MatchError(ContainSubstring("parsing kubectl version string")))
			})

			It("should return an error if kubectl version is not supported", func() {
				kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
					return exec.Command(filepath.Join("testdata", "fake-version"), oldClientVersionCommandOutput)
				})
				Expect(client.CheckKubectlVersion()).To(MatchError(ContainSubstring("kubectl version 1.9.0 was found at \"path\", minimum required version to use EKS is v1.10.0")))
			})

			It("should finish successfully", func() {
				kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
					return exec.Command(filepath.Join("testdata", "fake-version"), commandOutput)
				})
				Expect(client.CheckKubectlVersion()).NotTo(HaveOccurred())
			})
		})
	})
})

var commandOutput = `{
	"clientVersion": {
	  "major": "1",
	  "minor": "28",
	  "gitVersion": "v1.28.1",
	  "gitCommit": "8dc49c4b984b897d423aab4971090e1879eb4f23",
	  "gitTreeState": "clean",
	  "buildDate": "2023-08-24T11:16:29Z",
	  "goVersion": "go1.20.7",
	  "compiler": "gc",
	  "platform": "darwin/arm64"
	},
	"kustomizeVersion": "v5.0.4-0.20230601165947-6ce0bf390ce3",
	"serverVersion": {
	  "major": "1",
	  "minor": "24+",
	  "gitVersion": "v1.24.16-eks-2d98532",
	  "gitCommit": "af930c12e26ef9d1e8fac7e3532ff4bcc1b2b509",
	  "gitTreeState": "clean",
	  "buildDate": "2023-07-28T16:52:47Z",
	  "goVersion": "go1.20.6",
	  "compiler": "gc",
	  "platform": "linux/amd64"
	}
}`

var invalidCommandOutput = `{
	-
}`

var invalidClientVersionCommandOutput = `{
	"clientVersion": {
		"gitVersion": ""
	  }
}`

var oldClientVersionCommandOutput = `{	
	"clientVersion": {
		"gitVersion": "1.9.0"
	  }
}`
