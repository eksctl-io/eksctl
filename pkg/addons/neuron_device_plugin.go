package addons

import (
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
)

// NewNeuronDevicePlugin creates a new NeuronDevicePlugin
func NewNeuronDevicePlugin(rawClient kubernetes.RawClientInterface) *NeuronDevicePlugin {
	return &NeuronDevicePlugin{
		rawClient,
	}
}

// A NeuronDevicePlugin deploys the Neuron Device Plugin to a cluster
type NeuronDevicePlugin struct {
	rawClient kubernetes.RawClientInterface
}

func (n *NeuronDevicePlugin) applyRawResource(object runtime.Object) error {
	rawResource, err := n.rawClient.NewRawResource(object)
	if err != nil {
		return err
	}

	msg, err := rawResource.CreateOrReplace(false)
	if err != nil {
		return err
	}
	logger.Info(msg)
	return nil
}

func (n *NeuronDevicePlugin) applyResources(manifests []byte) error {
	list, err := kubernetes.NewList(manifests)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		if err := n.applyRawResource(item.Object); err != nil {
			return err
		}
	}
	return nil
}

// Deploy deploys the Neuron device plugin to the specified cluster
func (n *NeuronDevicePlugin) Deploy() error {
	return n.applyResources(mustGenerateAsset(neuronDevicePluginYamlBytes))
}
