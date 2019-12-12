package manager

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// FixClusterCompatibility adds any resources missing in the CloudFormation stack in order to support new features
// like Managed Nodegroups and Fargate
func (c *StackCollection) FixClusterCompatibility() error {
	logger.Info("checking cluster stack for missing resources")
	stack, err := c.DescribeClusterStack()
	if err != nil {
		return err
	}

	var (
		clusterDefaultSG string
		fargateRole      string
	)

	featureOutputs := map[string]outputs.Collector{
		// available on clusters created after Managed Nodes support was out
		outputs.ClusterDefaultSecurityGroup: func(v string) error {
			clusterDefaultSG = v
			return nil
		},
		// available on 1.14 clusters created after Fargate support was out
		outputs.FargatePodExecutionRoleARN: func(v string) error {
			fargateRole = v
			return nil
		},
	}

	if err := outputs.Collect(*stack, nil, featureOutputs); err != nil {
		return err
	}

	stackSupportsManagedNodes := false
	if clusterDefaultSG != "" {
		stackSupportsManagedNodes, err = c.hasManagedToUnmanagedSG()
		if err != nil {
			return err
		}
	}

	stackSupportsFargate := fargateRole != ""

	if stackSupportsManagedNodes && stackSupportsFargate {
		logger.Info("cluster stack has all required resources")
		return nil
	}

	if !stackSupportsManagedNodes {
		logger.Info("cluster stack is missing resources for Managed Nodegroups")
	}
	if !stackSupportsFargate {
		logger.Info("cluster stack is missing resources for Fargate")
	}
	logger.Info("adding missing resources to cluster stack")
	_, err = c.AppendNewClusterStackResource(false, true)
	return err
}

func (c *StackCollection) hasManagedToUnmanagedSG() (bool, error) {
	stackTemplate, err := c.GetStackTemplate(c.makeClusterStackName())
	if err != nil {
		return false, errors.Wrap(err, "error getting cluster stack template")
	}
	stackResources := gjson.Get(stackTemplate, resourcesRootPath)
	return builder.HasManagedNodesSG(&stackResources), nil
}
