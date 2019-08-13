package git

import (
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/git/executor"
)

var _ = Describe("GitClient", func() {

	var (
		fakeExecutor *executor.FakeExecutor
		gitClient    *Client
		tempCloneDir string
	)

	BeforeEach(func() {
		fakeExecutor = new(executor.FakeExecutor)
		gitClient = NewGitClientFromExecutor(fakeExecutor)
	})

	AfterEach(func() {
		deleteTempDir(tempCloneDir)
	})

	It("it can create a directory, clone the repo and delete it afterwards", func() {
		fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)
		deleteTempDir(tempCloneDir)

		var err error
		tempCloneDir, err = gitClient.CloneRepo("test-git-", "my-branch", "git@example.com:test/example-repo.git")

		// It called clone
		Expect(err).To(Not(HaveOccurred()))
		Expect(fakeExecutor.Dir).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Args).To(
			Equal([]string{"clone", "-b", "my-branch", "git@example.com:test/example-repo.git", tempCloneDir}))

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
		Expect(fakeExecutor.Args).To(
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
		Expect(fakeExecutor.Args).To(
			Equal([]string{"push"}))
	})
})

func deleteTempDir(tempDir string) {
	if tempDir != "" && strings.HasPrefix(tempDir, os.TempDir()) {
		_ = os.RemoveAll(tempDir)
	}
}
