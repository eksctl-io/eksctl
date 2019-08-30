package flux

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/instrumenta/kubeval/kubeval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/git"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/yaml"
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
	pki := &publicKeyInfrastructure{
		caCertificate:     []byte("caCertificateContent"),
		caKey:             []byte("caKeyContent"),
		serverCertificate: []byte("caCertificateContent"),
		serverKey:         []byte("serverKeyContent"),
		clientCertificate: []byte("clientCertificateContent"),
		clientKey:         []byte("clientKeyContent"),
	}
	pkiPaths := &publicKeyInfrastructurePaths{
		caKey:             filepath.Join(tmpDir, "ca-key.pem"),
		caCertificate:     filepath.Join(tmpDir, "ca.pem"),
		serverKey:         filepath.Join(tmpDir, "tiller-key.pem"),
		serverCertificate: filepath.Join(tmpDir, "tiller.pem"),
		clientKey:         filepath.Join(tmpDir, "flux-helm-operator-key.pem"),
		clientCertificate: filepath.Join(tmpDir, "flux-helm-operator.pem"),
	}
	err = pki.saveTo(pkiPaths)
	It("should not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	manifests, secrets, err := mockInstaller.getManifestsAndSecrets(pki, pkiPaths)
	It("should not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})
	It("should match expected lengths", func() {
		Expect(len(manifests)).To(Equal(13))
		Expect(len(secrets)).To(Equal(2))
	})

	manifestContents := [][]byte{}
	for _, manifest := range manifests {
		manifestContents = append(manifestContents, manifest)
	}
	for _, secret := range secrets {
		content, err := yaml.Marshal(secret)
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		manifestContents = append(manifestContents, content)
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
