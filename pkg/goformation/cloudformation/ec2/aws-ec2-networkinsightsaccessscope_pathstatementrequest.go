package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// NetworkInsightsAccessScope_PathStatementRequest AWS CloudFormation Resource (AWS::EC2::NetworkInsightsAccessScope.PathStatementRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsaccessscope-pathstatementrequest.html
type NetworkInsightsAccessScope_PathStatementRequest struct {

	// PacketHeaderStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsaccessscope-pathstatementrequest.html#cfn-ec2-networkinsightsaccessscope-pathstatementrequest-packetheaderstatement
	PacketHeaderStatement *NetworkInsightsAccessScope_PacketHeaderStatementRequest `json:"PacketHeaderStatement,omitempty"`

	// ResourceStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsaccessscope-pathstatementrequest.html#cfn-ec2-networkinsightsaccessscope-pathstatementrequest-resourcestatement
	ResourceStatement *NetworkInsightsAccessScope_ResourceStatementRequest `json:"ResourceStatement,omitempty"`

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
func (r *NetworkInsightsAccessScope_PathStatementRequest) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkInsightsAccessScope.PathStatementRequest"
}
