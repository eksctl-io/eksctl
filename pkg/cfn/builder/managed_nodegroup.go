package builder

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

// ManagedNodeGroupResourceSet defines the CloudFormation resources required for a managed nodegroup
type ManagedNodeGroupResourceSet struct {
	clusterConfig    *api.ClusterConfig
	clusterStackName string
	nodeGroup        *api.ManagedNodeGroup
	*resourceSet
}

// NewManagedNodeGroup creates a new ManagedNodeGroupResourceSet
func NewManagedNodeGroup(cluster *api.ClusterConfig, nodeGroup *api.ManagedNodeGroup, clusterStackName string) *ManagedNodeGroupResourceSet {
	return &ManagedNodeGroupResourceSet{
		clusterConfig:    cluster,
		clusterStackName: clusterStackName,
		nodeGroup:        nodeGroup,
		resourceSet:      newResourceSet(),
	}
}

// AddAllResources adds all required CloudFormation resources
func (m *ManagedNodeGroupResourceSet) AddAllResources() error {
	m.resourceSet.template.Description = fmt.Sprintf(
		"%s (SSH access: %v) %s",
		"EKS Managed Nodes",
		api.IsEnabled(m.nodeGroup.SSH.Allow),
		"[created by eksctl]")

	m.template.Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings

	var nodeRole *gfnt.Value
	if m.nodeGroup.IAM.InstanceRoleARN == "" {
		if err := createRole(m.resourceSet, m.nodeGroup.IAM, true); err != nil {
			return err
		}
		nodeRole = gfnt.MakeFnGetAttString(cfnIAMInstanceRoleName, "Arn")
	} else {
		nodeRole = gfnt.NewString(m.nodeGroup.IAM.InstanceRoleARN)
	}

	subnets, err := AssignSubnets(m.nodeGroup.AvailabilityZones, m.clusterStackName, m.clusterConfig, m.nodeGroup.PrivateNetworking)
	if err != nil {
		return err
	}

	scalingConfig := gfneks.Nodegroup_ScalingConfig{}
	if m.nodeGroup.MinSize != nil {
		scalingConfig.MinSize = gfnt.NewInteger(*m.nodeGroup.MinSize)
	}
	if m.nodeGroup.MaxSize != nil {
		scalingConfig.MaxSize = gfnt.NewInteger(*m.nodeGroup.MaxSize)
	}
	if m.nodeGroup.DesiredCapacity != nil {
		scalingConfig.DesiredSize = gfnt.NewInteger(*m.nodeGroup.DesiredCapacity)
	}
	managedResource := &gfneks.Nodegroup{
		ClusterName:   gfnt.NewString(m.clusterConfig.Metadata.Name),
		NodegroupName: gfnt.NewString(m.nodeGroup.Name),
		ScalingConfig: &scalingConfig,
		Subnets:       subnets,
		// Currently the API supports specifying only one instance type
		InstanceTypes: gfnt.NewStringSlice(m.nodeGroup.InstanceType),
		AmiType:       gfnt.NewString(getAMIType(m.nodeGroup.InstanceType)),
		NodeRole:      nodeRole,
		Labels:        m.nodeGroup.Labels,
		Tags:          m.nodeGroup.Tags,
	}

	if api.IsEnabled(m.nodeGroup.SSH.Allow) {
		managedResource.RemoteAccess = &gfneks.Nodegroup_RemoteAccess{
			Ec2SshKey:            gfnt.NewString(*m.nodeGroup.SSH.PublicKeyName),
			SourceSecurityGroups: gfnt.NewStringSlice(m.nodeGroup.SSH.SourceSecurityGroupIDs...),
		}
	}
	if m.nodeGroup.VolumeSize != nil {
		managedResource.DiskSize = gfnt.NewInteger(*m.nodeGroup.VolumeSize)
	}

	m.newResource("ManagedNodeGroup", managedResource)

	return nil
}

func getAMIType(instanceType string) string {
	if utils.IsGPUInstanceType(instanceType) {
		return eks.AMITypesAl2X8664Gpu
	}
	if utils.IsARMInstanceType(instanceType) {
		// TODO Upgrade SDK and use constant from the eks library
		return "AL2_ARM_64"
	}
	return eks.AMITypesAl2X8664
}

// RenderJSON implements the ResourceSet interface
func (m *ManagedNodeGroupResourceSet) RenderJSON() ([]byte, error) {
	return m.resourceSet.renderJSON()
}

// WithIAM implements the ResourceSet interface
func (m *ManagedNodeGroupResourceSet) WithIAM() bool {
	// eksctl does not support passing pre-created IAM instance roles to Managed Nodes,
	// so the IAM capability is always required
	return true
}

// WithNamedIAM implements the ResourceSet interface
func (m *ManagedNodeGroupResourceSet) WithNamedIAM() bool {
	return m.nodeGroup.IAM.InstanceRoleName != ""
}
