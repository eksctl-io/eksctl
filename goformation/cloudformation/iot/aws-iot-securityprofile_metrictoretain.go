package iot

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SecurityProfile_MetricToRetain AWS CloudFormation Resource (AWS::IoT::SecurityProfile.MetricToRetain)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metrictoretain.html
type SecurityProfile_MetricToRetain struct {

	// Metric AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metrictoretain.html#cfn-iot-securityprofile-metrictoretain-metric
	Metric *types.Value `json:"Metric,omitempty"`

	// MetricDimension AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metrictoretain.html#cfn-iot-securityprofile-metrictoretain-metricdimension
	MetricDimension *SecurityProfile_MetricDimension `json:"MetricDimension,omitempty"`

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
func (r *SecurityProfile_MetricToRetain) AWSCloudFormationType() string {
	return "AWS::IoT::SecurityProfile.MetricToRetain"
}
