//go:build integration
// +build integration

package matchers

import (
	"os/exec"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

const (
	flux2Namespace = "flux-system"
)

var (
	flux2Components = []string{"helm-controller", "kustomize-controller", "notification-controller", "source-controller"}
)

func AssertFluxPodsAbsentInKubernetes(kubeconfigPath, namespace string) {
	pods := fluxPods(kubeconfigPath, namespace)
	Expect(pods.Items).To(HaveLen(0))
}

func AssertFlux2PodsPresentInKubernetes(kubeconfigPath string) {
	assertLabelPresent(kubeconfigPath, flux2Namespace, "app", flux2Components)
}

func assertLabelPresent(kubeconfigPath, namespace, label string, components []string) {
	pods := fluxPods(kubeconfigPath, namespace)
	Expect(len(pods.Items)).To(Equal(len(components)))
	for i, component := range components {
		Expect(pods.Items[i].Labels[label]).To(Equal(component))
	}
}

func fluxPods(kubeconfigPath, namespace string) *corev1.PodList {
	output, err := kubectl("get", "pods", "--namespace", namespace, "--output", "json", "--kubeconfig", kubeconfigPath)
	Expect(err).ShouldNot(HaveOccurred())
	var pods corev1.PodList
	err = yaml.Unmarshal(output, &pods)
	Expect(err).ShouldNot(HaveOccurred())
	return &pods
}

func kubectl(args ...string) ([]byte, error) {
	kubectlCmd := exec.Command("kubectl", args...)
	return kubectlCmd.CombinedOutput()
}
