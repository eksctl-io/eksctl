package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncGraphQLApi_OpenIDConnectConfig AWS CloudFormation Resource (AWS::AppSync::GraphQLApi.OpenIDConnectConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-openidconnectconfig.html
type AWSAppSyncGraphQLApi_OpenIDConnectConfig struct {

	// AuthTTL AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-openidconnectconfig.html#cfn-appsync-graphqlapi-openidconnectconfig-authttl
	AuthTTL *Value `json:"AuthTTL,omitempty"`

	// ClientId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-openidconnectconfig.html#cfn-appsync-graphqlapi-openidconnectconfig-clientid
	ClientId *Value `json:"ClientId,omitempty"`

	// IatTTL AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-openidconnectconfig.html#cfn-appsync-graphqlapi-openidconnectconfig-iatttl
	IatTTL *Value `json:"IatTTL,omitempty"`

	// Issuer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-openidconnectconfig.html#cfn-appsync-graphqlapi-openidconnectconfig-issuer
	Issuer *Value `json:"Issuer,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncGraphQLApi_OpenIDConnectConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::GraphQLApi.OpenIDConnectConfig"
}

func (r *AWSAppSyncGraphQLApi_OpenIDConnectConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
