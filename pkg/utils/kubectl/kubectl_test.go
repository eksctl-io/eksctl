package kubectl_test

import (
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/utils/kubectl"
)

var _ = Describe("Kubectl", func() {
	var manager kubectl.KubernetesVersionManager
	// var genericError = "genericError"

	Context("GetClientVersion", func() {

		BeforeEach(func() {
			manager = kubectl.NewVersionManager()
		})
		AfterEach(func() {
			kubectl.SetExecCommand(exec.Command)
		})

		It("should return an error if kubectl call fails", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `fail`)
			})
			_, err := manager.ClientVersion()
			Expect(err).To(MatchError(ContainSubstring("error running `kubectl version`: exit status 1")))
		})

		It("should return an error if parsing the version fails", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), invalidCommandOutput)
			})
			_, err := manager.ClientVersion()
			Expect(err).To(MatchError(ContainSubstring("error parsing `kubectl version` output")))
		})

		It("should return the version successfully", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), commandOutput)
			})
			clientVersion, err := manager.ClientVersion()
			Expect(err).NotTo(HaveOccurred())
			Expect(clientVersion).To(Equal("v1.28.1"))
		})
	})

	Context("ServerVersion", func() {
		BeforeEach(func() {
			manager = kubectl.NewVersionManager()
		})
		AfterEach(func() {
			kubectl.SetExecCommand(exec.Command)
		})

		It("should return an error if kubectl call fails", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `fail`)
			})
			_, err := manager.ServerVersion([]string{}, []string{})
			Expect(err).To(MatchError(ContainSubstring("error running `kubectl version`: exit status 1")))
		})

		It("should return an error if parsing the version fails", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), invalidCommandOutput)
			})
			_, err := manager.ServerVersion([]string{}, []string{})
			Expect(err).To(MatchError(ContainSubstring("error parsing `kubectl version` output")))
		})

		It("should return the version successfully", func() {
			kubectl.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), commandOutput)
			})
			serverVersion, err := manager.ServerVersion([]string{}, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(serverVersion).To(Equal("v1.24.16-eks-2d98532"))
		})
	})

	Context("ValidateVersion", func() {
		BeforeEach(func() {
			manager = kubectl.NewVersionManager()
		})

		It("should return an error if it fails to parse kubectl version", func() {
			err := manager.ValidateVersion("", kubectl.Client)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("parsing kubernetes client version string")))
		})

		It("should return an error if kubectl version is not supported", func() {
			err := manager.ValidateVersion("v1.9.4", kubectl.Client)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("kubernetes client version v1.9.4 was found, minimum required version is v1.10.0")))
		})

		It("should finish successfully", func() {
			err := manager.ValidateVersion("v1.28.0", kubectl.Client)
			Expect(err).NotTo(HaveOccurred())
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
