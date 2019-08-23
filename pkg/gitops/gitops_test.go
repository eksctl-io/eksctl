package gitops

import (
	"context"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
)

type mockCloner struct {
	mock.Mock
}

func (m *mockCloner) CloneRepo(cloneDirPrefix string, branch string, gitURL string) (string, error) {
	args := m.Called(cloneDirPrefix, branch, gitURL)
	return args.String(0), args.Error(1)
}

var _ = Describe("gitops profile", func() {

	var (
		gitCloner *mockCloner
		memFs     afero.Fs
		io        afero.Afero
		testDir   string
		profile   *Profile
		outputDir string
		processor fileprocessor.FileProcessor
	)

	Context("generating a profile", func() {

		BeforeEach(func() {
			// In memory filesystem for the tests
			memFs = afero.NewMemMapFs()
			io = afero.Afero{Fs: memFs}

			// Create test data files instead of cloning
			testDir, _ = io.TempDir("", "test-dir-")
			createTestFiles(testDir, memFs)

			// mock git clone
			gitCloner = new(mockCloner)
			gitCloner.On("CloneRepo", mock.Anything, mock.Anything, mock.Anything).Return(testDir, nil)

			// output path
			outputDir, _ = io.TempDir("", "test-output-dir-")

			processor = &fileprocessor.GoTemplateProcessor{
				Params: fileprocessor.TemplateParameters{ClusterName: "test-cluster"},
			}
			profile = &Profile{
				Path: outputDir,
				GitOpts: git.Options{
					Branch: "master",
					URL:    "git@github.com:someorg/test-gitops-repo.git",
				},
				IO:        io,
				FS:        memFs,
				GitCloner: gitCloner,
				Processor: processor,
			}
		})

		AfterEach(func() {
			io.RemoveAll(testDir)
			io.RemoveAll(outputDir)
		})

		It("processes go templates and writes them to the output directory", func() {
			err := profile.Generate(context.Background())

			Expect(err).ToNot(HaveOccurred())
			template1, err := io.ReadFile(filepath.Join(outputDir, "a/good-template1.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(template1).To(MatchYAML([]byte("cluster: test-cluster")))

			template2, err := io.ReadFile(filepath.Join(outputDir, "a/b/good-template2.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(template2).To(MatchYAML([]byte("name: test-cluster")))
		})

		It("can load files and ignore .git/ files", func() {
			files, err := profile.loadFiles(testDir)

			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(4))
			Expect(files).To(ConsistOf(
				fileprocessor.File{
					Path: filepath.Join(testDir, "a/good-template1.yaml.tmpl"),
					Data: []byte("cluster: {{ .ClusterName }}"),
				},
				fileprocessor.File{
					Path: filepath.Join(testDir, "a/b/good-template2.yaml.tmpl"),
					Data: []byte("name: {{ .ClusterName }}"),
				},
				fileprocessor.File{
					Path: filepath.Join(testDir, "a/not-a-template2.yaml"),
					Data: []byte("somekey2: value2"),
				},
				fileprocessor.File{
					Path: filepath.Join(testDir, "not-a-template.yaml"),
					Data: []byte("somekey: value"),
				}))
		})

		Context("processing templates", func() {

			It("processes go templates and leaves the rest intact", func() {
				templateContent := []byte(`
apiVersion: v1
kind: Namespace
metadata:
  labels:
    name: {{ .ClusterName }}
  name: flux`)
				expectedProcessedTemplate := []byte(`
apiVersion: v1
kind: Namespace
metadata:
  labels:
    name: test-cluster
  name: flux`)
				pureYaml := []byte("this: is just yaml")
				inputFiles := []fileprocessor.File{
					{
						Data: templateContent,
						Path: "dir0/some-file.yaml.tmpl",
					},
					{
						Data: templateContent,
						Path: "dir0/dir1/some-file2.yaml.tmpl",
					},
					{
						Data: pureYaml,
						Path: "dir0/dir1/non-template.yaml",
					},
					{
						Data: templateContent,
						Path: "dir0/dir1/dir2/some-file3.yaml.tmpl",
					},
				}

				files, err := profile.processFiles(inputFiles, "dir0")

				Expect(err).ToNot(HaveOccurred())
				Expect(files).To(HaveLen(4))
				Expect(files).To(ConsistOf(
					fileprocessor.File{
						Path: "some-file.yaml",
						Data: expectedProcessedTemplate,
					},
					fileprocessor.File{
						Path: "dir1/some-file2.yaml",
						Data: expectedProcessedTemplate,
					},
					fileprocessor.File{
						Path: "dir1/non-template.yaml",
						Data: pureYaml,
					},
					fileprocessor.File{
						Path: "dir1/dir2/some-file3.yaml",
						Data: expectedProcessedTemplate,
					},
				))
			})
		})
	})
})

func createTestFiles(testDir string, memFs afero.Fs) {
	createFile(memFs, filepath.Join(testDir, "not-a-template.yaml"), "somekey: value")
	createFile(memFs, filepath.Join(testDir, "a/not-a-template2.yaml"), "somekey2: value2")
	createFile(memFs, filepath.Join(testDir, "a/good-template1.yaml.tmpl"), "cluster: {{ .ClusterName }}")
	createFile(memFs, filepath.Join(testDir, "a/b/good-template2.yaml.tmpl"), "name: {{ .ClusterName }}")
	memFs.Mkdir(".git", 0755)
	createFile(memFs, filepath.Join(testDir, ".git/some-git-file"), "this is a git file and should be ignored")
	createFile(memFs, filepath.Join(testDir, ".git/some-git-file.yaml"), "this is a git file and should be ignored")
	createFile(memFs, filepath.Join(testDir, ".git/some-git-file.yaml.tmpl"), "this is a git file and should be ignored")
}

func createFile(memFs afero.Fs, path string, content string) error {
	file, err := memFs.Create(path)
	if err != nil {
		return err
	}
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}
