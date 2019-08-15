// +build integration

package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	// Namespace is the default Kubernetes namespace under which to install Flux.
	Namespace = "flux"
	// Repository is the default testing Git repository.
	Repository = "git@github.com:eksctl-bot/my-gitops-repo.git"
	// Email is the default testing Git email.
	Email = "eksctl-bot@weave.works"
	// PrivateSSHKeyPath is the default SSH key to use for Git operations.
	PrivateSSHKeyPath = "~/.ssh/eksctl-bot_id_rsa"
	// Name is the default cluster name to test against.
	Name = "autoscaler"
	// Region is the default region to test against.
	Region = "ap-northeast-1"
)

func createBranch(branch string) (string, error) {
	cloneDir, err := ioutil.TempDir(os.TempDir(), "eksctl-install-flux-test-clone-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %s", err)
	}
	if err := gitWith(gitParams{Args: []string{"clone", "-b", "master", Repository, cloneDir}, Dir: cloneDir, Env: gitSSHCommand()}); err != nil {
		return "", err
	}
	if err := gitWith(gitParams{Args: []string{"checkout", "-b", branch}, Dir: cloneDir, Env: gitSSHCommand()}); err != nil {
		return "", err
	}
	if err := gitWith(gitParams{Args: []string{"push", "origin", branch}, Dir: cloneDir, Env: gitSSHCommand()}); err != nil {
		return "", err
	}
	return cloneDir, nil
}

func deleteBranch(branch, cloneDir string) error {
	defer os.RemoveAll(cloneDir)
	return gitWith(gitParams{Args: []string{"push", "origin", "--delete", branch}, Dir: cloneDir, Env: gitSSHCommand()})
}

func getBranch(branch string) (string, error) {
	cloneDir, err := ioutil.TempDir(os.TempDir(), "eksctl-install-flux-test-branch-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %s", err)
	}
	if err := gitWith(gitParams{Args: []string{"clone", "-b", branch, Repository, cloneDir}, Dir: cloneDir, Env: gitSSHCommand()}); err != nil {
		return "", err
	}
	return cloneDir, nil
}

func git(args ...string) error {
	return gitWith(gitParams{
		Args: args,
		Env:  gitSSHCommand(),
	})
}

type gitParams struct {
	Args []string
	Env  []string
	Dir  string
}

func gitWith(params gitParams) error {
	gitCmd := exec.Command("git", params.Args...)
	if params.Env != nil {
		gitCmd.Env = params.Env
	}
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if params.Dir != "" {
		gitCmd.Dir = params.Dir
	}
	return gitCmd.Run()
}

func gitSSHCommand() []string {
	return []string{fmt.Sprintf("GIT_SSH_COMMAND=ssh -i %s", PrivateSSHKeyPath)}
}

func assertFluxManifestsAbsentInGit(branch string) {
	dir, err := getBranch(branch)
	defer os.RemoveAll(dir)
	Expect(err).ShouldNot(HaveOccurred())
	assertDoesNotContainFluxDir(dir)
}

func assertFluxManifestsPresentInGit(branch string) {
	dir, err := getBranch(branch)
	defer os.RemoveAll(dir)
	Expect(err).ShouldNot(HaveOccurred())
	assertContainsFluxDir(dir)
	assertContainsFluxManifests(filepath.Join(dir, Namespace))
}

func assertContainsFluxDir(dir string) {
	fluxDirExists, err := fluxDirExists(dir)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(fluxDirExists).To(BeTrue(), "flux directory could not be found in %s", dir)
}

func assertDoesNotContainFluxDir(dir string) {
	fluxDirExists, err := fluxDirExists(dir)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(fluxDirExists).To(BeFalse(), "flux directory was unexpectedly found in %s", dir)
}

func fluxDirExists(dir string) (bool, error) {
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
		filePath := filepath.Join(dir, f.Name())
		if !f.IsDir() && f.Name() == "flux-account.yaml" {
			assertValidFluxAccountManifest(filePath)
		} else if !f.IsDir() && f.Name() == "flux-deployment.yaml" {
			assertValidFluxDeploymentManifest(filePath)
		} else if !f.IsDir() && f.Name() == "flux-secret.yaml" {
			assertValidFluxSecretManifest(filePath)
		} else if !f.IsDir() && f.Name() == "memcache-dep.yaml" {
			assertValidFluxMemcacheDeploymentManifest(filePath)
		} else if !f.IsDir() && f.Name() == "memcache-svc.yaml" {
			assertValidFluxMemcacheServiceManifest(filePath)
		} else if !f.IsDir() && f.Name() == "namespace.yaml" {
			assertValidFluxNamespaceManifest(filePath)
		} else if f.IsDir() {
			Fail(fmt.Sprintf("Unrecognized directory: %s", f.Name()))
		} else {
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
			Expect(container.Image).To(Equal("docker.io/fluxcd/flux:1.13.3"))
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
			Expect(container.Image).To(Equal("memcached:1.5.15"))
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

func assertFluxPodsAbsentInKubernetes() {
	pods := fluxPods()
	Expect(pods.Items).To(HaveLen(0))
}

func assertFluxPodsPresentInKubernetes() {
	pods := fluxPods()
	Expect(pods.Items).To(HaveLen(2))
	Expect(pods.Items[0].Labels["name"]).To(Equal("flux"))
	Expect(pods.Items[1].Labels["name"]).To(Equal("memcached"))
}

func fluxPods() *corev1.PodList {
	output, err := kubectl("get", "pods", "--namespace", "flux", "--output", "json")
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
