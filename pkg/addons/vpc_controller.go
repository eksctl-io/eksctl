package addons

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
)

const (
	VPCControllerInfoMessage = "you no longer need to install the VPC resource controller on Linux worker nodes to run " +
		"Windows workloads in EKS clusters. You can enable Windows IP address management on the EKS control plane via " +
		"a ConﬁgMap setting (see https://todo.com for details). eksctl will automatically patch the ConfigMap to enable " +
		"Windows IP address management when a Windows nodegroup is created. For existing clusters, you can enable it manually " +
		"and run `eksctl utils install-vpc-controllers` with the --delete ﬂag to remove the worker node installation of the VPC resource controller"
)

// VPCController deletes an existing installation of VPC controller from worker nodes.
type VPCController struct {
	RawClient *kubernetes.RawClient
	PlanMode  bool
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
	msg, err := r.DeleteSync(v.PlanMode)
	if err != nil {
		return errors.Wrapf(err, "error deleting resource %q", r.Info.String())
	}
	if msg != "" {
		logger.Info(msg)
	}
	return nil
}
