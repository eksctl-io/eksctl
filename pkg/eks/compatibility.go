package eks

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
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

	sharedClusterNodeSG := ""
	for _, x := range cluster.Outputs {
		if *x.OutputKey == outputs.ClusterSharedNodeSecurityGroup {
			sharedClusterNodeSG = *x.OutputValue
		}
	}

	if sharedClusterNodeSG == "" {
		return fmt.Errorf(
			"shared node security group missing, to fix this run 'eksctl utils update-cluster-stack --name=%s --region=%s'",
			cfg.Metadata.Name,
			cfg.Metadata.Region,
		)
	}

	return nil
}

// ValidateExistingNodeGroupsForCompatibility looks at each of the existing nodegroups and
// validates configuration, if it find issues it logs messages
func (c *ClusterProvider) ValidateExistingNodeGroupsForCompatibility(cfg *api.ClusterConfig, stackManager *manager.StackCollection) error {
	resourcesByNodeGroup, err := stackManager.DescribeNodeGroupStacksAndResources()
	if err != nil {
		return errors.Wrap(err, "getting resources for of all nodegroup stacks")
	}
	if len(resourcesByNodeGroup) == 0 {
		return nil
	}

	logger.Info("checking security group configuration for all nodegroups")
	incompatibleNodeGroups := []string{}
	for ng, resources := range resourcesByNodeGroup {
		compatible := false
		for _, x := range resources.Stack.Outputs {
			if *x.OutputKey == outputs.NodeGroupFeatureSharedSecurityGroup {
				compatible = true
			}
		}
		if !compatible {
			incompatibleNodeGroups = append(incompatibleNodeGroups, ng)
		}
	}

	numIncompatibleNodeGroups := len(incompatibleNodeGroups)
	if numIncompatibleNodeGroups == 0 {
		logger.Info("all security group nodegroups have up-to-date configuration")
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
