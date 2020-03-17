package defaultaddons

import (
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoadAsset return embedded manifest as a runtime.Object
func LoadAsset(name, ext string) (*metav1.List, error) {
	data, err := Asset(name + "." + ext)
	if err != nil {
		return nil, errors.Wrapf(err, "decoding embedded manifest for %q", name)
	}
	list, err := kubernetes.NewList(data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading individual resources from manifest for %q", name)
	}
	return list, nil
}
