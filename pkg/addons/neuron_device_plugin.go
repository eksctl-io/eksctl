package addons

import (
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
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

func watchDaemonSetReady(dsClientSet clientappsv1.DaemonSetInterface, dsName string) error {
	watcher, err := dsClientSet.Watch(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", dsName),
	})

	if err != nil {
		return err
	}

	defer watcher.Stop()
	timeout := 45 * time.Second

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return errors.New("failed waiting for daemon set: unexpected close of ResultChan")
			}
			switch event.Type {
			case watch.Added, watch.Modified:
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("timed out (after %v) waiting for certificate", timeout)
		}

	}
}

func (n *NeuronDevicePlugin) applyDeployment(manifest []byte) error {
	rawExtension, err := kubernetes.NewRawExtension(manifest)
	if err != nil {
		return err
	}
	daemonSet, ok := rawExtension.Object.(*appsv1.DaemonSet)
	if !ok {
		return &typeAssertionError{&appsv1.DaemonSet{}, rawExtension}
	}
	if err := useRegionalImage(&daemonSet.Spec.Template, n.region); err != nil {
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
	return watchDaemonSetReady(n.rawClient.ClientSet().AppsV1().DaemonSets(daemonSet.Namespace), daemonSet.Name)
}

// Deploy deploys the Neuron device plugin to the specified cluster
func (n *NeuronDevicePlugin) Deploy() error {
	return n.applyDeployment(mustGenerateAsset(neuronDevicePluginYamlBytes))
}
