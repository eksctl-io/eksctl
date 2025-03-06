package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VPNConnection_Phase2IntegrityAlgorithmsRequestListValue AWS CloudFormation Resource (AWS::EC2::VPNConnection.Phase2IntegrityAlgorithmsRequestListValue)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-phase2integrityalgorithmsrequestlistvalue.html
type VPNConnection_Phase2IntegrityAlgorithmsRequestListValue struct {

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-phase2integrityalgorithmsrequestlistvalue.html#cfn-ec2-vpnconnection-phase2integrityalgorithmsrequestlistvalue-value
	Value *types.Value `json:"Value,omitempty"`

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
func (r *VPNConnection_Phase2IntegrityAlgorithmsRequestListValue) AWSCloudFormationType() string {
	return "AWS::EC2::VPNConnection.Phase2IntegrityAlgorithmsRequestListValue"
}
