package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassConnectorDefinitionVersion_Connector AWS CloudFormation Resource (AWS::Greengrass::ConnectorDefinitionVersion.Connector)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinitionversion-connector.html
type AWSGreengrassConnectorDefinitionVersion_Connector struct {

	// ConnectorArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinitionversion-connector.html#cfn-greengrass-connectordefinitionversion-connector-connectorarn
	ConnectorArn *Value `json:"ConnectorArn,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinitionversion-connector.html#cfn-greengrass-connectordefinitionversion-connector-id
	Id *Value `json:"Id,omitempty"`

	// Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinitionversion-connector.html#cfn-greengrass-connectordefinitionversion-connector-parameters
	Parameters interface{} `json:"Parameters,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassConnectorDefinitionVersion_Connector) AWSCloudFormationType() string {
	return "AWS::Greengrass::ConnectorDefinitionVersion.Connector"
}

func (r *AWSGreengrassConnectorDefinitionVersion_Connector) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
