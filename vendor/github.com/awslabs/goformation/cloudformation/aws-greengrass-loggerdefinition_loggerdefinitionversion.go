package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassLoggerDefinition_LoggerDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::LoggerDefinition.LoggerDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-loggerdefinitionversion.html
type AWSGreengrassLoggerDefinition_LoggerDefinitionVersion struct {

	// Loggers AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-loggerdefinition-loggerdefinitionversion.html#cfn-greengrass-loggerdefinition-loggerdefinitionversion-loggers
	Loggers []AWSGreengrassLoggerDefinition_Logger `json:"Loggers,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassLoggerDefinition_LoggerDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::LoggerDefinition.LoggerDefinitionVersion"
}

func (r *AWSGreengrassLoggerDefinition_LoggerDefinitionVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
