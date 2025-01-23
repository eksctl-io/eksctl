package appflow

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Flow_DestinationConnectorProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.DestinationConnectorProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html
type Flow_DestinationConnectorProperties struct {

	// EventBridge AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-eventbridge
	EventBridge *Flow_EventBridgeDestinationProperties `json:"EventBridge,omitempty"`

	// LookoutMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-lookoutmetrics
	LookoutMetrics *Flow_LookoutMetricsDestinationProperties `json:"LookoutMetrics,omitempty"`

	// Redshift AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-redshift
	Redshift *Flow_RedshiftDestinationProperties `json:"Redshift,omitempty"`

	// S3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-s3
	S3 *Flow_S3DestinationProperties `json:"S3,omitempty"`

	// Salesforce AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-salesforce
	Salesforce *Flow_SalesforceDestinationProperties `json:"Salesforce,omitempty"`

	// Snowflake AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-snowflake
	Snowflake *Flow_SnowflakeDestinationProperties `json:"Snowflake,omitempty"`

	// Upsolver AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-upsolver
	Upsolver *Flow_UpsolverDestinationProperties `json:"Upsolver,omitempty"`

	// Zendesk AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-destinationconnectorproperties.html#cfn-appflow-flow-destinationconnectorproperties-zendesk
	Zendesk *Flow_ZendeskDestinationProperties `json:"Zendesk,omitempty"`

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
func (r *Flow_DestinationConnectorProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.DestinationConnectorProperties"
}
