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

// CreateDeployment creates a deployment in the given namespace.
func (test *Test) createDeployment(namespace string, d *appsv1.Deployment) error {
	test.Infof("creating deployment %s", d.Name)
	d.Namespace = namespace
	_, err := test.harness.kubeClient.AppsV1beta2().Deployments(namespace).Create(d)
	if err != nil {
		return errors.Wrapf(err, "failed to create deployment %s", d.Name)
	}
	return nil
}

// CreateDeployment creates a deployment in the given namespace.
func (test *Test) CreateDeployment(namespace string, d *appsv1.Deployment) {
	err := test.createDeployment(namespace, d)
	test.err(err)
}

func (test *Test) loadDeployment(manifestPath string) (*appsv1.Deployment, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, errors.Wrapf(err, "failed to decode deployment %s", manifestPath)
	}

	return &dep, nil
}

// LoadDeployment loads a deployment from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestsDirectory.
func (test *Test) LoadDeployment(manifestPath string) *appsv1.Deployment {
	dep, err := test.loadDeployment(manifestPath)
	test.err(err)
	return dep
}

func (test *Test) createDeploymentFromFile(namespace string, manifestPath string) (*appsv1.Deployment, error) {
	d, err := test.loadDeployment(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createDeployment(namespace, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CreateDeploymentFromFile creates a deployment from a manifest file in the given namespace.
func (test *Test) CreateDeploymentFromFile(namespace string, manifestPath string) *appsv1.Deployment {
	d, err := test.createDeploymentFromFile(namespace, manifestPath)
	test.err(err)
	return d
}

// GetDeployment returns Deployment if it exists or error if it doesn't.
func (test *Test) GetDeployment(ns, name string) (*appsv1.Deployment, error) {
	d, err := test.harness.kubeClient.AppsV1beta2().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return d, nil
}

// waitForDeploymentReady waits until all replica pods are running and ready.
func (test *Test) waitForDeploymentReady(d *appsv1.Deployment, timeout time.Duration) error {
	numReady := int32(0)

	test.Infof("waiting for deployment %s to be ready", d.Name)

	return wait.Poll(time.Second, timeout, func() (bool, error) {
		current, err := test.GetDeployment(d.Namespace, d.Name)
		if err != nil {
			return false, err
		}

		if numReady != current.Status.ReadyReplicas {
			numReady = current.Status.AvailableReplicas
			test.Debugf("%s number of replicas: %d/%d", d.Name, numReady, *d.Spec.Replicas)
		}
		if current.Status.ReadyReplicas == int32(*d.Spec.Replicas) {
			return true, nil
		}

		return false, nil
	})
}

// WaitForDeploymentReady waits until all replica pods are running and ready.
func (test *Test) WaitForDeploymentReady(d *appsv1.Deployment, timeout time.Duration) {
	err := test.waitForDeploymentReady(d, timeout)
	test.err(err)
}

// deleteDeployment deletes a deployment in the given namespace.
func (test *Test) deleteDeployment(d *appsv1.Deployment) error {
	test.Infof("deleting deployment %s ", d.Name)

	d, err := test.GetDeployment(d.Namespace, d.Name)
	if err != nil {
		return err
	}

	zero := int32(0)
	d.Spec.Replicas = &zero

	d, err = test.harness.kubeClient.AppsV1beta2().Deployments(d.Namespace).Update(d)
	if err != nil {
		return err
	}
	return test.harness.kubeClient.AppsV1beta2().Deployments(d.Namespace).Delete(d.Name, &metav1.DeleteOptions{})
}

// DeleteDeployment deletes a deployment in the given namespace.
func (test *Test) DeleteDeployment(d *appsv1.Deployment) {
	test.err(test.deleteDeployment(d))
}

// waitForDeploymentDeleted waits until a deleted deployment has disappeared from the cluster.
func (test *Test) waitForDeploymentDeleted(d *appsv1.Deployment, timeout time.Duration) error {
	test.Infof("waiting for deployment %s to be deleted", d.Name)

	return wait.Poll(time.Second, timeout, func() (bool, error) {

		_, err := test.GetDeployment(d.Namespace, d.Name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}

			return false, err
		}

		return false, nil
	})
}

// WaitForDeploymentDeleted waits until a deleted deployment has disappeared from the cluster.
func (test *Test) WaitForDeploymentDeleted(d *appsv1.Deployment, timeout time.Duration) {
	test.err(test.waitForDeploymentDeleted(d, timeout))
}
