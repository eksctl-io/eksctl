package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassLoggerDefinitionVersion_Logger AWS CloudFormation Resource (AWS::Greengrass::LoggerDefinitionVersion.Logger)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinitionversion-logger.html
type AWSGreengrassLoggerDefinitionVersion_Logger struct {

	// Component AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinitionversion-logger.html#cfn-greengrass-loggerdefinitionversion-logger-component
	Component *Value `json:"Component,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinitionversion-logger.html#cfn-greengrass-loggerdefinitionversion-logger-id
	Id *Value `json:"Id,omitempty"`

	// Level AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinitionversion-logger.html#cfn-greengrass-loggerdefinitionversion-logger-level
	Level *Value `json:"Level,omitempty"`

	// Space AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinitionversion-logger.html#cfn-greengrass-loggerdefinitionversion-logger-space
	Space *Value `json:"Space,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinitionversion-logger.html#cfn-greengrass-loggerdefinitionversion-logger-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassLoggerDefinitionVersion_Logger) AWSCloudFormationType() string {
	return "AWS::Greengrass::LoggerDefinitionVersion.Logger"
}

func (r *AWSGreengrassLoggerDefinitionVersion_Logger) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
