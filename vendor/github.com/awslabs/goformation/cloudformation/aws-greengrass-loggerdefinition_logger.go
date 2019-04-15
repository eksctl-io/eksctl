package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassLoggerDefinition_Logger AWS CloudFormation Resource (AWS::Greengrass::LoggerDefinition.Logger)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-logger.html
type AWSGreengrassLoggerDefinition_Logger struct {

	// Component AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-logger.html#cfn-greengrass-loggerdefinition-logger-component
	Component *Value `json:"Component,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-logger.html#cfn-greengrass-loggerdefinition-logger-id
	Id *Value `json:"Id,omitempty"`

	// Level AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-logger.html#cfn-greengrass-loggerdefinition-logger-level
	Level *Value `json:"Level,omitempty"`

	// Space AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-logger.html#cfn-greengrass-loggerdefinition-logger-space
	Space *Value `json:"Space,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-logger.html#cfn-greengrass-loggerdefinition-logger-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassLoggerDefinition_Logger) AWSCloudFormationType() string {
	return "AWS::Greengrass::LoggerDefinition.Logger"
}

func (r *AWSGreengrassLoggerDefinition_Logger) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
