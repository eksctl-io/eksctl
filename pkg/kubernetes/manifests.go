package kubernetes

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

// ConcatManifests joins the provided manifests (as byte arrays) into one single
// manifest. This can be useful to only have one I/O operation with Kubernetes
// down the line, when trying to apply these manifests.
func ConcatManifests(manifests ...[]byte) []byte {
	return bytes.Join(manifests, separator)
}

var separator = []byte("---\n")

// NewRawExtensions decodes the provided manifest's bytes into "raw extension"
// Kubernetes objects. These can then be passed to NewRawResource.
func NewRawExtensions(manifest []byte) ([]runtime.RawExtension, error) {
	objects := []runtime.RawExtension{}
	list, err := NewList(manifest)
	if err != nil {
		return nil, err
	}
	objects = append(objects, list.Items...)
	return objects, nil
}

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

// NewRawExtension decodes the provided manifest into runtime.RawExtension
func NewRawExtension(manifest []byte) (runtime.RawExtension, error) {
	list, err := NewList(manifest)
	if err != nil {
		return runtime.RawExtension{}, err
	}
	if len(list.Items) != 1 {
		return runtime.RawExtension{}, fmt.Errorf("expected to decode a single item; got %d items", len(list.Items))
	}
	return list.Items[0], nil
}

// AppendFlattened will append newItem to list; making sure that raw newItem is decoded
// and flattened with another list
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
