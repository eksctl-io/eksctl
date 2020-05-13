package addons

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

// NewNeuronDevicePlugin creates a new NeuronDevicePlugin
func NewNeuronDevicePlugin(rawClient kubernetes.RawClientInterface, region string, planMode bool) *NeuronDevicePlugin {
	return &NeuronDevicePlugin{
		rawClient,
		region,
		planMode,
	}
}

// A NeuronDevicePlugin deploys the Neuron Device Plugin to a cluster
type NeuronDevicePlugin struct {
	rawClient kubernetes.RawClientInterface
	region    string
	planMode  bool
}

// Currently only us-east-1 and us-west-2 are available
// and both use this AWS account
const neuronResourceAccount = "790709498068"

// useRegionalImage is specific to the neuron device plugin
// we assume that only us-east-1 and us-west-2 support inferentia for now
func useRegionalImage(spec *v1.PodTemplateSpec, region string) error {
	imageFormat := spec.Spec.Containers[0].Image
	dnsSuffix, err := awsDNSSuffixForRegion(region)
	if err != nil {
		return err
	}
	regionalImage := fmt.Sprintf(imageFormat, neuronResourceAccount, region, dnsSuffix)
	spec.Spec.Containers[0].Image = regionalImage
	return nil
}

func (n *NeuronDevicePlugin) applyDeployment(manifest []byte) error {
	rawExtension, err := kubernetes.NewRawExtension(manifest)
	if err != nil {
		return err
	}
	deployment, ok := rawExtension.Object.(*appsv1.DaemonSet)
	if !ok {
		return &typeAssertionError{&appsv1.DaemonSet{}, rawExtension}
	}
	if err := useRegionalImage(&deployment.Spec.Template, n.region); err != nil {
		return err
	}

	rawResource, err := n.rawClient.NewRawResource(rawExtension.Object)
	if err != nil {
		return err
	}
	msg, err := rawResource.CreateOrReplace(n.planMode)
	if err != nil {
		return err
	}
	logger.Info(msg)
	return nil
}

// Deploy deploys the Neuron device plugin to the specified cluster
func (n *NeuronDevicePlugin) Deploy() error {
	return n.applyDeployment(mustGenerateAsset(neuronDevicePluginYamlBytes))
}
