package events

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Rule_HttpParameters AWS CloudFormation Resource (AWS::Events::Rule.HttpParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-httpparameters.html
type Rule_HttpParameters struct {

	// HeaderParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-httpparameters.html#cfn-events-rule-httpparameters-headerparameters
	HeaderParameters map[string]*types.Value `json:"HeaderParameters,omitempty"`

	// PathParameterValues AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-httpparameters.html#cfn-events-rule-httpparameters-pathparametervalues
	PathParameterValues *types.Value `json:"PathParameterValues,omitempty"`

	// QueryStringParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-httpparameters.html#cfn-events-rule-httpparameters-querystringparameters
	QueryStringParameters map[string]*types.Value `json:"QueryStringParameters,omitempty"`

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
func (r *Rule_HttpParameters) AWSCloudFormationType() string {
	return "AWS::Events::Rule.HttpParameters"
}
