// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
package git_test

import (
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/executor/fakes"
	"github.com/weaveworks/eksctl/pkg/git"
)

var _ = Describe("git", func() {
	Describe("Client", func() {
		var (
			fakeExecutor *fakes.FakeExecutor
			gitClient    *git.Client
			tempCloneDir string
		)

		BeforeEach(func() {
			fakeExecutor = new(fakes.FakeExecutor)
			gitClient = git.NewGitClientFromExecutor(fakeExecutor)
		})

		AfterEach(func() {
			deleteTempDir(tempCloneDir)
		})

		It("it can create a directory, clone the repo and delete it afterwards", func() {
			deleteTempDir(tempCloneDir)

			var err error
			options := git.CloneOptions{
				Branch: "my-branch",
				URL:    "git@example.com:test/example-repo.git",
			}
			tempCloneDir, err = gitClient.CloneRepoInTmpDir("test-git-", options)

			// It called clone and checkout on the branch
			Expect(err).To(Not(HaveOccurred()))
			Expect(fakeExecutor.ExecInDirCallCount()).To(Equal(2))

			_, _, receivedArgs := fakeExecutor.ExecInDirArgsForCall(0)
			Expect(receivedArgs).To(Equal([]string{"clone", "git@example.com:test/example-repo.git", tempCloneDir}))

			_, receivedDir, receivedArgs := fakeExecutor.ExecInDirArgsForCall(1)
			Expect(receivedArgs).To(Equal([]string{"checkout", "my-branch"}))
			Expect(receivedDir).To(Equal(tempCloneDir))

			// The directory was created
			_, err = os.Stat(tempCloneDir)
			Expect(err).ToNot(HaveOccurred())

			// It can delete it
			err = gitClient.DeleteLocalRepo()

			Expect(err).ToNot(HaveOccurred())
			_, err = os.Stat(tempCloneDir)
			Expect(err).To(HaveOccurred())
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		It("can add files", func() {
			err := gitClient.Add("file1", "file2")
			Expect(err).To(Not(HaveOccurred()))

			_, _, receivedArgs := fakeExecutor.ExecInDirArgsForCall(0)
			Expect(receivedArgs).To(Equal([]string{"add", "--", "file1", "file2"}))
		})

		It("can make commits", func() {
			fakeExecutor.ExecInDirCalls(func(arg1, arg2 string, array ...string) error {
				if array[0] == "diff" {
					return &exec.ExitError{}
				}
				return nil
			})

			err := gitClient.Commit("test commit", "test-user", "test-user@example.com")
			Expect(err).To(Not(HaveOccurred()))

			receivedCommand, _, receivedArgs := fakeExecutor.ExecInDirArgsForCall(0)
			Expect(receivedCommand).To(Equal("git"))
			Expect(receivedArgs).To(Equal([]string{"diff", "--cached", "--quiet"}))

			receivedCommand, _, receivedArgs = fakeExecutor.ExecInDirArgsForCall(1)
			Expect(receivedCommand).To(Equal("git"))
			Expect(receivedArgs).To(Equal([]string{"config", "user.email", "test-user@example.com"}))

			receivedCommand, _, receivedArgs = fakeExecutor.ExecInDirArgsForCall(2)
			Expect(receivedCommand).To(Equal("git"))
			Expect(receivedArgs).To(Equal([]string{"config", "user.name", "test-user"}))

			receivedCommand, _, receivedArgs = fakeExecutor.ExecInDirArgsForCall(3)
			Expect(receivedCommand).To(Equal("git"))
			Expect(receivedArgs).To(Equal([]string{"commit", "-m", "test commit", "--author=test-user <test-user@example.com>"}))
		})

		It("can push", func() {
			err := gitClient.Push()
			Expect(err).To(Not(HaveOccurred()))

			_, _, receivedArgs := fakeExecutor.ExecInDirArgsForCall(0)
			Expect(receivedArgs).To(Equal([]string{"config", "push.default", "current"}))

			_, _, receivedArgs = fakeExecutor.ExecInDirArgsForCall(1)
			Expect(receivedArgs).To(Equal([]string{"push"}))
		})
	})

	Describe("RepoName", func() {
		It("can parse the repository name from a URL", func() {
			name, err := git.RepoName("git@github.com:weaveworks/eksctl.git")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("eksctl"))

			name, err = git.RepoName("git@github.com:weaveworks/sock-shop.git")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("sock-shop"))

			name, err = git.RepoName("https://example.com/department1/team1/some-repo-name.git")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("some-repo-name"))

			name, err = git.RepoName("https://github.com/department1/team2/another-repo-name")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("another-repo-name"))
		})
	})

	Describe("IsGitURL", func() {
		It("can determine if a string is a git URL", func() {
			Expect(git.IsGitURL("git@github.com:weaveworks/eksctl.git")).To(BeTrue())
			Expect(git.IsGitURL("https://github.com/weaveworks/eksctl.git")).To(BeTrue())
			Expect(git.IsGitURL("https://username@secr3t:my-repo.example.com:8080/weaveworks/eksctl.git")).To(BeTrue())

			Expect(git.IsGitURL("git@github")).To(BeFalse())
			Expect(git.IsGitURL("https://")).To(BeFalse())
			Expect(git.IsGitURL("app-dev")).To(BeFalse())
		})
	})

	Describe("ValidateURL", func() {
		It("returns an error on empty Git URL", func() {
			Expect(git.ValidateURL("")).To(MatchError("empty Git URL"))
		})

		It("returns an error on invalid Git URL", func() {
			Expect(git.ValidateURL("https://")).To(MatchError("invalid Git URL"))
		})

		It("returns an error on HTTPS Git URL", func() {
			Expect(git.ValidateURL("https://github.com/eksctl-bot/my-gitops-repo.git")).
				To(MatchError("got a HTTP(S) Git URL, but eksctl currently only supports SSH Git URLs"))
		})

		It("succeeds when a SSH Git URL is provided", func() {
			Expect(git.ValidateURL("git@github.com:eksctl-bot/my-gitops-repo.git")).NotTo(HaveOccurred())
		})
	})

	Describe("Validate", func() {
		It("returns an error on invalid Git URL", func() {
			Expect(git.ValidateURL("https://")).To(MatchError("invalid Git URL"))
		})

		It("returns an error on HTTPS Git URL", func() {
			Expect(git.ValidateURL("https://github.com/eksctl-bot/my-gitops-repo.git")).
				To(MatchError("got a HTTP(S) Git URL, but eksctl currently only supports SSH Git URLs"))
		})

		It("succeeds when a SSH Git URL and an email address is provided", func() {
			Expect(git.ValidateURL("git@github.com:eksctl-bot/my-gitops-repo.git")).NotTo(HaveOccurred())
		})

		It("returns an error on non-existing file path", func() {
			Expect(git.ValidatePrivateSSHKeyPath("/path/to/non/existing/file/id_rsa")).
				To(MatchError("invalid path to private SSH key: /path/to/non/existing/file/id_rsa"))
		})

		It("succeeds when a valid path is provided", func() {
			privateSSHKey, err := os.CreateTemp("", "fake_id_rsa")
			Expect(err).To(Not(HaveOccurred()))
			defer os.Remove(privateSSHKey.Name()) // clean up

			Expect(git.ValidatePrivateSSHKeyPath(privateSSHKey.Name())).
				To(Not(HaveOccurred()))
		})
	})
})

func deleteTempDir(tempDir string) {
	if tempDir != "" && strings.HasPrefix(tempDir, os.TempDir()) {
		_ = os.RemoveAll(tempDir)
	}
}
