package finspace

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Environment_FederationParameters AWS CloudFormation Resource (AWS::FinSpace::Environment.FederationParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-finspace-environment-federationparameters.html
type Environment_FederationParameters struct {

	// ApplicationCallBackURL AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-finspace-environment-federationparameters.html#cfn-finspace-environment-federationparameters-applicationcallbackurl
	ApplicationCallBackURL *types.Value `json:"ApplicationCallBackURL,omitempty"`

	// AttributeMap AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-finspace-environment-federationparameters.html#cfn-finspace-environment-federationparameters-attributemap
	AttributeMap interface{} `json:"AttributeMap,omitempty"`

	// FederationProviderName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-finspace-environment-federationparameters.html#cfn-finspace-environment-federationparameters-federationprovidername
	FederationProviderName *types.Value `json:"FederationProviderName,omitempty"`

	// FederationURN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-finspace-environment-federationparameters.html#cfn-finspace-environment-federationparameters-federationurn
	FederationURN *types.Value `json:"FederationURN,omitempty"`

	// SamlMetadataDocument AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-finspace-environment-federationparameters.html#cfn-finspace-environment-federationparameters-samlmetadatadocument
	SamlMetadataDocument *types.Value `json:"SamlMetadataDocument,omitempty"`

	// SamlMetadataURL AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-finspace-environment-federationparameters.html#cfn-finspace-environment-federationparameters-samlmetadataurl
	SamlMetadataURL *types.Value `json:"SamlMetadataURL,omitempty"`

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
func (r *Environment_FederationParameters) AWSCloudFormationType() string {
	return "AWS::FinSpace::Environment.FederationParameters"
}
