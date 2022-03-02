package addons

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/assetutil"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	
	// For go:embed
	_ "embed"
)

//go:embed assets/spot-ocean-controller.yaml
var spotOceanControllerYamlBytes []byte

// SpotOceanController deploys Spot Ocean controller to a cluster.
type SpotOceanController struct {
	rawClient   kubernetes.RawClientInterface
	clusterSpec *api.ClusterConfig
	plan        bool
}

// NewSpotOceanController creates a new Spot Ocean controller.
func NewSpotOceanController(rawClient kubernetes.RawClientInterface,
	clusterSpec *api.ClusterConfig, plan bool) *SpotOceanController {
	return &SpotOceanController{
		rawClient:   rawClient,
		clusterSpec: clusterSpec,
		plan:        plan,
	}
}

// Deploy deploys the Spot Ocean controller.
func (x *SpotOceanController) Deploy() (err error) {
	logger.Debug("ocean: installing controller")

	defer func() {
		if r := recover(); r != nil {
			if ae, ok := r.(*assetutil.Error); ok {
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

	// Deploy a Secret to store the controller credentials.
	if err := x.applySecret(); err != nil {
		return fmt.Errorf("error creating secret: %w", err)
	}

	// Deploy the controller and its resources (RBAC, etc.).
	if err := x.applyResources(spotOceanControllerYamlBytes); err != nil {
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
	o := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spotinst-kubernetes-cluster-controller-config",
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"spotinst.cluster-identifier": x.clusterSpec.Metadata.Name,
		},
	}

	return x.applyRawResource(o)
}

func (x *SpotOceanController) applySecret() error {
	config := spotinst.DefaultConfig()
	c, err := config.Credentials.Get()
	if err != nil {
		return err
	}

	o := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spotinst-kubernetes-cluster-controller",
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"token":   []byte(c.Token),
			"account": []byte(c.Account),
		},
	}

	return x.applyRawResource(o)
}
