package events

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// Connection_AuthParameters AWS CloudFormation Resource (AWS::Events::Connection.AuthParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-authparameters.html
type Connection_AuthParameters struct {

	// ApiKeyAuthParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-authparameters.html#cfn-events-connection-authparameters-apikeyauthparameters
	ApiKeyAuthParameters *Connection_ApiKeyAuthParameters `json:"ApiKeyAuthParameters,omitempty"`

	// BasicAuthParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-authparameters.html#cfn-events-connection-authparameters-basicauthparameters
	BasicAuthParameters *Connection_BasicAuthParameters `json:"BasicAuthParameters,omitempty"`

	// ConnectivityParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-authparameters.html#cfn-events-connection-authparameters-connectivityparameters
	ConnectivityParameters *Connection_ConnectivityParameters `json:"ConnectivityParameters,omitempty"`

	// InvocationHttpParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-authparameters.html#cfn-events-connection-authparameters-invocationhttpparameters
	InvocationHttpParameters *Connection_ConnectionHttpParameters `json:"InvocationHttpParameters,omitempty"`

	// OAuthParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-authparameters.html#cfn-events-connection-authparameters-oauthparameters
	OAuthParameters *Connection_OAuthParameters `json:"OAuthParameters,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Connection_AuthParameters) AWSCloudFormationType() string {
	return "AWS::Events::Connection.AuthParameters"
}
