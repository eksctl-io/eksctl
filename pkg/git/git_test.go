package git

import (
	"github.com/docker/docker/pkg/ioutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/git/executor"
	"os"
	"os/exec"
	"strings"
)

var _ = Describe("GitClient", func() {

	var (
		fakeExecutor *executor.FakeExecutor
		gitClient    *Client
		tempCloneDir string
	)

	BeforeEach(func() {
		tempCloneDir, _ = ioutils.TempDir("", "git_test-")
		fakeExecutor = new(executor.FakeExecutor)
		gitClient = NewGitClientFromExecutor(tempCloneDir, "test-user", "test-user@example.com", fakeExecutor)
	})

	AfterEach(func() {
		deleteTempDir(tempCloneDir)
	})

	It("can clone repo on the given directory", func() {
		fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)

		_, err := gitClient.CloneRepo("my-branch", "git@example.com:test/example-repo.git")

		Expect(err).To(Not(HaveOccurred()))
		Expect(fakeExecutor.Dir).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Args).To(
			Equal([]string{"clone", "-b", "my-branch", "git@example.com:test/example-repo.git", tempCloneDir}))
	})

	It("creates the directory when it doesn't exist", func() {
		fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)
		deleteTempDir(tempCloneDir)

		_, err := gitClient.CloneRepo("my-branch", "git@example.com:test/example-repo.git")

		// It called clone
		Expect(err).To(Not(HaveOccurred()))
		Expect(fakeExecutor.Dir).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Args).To(
			Equal([]string{"clone", "-b", "my-branch", "git@example.com:test/example-repo.git", tempCloneDir}))

		// The directory was created
		_, err = os.Stat(tempCloneDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("can delete the clone directory", func() {
		err := gitClient.DeleteLocalRepo()

		Expect(err).ToNot(HaveOccurred())
		_, err = os.Stat(tempCloneDir)
		Expect(err).To(HaveOccurred())
		Expect(os.IsNotExist(err)).To(BeTrue())

	})

	It("can add files", func() {
		fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)

		err := gitClient.Add("file1", "file2")

		Expect(err).To(Not(HaveOccurred()))
		Expect(fakeExecutor.Dir).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Args).To(
			Equal([]string{"add", "--", "file1", "file2"}))
	})

	It("can make commits", func() {
		fakeExecutor.On("Exec", mock.Anything, mock.Anything, mock.MatchedBy(func(args []string) bool {
			return args[0] == "diff"
		})).Return(&exec.ExitError{})
		fakeExecutor.On("Exec", mock.Anything, mock.Anything, mock.MatchedBy(func(args []string) bool {
			return args[0] == "commit"
		})).Return(nil)

		err := gitClient.Commit("test commit")

		Expect(err).To(Not(HaveOccurred()))
		Expect(fakeExecutor.Calls[0].Arguments[0]).To(Equal("git"))
		Expect(fakeExecutor.Calls[0].Arguments[1]).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Calls[0].Arguments[2]).To(Equal([]string{"diff", "--cached", "--quiet"}))

		Expect(fakeExecutor.Calls[1].Arguments[0]).To(Equal("git"))
		Expect(fakeExecutor.Calls[1].Arguments[1]).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Calls[1].Arguments[2]).To(
			Equal([]string{"commit", "-m", "test commit", "--author=test-user <test-user@example.com>"}))
	})

	It("can push", func() {
		fakeExecutor.On("Exec", "git", mock.Anything, mock.Anything).Return(nil)

		err := gitClient.Push()

		Expect(err).To(Not(HaveOccurred()))
		Expect(fakeExecutor.Dir).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Args).To(
			Equal([]string{"push"}))
	})
})

func deleteTempDir(tempDir string) {
	if strings.HasPrefix(tempDir, os.TempDir()) && len(tempDir) > len(os.TempDir())+1 {
		_ = os.RemoveAll(tempDir)
	}
}
