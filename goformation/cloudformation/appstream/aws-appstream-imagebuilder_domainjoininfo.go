package appstream

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ImageBuilder_DomainJoinInfo AWS CloudFormation Resource (AWS::AppStream::ImageBuilder.DomainJoinInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-domainjoininfo.html
type ImageBuilder_DomainJoinInfo struct {

	// DirectoryName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-domainjoininfo.html#cfn-appstream-imagebuilder-domainjoininfo-directoryname
	DirectoryName *types.Value `json:"DirectoryName,omitempty"`

	// OrganizationalUnitDistinguishedName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-domainjoininfo.html#cfn-appstream-imagebuilder-domainjoininfo-organizationalunitdistinguishedname
	OrganizationalUnitDistinguishedName *types.Value `json:"OrganizationalUnitDistinguishedName,omitempty"`

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
func (r *ImageBuilder_DomainJoinInfo) AWSCloudFormationType() string {
	return "AWS::AppStream::ImageBuilder.DomainJoinInfo"
}
