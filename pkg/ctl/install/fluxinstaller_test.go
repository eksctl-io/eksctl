package install

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/instrumenta/kubeval/kubeval"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/yaml"
)

func TestGetManifestsAndSecrets(t *testing.T) {
	mockOpts := &installFluxOpts{
		gitURL:          "git@github.com/foo/bar.git",
		gitBranch:       "gitbranch",
		gitPaths:        []string{"gitpath/"},
		gitLabel:        "gitlabel",
		gitUser:         "gituser",
		gitEmail:        "gitemail@example.com",
		gitFluxPath:     "fluxpath/",
		namespace:       "fluxnamespace",
		tillerNamespace: "tillernamespace",
		noHelmOp:        false,
		noTiller:        false,
		tillerHost:      "tillerhost",
	}
	mockInstaller := &fluxInstaller{
		opts:         mockOpts,
		helmOpTLS:    true,
		k8sClientSet: fake.NewSimpleClientset(),
	}
	tmpDir, err := ioutil.TempDir(os.TempDir(), "getmanifestsandsecrets")
	if err != nil {
		t.Fatal(err)
	}
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
	if err := pki.saveTo(pkiPaths); err != nil {
		t.Fatal(err)
	}

	manifests, secrets, err := mockInstaller.getManifestsAndSecrets(pki, pkiPaths)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifests) != 14 {
		t.Fatalf("unexpected number of manifest files: %d", len(manifests))
	}
	if len(secrets) != 2 {
		t.Fatalf("unexpected number of secrets: %d", len(secrets))
	}

	manifestContents := [][]byte{}
	for _, manifest := range manifests {
		manifestContents = append(manifestContents, manifest)
	}
	for _, secret := range secrets {
		content, err := yaml.Marshal(secret)
		if err != nil {
			t.Fatalf("failed to serialize secret: %s", err)
		}
		manifestContents = append(manifestContents, content)
	}

	config := &kubeval.Config{
		IgnoreMissingSchemas: true,
		KubernetesVersion:    "master",
	}
	for _, content := range manifestContents {
		validationResults, err := kubeval.Validate(content, config)
		if err != nil {
			t.Fatalf("failed to validate manifest: %s\ncontents:\n%s", err, string(content))
		}
		for _, result := range validationResults {
			if len(result.Errors) > 0 {
				t.Errorf("found problems with manifest (Kind %s):\ncontent:\n%s\nerrors: %s",
					result.Kind,
					string(content),
					result.Errors)
			}
		}
	}
}
