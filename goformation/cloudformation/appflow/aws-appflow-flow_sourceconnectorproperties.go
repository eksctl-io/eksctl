package appflow

import (
	"goformation/v4/cloudformation/policies"
)

// Flow_SourceConnectorProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.SourceConnectorProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html
type Flow_SourceConnectorProperties struct {

	// Amplitude AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-amplitude
	Amplitude *Flow_AmplitudeSourceProperties `json:"Amplitude,omitempty"`

	// Datadog AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-datadog
	Datadog *Flow_DatadogSourceProperties `json:"Datadog,omitempty"`

	// Dynatrace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-dynatrace
	Dynatrace *Flow_DynatraceSourceProperties `json:"Dynatrace,omitempty"`

	// GoogleAnalytics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-googleanalytics
	GoogleAnalytics *Flow_GoogleAnalyticsSourceProperties `json:"GoogleAnalytics,omitempty"`

	// InforNexus AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-infornexus
	InforNexus *Flow_InforNexusSourceProperties `json:"InforNexus,omitempty"`

	// Marketo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-marketo
	Marketo *Flow_MarketoSourceProperties `json:"Marketo,omitempty"`

	// S3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-s3
	S3 *Flow_S3SourceProperties `json:"S3,omitempty"`

	// SAPOData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-sapodata
	SAPOData *Flow_SAPODataSourceProperties `json:"SAPOData,omitempty"`

	// Salesforce AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-salesforce
	Salesforce *Flow_SalesforceSourceProperties `json:"Salesforce,omitempty"`

	// ServiceNow AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-servicenow
	ServiceNow *Flow_ServiceNowSourceProperties `json:"ServiceNow,omitempty"`

	// Singular AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-singular
	Singular *Flow_SingularSourceProperties `json:"Singular,omitempty"`

	// Slack AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-slack
	Slack *Flow_SlackSourceProperties `json:"Slack,omitempty"`

	// Trendmicro AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-trendmicro
	Trendmicro *Flow_TrendmicroSourceProperties `json:"Trendmicro,omitempty"`

	// Veeva AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-veeva
	Veeva *Flow_VeevaSourceProperties `json:"Veeva,omitempty"`

	// Zendesk AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-sourceconnectorproperties.html#cfn-appflow-flow-sourceconnectorproperties-zendesk
	Zendesk *Flow_ZendeskSourceProperties `json:"Zendesk,omitempty"`

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
func (r *Flow_SourceConnectorProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.SourceConnectorProperties"
}
