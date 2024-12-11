package apigateway

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// UsagePlan_QuotaSettings AWS CloudFormation Resource (AWS::ApiGateway::UsagePlan.QuotaSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-usageplan-quotasettings.html
type UsagePlan_QuotaSettings struct {

	// Limit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-usageplan-quotasettings.html#cfn-apigateway-usageplan-quotasettings-limit
	Limit *types.Value `json:"Limit,omitempty"`

	// Offset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-usageplan-quotasettings.html#cfn-apigateway-usageplan-quotasettings-offset
	Offset *types.Value `json:"Offset,omitempty"`

	// Period AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-usageplan-quotasettings.html#cfn-apigateway-usageplan-quotasettings-period
	Period *types.Value `json:"Period,omitempty"`

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
func (r *UsagePlan_QuotaSettings) AWSCloudFormationType() string {
	return "AWS::ApiGateway::UsagePlan.QuotaSettings"
}
