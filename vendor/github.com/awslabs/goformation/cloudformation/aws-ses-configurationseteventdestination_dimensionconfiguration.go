package cloudformation

import (
	"encoding/json"
)

// AWSSESConfigurationSetEventDestination_DimensionConfiguration AWS CloudFormation Resource (AWS::SES::ConfigurationSetEventDestination.DimensionConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html
type AWSSESConfigurationSetEventDestination_DimensionConfiguration struct {

	// DefaultDimensionValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html#cfn-ses-configurationseteventdestination-dimensionconfiguration-defaultdimensionvalue
	DefaultDimensionValue *Value `json:"DefaultDimensionValue,omitempty"`

	// DimensionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html#cfn-ses-configurationseteventdestination-dimensionconfiguration-dimensionname
	DimensionName *Value `json:"DimensionName,omitempty"`

	// DimensionValueSource AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html#cfn-ses-configurationseteventdestination-dimensionconfiguration-dimensionvaluesource
	DimensionValueSource *Value `json:"DimensionValueSource,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESConfigurationSetEventDestination_DimensionConfiguration) AWSCloudFormationType() string {
	return "AWS::SES::ConfigurationSetEventDestination.DimensionConfiguration"
}

func (r *AWSSESConfigurationSetEventDestination_DimensionConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
