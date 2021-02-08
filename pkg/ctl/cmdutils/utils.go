package cmdutils

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesClientAndConfigFrom returns a Kubernetes client set and REST
// configuration object for the currently configured cluster.
func KubernetesClientAndConfigFrom(cmd *Cmd) (*kubernetes.Clientset, *rest.Config, error) {
	ctl, err := cmd.NewCtl()
	if err != nil {
		return nil, nil, err
	}
	cfg := cmd.ClusterConfig
	if ok, err := ctl.CanOperate(cfg); !ok {
		return nil, nil, err
	}
	kubernetesClientConfigs, err := ctl.NewClient(cfg)
	if err != nil {
		return nil, nil, err
	}
	k8sConfig := kubernetesClientConfigs.Config
	k8sRestConfig, err := clientcmd.NewDefaultClientConfig(*k8sConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot create Kubernetes client configuration")
	}
	k8sClientSet, err := kubernetes.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, nil, errors.Errorf("cannot create Kubernetes client set: %s", err)
	}
	return k8sClientSet, k8sRestConfig, nil
}
