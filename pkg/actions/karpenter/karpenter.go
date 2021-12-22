package karpenter

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/kris-nova/logger"
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers/helm"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

// Installer contains all necessary dependencies for the Karpenter Install tasks and others.
type Installer struct {
	StackManager       manager.StackManager
	CTL                *eks.ClusterProvider
	Config             *api.ClusterConfig
	Wait               WaitFunc
	KarpenterInstaller karpenter.ChartInstaller
	ClientSet          kubernetes.Interface
	OIDC               *iamoidc.OpenIDConnectManager
}

type WaitFunc func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error

// NewInstaller creates a new Karpenter installer.
func NewInstaller(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager, clientSet kubeclient.Interface) (*Installer, error) {
	helmInstaller, err := helm.NewInstaller(helm.Options{
		Namespace: karpenter.DefaultNamespace,
	})
	if err != nil {
		return nil, err
	}
	karpenterInstaller := karpenter.NewKarpenterInstaller(karpenter.Options{
		HelmInstaller: helmInstaller,
		Namespace:     karpenter.DefaultNamespace,
		ClusterConfig: cfg,
	})
	oidc, err := ctl.NewOpenIDConnectManager(cfg)
	if err != nil {
		return nil, err
	}

	oidcProviderExists, err := oidc.CheckProviderExists()
	if err != nil {
		return nil, err
	}

	if !oidcProviderExists {
		logger.Warning("no IAM OIDC provider associated with cluster, try 'eksctl utils associate-iam-oidc-provider --region=%s --cluster=%s'", cfg.Metadata.Region, cfg.Metadata.Name)
	}

	return &Installer{
		StackManager:       stackManager,
		CTL:                ctl,
		Config:             cfg,
		Wait:               waiters.Wait,
		KarpenterInstaller: karpenterInstaller,
		ClientSet:          clientSet,
		OIDC:               oidc,
	}, nil
}

func doTasks(taskTree *tasks.TaskTree) error {
	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		logger.Info("%d error(s) occurred while installing Karpenter, you may wish to check your Cluster for further information", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to install Karpenter on cluster")
	}
	return nil
}
