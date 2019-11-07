package builder

import (
	"encoding/json"
	"fmt"

	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type ManagedNodeGroupResourceSet struct {
	clusterConfig    *v1alpha5.ClusterConfig
	clusterStackName string
	nodeGroup        *v1alpha5.ManagedNodeGroup
	*resourceSet
}

// This type exists because goformation does not support Managed Nodes (yet)
// Rather than setting all field types to *gfn.Value, the types are conveniently chosen
// to allow using values without requiring any conversion
type managedNodeGroup struct {
	ClusterName   string              `json:"ClusterName"`
	NodegroupName string              `json:"NodegroupName"`
	ScalingConfig *scalingConfig      `json:"ScalingConfig,omitempty"`
	DiskSize      *int64              `json:"DiskSize,omitempty"`
	Subnets       interface{}         `json:"Subnets"`
	InstanceTypes []string            `json:"InstanceTypes"`
	AmiType       string              `json:"AmiType,omitempty"`
	RemoteAccess  *remoteAccessConfig `json:"RemoteAccess,omitempty"`
	NodeRole      *gfn.Value          `json:"NodeRole"`
	Labels        map[string]string   `json:"Labels,omitempty"`
	Tags          map[string]string   `json:"Tags,omitempty"`
}

type scalingConfig struct {
	MinSize     *int `json:"MinSize,omitempty"`
	MaxSize     *int `json:"MaxSize,omitempty"`
	DesiredSize *int `json:"DesiredSize,omitempty"`
}

type remoteAccessConfig struct {
	Ec2SshKey            *string      `json:"Ec2SshKey,omitempty"`
	SourceSecurityGroups []*gfn.Value `json:"SourceSecurityGroups,omitempty"`
}

// TODO consider using the Template.Resource interface

// MarshalJSON returns the JSON encoding for this CloudFormation resource
func (e *managedNodeGroup) MarshalJSON() ([]byte, error) {
	type Properties managedNodeGroup
	return json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       "Dev::EKS::Nodegroup",
		Properties: Properties(*e),
	})

}

// NewManagedNodeGroup creates a new ManagedNodeGroupResourceSet
func NewManagedNodeGroup(cluster *v1alpha5.ClusterConfig, nodeGroup *v1alpha5.ManagedNodeGroup, clusterStackName string) *ManagedNodeGroupResourceSet {
	return &ManagedNodeGroupResourceSet{
		clusterConfig:    cluster,
		clusterStackName: clusterStackName,
		nodeGroup:        nodeGroup,
		resourceSet:      newResourceSet(),
	}
}

func (m *ManagedNodeGroupResourceSet) AddAllResources() error {
	iamHelper := NewIAMHelper(m.resourceSet, m.nodeGroup.IAM)
	iamHelper.CreateRole()

	var remoteAccess *remoteAccessConfig
	if v1alpha5.IsEnabled(m.nodeGroup.SSH.Allow) {
		remoteAccess = &remoteAccessConfig{
			Ec2SshKey: m.nodeGroup.SSH.PublicKeyName,
			// FIXME
			SourceSecurityGroups: []*gfn.Value{sgSourceAnywhereIPv4, sgSourceAnywhereIPv6},
		}
	}

	m.newResource("ManagedNodeGroup", &managedNodeGroup{
		ClusterName:   m.clusterConfig.Metadata.Name,
		NodegroupName: m.nodeGroup.Name,
		ScalingConfig: &scalingConfig{
			MinSize:     m.nodeGroup.MinSize,
			MaxSize:     m.nodeGroup.MaxSize,
			DesiredSize: m.nodeGroup.DesiredCapacity,
		},
		DiskSize: m.nodeGroup.VolumeSize,
		// Only public subnets are supported at launch
		Subnets: AssignSubnets(m.nodeGroup.AvailabilityZones, m.clusterStackName, m.clusterConfig, false),
		// Currently the API supports specifying only one instance type
		InstanceTypes: []string{m.nodeGroup.InstanceType},
		AmiType:       m.nodeGroup.AMIType,
		RemoteAccess:  remoteAccess,
		// ManagedNodeGroup.IAM.InstanceRoleARN is not supported, so this field is always retrieved from the
		// CFN resource
		NodeRole: gfn.MakeFnGetAttString(fmt.Sprintf("%s.%s", cfnIAMInstanceRoleName, "Arn")),
		Labels:   m.nodeGroup.Labels,
		Tags:     m.nodeGroup.Tags,
	})

	return nil
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
