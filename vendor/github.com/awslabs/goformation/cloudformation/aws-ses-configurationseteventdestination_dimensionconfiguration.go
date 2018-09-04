package cloudformation

// AWSSESConfigurationSetEventDestination_DimensionConfiguration AWS CloudFormation Resource (AWS::SES::ConfigurationSetEventDestination.DimensionConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html
type AWSSESConfigurationSetEventDestination_DimensionConfiguration struct {

	// DefaultDimensionValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html#cfn-ses-configurationseteventdestination-dimensionconfiguration-defaultdimensionvalue
	DefaultDimensionValue *StringIntrinsic `json:"DefaultDimensionValue,omitempty"`

	// DimensionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html#cfn-ses-configurationseteventdestination-dimensionconfiguration-dimensionname
	DimensionName *StringIntrinsic `json:"DimensionName,omitempty"`

	// DimensionValueSource AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-configurationseteventdestination-dimensionconfiguration.html#cfn-ses-configurationseteventdestination-dimensionconfiguration-dimensionvaluesource
	DimensionValueSource *StringIntrinsic `json:"DimensionValueSource,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESConfigurationSetEventDestination_DimensionConfiguration) AWSCloudFormationType() string {
	return "AWS::SES::ConfigurationSetEventDestination.DimensionConfiguration"
}
