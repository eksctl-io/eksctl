package harness

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (test *Test) createService(namespace string, service *v1.Service) error {
	service.Namespace = namespace
	if _, err := test.harness.kubeClient.CoreV1().Services(namespace).Create(service); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create service %v ", service.Name))
	}
	return nil
}

// CreateService creates a service in the given namespace.
func (test *Test) CreateService(namespace string, service *v1.Service) {
	err := test.createService(namespace, service)
	test.err(err)
}

func (test *Test) loadService(manifestPath string) (*v1.Service, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := v1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, errors.Wrapf(err, "failed to decode service %s", manifestPath)
	}

	return &dep, nil
}

// LoadService loads a service from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestsDirectory.
func (test *Test) LoadService(manifestPath string) *v1.Service {
	dep, err := test.loadService(manifestPath)
	test.err(err)
	return dep
}

func (test *Test) createServiceFromFile(namespace string, manifestPath string) (*v1.Service, error) {
	s, err := test.loadService(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createService(namespace, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateServiceFromFile creates a service from a manifest file in the given namespace.
func (test *Test) CreateServiceFromFile(namespace string, manifestPath string) *v1.Service {
	d, err := test.createServiceFromFile(namespace, manifestPath)
	test.err(err)
	return d
}

func (test *Test) waitForServiceReady(service *v1.Service) error {
	test.Infof("waiting for service %s to be ready", service.Name)
	err := wait.Poll(time.Second, time.Minute*5, func() (bool, error) {
		endpoints, err := test.getEndpoints(service.Namespace, service.Name)
		if err != nil {
			return false, err
		}
		if len(endpoints.Subsets) != 0 && len(endpoints.Subsets[0].Addresses) > 0 {
			return true, nil
		}
		return false, nil
	})
	return err
}

// WaitForServiceReady will wait until at least one endpoint backing up the service is ready.
func (test *Test) WaitForServiceReady(service *v1.Service) {
	test.err(test.waitForServiceReady(service))
}

func (test *Test) deleteService(service *v1.Service) error {
	if err := test.harness.kubeClient.CoreV1().Services(service.Namespace).Delete(service.Name, nil); err != nil {
		return errors.Wrap(err, fmt.Sprintf("deleting service %v failed", service.Name))
	}
	return nil
}

// DeleteService deletes a service.
func (test *Test) DeleteService(service *v1.Service) {
	err := test.deleteService(service)
	test.err(err)
}

func (test *Test) waitForServiceDeleted(service *v1.Service) error {
	test.Infof("waiting for service %s to be deleted", service.Name)

	err := wait.Poll(5*time.Second, time.Minute, func() (bool, error) {
		_, err := test.getEndpoints(service.Namespace, service.Name)
		if err != nil {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return errors.Wrap(err, "waiting for service to go away failed")
	}

	return nil
}

// WaitForServiceDeleted waits until deleted service has disappeared from the cluster.
func (test *Test) WaitForServiceDeleted(service *v1.Service) {
	test.err(test.waitForServiceDeleted(service))
}

func (test *Test) getEndpoints(namespace, serviceName string) (*v1.Endpoints, error) {
	endpoints, err := test.harness.kubeClient.CoreV1().Endpoints(namespace).Get(serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to request endpoints for service %v", serviceName))
	}
	return endpoints, nil
}
