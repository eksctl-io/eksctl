package kubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/weaveworks/eksctl/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func TestNamespace(t *testing.T) {
	ns := Namespace("flux")
	assert.Equal(t, "v1", ns.APIVersion)
	assert.Equal(t, "Namespace", ns.Kind)
	assert.Equal(t, "flux", ns.Name)
	assert.Equal(t, "flux", ns.Labels["name"])
	bytes, err := yaml.Marshal(ns)
	assert.NoError(t, err)
	assert.Equal(t, `apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: null
  labels:
    name: flux
  name: flux
spec: {}
status: {}
`, string(bytes))
}

func TestNamespaceYAML(t *testing.T) {
	nsBytes := NamespaceYAML("flux")
	var ns corev1.Namespace
	err := yaml.Unmarshal(nsBytes, &ns)
	assert.NoError(t, err)
	assert.Equal(t, "v1", ns.APIVersion)
	assert.Equal(t, "Namespace", ns.Kind)
	assert.Equal(t, "flux", ns.Name)
	assert.Equal(t, "flux", ns.Labels["name"])
	assert.Equal(t, ns, *Namespace("flux"))
}
