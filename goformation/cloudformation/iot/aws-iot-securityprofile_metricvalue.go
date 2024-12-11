package iot

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SecurityProfile_MetricValue AWS CloudFormation Resource (AWS::IoT::SecurityProfile.MetricValue)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metricvalue.html
type SecurityProfile_MetricValue struct {

	// Cidrs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metricvalue.html#cfn-iot-securityprofile-metricvalue-cidrs
	Cidrs *types.Value `json:"Cidrs,omitempty"`

	// Count AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metricvalue.html#cfn-iot-securityprofile-metricvalue-count
	Count *types.Value `json:"Count,omitempty"`

	// Number AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metricvalue.html#cfn-iot-securityprofile-metricvalue-number
	Number *types.Value `json:"Number,omitempty"`

	// Numbers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metricvalue.html#cfn-iot-securityprofile-metricvalue-numbers
	Numbers *types.Value `json:"Numbers,omitempty"`

	// Ports AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metricvalue.html#cfn-iot-securityprofile-metricvalue-ports
	Ports *types.Value `json:"Ports,omitempty"`

	// Strings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-securityprofile-metricvalue.html#cfn-iot-securityprofile-metricvalue-strings
	Strings *types.Value `json:"Strings,omitempty"`

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
func (r *SecurityProfile_MetricValue) AWSCloudFormationType() string {
	return "AWS::IoT::SecurityProfile.MetricValue"
}
