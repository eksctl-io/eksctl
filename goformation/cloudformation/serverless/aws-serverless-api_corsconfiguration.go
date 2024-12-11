package serverless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Api_CorsConfiguration AWS CloudFormation Resource (AWS::Serverless::Api.CorsConfiguration)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration
type Api_CorsConfiguration struct {

	// AllowCredentials AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration
	AllowCredentials *types.Value `json:"AllowCredentials,omitempty"`

	// AllowHeaders AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration
	AllowHeaders *types.Value `json:"AllowHeaders,omitempty"`

	// AllowMethods AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration
	AllowMethods *types.Value `json:"AllowMethods,omitempty"`

	// AllowOrigin AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration
	AllowOrigin *types.Value `json:"AllowOrigin,omitempty"`

	// MaxAge AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration
	MaxAge *types.Value `json:"MaxAge,omitempty"`

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
func (r *Api_CorsConfiguration) AWSCloudFormationType() string {
	return "AWS::Serverless::Api.CorsConfiguration"
}
