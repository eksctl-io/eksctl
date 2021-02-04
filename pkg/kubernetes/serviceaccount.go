package kubernetes

import (
	"context"

	"github.com/pkg/errors"
	"github.com/weaveworks/logger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewServiceAccount creates a corev1.ServiceAccount object using the provided meta.
func NewServiceAccount(meta metav1.ObjectMeta) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: meta,
	}
}

// CheckServiceAccountExists check if a serviceaccount with a given name already exists, and
// returns boolean or an error
func CheckServiceAccountExists(clientSet Interface, meta metav1.ObjectMeta) (bool, error) {
	name := meta.Namespace + "/" + meta.Name
	_, err := clientSet.CoreV1().ServiceAccounts(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
	if err == nil {
		return true, nil
	}
	if !apierrors.IsNotFound(err) {
		return false, errors.Wrapf(err, "checking whether serviceaccount %q exists", name)
	}
	return false, nil
}

// MaybeCreateServiceAccountOrUpdateMetadata will only create serviceaccount with the given name if
// it doesn't already exist, it will also create namespace if needed; if serviceaccount exists, new
// labels and annotations will get added, all user-set label and annotation keys that are not set in
// meta will be retained
func MaybeCreateServiceAccountOrUpdateMetadata(clientSet Interface, meta metav1.ObjectMeta) error {
	name := meta.Namespace + "/" + meta.Name

	if err := MaybeCreateNamespace(clientSet, meta.Namespace); err != nil {
		return err
	}
	exists, err := CheckServiceAccountExists(clientSet, meta)
	if err != nil {
		return err
	}
	if !exists {
		_, err = clientSet.CoreV1().ServiceAccounts(meta.Namespace).Create(context.TODO(), NewServiceAccount(meta), metav1.CreateOptions{})
		if err != nil {
			return err
		}
		logger.Info("created serviceaccount %q", name)
		return nil
	}

	logger.Info("serviceaccount %q already exists", name)

	current, err := clientSet.CoreV1().ServiceAccounts(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	updateRequired := false

	mergeMetadata := func(src, dst map[string]string) {
		for key, value := range src {
			currentValue, ok := dst[key]
			updateRequired = !ok || ok && currentValue != value
			dst[key] = value
		}
	}

	if current.Annotations == nil {
		current.Annotations = make(map[string]string)
	}
	mergeMetadata(meta.Annotations, current.Annotations)

	if current.Labels == nil {
		current.Labels = make(map[string]string)
	}
	mergeMetadata(meta.Labels, current.Labels)

	if !updateRequired {
		logger.Info("serviceaccount %q is already up-to-date", name)
		return nil
	}
	_, err = clientSet.CoreV1().ServiceAccounts(meta.Namespace).Update(context.TODO(), current, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	logger.Info("updated serviceaccount %q", name)
	return nil
}

// MaybeDeleteServiceAccount will only delete the serviceaccount if it exists
func MaybeDeleteServiceAccount(clientSet Interface, meta metav1.ObjectMeta) error {
	name := meta.Namespace + "/" + meta.Name
	exists, err := CheckServiceAccountExists(clientSet, meta)
	if err != nil {
		return err
	}
	if !exists {
		logger.Info("serviceaccount %q was already deleted", name)
		return nil
	}
	err = clientSet.CoreV1().ServiceAccounts(meta.Namespace).Delete(context.TODO(), meta.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	logger.Info("deleted serviceaccount %q", name)
	return nil
}
