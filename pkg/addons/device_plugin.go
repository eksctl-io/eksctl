package addons

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/assetutil"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/logger"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

const (
	// Currently only us-east-1 and us-west-2 are available
	// and both use this AWS account
	neuronResourceAccount = "790709498068"
)

// useRegionalImage is specific to the neuron device plugin
// we assume that only us-east-1 and us-west-2 support inferentia for now
func useRegionalImage(spec *v1.PodTemplateSpec, region string, account string) error {
	imageFormat := spec.Spec.Containers[0].Image
	dnsSuffix, err := awsDNSSuffixForRegion(region)
	if err != nil {
		return err
	}
	regionalImage := fmt.Sprintf(imageFormat, account, region, dnsSuffix)
	spec.Spec.Containers[0].Image = regionalImage
	return nil
}

func watchDaemonSetReady(dsClientSet clientappsv1.DaemonSetInterface, dsName string) error {
	watcher, err := dsClientSet.Watch(context.TODO(), metav1.ListOptions{
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

type MkDevicePlugin func(rawClient kubernetes.RawClientInterface, region string, planMode bool) DevicePlugin

type DevicePlugin interface {
	RawClient() kubernetes.RawClientInterface
	PlanMode() bool
	Manifest() []byte
	SetImage(t *v1.PodTemplateSpec) error
	Deploy() error
}

func applyDevicePlugin(dp DevicePlugin) error {
	list, err := kubernetes.NewList(dp.Manifest())
	if err != nil {
		return errors.Wrap(err, "creating list from device plugin manifest")
	}

	rawClient := dp.RawClient()
	for _, rawObj := range list.Items {
		rawResource, err := rawClient.NewRawResource(rawObj.Object)
		if err != nil {
			return errors.Wrap(err, "creating raw resource from list item")
		}
		switch rawResource.GVK.Kind {
		case "DaemonSet":
			daemonSet, ok := rawResource.Info.Object.(*appsv1.DaemonSet)
			if !ok {
				return &typeAssertionError{&appsv1.DaemonSet{}, rawResource}
			}
			if err := dp.SetImage(&daemonSet.Spec.Template); err != nil {
				return errors.Wrap(err, "setting image of device plugin daemonset")
			}

			msg, err := rawResource.CreateOrReplace(dp.PlanMode())
			if err != nil {
				return errors.Wrap(err, "calling create or replace on raw device plugin daemonset")
			}
			logger.Info(msg)
			if err := watchDaemonSetReady(dp.RawClient().ClientSet().AppsV1().DaemonSets(daemonSet.Namespace), daemonSet.Name); err != nil {
				return errors.Wrap(err, "waiting for device plugin daemonset to become ready")
			}
		default:
			status, err := rawResource.CreateOrReplace(dp.PlanMode())
			if err != nil {
				return errors.Wrap(err, "calling create or replace on raw device plugin rawResource")
			}
			logger.Info(status)
		}
	}
	return nil
}

// NewNeuronDevicePlugin creates a new NeuronDevicePlugin
func NewNeuronDevicePlugin(rawClient kubernetes.RawClientInterface, region string, planMode bool) DevicePlugin {
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

func (n *NeuronDevicePlugin) RawClient() kubernetes.RawClientInterface {
	return n.rawClient
}

func (n *NeuronDevicePlugin) PlanMode() bool {
	return n.planMode
}

func (n *NeuronDevicePlugin) Manifest() []byte {
	return assetutil.MustLoad(neuronDevicePluginYamlBytes)
}

func (n *NeuronDevicePlugin) SetImage(t *v1.PodTemplateSpec) error {
	return useRegionalImage(t, n.region, neuronResourceAccount)
}

// Deploy deploys the Neuron device plugin to the specified cluster
func (n *NeuronDevicePlugin) Deploy() error {
	return applyDevicePlugin(n)
}

// NewNvidiaDevicePlugin creates a new NvidiaDevicePlugin
func NewNvidiaDevicePlugin(rawClient kubernetes.RawClientInterface, region string, planMode bool) DevicePlugin {
	return &NvidiaDevicePlugin{
		rawClient,
		region,
		planMode,
	}
}

// A NvidiaDevicePlugin deploys the Nvidia Device Plugin to a cluster
type NvidiaDevicePlugin struct {
	rawClient kubernetes.RawClientInterface
	region    string
	planMode  bool
}

func (n *NvidiaDevicePlugin) RawClient() kubernetes.RawClientInterface {
	return n.rawClient
}

func (n *NvidiaDevicePlugin) PlanMode() bool {
	return n.planMode
}

func (n *NvidiaDevicePlugin) SetImage(t *v1.PodTemplateSpec) error {
	return nil
}

func (n *NvidiaDevicePlugin) Manifest() []byte {
	return assetutil.MustLoad(nvidiaDevicePluginYamlBytes)
}

// Deploy deploys the Nvidia device plugin to the specified cluster
func (n *NvidiaDevicePlugin) Deploy() error {
	return applyDevicePlugin(n)
}
