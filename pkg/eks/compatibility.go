package eks

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ValidateClusterForCompatibility looks at the cluster stack and check if it's
// compatible with current nodegroup configuration, if it find issues it returns an error
func (c *ClusterProvider) ValidateClusterForCompatibility(ctx context.Context, cfg *api.ClusterConfig, stackManager manager.StackManager) error {
	cluster, err := stackManager.DescribeClusterStack(ctx)
	if err != nil {
		return fmt.Errorf("getting cluster stack: %w", err)
	}

	err = outputs.Collect(*cluster,
		map[string]outputs.Collector{
			outputs.ClusterSharedNodeSecurityGroup: func(v string) error {
				logger.Debug("ClusterSharedNodeSecurityGroup = %s", v)
				return nil
			},
		},
		nil,
	)

	if err != nil {
		logger.Debug("err = %s", err.Error())
		return fmt.Errorf(
			"shared node security group missing, to fix this run 'eksctl update cluster --name=%s --region=%s'",
			cfg.Metadata.Name,
			cfg.Metadata.Region,
		)
	}

	return nil
}

func isNodeGroupCompatible(name string, info manager.StackInfo) (bool, error) {
	nodeGroupType, err := manager.GetNodeGroupType(info.Stack.Tags)
	if err != nil {
		return false, err
	}
	if nodeGroupType == api.NodeGroupTypeManaged {
		// Managed Nodegroups have a cluster security group attached by default
		return true, nil
	}

	hasSharedSecurityGroupFlag := false
	usesSharedSecurityGroup := false
	hasLocalSecurityGroupFlag := false

	_ = outputs.Collect(*info.Stack,
		nil,
		map[string]outputs.Collector{
			outputs.NodeGroupFeatureSharedSecurityGroup: func(v string) error {
				hasSharedSecurityGroupFlag = true
				switch v {
				case "true":
					usesSharedSecurityGroup = true
				case "false":
					usesSharedSecurityGroup = false
				}
				return nil
			},
			outputs.NodeGroupFeatureLocalSecurityGroup: func(v string) error {
				hasLocalSecurityGroupFlag = true
				return nil
			},
		},
	)

	if !hasSharedSecurityGroupFlag {
		// if it doesn't have `outputs.NodeGroupFeatureSharedSecurityGroup` flags at all,
		// it must be incompatible
		return false, nil
	}

	if !hasLocalSecurityGroupFlag && !usesSharedSecurityGroup {
		// when `outputs.NodeGroupFeatureSharedSecurityGroup` was added in 0.1.19, v1alpha3 didn't set
		// `ng.SharedSecurityGroup=true` by default, and (technically) it implies the nodegroup maybe compatible,
		// however users are unaware of that API v1alpha3 was broken this way, so we need this warning;
		// as `outputs.NodeGroupFeatureLocalSecurityGroup` was added in 0.1.20/v1alpha4, it can be used to
		// infer use of v1alpha3 and thereby warn the user that their cluster may had been misconfigured
		logger.Warning("looks like nodegroup %q was created using v1alpha3 API and is not using shared SG", name)
		logger.Warning("if you didn't disable shared SG and expect that pods running on %q should have access to all other pods", name)
		logger.Warning("then you should replace nodegroup %q or patch the configuration", name)
	}

	return true, nil
}
