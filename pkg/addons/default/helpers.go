package defaultaddons

import (
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoadAsset return embedded manifest as a runtime.Object
func newList(data []byte) (*metav1.List, error) {
	list, err := kubernetes.NewList(data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading individual resources from manifest")
	}
	return list, nil
}
