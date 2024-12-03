package windows

import (
	"context"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/kris-nova/logger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	vpcCNIName       = "amazon-vpc-cni"
	vpcCNINamespace  = metav1.NamespaceSystem
	windowsIPAMField = "enable-windows-ipam"
)

// IPAM enables Windows IPAM in the VPC CNI ConfigMap.
type IPAM struct {
	Clientset kubernetes.Interface
}

// Enable enables Windows IPAM in the VPC CNI ConfigMap.
func (w *IPAM) Enable(ctx context.Context) error {
	configMaps := w.Clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem)
	vpcCNIConfig, err := configMaps.Get(ctx, vpcCNIName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("error getting ConfigMap %q: %w", vpcCNIName, err)
		}
		return createConfigMap(ctx, configMaps)
	}

	if val, ok := vpcCNIConfig.Data[windowsIPAMField]; ok && val == "true" {
		logger.Info("Windows IPAM is already enabled")
		return nil
	}

	patch, err := createPatch(vpcCNIConfig)
	if err != nil {
		return fmt.Errorf("error creating merge patch: %w", err)
	}

	_, err = configMaps.Patch(ctx, vpcCNIName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch resource %q: %w", vpcCNIName, err)
	}
	return nil
}

func createPatch(cm *corev1.ConfigMap) ([]byte, error) {
	oldData, err := json.Marshal(cm)
	if err != nil {
		return nil, err
	}
	cm.Data[windowsIPAMField] = "true"
	modifiedData, err := json.Marshal(cm)
	if err != nil {
		return nil, err
	}
	return jsonpatch.CreateMergePatch(oldData, modifiedData)
}

func createConfigMap(ctx context.Context, configMaps corev1client.ConfigMapInterface) error {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vpcCNIName,
			Namespace: vpcCNINamespace,
		},
		Data: map[string]string{
			windowsIPAMField: "true",
		},
	}
	_, err := configMaps.Create(ctx, cm, metav1.CreateOptions{})
	return err
}
