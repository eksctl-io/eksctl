package defaultaddons

import (
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

// LoadAsset return embedded manifest as a runtime.Object
func LoadAsset(name string) (runtime.Object, error) {
	data, err := Asset(name + ".yaml")
	if err != nil {
		return nil, errors.Wrapf(err, "decoding embedded manifest for %q", name)
	}

	obj, err := runtime.Decode(scheme.Codecs.UniversalDeserializer(), data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading embedded manifest for %q", name)
	}

	return obj, nil
}
