package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Asset_EgressEndpoint AWS CloudFormation Resource (AWS::MediaPackage::Asset.EgressEndpoint)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-asset-egressendpoint.html
type Asset_EgressEndpoint struct {

	// PackagingConfigurationId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-asset-egressendpoint.html#cfn-mediapackage-asset-egressendpoint-packagingconfigurationid
	PackagingConfigurationId *types.Value `json:"PackagingConfigurationId,omitempty"`

	// Url AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-asset-egressendpoint.html#cfn-mediapackage-asset-egressendpoint-url
	Url *types.Value `json:"Url,omitempty"`

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
func (r *Asset_EgressEndpoint) AWSCloudFormationType() string {
	return "AWS::MediaPackage::Asset.EgressEndpoint"
}
