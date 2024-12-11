package events

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Rule_BatchParameters AWS CloudFormation Resource (AWS::Events::Rule.BatchParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-batchparameters.html
type Rule_BatchParameters struct {

	// ArrayProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-batchparameters.html#cfn-events-rule-batchparameters-arrayproperties
	ArrayProperties *Rule_BatchArrayProperties `json:"ArrayProperties,omitempty"`

	// JobDefinition AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-batchparameters.html#cfn-events-rule-batchparameters-jobdefinition
	JobDefinition *types.Value `json:"JobDefinition,omitempty"`

	// JobName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-batchparameters.html#cfn-events-rule-batchparameters-jobname
	JobName *types.Value `json:"JobName,omitempty"`

	// RetryStrategy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-batchparameters.html#cfn-events-rule-batchparameters-retrystrategy
	RetryStrategy *Rule_BatchRetryStrategy `json:"RetryStrategy,omitempty"`

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
func (r *Rule_BatchParameters) AWSCloudFormationType() string {
	return "AWS::Events::Rule.BatchParameters"
}
