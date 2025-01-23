package iotanalytics

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Dataset_LateDataRuleConfiguration AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.LateDataRuleConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-latedataruleconfiguration.html
type Dataset_LateDataRuleConfiguration struct {

	// DeltaTimeSessionWindowConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-latedataruleconfiguration.html#cfn-iotanalytics-dataset-latedataruleconfiguration-deltatimesessionwindowconfiguration
	DeltaTimeSessionWindowConfiguration *Dataset_DeltaTimeSessionWindowConfiguration `json:"DeltaTimeSessionWindowConfiguration,omitempty"`

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
func (r *Dataset_LateDataRuleConfiguration) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.LateDataRuleConfiguration"
}
