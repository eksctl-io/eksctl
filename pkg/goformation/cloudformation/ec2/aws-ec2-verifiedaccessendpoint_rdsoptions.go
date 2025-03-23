package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VerifiedAccessEndpoint_RdsOptions AWS CloudFormation Resource (AWS::EC2::VerifiedAccessEndpoint.RdsOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html
type VerifiedAccessEndpoint_RdsOptions struct {

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html#cfn-ec2-verifiedaccessendpoint-rdsoptions-port
	Port *types.Value `json:"Port,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html#cfn-ec2-verifiedaccessendpoint-rdsoptions-protocol
	Protocol *types.Value `json:"Protocol,omitempty"`

	// RdsDbClusterArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html#cfn-ec2-verifiedaccessendpoint-rdsoptions-rdsdbclusterarn
	RdsDbClusterArn *types.Value `json:"RdsDbClusterArn,omitempty"`

	// RdsDbInstanceArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html#cfn-ec2-verifiedaccessendpoint-rdsoptions-rdsdbinstancearn
	RdsDbInstanceArn *types.Value `json:"RdsDbInstanceArn,omitempty"`

	// RdsDbProxyArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html#cfn-ec2-verifiedaccessendpoint-rdsoptions-rdsdbproxyarn
	RdsDbProxyArn *types.Value `json:"RdsDbProxyArn,omitempty"`

	// RdsEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html#cfn-ec2-verifiedaccessendpoint-rdsoptions-rdsendpoint
	RdsEndpoint *types.Value `json:"RdsEndpoint,omitempty"`

	// SubnetIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-rdsoptions.html#cfn-ec2-verifiedaccessendpoint-rdsoptions-subnetids
	SubnetIds *types.Value `json:"SubnetIds,omitempty"`

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
func (r *VerifiedAccessEndpoint_RdsOptions) AWSCloudFormationType() string {
	return "AWS::EC2::VerifiedAccessEndpoint.RdsOptions"
}
