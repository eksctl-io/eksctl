package flux

import (
	"io/ioutil"
	"os"

	"github.com/instrumenta/kubeval/kubeval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Installer", func() {
	mockOpts := &api.Git{
		Repo: &api.Repo{
			URL:      "git@github.com/foo/bar.git",
			Branch:   "gitbranch",
			User:     "gituser",
			Email:    "gitemail@example.com",
			Paths:    []string{"gitpath/"},
			FluxPath: "fluxpath/",
		},
		Operator: api.Operator{
			Label:     "gitlabel",
			Namespace: "fluxnamespace",
			WithHelm:  api.Enabled(),
		},
	}
	mockInstaller := &Installer{
		cluster:      &api.ClusterMeta{Name: "cluster-1", Region: "us-west-2"},
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
