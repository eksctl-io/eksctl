package cluster

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/logger"
)

type Cluster interface {
	Upgrade(dryRun bool) error
	Delete(waitInterval time.Duration, wait bool) error
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) (Cluster, error) {
	clusterExists := true
	err := ctl.RefreshClusterStatus(cfg)
	if err != nil {
		if awsError, ok := errors.Unwrap(errors.Unwrap(err)).(awserr.Error); ok &&
			awsError.Code() == awseks.ErrCodeResourceNotFoundException {
			clusterExists = false
		} else {
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

	if !clusterExists {
		return nil, fmt.Errorf("cluster %q does not exist", cfg.Metadata.Name)
	}

	logger.Debug("cluster %q was not created by eksctl", cfg.Metadata.Name)

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return nil, err
	}
	return NewUnownedCluster(cfg, ctl, clientSet, stackManager), nil
}
