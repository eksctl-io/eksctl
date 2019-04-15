package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassCoreDefinition_CoreDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::CoreDefinition.CoreDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-coredefinition-coredefinitionversion.html
type AWSGreengrassCoreDefinition_CoreDefinitionVersion struct {

	// Cores AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-coredefinition-coredefinitionversion.html#cfn-greengrass-coredefinition-coredefinitionversion-cores
	Cores []AWSGreengrassCoreDefinition_Core `json:"Cores,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassCoreDefinition_CoreDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::CoreDefinition.CoreDefinitionVersion"
}

func (r *AWSGreengrassCoreDefinition_CoreDefinitionVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
