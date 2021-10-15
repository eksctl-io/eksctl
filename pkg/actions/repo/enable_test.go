// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
package repo_test

import (
	"os"

	"github.com/instrumenta/kubeval/kubeval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/repo"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Installer", func() {
	var (
		installer *repo.Installer
		tmpDir    string
		mockOpts  *api.Git
	)

	BeforeEach(func() {
		mockOpts = &api.Git{
			Repo: &api.Repo{
				URL:      "git@github.com/foo/bar.git",
				Branch:   "gitbranch",
				User:     "gituser",
				Email:    "gitemail@example.com",
				Paths:    []string{"gitpath/"},
				FluxPath: "fluxpath/",
			},
			Operator: api.Operator{
				Label:                      "gitlabel",
				Namespace:                  "fluxnamespace",
				WithHelm:                   api.Enabled(),
				AdditionalFluxArgs:         []string{"--git-poll-interval=30s"},
				AdditionalHelmOperatorArgs: []string{"--log-format=json"},
			},
		}

		installer = &repo.Installer{
			Opts:         mockOpts,
			K8sClientSet: fake.NewSimpleClientset(),
		}

		var err error
		tmpDir, err = os.MkdirTemp(os.TempDir(), "getmanifestsandsecrets")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	Context("GetManifests", func() {
		var (
			err       error
			manifests map[string][]byte
		)

		BeforeEach(func() {
			manifests, err = installer.GetManifests()
		})

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should match expected lengths", func() {
			Expect(len(manifests)).To(Equal(9))
		})

		It("should not error", func() {
			manifestContents := [][]byte{}
			for _, manifest := range manifests {
				manifestContents = append(manifestContents, manifest)
			}

			config := &kubeval.Config{
				IgnoreMissingSchemas: true,
				KubernetesVersion:    "master",
			}
			for _, content := range manifestContents {
				validationResults, err := kubeval.Validate(content, config)
				Expect(err).NotTo(HaveOccurred())
				for _, result := range validationResults {
					Expect(len(result.Errors)).To(Equal(0))
				}
			}
		})
	})
})
