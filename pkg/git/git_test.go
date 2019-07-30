package git

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/git/executor"
	"os"
)

var _ = Describe("GitClient", func() {

	var (
		fakeExecutor *executor.FakeExecutor
		gitClient    *Client
		tempCloneDir string
	)

	BeforeEach(func() {
		tempCloneDir = os.TempDir()
		fakeExecutor = new(executor.FakeExecutor)
		gitClient = NewGitClientFromExecutor(tempCloneDir, "test-user", "test-user@example.com", fakeExecutor)

	})

	It("can clone repo", func() {
		fakeExecutor.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		_, err := gitClient.CloneRepo("my-branch", "git@example.com:test/example-repo.git")

		Expect(err).To(Not(HaveOccurred()))
		Expect(fakeExecutor.Command).To(Equal("git"))
		Expect(fakeExecutor.Dir).To(Equal(tempCloneDir))
		Expect(fakeExecutor.Args).To(
			Equal([]string{"clone", "-b", "my-branch", "git@example.com:test/example-repo.git", tempCloneDir}))
	})

	//It("can add files", func() {
	//	Expect(nil).To(BeNil())
	//})
	//It("can make commits", func() {
	//	Expect(nil).To(BeNil())
	//})
	//It("can push", func() {
	//	Expect(nil).To(BeNil())
	//})
	//It("can delete a cloned repo", func() {
	//	Expect(nil).To(BeNil())
	//})
})
