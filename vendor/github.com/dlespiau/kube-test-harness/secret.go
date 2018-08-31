package harness

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (test *Test) createSecret(namespace string, secret *v1.Secret) error {
	secret.Namespace = namespace
	if _, err := test.harness.kubeClient.CoreV1().Secrets(namespace).Create(secret); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create secret %v ", secret.Name))
	}
	return nil
}

// CreateSecret creates a secret in the given namespace.
func (test *Test) CreateSecret(namespace string, secret *v1.Secret) {
	err := test.createSecret(namespace, secret)
	test.err(err)
}

func (test *Test) loadSecret(manifestPath string) (*v1.Secret, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := v1.Secret{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, errors.Wrapf(err, "failed to decode secret %s", manifestPath)
	}

	return &dep, nil
}

// LoadSecret loads a secret from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestsDirectory.
func (test *Test) LoadSecret(manifestPath string) *v1.Secret {
	dep, err := test.loadSecret(manifestPath)
	test.err(err)
	return dep
}

func (test *Test) createSecretFromFile(namespace string, manifestPath string) (*v1.Secret, error) {
	s, err := test.loadSecret(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createSecret(namespace, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateSecretFromFile creates a secret from a manifest file in the given namespace.
func (test *Test) CreateSecretFromFile(namespace string, manifestPath string) *v1.Secret {
	d, err := test.createSecretFromFile(namespace, manifestPath)
	test.err(err)
	return d
}

func (test *Test) deleteSecret(secret *v1.Secret) error {
	if err := test.harness.kubeClient.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, nil); err != nil {
		return errors.Wrap(err, fmt.Sprintf("deleting secret %v failed", secret.Name))
	}
	return nil
}

// DeleteSecret deletes a secret.
func (test *Test) DeleteSecret(secret *v1.Secret) {
	err := test.deleteSecret(secret)
	test.err(err)
}

// GetSecret returns a Secret object if it exists or error.
func (test *Test) GetSecret(ns, name string) (*v1.Secret, error) {
	s, err := test.harness.kubeClient.CoreV1().Secrets(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return s, nil
}

// WaitForSecretReady waits until Secret is created, otherwise times out.
func (test *Test) WaitForSecretReady(secret *v1.Secret, timeout time.Duration) {
	err := test.waitForSecretReady(secret.Namespace, secret.Name, timeout)
	test.err(err)
}

func (test *Test) waitForSecretReady(ns, name string, timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		_, err := test.GetSecret(ns, name)
		if err != nil {
			return false, err
		}

		return true, nil
	})
}
