package ec2

import (
	"github.com/awslabs/goformation/v4/cloudformation/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// LocalGatewayRouteTableVPCAssociation_Tags AWS CloudFormation Resource (AWS::EC2::LocalGatewayRouteTableVPCAssociation.Tags)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-localgatewayroutetablevpcassociation-tags.html
type LocalGatewayRouteTableVPCAssociation_Tags struct {

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-localgatewayroutetablevpcassociation-tags.html#cfn-ec2-localgatewayroutetablevpcassociation-tags-tags
	Tags []cloudformation.Tag `json:"Tags,omitempty"`

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
func (r *LocalGatewayRouteTableVPCAssociation_Tags) AWSCloudFormationType() string {
	return "AWS::EC2::LocalGatewayRouteTableVPCAssociation.Tags"
}
