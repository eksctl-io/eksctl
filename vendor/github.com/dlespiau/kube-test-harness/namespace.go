package harness

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (test *Test) createNamespace(name string) (*v1.Namespace, error) {
	test.Infof("creating namespace %s", name)

	namespace, err := test.harness.kubeClient.Core().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create namespace with name %v", name))
	}
	return namespace, nil
}

// CreateNamespace creates a new namespace.
func (test *Test) CreateNamespace(name string) {

	_, err := test.createNamespace(name)
	test.err(err)

	test.addNamespace(name)

	test.addFinalizer(func() error {
		if err := test.deleteNamespace(name); err != nil {
			return err
		}
		return nil
	})
}

func (test *Test) deleteNamespace(name string) error {
	test.Infof("deleting namespace %s", name)

	test.removeNamespace(name)

	return test.harness.kubeClient.Core().Namespaces().Delete(name, nil)
}
