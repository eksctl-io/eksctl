package applicationinsights

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Application_JMXPrometheusExporter AWS CloudFormation Resource (AWS::ApplicationInsights::Application.JMXPrometheusExporter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-jmxprometheusexporter.html
type Application_JMXPrometheusExporter struct {

	// HostPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-jmxprometheusexporter.html#cfn-applicationinsights-application-jmxprometheusexporter-hostport
	HostPort *types.Value `json:"HostPort,omitempty"`

	// JMXURL AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-jmxprometheusexporter.html#cfn-applicationinsights-application-jmxprometheusexporter-jmxurl
	JMXURL *types.Value `json:"JMXURL,omitempty"`

	// PrometheusPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-jmxprometheusexporter.html#cfn-applicationinsights-application-jmxprometheusexporter-prometheusport
	PrometheusPort *types.Value `json:"PrometheusPort,omitempty"`

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
func (r *Application_JMXPrometheusExporter) AWSCloudFormationType() string {
	return "AWS::ApplicationInsights::Application.JMXPrometheusExporter"
}
