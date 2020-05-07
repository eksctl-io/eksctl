// +build integration

package matchers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/integration/utilities/git"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	// Namespace is the default Kubernetes namespace under which to install Flux.
	Namespace = "flux"
)

// AssertFluxManifestsAbsentInGit asserts expected Flux manifests are not present in Git.
func AssertFluxManifestsAbsentInGit(branch, privateSSHKeyPath string) {
	dir, err := git.GetBranch(branch, privateSSHKeyPath)
	defer os.RemoveAll(dir)
	Expect(err).ShouldNot(HaveOccurred())
	assertDoesNotContainFluxDir(dir)
}

// AssertFluxManifestsPresentInGit asserts expected Flux manifests are present in Git.
func AssertFluxManifestsPresentInGit(branch, privateSSHKeyPath string) {
	dir, err := git.GetBranch(branch, privateSSHKeyPath)
	defer os.RemoveAll(dir)
	Expect(err).ShouldNot(HaveOccurred())
	assertContainsFluxDir(dir)
	assertContainsFluxManifests(filepath.Join(dir, Namespace))
}

func assertContainsFluxDir(dir string) {
	fluxDirExists, err := dirExists(dir)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(fluxDirExists).To(BeTrue(), "flux directory could not be found in %s", dir)
}

func assertDoesNotContainFluxDir(dir string) {
	fluxDirExists, err := dirExists(dir)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(fluxDirExists).To(BeFalse(), "flux directory was unexpectedly found in %s", dir)
}

func dirExists(dir string) (bool, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, f := range files {
		if f.Name() == Namespace && f.IsDir() {
			return true, nil
		}
	}
	return false, nil
}

func assertContainsFluxManifests(dir string) {
	// We could have stricter validation by comparing objects in Kubernetes and
	// in the manifests, and ensuring they are equal. However, this may not be
	// very easy to achieve, especially for values which are defaulted by the
	// API server. Hence, for now, we simply ensure that all files & objects
	// are present, and that the main fields of these objects match expected
	// values.
	files, err := ioutil.ReadDir(dir)
	Expect(err).ShouldNot(HaveOccurred())
	for _, f := range files {
		if f.IsDir() {
			Fail(fmt.Sprintf("Unrecognized directory: %s", f.Name()))
		}
		filePath := filepath.Join(dir, f.Name())
		switch f.Name() {
		// Flux resources:
		case "flux-account.yaml":
			assertValidFluxAccountManifest(filePath)
		case "flux-deployment.yaml":
			assertValidFluxDeploymentManifest(filePath)
		case "flux-namespace.yaml":
			assertValidFluxNamespaceManifest(filePath)
		case "flux-secret.yaml":
			assertValidFluxSecretManifest(filePath)
		case "memcache-dep.yaml":
			assertValidFluxMemcacheDeploymentManifest(filePath)
		case "memcache-svc.yaml":
			assertValidFluxMemcacheServiceManifest(filePath)
		// Helm operator resources:
		case "flux-helm-operator-account.yaml":
			assertValidFluxHelmOperatorAccount(filePath)
		case "flux-helm-release-crd.yaml":
			assertValidFluxHelmReleaseCRD(filePath)
		case "helm-operator-deployment.yaml":
			assertValidHelmOperatorDeployment(filePath)
		default:
			Fail(fmt.Sprintf("Unrecognized file: %s", f.Name()))
		}
	}
}

