package events

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// Connection_ConnectionHttpParameters AWS CloudFormation Resource (AWS::Events::Connection.ConnectionHttpParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-connectionhttpparameters.html
type Connection_ConnectionHttpParameters struct {

	// BodyParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-connectionhttpparameters.html#cfn-events-connection-connectionhttpparameters-bodyparameters
	BodyParameters []Connection_Parameter `json:"BodyParameters,omitempty"`

	// HeaderParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-connectionhttpparameters.html#cfn-events-connection-connectionhttpparameters-headerparameters
	HeaderParameters []Connection_Parameter `json:"HeaderParameters,omitempty"`

	// QueryStringParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-connection-connectionhttpparameters.html#cfn-events-connection-connectionhttpparameters-querystringparameters
	QueryStringParameters []Connection_Parameter `json:"QueryStringParameters,omitempty"`

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
func (r *Connection_ConnectionHttpParameters) AWSCloudFormationType() string {
	return "AWS::Events::Connection.ConnectionHttpParameters"
}
