package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// IPAMPool_SourceResource AWS CloudFormation Resource (AWS::EC2::IPAMPool.SourceResource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipampool-sourceresource.html
type IPAMPool_SourceResource struct {

	// ResourceId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipampool-sourceresource.html#cfn-ec2-ipampool-sourceresource-resourceid
	ResourceId *types.Value `json:"ResourceId,omitempty"`

	// ResourceOwner AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipampool-sourceresource.html#cfn-ec2-ipampool-sourceresource-resourceowner
	ResourceOwner *types.Value `json:"ResourceOwner,omitempty"`

	// ResourceRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipampool-sourceresource.html#cfn-ec2-ipampool-sourceresource-resourceregion
	ResourceRegion *types.Value `json:"ResourceRegion,omitempty"`

	// ResourceType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipampool-sourceresource.html#cfn-ec2-ipampool-sourceresource-resourcetype
	ResourceType *types.Value `json:"ResourceType,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *IPAMPool_SourceResource) AWSCloudFormationType() string {
	return "AWS::EC2::IPAMPool.SourceResource"
}
