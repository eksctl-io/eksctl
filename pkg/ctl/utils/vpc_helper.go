package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kris-nova/logger"

	"golang.org/x/exp/slices"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// A VPCConfigUpdater updates a cluster's VPC config.
type VPCConfigUpdater interface {
	UpdateClusterConfig(ctx context.Context, input *eks.UpdateClusterConfigInput) error
}

// A VPCHelper is a helper for updating a cluster's VPC config.
type VPCHelper struct {
	// VPCUpdater updates the VPC config.
	VPCUpdater VPCConfigUpdater
	// ClusterMeta holds the cluster metadata.
	ClusterMeta *api.ClusterMeta
	// Cluster holds the current cluster state.
	Cluster *ekstypes.Cluster
	// PlanMode configures the plan mode.
	PlanMode bool
}

// UpdateClusterVPCConfig updates the cluster endpoints and public access CIDRs.
func (v *VPCHelper) UpdateClusterVPCConfig(ctx context.Context, vpc *api.ClusterVPC) error {
	if v.Cluster.OutpostConfig != nil {
		return api.ErrUnsupportedLocalCluster
	}
	if vpc.ClusterEndpoints != nil {
		if err := v.updateEndpointAccess(ctx, *vpc.ClusterEndpoints); err != nil {
			return err
		}
	}
	if vpc.PublicAccessCIDRs != nil {
		if err := v.updatePublicAccessCIDRs(ctx, vpc); err != nil {
			return err
		}
	}
	if vpc.ControlPlaneSubnetIDs != nil || vpc.ControlPlaneSecurityGroupIDs != nil {
		if err := v.updateSubnetsSecurityGroups(ctx, vpc); err != nil {
			return err
		}
	}
	cmdutils.LogPlanModeWarning(v.PlanMode)
	return nil
}

func (v *VPCHelper) updateSubnetsSecurityGroups(ctx context.Context, vpc *api.ClusterVPC) error {
	current := v.Cluster.ResourcesVpcConfig
	hasUpdate := false
	vpcUpdate := &ekstypes.VpcConfigRequest{
		SubnetIds:        current.SubnetIds,
		SecurityGroupIds: current.SecurityGroupIds,
	}

	compareValues := func(currentValues, newValues []string, resourceName string, updateFn func()) {
		if !slices.Equal(currentValues, newValues) {
			updateFn()
			hasUpdate = true
			cmdutils.LogIntendedAction(v.PlanMode, "update %s for cluster %q in %q to: %v", resourceName,
				v.ClusterMeta.Name, v.ClusterMeta.Region, newValues)
		} else {
			logger.Success("%s for cluster %q in %q are already up-to-date", resourceName, v.ClusterMeta.Name, v.ClusterMeta.Region)
		}
	}
	if vpc.ControlPlaneSubnetIDs != nil {
		compareValues(current.SubnetIds, vpc.ControlPlaneSubnetIDs, "control plane subnet IDs", func() {
			vpcUpdate.SubnetIds = vpc.ControlPlaneSubnetIDs
		})
	}

	if vpc.ControlPlaneSecurityGroupIDs != nil {
		compareValues(current.SecurityGroupIds, vpc.ControlPlaneSecurityGroupIDs, "control plane security group IDs", func() {
			vpcUpdate.SecurityGroupIds = vpc.ControlPlaneSecurityGroupIDs
		})
	}

	if v.PlanMode || !hasUpdate {
		return nil
	}
	if err := v.updateVPCConfig(ctx, vpcUpdate); err != nil {
		return err
	}
	cmdutils.LogCompletedAction(false, "control plane subnets and security groups for cluster %q in %q have been updated to: "+
		"controlPlaneSubnetIDs=%v, controlPlaneSecurityGroupIDs=%v", v.ClusterMeta.Name, v.ClusterMeta.Region, vpcUpdate.SubnetIds, vpcUpdate.SecurityGroupIds)

	return nil
}

func (v *VPCHelper) updateEndpointAccess(ctx context.Context, desired api.ClusterEndpoints) error {
	current := v.Cluster.ResourcesVpcConfig
	if desired.PublicAccess == nil {
		desired.PublicAccess = aws.Bool(current.EndpointPublicAccess)
	}
	if desired.PrivateAccess == nil {
		desired.PrivateAccess = aws.Bool(current.EndpointPrivateAccess)
	}
	if *desired.PublicAccess == current.EndpointPublicAccess && *desired.PrivateAccess == current.EndpointPrivateAccess {
		logger.Success("Kubernetes API endpoint access for cluster %q in %q is already up-to-date",
			v.ClusterMeta.Name, v.ClusterMeta.Region)
		return nil
	}

	cmdutils.LogIntendedAction(
		v.PlanMode, "update Kubernetes API endpoint access for cluster %q in %q to: privateAccess=%v, publicAccess=%v",
		v.ClusterMeta.Name, v.ClusterMeta.Region, *desired.PrivateAccess, *desired.PublicAccess)
	if api.PrivateOnly(&desired) {
		logger.Warning(api.ErrClusterEndpointPrivateOnly.Error())
	}
	if v.PlanMode {
		return nil
	}
	endpointUpdate := &ekstypes.VpcConfigRequest{
		EndpointPrivateAccess: desired.PrivateAccess,
		EndpointPublicAccess:  desired.PublicAccess,
	}
	if err := v.updateVPCConfig(ctx, endpointUpdate); err != nil {
		return err
	}
	cmdutils.LogCompletedAction(
		false,
		"Kubernetes API endpoint access for cluster %q in %q has been updated to: "+
			"privateAccess=%v, publicAccess=%v",
		v.ClusterMeta.Name, v.ClusterMeta.Region, *desired.PrivateAccess, *desired.PublicAccess)
	return nil
}

func (v *VPCHelper) updatePublicAccessCIDRs(ctx context.Context, vpc *api.ClusterVPC) error {
	if cidrsEqual(v.Cluster.ResourcesVpcConfig.PublicAccessCidrs, vpc.PublicAccessCIDRs) {
		logger.Success("public access CIDRs for cluster %q in %q are already up-to-date",
			v.ClusterMeta.Name, v.ClusterMeta.Region)
		return nil
	}

	logger.Info("current public access CIDRs: %v", v.Cluster.ResourcesVpcConfig.PublicAccessCidrs)
	cmdutils.LogIntendedAction(
		v.PlanMode, "update public access CIDRs for cluster %q in %q to: %v",
		v.ClusterMeta.Name, v.ClusterMeta.Region, vpc.PublicAccessCIDRs)

	if v.PlanMode {
		return nil
	}

	if err := v.updateVPCConfig(ctx, &ekstypes.VpcConfigRequest{
		PublicAccessCidrs: vpc.PublicAccessCIDRs,
	}); err != nil {
		return fmt.Errorf("error updating CIDRs for public access: %w", err)
	}
	cmdutils.LogCompletedAction(
		false,
		"public access CIDRs for cluster %q in %q have been updated to: %v",
		v.ClusterMeta.Name, v.ClusterMeta.Region, vpc.PublicAccessCIDRs)
	return nil
}

func (v *VPCHelper) updateVPCConfig(ctx context.Context, vpcConfig *ekstypes.VpcConfigRequest) error {
	return v.VPCUpdater.UpdateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
		Name:               v.Cluster.Name,
		ResourcesVpcConfig: vpcConfig,
	})
}

func cidrsEqual(currentValues, newValues []string) bool {
	if len(newValues) == 0 && len(currentValues) == 1 && currentValues[0] == "0.0.0.0/0" {
		return true
	}
	return slices.Equal(currentValues, newValues)
}
