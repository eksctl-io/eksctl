package addons

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/credentials"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// SpotOceanController deploys Spot Ocean controller to a cluster.
type SpotOceanController struct {
	rawClient   kubernetes.RawClientInterface
	clusterSpec *api.ClusterConfig
	plan        bool
	profile     string
}

// NewSpotOceanController creates a new Spot Ocean controller.
func NewSpotOceanController(rawClient kubernetes.RawClientInterface,
	clusterSpec *api.ClusterConfig, plan bool, profile string) *SpotOceanController {

	return &SpotOceanController{
		rawClient:   rawClient,
		clusterSpec: clusterSpec,
		plan:        plan,
		profile:     profile,
	}
}

// Deploy deploys the Spot Ocean controller.
func (x *SpotOceanController) Deploy() (err error) {
	logger.Debug("ocean: installing controller")

	defer func() {
		if r := recover(); r != nil {
			if ae, ok := r.(*assetError); ok {
				err = ae
			} else {
				panic(r)
			}
		}
	}()

	// Deploy a ConfigMap to store the controller configuration.
	if err := x.applyConfigMap(); err != nil {
		return fmt.Errorf("error creating configmap: %w", err)
	}

	// Deploy the controller and its resources (RBAC, etc.).
	if err := x.applyResources(mustGenerateAsset(spotOceanControllerYamlBytes)); err != nil {
		return fmt.Errorf("error creating resources: %w", err)
	}

	return nil
}

func (x *SpotOceanController) applyResources(manifests []byte) error {
	list, err := kubernetes.NewList(manifests)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		if err := x.applyRawResource(item.Object); err != nil {
			return err
		}
	}

	return nil
}

func (x *SpotOceanController) applyRawResource(object runtime.Object) error {
	rawResource, err := x.rawClient.NewRawResource(object)
	if err != nil {
		return err
	}

	msg, err := rawResource.CreateOrReplace(x.plan)
	if err != nil {
		return err
	}

	logger.Debug(msg)
	return nil
}

func (x *SpotOceanController) applyConfigMap() error {
	config := spotinst.DefaultConfig()
	config.WithCredentials(credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.FileProvider{Profile: x.profile},
	}...))

	c, err := config.Credentials.Get()
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spotinst-kubernetes-cluster-controller-config",
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"spotinst.token":              c.Token,
			"spotinst.account":            c.Account,
			"spotinst.cluster-identifier": x.clusterSpec.Metadata.Name,
		},
	}

	return x.applyRawResource(cm)
}
