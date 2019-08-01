package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Namespace creates a corev1.Namespace object with the standard labels set
// using the provided name.
func Namespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"name": name},
		},
	}
}

const namespaceTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    name: %s
  name: %s
`

// NamespaceYAML returns a YAML string for a Kubernetes Namespace object with
// the standard labels set using the provided name.
// N.B.: Kubernetes' serializers are not used as unnecessary fields are being
// generated, e.g.: spec, status, creatimeTimestamp.
func NamespaceYAML(name string) []byte {
	return []byte(fmt.Sprintf(namespaceTemplate, name, name))
}
