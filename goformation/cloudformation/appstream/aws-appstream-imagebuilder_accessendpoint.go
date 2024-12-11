package appstream

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ImageBuilder_AccessEndpoint AWS CloudFormation Resource (AWS::AppStream::ImageBuilder.AccessEndpoint)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-accessendpoint.html
type ImageBuilder_AccessEndpoint struct {

	// EndpointType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-accessendpoint.html#cfn-appstream-imagebuilder-accessendpoint-endpointtype
	EndpointType *types.Value `json:"EndpointType,omitempty"`

	// VpceId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-accessendpoint.html#cfn-appstream-imagebuilder-accessendpoint-vpceid
	VpceId *types.Value `json:"VpceId,omitempty"`

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
func (r *ImageBuilder_AccessEndpoint) AWSCloudFormationType() string {
	return "AWS::AppStream::ImageBuilder.AccessEndpoint"
}
