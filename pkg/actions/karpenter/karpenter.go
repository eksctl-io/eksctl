package karpenter

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/aws/aws-sdk-go/aws/request"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Installer struct {
	stackManager       manager.StackManager
	ctl                *eks.ClusterProvider
	cfg                *api.ClusterConfig
	clientSet          kubernetes.Interface
	wait               WaitFunc
	kubeProvider       eks.KubeProvider
	karpenterInstaller karpenter.InstallKarpenter
}

type WaitFunc func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error

// NewInstaller creates a new Karpenter installer.
func NewInstaller(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface, karpenter karpenter.InstallKarpenter) *Installer {
	return &Installer{
		stackManager:       ctl.NewStackManager(cfg),
		ctl:                ctl,
		cfg:                cfg,
		clientSet:          clientSet,
		wait:               waiters.Wait,
		kubeProvider:       ctl,
		karpenterInstaller: karpenter,
	}
}
