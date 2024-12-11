package ssmincidents

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResponsePlan_SsmAutomation AWS CloudFormation Resource (AWS::SSMIncidents::ResponsePlan.SsmAutomation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-ssmautomation.html
type ResponsePlan_SsmAutomation struct {

	// DocumentName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-ssmautomation.html#cfn-ssmincidents-responseplan-ssmautomation-documentname
	DocumentName *types.Value `json:"DocumentName,omitempty"`

	// DocumentVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-ssmautomation.html#cfn-ssmincidents-responseplan-ssmautomation-documentversion
	DocumentVersion *types.Value `json:"DocumentVersion,omitempty"`

	// Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-ssmautomation.html#cfn-ssmincidents-responseplan-ssmautomation-parameters
	Parameters []ResponsePlan_SsmParameter `json:"Parameters,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-ssmautomation.html#cfn-ssmincidents-responseplan-ssmautomation-rolearn
	RoleArn *types.Value `json:"RoleArn,omitempty"`

	// TargetAccount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-ssmautomation.html#cfn-ssmincidents-responseplan-ssmautomation-targetaccount
	TargetAccount *types.Value `json:"TargetAccount,omitempty"`

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
func (r *ResponsePlan_SsmAutomation) AWSCloudFormationType() string {
	return "AWS::SSMIncidents::ResponsePlan.SsmAutomation"
}
