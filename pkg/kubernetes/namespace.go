package kubernetes

import (
	"fmt"
	"math"
	"strings"
	"time"

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

// DeleteNamespace deletes the provided namespace and all resources namespaced
// under it. It waits until all "common" (pods, deployments, services, etc.)
// resources are deleted before returning.
// N.B.: this method does NOT delete non-namespaced resources, e.g.:
//       ClusterRoles or ClusterRoleBindings.
func DeleteNamespace(clientSet Interface, namespace string) error {
	if err := clientSet.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{}); err != nil {
		return errors.Wrapf(err, "failed to delete namespace %s", namespace)
	}
	logger.Info("deleting namespace %s", namespace)
	// There may be race conditions if we return right after asynchronously
	// deleting the namespace, hence we wait for commonly used resources.
	// This isn't perfect in any way, but given there is no way to directly
	// list all namespaced resources without using kubectl [1] this will
	// hopefully do in most cases.
	// [1]: kubectl api-resources --verbs=list --namespaced -o name \
	//      | xargs -n 1 kubectl get --show-kind --ignore-not-found -n <namespace>
	return waitForCommonNamespacedResourcesDeletion(clientSet, namespace)
}

func waitForCommonNamespacedResourcesDeletion(clientSet Interface, namespace string) error {
	if err := waitForPodsDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForDaemonSetsDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForDeploymentsDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForExtensionDeploymentsDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForServicesDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForConfigMapsDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForSecretsDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForServiceAccountsDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForRolesDeletion(clientSet, namespace); err != nil {
		return err
	}
	if err := waitForNamespaceDeletion(clientSet, namespace); err != nil {
		return err
	}
	logger.Info("deleted namespace %s", namespace)
	return nil
}

func waitForPodsDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list pods under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for pods under namespace %s to be deleted", namespace)
	}
}

func waitForDaemonSetsDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list daemonsets under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for daemonsets under namespace %s to be deleted", namespace)
	}
}

func waitForDeploymentsDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list corev1/deployments under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for corev1/deployments under namespace %s to be deleted", namespace)
	}
}

func waitForExtensionDeploymentsDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.ExtensionsV1beta1().Deployments(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list extensionsv1beta1/deployments under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for extensionsv1beta1/deployments under namespace %s to be deleted", namespace)
	}
}

func waitForServicesDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list services under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for services under namespace %s to be deleted", namespace)
	}
}

func waitForConfigMapsDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list config maps under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for config maps under namespace %s to be deleted", namespace)
	}
}

func waitForSecretsDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list secrets under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for secrets under namespace %s to be deleted", namespace)
	}
}

func waitForServiceAccountsDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.CoreV1().ServiceAccounts(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list service accounts under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for service accounts under namespace %s to be deleted", namespace)
	}
}

func waitForRolesDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		list, err := clientSet.RbacV1beta1().Roles(namespace).List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to list roles under namespace %s", namespace)
		}
		if len(list.Items) == 0 {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for roles under namespace %s to be deleted", namespace)
	}
}

func waitForNamespaceDeletion(clientSet Interface, namespace string) error {
	attempt := 0
	for {
		exists, err := CheckNamespaceExists(clientSet, namespace)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
		time.Sleep(exponentialBackOff(attempt) * time.Second)
		attempt++
		logger.Info("waiting for namespace %s to be deleted", namespace)
	}
}

func exponentialBackOff(attempt int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt)))
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
		if err != nil {
			return err
		}
		logger.Info("created namespace %q", name)
	}
	return nil
}
