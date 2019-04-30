package eks

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ValidateClusterForCompatibility looks at the cluster stack and check if it's
// compatible with current nodegroup configuration, if it find issues it returns an error
func (c *ClusterProvider) ValidateClusterForCompatibility(cfg *api.ClusterConfig, stackManager *manager.StackCollection) error {
	cluster, err := stackManager.DescribeClusterStack()
	if err != nil {
		return errors.Wrap(err, "getting cluster stacks")
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

func isNodeGroupCompatible(name string, info manager.StackInfo) bool {
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
		return false
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

	return true
}

// ValidateExistingNodeGroupsForCompatibility looks at each of the existing nodegroups and
// validates configuration, if it find issues it logs messages
func (c *ClusterProvider) ValidateExistingNodeGroupsForCompatibility(cfg *api.ClusterConfig, stackManager *manager.StackCollection) error {
	infoByNodeGroup, err := stackManager.DescribeNodeGroupStacksAndResources()
	if err != nil {
		return errors.Wrap(err, "getting resources for of all nodegroup stacks")
	}
	if len(infoByNodeGroup) == 0 {
		return nil
	}

	logger.Info("checking security group configuration for all nodegroups")
	incompatibleNodeGroups := []string{}
	for ng, info := range infoByNodeGroup {
		if stackManager.StackStatusIsNotTransitional(info.Stack) {
			if isNodeGroupCompatible(ng, info) {
				logger.Debug("nodegroup %q is compatible", ng)
			} else {
				logger.Debug("nodegroup %q is incompatible", ng)
				incompatibleNodeGroups = append(incompatibleNodeGroups, ng)
			}
		}
	}

	numIncompatibleNodeGroups := len(incompatibleNodeGroups)
	if numIncompatibleNodeGroups == 0 {
		logger.Info("all nodegroups have up-to-date configuration")
		return nil
	}

	logger.Critical("found %d nodegroup(s) (%s) without shared security group, cluster networking maybe be broken",
		numIncompatibleNodeGroups, strings.Join(incompatibleNodeGroups, ", "))
	logger.Critical("it's recommended to create new nodegroups, then delete old ones")
	if cfg.VPC.SharedNodeSecurityGroup != "" {
		logger.Critical("as a temporary fix, you can patch the configuration and add each of these nodegroup(s) to %q",
			cfg.VPC.SharedNodeSecurityGroup)
	}

	return nil
}
