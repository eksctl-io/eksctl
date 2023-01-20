package cmdutils

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

// KubernetesClientAndConfigFrom returns a Kubernetes client set and REST
// configuration object for the currently configured cluster.
func KubernetesClientAndConfigFrom(cmd *Cmd) (kubernetes.Interface, error) {
	ctl, err := cmd.NewProviderForExistingCluster(context.TODO())
	if err != nil {
		return nil, err
	}
	cfg := cmd.ClusterConfig
	if ok, err := ctl.CanOperate(cfg); !ok {
		return nil, err
	}
	k8sClientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return nil, err
	}
	return k8sClientSet, nil
}
