package vpc

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/fake_vpc_importer.go . Importer
type Importer interface {
	VPC() *gfnt.Value
	ClusterSecurityGroup() *gfnt.Value
	ControlPlaneSecurityGroup() *gfnt.Value
	SharedNodeSecurityGroup() *gfnt.Value
	SecurityGroups() gfnt.Slice
	SubnetsPublic() *gfnt.Value
	SubnetsPrivate() *gfnt.Value
}

// StackConfigImporter returns VPC info based on the Cluster Stack
type StackConfigImporter struct {
	clusterStackName string
}

// NewStackConfigImporter creates a new StackConfigImporter instance
func NewStackConfigImporter(clusterStackName string) *StackConfigImporter {
	return &StackConfigImporter{
		clusterStackName: clusterStackName,
	}
}

// VPC returns a gfnt value based on the cluster stack name and the VPC from the
// cluster stack output
func (si *StackConfigImporter) VPC() *gfnt.Value {
	return makeImportValue(si.clusterStackName, outputs.ClusterVPC)
}

// ClusterSecurityGroup returns a gfnt value based on the cluster stack name
// and the default security group from the cluster stack output
func (si *StackConfigImporter) ClusterSecurityGroup() *gfnt.Value {
	return makeImportValue(si.clusterStackName, outputs.ClusterDefaultSecurityGroup)
}

// ControlPlaneSecurityGroup returns a gfnt value based on the cluster stack name
// and the control plane security group from the cluster stack output
func (si *StackConfigImporter) ControlPlaneSecurityGroup() *gfnt.Value {
	return makeImportValue(si.clusterStackName, outputs.ClusterSecurityGroup)
}

// SharedNodeSecurityGroup returns a gfnt value based on the cluster stack name
// and the shared node security group from the cluster stack output
func (si *StackConfigImporter) SharedNodeSecurityGroup() *gfnt.Value {
	return makeImportValue(si.clusterStackName, outputs.ClusterSharedNodeSecurityGroup)
}

// SecurityGroups returns a gfnt slice based on the cluster stack name
// and the default security group from the cluster stack output
func (si *StackConfigImporter) SecurityGroups() gfnt.Slice {
	return gfnt.Slice{si.ClusterSecurityGroup()}
}

// SubnetsPublic returns a gfnt value based on the cluster stack name
// and the public subnets from the cluster stack output
func (si *StackConfigImporter) SubnetsPublic() *gfnt.Value {
	return gfnt.MakeFnSplit(",", makeImportValue(si.clusterStackName, outputs.ClusterSubnetsPublic))
}

// SubnetsPrivate returns a gfnt value based on the cluster stack name
// and the public subnets from the cluster stack output
func (si *StackConfigImporter) SubnetsPrivate() *gfnt.Value {
	return gfnt.MakeFnSplit(",", makeImportValue(si.clusterStackName, outputs.ClusterSubnetsPrivate))
}

func makeImportValue(stackName, output string) *gfnt.Value {
	return gfnt.MakeFnImportValueString(fmt.Sprintf("%s::%s", stackName, output))
}

// SpecConfigImporter returns VPC info based on the ClusterConfig Spec
type SpecConfigImporter struct {
	clusterSecurityGroup string
	vpc                  *api.ClusterVPC
}

// NewSpecConfigImporter creates a new SpecConfigImporter instance
func NewSpecConfigImporter(securityGroup string, vpc *api.ClusterVPC) *SpecConfigImporter {
	return &SpecConfigImporter{
		clusterSecurityGroup: securityGroup,
		vpc:                  vpc,
	}
}

// VPC returns the gfnt value of the cluster config VPC ID
func (si *SpecConfigImporter) VPC() *gfnt.Value {
	return gfnt.NewString(si.vpc.ID)
}

// ClusterSecurityGroup returns the gfnt value of the default cluser security group
func (si *SpecConfigImporter) ClusterSecurityGroup() *gfnt.Value {
	return gfnt.NewString(si.clusterSecurityGroup)
}

// ControlPlaneSecurityGroup returns the gfnt value of the cluster config VPC
// securityGroup
func (si *SpecConfigImporter) ControlPlaneSecurityGroup() *gfnt.Value {
	return gfnt.NewString(si.vpc.SecurityGroup)
}

// SharedNodeSecurityGroup returns the gfnt value of the cluster config VPC
// sharedNodeSecurityGroup if it is set. If not, it returns the default
// cluster security group
func (si *SpecConfigImporter) SharedNodeSecurityGroup() *gfnt.Value {
	if si.vpc.SharedNodeSecurityGroup != "" {
		return gfnt.NewString(si.vpc.SharedNodeSecurityGroup)
	}
	return si.ClusterSecurityGroup()
}

// SecurityGroups returns a gfnt slice of the ClusterSecurityGroup
func (si *SpecConfigImporter) SecurityGroups() gfnt.Slice {
	return gfnt.Slice{si.ClusterSecurityGroup()}
}

// SubnetsPublic returns a gfnt string slice of the Public subnets from the
// cluster config VPC subnets spec
func (si *SpecConfigImporter) SubnetsPublic() *gfnt.Value {
	return gfnt.NewStringSlice(si.vpc.Subnets.Public.WithIDs()...)
}

// SubnetsPrivate returns a gfnt string slice of the Private subnets from the
// cluster config VPC subnets spec
func (si *SpecConfigImporter) SubnetsPrivate() *gfnt.Value {
	return gfnt.NewStringSlice(si.vpc.Subnets.Private.WithIDs()...)
}
