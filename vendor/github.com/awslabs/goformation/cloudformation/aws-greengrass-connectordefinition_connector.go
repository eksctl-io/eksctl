package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassConnectorDefinition_Connector AWS CloudFormation Resource (AWS::Greengrass::ConnectorDefinition.Connector)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinition-connector.html
type AWSGreengrassConnectorDefinition_Connector struct {

	// ConnectorArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinition-connector.html#cfn-greengrass-connectordefinition-connector-connectorarn
	ConnectorArn *Value `json:"ConnectorArn,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinition-connector.html#cfn-greengrass-connectordefinition-connector-id
	Id *Value `json:"Id,omitempty"`

	// Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-connectordefinition-connector.html#cfn-greengrass-connectordefinition-connector-parameters
	Parameters interface{} `json:"Parameters,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassConnectorDefinition_Connector) AWSCloudFormationType() string {
	return "AWS::Greengrass::ConnectorDefinition.Connector"
}

func (r *AWSGreengrassConnectorDefinition_Connector) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
