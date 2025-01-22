package serverless

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Function_DestinationConfig AWS CloudFormation Resource (AWS::Serverless::Function.DestinationConfig)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#destination-config-object
type Function_DestinationConfig struct {

	// OnFailure AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#destination-config-object
	OnFailure *Function_OnFailure `json:"OnFailure,omitempty"`

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
func (r *Function_DestinationConfig) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.DestinationConfig"
}
