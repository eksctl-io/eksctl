package cluster

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/logger"
	"time"
)

type Cluster interface {
	Upgrade(dryRun bool) error
	Delete(waitInterval time.Duration, wait bool) error
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) (Cluster, error) {
	err := ctl.RefreshClusterStatus(cfg)

	if err != nil  {
		awsError, ok :=  err.(awserr.Error)
		//if the cluster doesn't exist it might still have stacks to delete
		if !ok || awsError.Code() != awseks.ErrCodeResourceNotFoundException {
			return nil, err
		}
	}

	stackManager := ctl.NewStackManager(cfg)
	hasClusterStack, err := stackManager.HasClusterStack()
	if err != nil {
		return nil, err
	}

	if hasClusterStack {
		logger.Debug("Cluster %q was created by eksctl", cfg.Metadata.Name)
		return NewOwnedCluster(cfg, ctl, stackManager)
	}
	logger.Debug("Cluster %q was not created by eksctl", cfg.Metadata.Name)

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return nil, err
	}
	return NewUnownedCluster(cfg, ctl, clientSet, stackManager), nil
}
