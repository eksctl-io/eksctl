package kubernetes

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewNamespace creates a corev1.Namespace object using the provided name.
func NewNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// NewNamespaceYAML returns a YAML string for a Kubernetes Namespace object.
// N.B.: Kubernetes' serializers are not used as unnecessary fields are being
// generated, e.g.: spec, status, creatimeTimestamp.
func NewNamespaceYAML(name string) []byte {
	nsFmt := strings.Join(
		[]string{
			"---",
			"apiVersion: v1",
			"kind: Namespace",
			"metadata: {name: %s}",
		},
		"\n")

	return []byte(fmt.Sprintf(nsFmt, name))
}

// CheckNamespaceExists check if a namespace with a given name already exists, and
// returns boolean or an error
func CheckNamespaceExists(clientSet Interface, name string) (bool, error) {
	_, err := clientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err == nil {
		return true, nil
	}
	if !apierrors.IsNotFound(err) {
		return false, errors.Wrapf(err, "checking whether namespace %q exists", name)
	}
	return false, nil
}

// MaybeCreateNamespace will only create namespace with the given name if it doesn't
// already exist
func MaybeCreateNamespace(clientSet Interface, name string) error {
	exists, err := CheckNamespaceExists(clientSet, name)
	if err != nil {
		return err
	}
	if !exists {
		_, err = clientSet.CoreV1().Namespaces().Create(NewNamespace(name))
		if apierrors.IsAlreadyExists(err) {
			logger.Debug("ignoring failed creation of existing namespace %q", name)
			return nil
		} else if err != nil {
			return err
		}
		logger.Info("created namespace %q", name)
	}
	return nil
}
