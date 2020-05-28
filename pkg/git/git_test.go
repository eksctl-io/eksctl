package git_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/git/executor"
)

var _ = Describe("git", func() {
	Describe("Client", func() {
		var (
			fakeExecutor *executor.FakeExecutor
			gitClient    *git.Client
			tempCloneDir string
		)

		BeforeEach(func() {
			fakeExecutor = new(executor.FakeExecutor)
			gitClient = git.NewGitClientFromExecutor(fakeExecutor)
		})

		AfterEach(func() {
			deleteTempDir(tempCloneDir)
		})

		It("it can create a directory, clone the repo and delete it afterwards", func() {
			fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)
			deleteTempDir(tempCloneDir)

			var err error
			options := git.CloneOptions{
				Branch: "my-branch",
				URL:    "git@example.com:test/example-repo.git",
			}
			tempCloneDir, err = gitClient.CloneRepoInTmpDir("test-git-", options)

			// It called clone and checkout on the branch
			Expect(err).To(Not(HaveOccurred()))
			Expect(len(fakeExecutor.Calls)).To(Equal(2))
			Expect(fakeExecutor.Calls[0].Arguments[2]).To(Equal([]string{"clone", "git@example.com:test/example-repo.git", tempCloneDir}))
			Expect(fakeExecutor.Calls[1].Arguments[2]).To(Equal([]string{"checkout", "my-branch"}))
			// The directory for the checkout command was the cloned directory
			Expect(fakeExecutor.Calls[1].Arguments[1]).To(Equal(tempCloneDir))

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
			fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)

			err := gitClient.Add("file1", "file2")

			Expect(err).To(Not(HaveOccurred()))
			Expect(fakeExecutor.Calls[0].Arguments[2]).To(
				Equal([]string{"add", "--", "file1", "file2"}))
		})

		It("can make commits", func() {
			fakeExecutor.On("Exec", mock.Anything, mock.Anything, mock.MatchedBy(func(args []string) bool {
				return args[0] == "diff"
			})).Return(&exec.ExitError{})
			fakeExecutor.On("Exec", mock.Anything, mock.Anything, mock.MatchedBy(func(args []string) bool {
				return args[0] == "config"
			})).Return(nil)
			fakeExecutor.On("Exec", mock.Anything, mock.Anything, mock.MatchedBy(func(args []string) bool {
				return args[0] == "config"
			})).Return(nil)
			fakeExecutor.On("Exec", mock.Anything, mock.Anything, mock.MatchedBy(func(args []string) bool {
				return args[0] == "commit"
			})).Return(nil)

			err := gitClient.Commit("test commit", "test-user", "test-user@example.com")

			Expect(err).To(Not(HaveOccurred()))
			Expect(fakeExecutor.Calls[0].Arguments[0]).To(Equal("git"))
			Expect(fakeExecutor.Calls[0].Arguments[2]).To(Equal([]string{"diff", "--cached", "--quiet"}))

			Expect(fakeExecutor.Calls[1].Arguments[0]).To(Equal("git"))
			Expect(fakeExecutor.Calls[1].Arguments[2]).To(
				Equal([]string{"config", "user.email", "test-user@example.com"}))

			Expect(fakeExecutor.Calls[2].Arguments[0]).To(Equal("git"))
			Expect(fakeExecutor.Calls[2].Arguments[2]).To(
				Equal([]string{"config", "user.name", "test-user"}))

			Expect(fakeExecutor.Calls[3].Arguments[0]).To(Equal("git"))
			Expect(fakeExecutor.Calls[3].Arguments[2]).To(
				Equal([]string{"commit", "-m", "test commit", "--author=test-user <test-user@example.com>"}))
		})

		It("can push", func() {
			fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)

			err := gitClient.Push()

			Expect(err).To(Not(HaveOccurred()))
			Expect(fakeExecutor.Calls[0].Arguments[2]).To(
				Equal([]string{"config", "push.default", "current"}))
			Expect(fakeExecutor.Calls[1].Arguments[2]).To(
				Equal([]string{"push"}))
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
			privateSSHKey, err := ioutil.TempFile("", "fake_id_rsa")
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
