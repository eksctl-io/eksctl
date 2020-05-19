package flux

import (
	"io/ioutil"
	"os"

	"github.com/instrumenta/kubeval/kubeval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/git"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Installer", func() {
	mockOpts := &InstallOpts{
		GitOptions: git.Options{
			URL:    "git@github.com/foo/bar.git",
			Branch: "gitbranch",
			User:   "gituser",
			Email:  "gitemail@example.com",
		},
		GitPaths:    []string{"gitpath/"},
		GitLabel:    "gitlabel",
		GitFluxPath: "fluxpath/",
		Namespace:   "fluxnamespace",
		WithHelm:    true,
	}
	mockInstaller := &Installer{
		opts:         mockOpts,
		k8sClientSet: fake.NewSimpleClientset(),
	}
	tmpDir, err := ioutil.TempDir(os.TempDir(), "getmanifestsandsecrets")
	It("should not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	defer os.RemoveAll(tmpDir)
	It("should not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	manifests, err := mockInstaller.getManifests()
	It("should not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})
	It("should match expected lengths", func() {
		Expect(len(manifests)).To(Equal(9))
	})

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
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		for _, result := range validationResults {
			It("should not error", func() {
				Expect(len(result.Errors)).To(Equal(0))
			})
		}
	}
})
