package addons

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	vpcControllerNamespace = metav1.NamespaceSystem
)

// VPCController deletes an existing installation of VPC controller from worker nodes.
type VPCController struct {
	RawClient kubernetes.RawClientInterface
}

// Delete deletes the resources for VPC controller.
func (v *VPCController) Delete() error {
	vpcControllerMetadata, err := vpcControllerMetadataYamlBytes()
	if err != nil {
		return errors.Wrap(err, "unexpected error loading manifests")
	}
	list, err := kubernetes.NewList(vpcControllerMetadata)
	if err != nil {
		return errors.Wrap(err, "unexpected error parsing manifests")
	}
	for _, item := range list.Items {
		if err := v.deleteResource(item.Object); err != nil {
			return err
		}
	}
	return nil
}

func (v *VPCController) deleteResource(o runtime.Object) error {
	r, err := v.RawClient.NewRawResource(o)
	if err != nil {
		return errors.Wrap(err, "unexpected error creating raw resource")
	}
	msg, err := r.DeleteSync()
	if err != nil {
		return errors.Wrapf(err, "error deleting resource %q", r.Info.String())
	}
	logger.Info(msg)
	return nil
}
