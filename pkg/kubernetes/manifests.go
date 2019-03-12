package kubernetes

import (
	"bytes"
	"io"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

// NewList decoded data into a list of Kubernetes resources
func NewList(data []byte) (*metav1.List, error) {
	list := metav1.List{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer(data), 4096)

	for {
		obj := new(runtime.RawExtension)
		err := decoder.Decode(obj)
		if err != nil {
			if err == io.EOF {
				return &list, nil
			}
			return nil, err
		}
		if err := AppendFlattened(&list, *obj); err != nil {
			return nil, err
		}
	}
}

// AppendFlattened will appaned newItem to list; making sure that raw newItem is decoded
// and flattended when its another list
func AppendFlattened(components *metav1.List, component runtime.RawExtension) error {
	if component.Object != nil {
		gvk := component.Object.GetObjectKind().GroupVersionKind()
		if strings.HasSuffix(gvk.Kind, "List") {
			// must use *corev1.List, as it cannot be converted to *metav1.List
			newList := component.Object.(*corev1.List)
			for _, item := range (*newList).Items {
				// we attempt to recurse here, but most likely
				// we will have to try decoding component.Raw
				if err := AppendFlattened(components, item); err != nil {
					return err
				}
			}
			return nil
		}
		components.Items = append(components.Items, component)
		return nil
	}
	obj, err := runtime.Decode(scheme.Codecs.UniversalDeserializer(), component.Raw)
	if err != nil {
		return errors.Wrapf(err, "decoding object")
	}
	return AppendFlattened(components, runtime.RawExtension{Object: obj})
}
