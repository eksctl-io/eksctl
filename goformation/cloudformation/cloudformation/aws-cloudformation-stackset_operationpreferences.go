package cloudformation

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StackSet_OperationPreferences AWS CloudFormation Resource (AWS::CloudFormation::StackSet.OperationPreferences)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-operationpreferences.html
type StackSet_OperationPreferences struct {

	// FailureToleranceCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-operationpreferences.html#cfn-cloudformation-stackset-operationpreferences-failuretolerancecount
	FailureToleranceCount *types.Value `json:"FailureToleranceCount,omitempty"`

	// FailureTolerancePercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-operationpreferences.html#cfn-cloudformation-stackset-operationpreferences-failuretolerancepercentage
	FailureTolerancePercentage *types.Value `json:"FailureTolerancePercentage,omitempty"`

	// MaxConcurrentCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-operationpreferences.html#cfn-cloudformation-stackset-operationpreferences-maxconcurrentcount
	MaxConcurrentCount *types.Value `json:"MaxConcurrentCount,omitempty"`

	// MaxConcurrentPercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-operationpreferences.html#cfn-cloudformation-stackset-operationpreferences-maxconcurrentpercentage
	MaxConcurrentPercentage *types.Value `json:"MaxConcurrentPercentage,omitempty"`

	// RegionConcurrencyType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-operationpreferences.html#cfn-cloudformation-stackset-operationpreferences-regionconcurrencytype
	RegionConcurrencyType *types.Value `json:"RegionConcurrencyType,omitempty"`

	// RegionOrder AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-operationpreferences.html#cfn-cloudformation-stackset-operationpreferences-regionorder
	RegionOrder *types.Value `json:"RegionOrder,omitempty"`

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
func (r *StackSet_OperationPreferences) AWSCloudFormationType() string {
	return "AWS::CloudFormation::StackSet.OperationPreferences"
}
