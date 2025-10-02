package efa

import (
	"fmt"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfnt "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
	"github.com/weaveworks/eksctl/pkg/utils/version"
)

// IsBuiltInSupported returns true if the Kubernetes version supports built-in EFA in the default security group
func IsBuiltInSupported(kubernetesVersion string) (bool, error) {
	supported, err := version.IsMinVersion(api.EFABuiltInSupportVersion, kubernetesVersion)
	if err != nil {
		return false, fmt.Errorf("failed to determine EFA built-in support for Kubernetes version %q (minimum required: %s): %w",
			kubernetesVersion, api.EFABuiltInSupportVersion, err)
	}
	return supported, nil
}

// SecurityGroupConfig holds the configuration needed for EFA security group creation
type SecurityGroupConfig struct {
	ClusterVersion string
	ClusterName    string
	NodeGroupName  string
	VPCID          *gfnt.Value
	Description    string
}

// ProcessSecurityGroup handles the common EFA security group logic
// Returns the security group (nil if built-in EFA is used) and any error
func ProcessSecurityGroup(config SecurityGroupConfig, addEFASecurityGroupFunc func(*gfnt.Value, string, string) *gfnt.Value) (*gfnt.Value, error) {
	supported, err := IsBuiltInSupported(config.ClusterVersion)
	if err != nil {
		logger.Warning("failed to parse Kubernetes version %s for EFA configuration: %v; falling back to custom EFA security group creation", config.ClusterVersion, err)
		// Fall back to creating custom EFA security group when version parsing fails
		supported = false
	}

	if !supported {
		logger.Info("creating custom EFA security group for nodegroup %s with Kubernetes %s (EFA built-in support requires version 1.33+)", config.NodeGroupName, config.ClusterVersion)
		efaSG := addEFASecurityGroupFunc(config.VPCID, config.ClusterName, config.Description)
		if efaSG == nil {
			return nil, fmt.Errorf("failed to create EFA security group for nodegroup %s with Kubernetes %s: invalid VPC ID or cluster configuration", config.NodeGroupName, config.ClusterVersion)
		}
		logger.Info("successfully created custom EFA security group for nodegroup %s", config.NodeGroupName)
		return efaSG, nil
	}

	logger.Info("using built-in EFA support in default security group for nodegroup %s with Kubernetes %s (no custom EFA security group needed)", config.NodeGroupName, config.ClusterVersion)
	return nil, nil
}
