package batch

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SchedulingPolicy_ShareAttributes AWS CloudFormation Resource (AWS::Batch::SchedulingPolicy.ShareAttributes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-schedulingpolicy-shareattributes.html
type SchedulingPolicy_ShareAttributes struct {

	// ShareIdentifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-schedulingpolicy-shareattributes.html#cfn-batch-schedulingpolicy-shareattributes-shareidentifier
	ShareIdentifier *types.Value `json:"ShareIdentifier,omitempty"`

	// WeightFactor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-schedulingpolicy-shareattributes.html#cfn-batch-schedulingpolicy-shareattributes-weightfactor
	WeightFactor *types.Value `json:"WeightFactor,omitempty"`

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
func (r *SchedulingPolicy_ShareAttributes) AWSCloudFormationType() string {
	return "AWS::Batch::SchedulingPolicy.ShareAttributes"
}
