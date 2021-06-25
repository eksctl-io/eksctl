package kubernetes

import (
	"context"
	"encoding/json"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const kmsAnnotation = "eksctl.io/kms-encryption-timestamp"

// RefreshSecrets updates all secrets to apply KMS encryption
func RefreshSecrets(ctx context.Context, c v1.CoreV1Interface) error {
	var cont string
	for {
		list, err := c.Secrets(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
			Continue: cont,
		})
		if err != nil {
			return errors.Wrap(err, "error listing resources")
		}
		for _, secret := range list.Items {
			if err := refreshSecret(ctx, c, secret); err != nil {
				return errors.Wrapf(err, "error updating secret %q", secret.Name)
			}
		}
		if cont = list.Continue; cont == "" {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

		}
	}
	return nil
}
func createPatch(o runtime.Object, annotationName string) ([]byte, error) {
	metaAccessor := meta.NewAccessor()
	oldData, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	annotations, err := metaAccessor.Annotations(o)
	if err != nil {
		return nil, err
	}
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[annotationName] = time.Now().Format(time.RFC3339)
	if err := metaAccessor.SetAnnotations(o, annotations); err != nil {
		return nil, err
	}
	modifiedData, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return jsonpatch.CreateMergePatch(oldData, modifiedData)
}

func refreshSecret(ctx context.Context, c v1.CoreV1Interface, s corev1.Secret) error {
	patch, err := createPatch(&s, kmsAnnotation)
	if err != nil {
		return errors.Wrap(err, "unexpected error creating JSON patch")
	}
	if _, err := c.Secrets(s.Namespace).Patch(ctx, s.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{}); err != nil {
		return errors.Wrap(err, "error updating secret")
	}
	return nil
}
