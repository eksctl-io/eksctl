package cluster

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Cluster interface {
	Upgrade(dryRun bool) error
	Delete(waitInterval time.Duration, wait, force, disableNodegroupEviction bool) error
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) (Cluster, error) {
	clusterExists := true
	if err := ctl.RefreshClusterStatusIfStale(cfg); err != nil {
		if awsError, ok := errors.Unwrap(errors.Unwrap(err)).(awserr.Error); ok &&
			awsError.Code() == awseks.ErrCodeResourceNotFoundException {
			clusterExists = false
		} else {
			return nil, err
		}
	}

	stackManager := ctl.NewStackManager(cfg)
	clusterStack, err := stackManager.GetClusterStackIfExists()
	if err != nil {
		return nil, err
	}

	if clusterStack != nil {
		logger.Debug("cluster %q was created by eksctl", cfg.Metadata.Name)
		return NewOwnedCluster(cfg, ctl, clusterStack, stackManager), nil
	}

	if !clusterExists {
		return nil, fmt.Errorf("cluster %q does not exist", cfg.Metadata.Name)
	}

	logger.Debug("cluster %q was not created by eksctl", cfg.Metadata.Name)

	return NewUnownedCluster(cfg, ctl, stackManager), nil
}
