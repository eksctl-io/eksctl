package gitops

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"strings"
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

			profile = &Profile{
				Path: outputDir,
				GitOpts: GitOptions{
					Branch: "master",
					URL:    "git@github.com:someorg/test-gitops-repo.git",
				},
				IO:        io,
				Fs:        memFs,
				gitCloner: gitCloner,
				Params:    TemplateParameters{ClusterName: "test-cluster"},
			}
		})

		AfterEach(func() {
			io.RemoveAll(testDir)
			io.RemoveAll(outputDir)
		})

		It("process the templates and writes them to the output directory", func() {
			err := profile.Generate(context.Background())

			Expect(err).ToNot(HaveOccurred())
			template1, err := io.ReadFile(filepath.Join(outputDir, "/a/good-template1.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(template1).To(Equal([]byte("cluster: test-cluster")))

			template2, err := io.ReadFile(filepath.Join(outputDir, "a/b/good-template2.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(template2).To(Equal([]byte("name: test-cluster")))
		})

		It("loads only .templ files", func() {
			files, err := profile.loadFiles(testDir)

			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(2))
			Expect(files).To(ConsistOf(
				ManifestFile{
					Name: filepath.Join(testDir, "/a/good-template1.yaml.templ"),
					Data: []byte("cluster: {{ .ClusterName }}"),
				},
				ManifestFile{
					Name: filepath.Join(testDir, "a/b/good-template2.yaml.templ"),
					Data: []byte("name: {{ .ClusterName }}"),
				}))
		})

		Context("processing templates", func() {

			It("loads only .templ files", func() {
				params := TemplateParameters{
					ClusterName: "fancy-unicorn",
				}
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
    name: fancy-unicorn
  name: flux`)
				templates := []ManifestFile{
					{
						Data: templateContent,
						Name: "dir0/some-file.yaml.templ",
					},
					{
						Data: templateContent,
						Name: "dir0/dir1/some-file2.yaml.templ",
					},
					{
						Data: templateContent,
						Name: "dir0/dir1/dir2/some-file3.yaml.templ",
					},
				}

				files, err := processFiles(templates, params, "dir0")

				Expect(err).ToNot(HaveOccurred())
				Expect(files).To(HaveLen(3))
				Expect(files).To(ConsistOf(
					ManifestFile{
						Name: "some-file.yaml",
						Data: expectedProcessedTemplate,
					},
					ManifestFile{
						Name: "dir1/some-file2.yaml",
						Data: expectedProcessedTemplate,
					},
					ManifestFile{
						Name: "dir1/dir2/some-file3.yaml",
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
	createFile(memFs, filepath.Join(testDir, "a/good-template1.yaml.templ"), "cluster: {{ .ClusterName }}")
	createFile(memFs, filepath.Join(testDir, "a/b/good-template2.yaml.templ"), "name: {{ .ClusterName }}")
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

func deleteTempDir(tempDir string) {
	if tempDir != "" && strings.HasPrefix(tempDir, os.TempDir()) {
		_ = os.RemoveAll(tempDir)
	}
}
