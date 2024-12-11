package ssmincidents

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResponsePlan_IncidentTemplate AWS CloudFormation Resource (AWS::SSMIncidents::ResponsePlan.IncidentTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-incidenttemplate.html
type ResponsePlan_IncidentTemplate struct {

	// DedupeString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-incidenttemplate.html#cfn-ssmincidents-responseplan-incidenttemplate-dedupestring
	DedupeString *types.Value `json:"DedupeString,omitempty"`

	// Impact AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-incidenttemplate.html#cfn-ssmincidents-responseplan-incidenttemplate-impact
	Impact *types.Value `json:"Impact"`

	// NotificationTargets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-incidenttemplate.html#cfn-ssmincidents-responseplan-incidenttemplate-notificationtargets
	NotificationTargets []ResponsePlan_NotificationTargetItem `json:"NotificationTargets,omitempty"`

	// Summary AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-incidenttemplate.html#cfn-ssmincidents-responseplan-incidenttemplate-summary
	Summary *types.Value `json:"Summary,omitempty"`

	// Title AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmincidents-responseplan-incidenttemplate.html#cfn-ssmincidents-responseplan-incidenttemplate-title
	Title *types.Value `json:"Title,omitempty"`

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
func (r *ResponsePlan_IncidentTemplate) AWSCloudFormationType() string {
	return "AWS::SSMIncidents::ResponsePlan.IncidentTemplate"
}
