package harness

import (
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1beta2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// createDaemonSet creates a daemonset in the given namespace.
func (test *Test) createDaemonSet(namespace string, d *appsv1.DaemonSet) error {
	test.Infof("creating daemonset %s", d.Name)
	d.Namespace = namespace
	_, err := test.harness.kubeClient.AppsV1beta2().DaemonSets(namespace).Create(d)
	if err != nil {
		return errors.Wrapf(err, "failed to create daemonset %s", d.Name)
	}
	return nil
}

// CreateDaemonSet creates a daemonset in the given namespace.
func (test *Test) CreateDaemonSet(namespace string, d *appsv1.DaemonSet) {
	err := test.createDaemonSet(namespace, d)
	test.err(err)
}

func (test *Test) loadDaemonSet(manifestPath string) (*appsv1.DaemonSet, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := appsv1.DaemonSet{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, errors.Wrapf(err, "failed to decode daemonset %s", manifestPath)
	}

	return &dep, nil
}

// LoadDaemonSet loads a daemonset from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestsDirectory.
func (test *Test) LoadDaemonSet(manifestPath string) *appsv1.DaemonSet {
	dep, err := test.loadDaemonSet(manifestPath)
	test.err(err)
	return dep
}

func (test *Test) createDaemonSetFromFile(namespace string, manifestPath string) (*appsv1.DaemonSet, error) {
	d, err := test.loadDaemonSet(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createDaemonSet(namespace, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CreateDaemonSetFromFile creates a daemonset from a manifest file in the given namespace.
func (test *Test) CreateDaemonSetFromFile(namespace string, manifestPath string) *appsv1.DaemonSet {
	d, err := test.createDaemonSetFromFile(namespace, manifestPath)
	test.err(err)
	return d
}

// GetDaemonSet returns daemonset if it exists or error if it doesn't.
func (test *Test) GetDaemonSet(ns, name string) (*appsv1.DaemonSet, error) {
	d, err := test.harness.kubeClient.AppsV1beta2().DaemonSets(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return d, nil
}

// waitForDaemonSetReady waits until all replica pods are running and ready.
func (test *Test) waitForDaemonSetReady(d *appsv1.DaemonSet, timeout time.Duration) error {

	test.Infof("waiting for daemonset %s to be ready", d.Name)

	return wait.Poll(time.Second, timeout, func() (bool, error) {
		current, err := test.GetDaemonSet(d.Namespace, d.Name)
		if err != nil {
			return false, err
		}

		if current.Status.DesiredNumberScheduled == current.Status.NumberReady {
			return true, nil
		}

		return false, nil
	})
}

// WaitForDaemonSetReady waits until all replica pods are running and ready.
func (test *Test) WaitForDaemonSetReady(d *appsv1.DaemonSet, timeout time.Duration) {
	err := test.waitForDaemonSetReady(d, timeout)
	test.err(err)
}

// deleteDaemonSet deletes a daemonset in the given namespace.
func (test *Test) deleteDaemonSet(d *appsv1.DaemonSet) error {
	test.Infof("deleting daemonset %s ", d.Name)

	d, err := test.GetDaemonSet(d.Namespace, d.Name)
	if err != nil {
		return err
	}

	return test.harness.kubeClient.AppsV1beta2().DaemonSets(d.Namespace).Delete(d.Name, &metav1.DeleteOptions{})
}

// DeleteDaemonSet deletes a daemonset in the given namespace.
func (test *Test) DeleteDaemonSet(d *appsv1.DaemonSet) {
	test.err(test.deleteDaemonSet(d))
}

// waitForDaemonSetDeleted waits until a deleted daemonset has disappeared from the cluster.
func (test *Test) waitForDaemonSetDeleted(d *appsv1.DaemonSet, timeout time.Duration) error {
	test.Infof("waiting for daemonset %s to be deleted", d.Name)

	return wait.Poll(time.Second, timeout, func() (bool, error) {

		_, err := test.GetDaemonSet(d.Namespace, d.Name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}

			return false, err
		}

		return false, nil
	})
}

// WaitForDaemonSetDeleted waits until a deleted daemonset has disappeared from the cluster.
func (test *Test) WaitForDaemonSetDeleted(d *appsv1.DaemonSet, timeout time.Duration) {
	test.err(test.waitForDaemonSetDeleted(d, timeout))
}
