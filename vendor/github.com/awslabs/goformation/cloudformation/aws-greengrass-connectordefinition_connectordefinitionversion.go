package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassConnectorDefinition_ConnectorDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::ConnectorDefinition.ConnectorDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinition-connectordefinitionversion.html
type AWSGreengrassConnectorDefinition_ConnectorDefinitionVersion struct {

	// Connectors AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinition-connectordefinitionversion.html#cfn-greengrass-connectordefinition-connectordefinitionversion-connectors
	Connectors []AWSGreengrassConnectorDefinition_Connector `json:"Connectors,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassConnectorDefinition_ConnectorDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::ConnectorDefinition.ConnectorDefinitionVersion"
}

func (r *AWSGreengrassConnectorDefinition_ConnectorDefinitionVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