func assertValidFluxAccountManifest(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(3))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "ServiceAccount" {
			sa, ok := item.Object.(*corev1.ServiceAccount)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(sa.Kind).To(Equal("ServiceAccount"))
			Expect(sa.Namespace).To(Equal(Namespace))
			Expect(sa.Name).To(Equal("flux"))
			Expect(sa.Labels["name"]).To(Equal("flux"))
		} else if gvk.Version == "v1beta1" && gvk.Kind == "ClusterRole" {
			cr, ok := item.Object.(*rbacv1beta1.ClusterRole)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(cr.Kind).To(Equal("ClusterRole"))
			Expect(cr.Name).To(Equal("flux"))
			Expect(cr.Labels["name"]).To(Equal("flux"))
		} else if gvk.Version == "v1beta1" && gvk.Kind == "ClusterRoleBinding" {
			crb, ok := item.Object.(*rbacv1beta1.ClusterRoleBinding)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(crb.Kind).To(Equal("ClusterRoleBinding"))
			Expect(crb.Name).To(Equal("flux"))
			Expect(crb.Labels["name"]).To(Equal("flux"))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidFluxDeploymentManifest(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(1))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "Deployment" {
			deployment, ok := item.Object.(*appsv1.Deployment)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(deployment.Kind).To(Equal("Deployment"))

			Expect(deployment.Namespace).To(Equal(Namespace))
			Expect(deployment.Name).To(Equal("flux"))
			Expect(*deployment.Spec.Replicas).To(Equal(int32(1)))
			Expect(deployment.Spec.Template.Labels["name"]).To(Equal("flux"))
			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Name).To(Equal("flux"))
			Expect(container.Image).To(Equal("docker.io/fluxcd/flux:1.19.0"))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidFluxSecretManifest(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(1))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "Secret" {
			secret, ok := item.Object.(*corev1.Secret)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(secret.Kind).To(Equal("Secret"))
			Expect(secret.Namespace).To(Equal(Namespace))
			Expect(secret.Name).To(Equal("flux-git-deploy"))
			Expect(secret.Type).To(Equal(corev1.SecretTypeOpaque))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidFluxMemcacheDeploymentManifest(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(1))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "Deployment" {
			deployment, ok := item.Object.(*appsv1.Deployment)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(deployment.Kind).To(Equal("Deployment"))
			Expect(deployment.Namespace).To(Equal(Namespace))
			Expect(deployment.Name).To(Equal("memcached"))
			Expect(*deployment.Spec.Replicas).To(Equal(int32(1)))
			Expect(deployment.Spec.Template.Labels["name"]).To(Equal("memcached"))
			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Name).To(Equal("memcached"))
			Expect(container.Image).To(Equal("memcached:1.5.20"))
			Expect(container.Ports).To(HaveLen(1))
			Expect(container.Ports[0].ContainerPort).To(Equal(int32(11211)))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidFluxMemcacheServiceManifest(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(1))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "Service" {
			service, ok := item.Object.(*corev1.Service)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(service.Kind).To(Equal("Service"))

			Expect(service.Namespace).To(Equal(Namespace))
			Expect("memcached").To(Equal(service.Name))
			Expect(service.Spec.Ports).To(HaveLen(1))
			Expect(service.Spec.Ports[0].Name).To(Equal("memcached"))
			Expect(service.Spec.Ports[0].Port).To(Equal(int32(11211)))
			Expect(service.Spec.Selector).To(HaveLen(1))
			Expect(service.Spec.Selector["name"]).To(Equal("memcached"))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidFluxNamespaceManifest(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(1))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "Namespace" {
			ns, ok := item.Object.(*corev1.Namespace)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(ns.Kind).To(Equal("Namespace"))
			Expect(ns.Name).To(Equal(Namespace))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidFluxHelmOperatorAccount(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(3))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "ServiceAccount" {
			sa, ok := item.Object.(*corev1.ServiceAccount)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(sa.Kind).To(Equal("ServiceAccount"))
			Expect(sa.Namespace).To(Equal(Namespace))
			Expect(sa.Name).To(Equal("helm-operator"))
			Expect(sa.Labels["name"]).To(Equal("helm-operator"))
		} else if gvk.Version == "v1beta1" && gvk.Kind == "ClusterRole" {
			cr, ok := item.Object.(*rbacv1beta1.ClusterRole)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(cr.Kind).To(Equal("ClusterRole"))
			Expect(cr.Name).To(Equal("helm-operator"))
			Expect(cr.Labels["name"]).To(Equal("helm-operator"))
		} else if gvk.Version == "v1beta1" && gvk.Kind == "ClusterRoleBinding" {
			crb, ok := item.Object.(*rbacv1beta1.ClusterRoleBinding)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(crb.Kind).To(Equal("ClusterRoleBinding"))
			Expect(crb.Name).To(Equal("helm-operator"))
			Expect(crb.Labels["name"]).To(Equal("helm-operator"))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidFluxHelmReleaseCRD(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(1))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1beta1" && gvk.Kind == "CustomResourceDefinition" {
			crd, ok := item.Object.(*apiextensionsv1beta1.CustomResourceDefinition)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(crd.Kind).To(Equal("CustomResourceDefinition"))
			Expect(crd.Name).To(Equal("helmreleases.helm.fluxcd.io"))
			Expect(crd.Spec.Group).To(Equal("helm.fluxcd.io"))
			Expect(crd.Spec.Names.Kind).To(Equal("HelmRelease"))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func assertValidHelmOperatorDeployment(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	list, err := kubernetes.NewRawExtensions(bytes)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(list).To(HaveLen(1))
	for _, item := range list {
		gvk := item.Object.GetObjectKind().GroupVersionKind()
		if gvk.Version == "v1" && gvk.Kind == "Deployment" {
			deployment, ok := item.Object.(*appsv1.Deployment)
			Expect(ok).To(BeTrue(), "Failed to convert object of type %T to %s", item.Object, gvk.Kind)
			Expect(deployment.Kind).To(Equal("Deployment"))
			Expect(deployment.Namespace).To(Equal(Namespace))
			Expect(deployment.Name).To(Equal("helm-operator"))
			Expect(*deployment.Spec.Replicas).To(Equal(int32(1)))
			Expect(deployment.Spec.Template.Labels["name"]).To(Equal("helm-operator"))
			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Name).To(Equal("helm-operator"))
			Expect(container.Image).To(Equal("docker.io/fluxcd/helm-operator:1.0.0"))
		} else {
			Fail(fmt.Sprintf("Unsupported Kubernetes object. Got %s object with version %s in: %s", gvk.Kind, gvk.Version, fileName))
		}
	}
}

func AssertFluxPodsAbsentInKubernetes(kubeconfigPath string) {
	pods := fluxPods(kubeconfigPath)
	Expect(pods.Items).To(HaveLen(0))
}

func AssertFluxPodsPresentInKubernetes(kubeconfigPath string) {
	pods := fluxPods(kubeconfigPath)
	Expect(pods.Items).To(HaveLen(3))
	Expect(pods.Items[0].Labels["name"]).To(Equal("flux"))
	Expect(pods.Items[1].Labels["name"]).To(Equal("helm-operator"))
	Expect(pods.Items[2].Labels["name"]).To(Equal("memcached"))
}

func fluxPods(kubeconfigPath string) *corev1.PodList {
	output, err := kubectl("get", "pods", "--namespace", "flux", "--output", "json", "--kubeconfig", kubeconfigPath)
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
