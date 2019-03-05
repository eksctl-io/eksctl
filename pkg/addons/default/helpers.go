package defaultaddons

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

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func init() {
	apiextensionsv1beta1.AddToScheme(scheme.Scheme)
}

// LoadAsset return embedded manifest as a runtime.Object
// TODO: we certainly need tests for this
func LoadAsset(name, ext string) (*metav1.List, error) {
	data, err := Asset(name + "." + ext)
	if err != nil {
		return nil, errors.Wrapf(err, "decoding embedded manifest for %q", name)
	}

	list := metav1.List{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer(data), 4096)

	for {
		obj := new(runtime.RawExtension)
		err := decoder.Decode(obj)
		if err != nil {
			if err == io.EOF {
				return &list, nil
			}
			return nil, errors.Wrapf(err, "loading individual resources from manifest for %q", name)
		}
		// obj.Object, err = runtime.Decode(scheme.Codecs.UniversalDeserializer(), obj.Raw)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "converting object")
		// }
		// list.Items = append(list.Items, *obj)
		if err := listAppendFlattened(&list, *obj); err != nil {
			return nil, err
		}
	}
}

// this was copied from kubegen
func listAppendFlattened(components *metav1.List, component runtime.RawExtension) error {
	if component.Object != nil {
		if strings.HasSuffix(component.Object.GetObjectKind().GroupVersionKind().Kind, "List") {
			// must use corev1, as it panics on obj.(*metav1.List) with
			// an amusing error message saying that *v1.List is not *v1.List
			// TODO: test this to find out if the conversion is still required
			list := component.Object.(*corev1.List)
			for _, item := range (*list).Items {
				// we attempt to recurse here, but most likely
				// we will have to try decoding component.Raw
				if err := listAppendFlattened(components, item); err != nil {
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
		return errors.Wrapf(err, "converting object")
	}
	return listAppendFlattened(components, runtime.RawExtension{Object: obj})
}
