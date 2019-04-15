package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamStack_StorageConnector AWS CloudFormation Resource (AWS::AppStream::Stack.StorageConnector)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-storageconnector.html
type AWSAppStreamStack_StorageConnector struct {

	// ConnectorType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-storageconnector.html#cfn-appstream-stack-storageconnector-connectortype
	ConnectorType *Value `json:"ConnectorType,omitempty"`

	// Domains AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-storageconnector.html#cfn-appstream-stack-storageconnector-domains
	Domains []*Value `json:"Domains,omitempty"`

	// ResourceIdentifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-storageconnector.html#cfn-appstream-stack-storageconnector-resourceidentifier
	ResourceIdentifier *Value `json:"ResourceIdentifier,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamStack_StorageConnector) AWSCloudFormationType() string {
	return "AWS::AppStream::Stack.StorageConnector"
}

func (r *AWSAppStreamStack_StorageConnector) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
