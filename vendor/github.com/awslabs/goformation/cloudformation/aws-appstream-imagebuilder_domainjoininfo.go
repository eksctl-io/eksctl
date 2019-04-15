package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamImageBuilder_DomainJoinInfo AWS CloudFormation Resource (AWS::AppStream::ImageBuilder.DomainJoinInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-domainjoininfo.html
type AWSAppStreamImageBuilder_DomainJoinInfo struct {

	// DirectoryName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-domainjoininfo.html#cfn-appstream-imagebuilder-domainjoininfo-directoryname
	DirectoryName *Value `json:"DirectoryName,omitempty"`

	// OrganizationalUnitDistinguishedName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-domainjoininfo.html#cfn-appstream-imagebuilder-domainjoininfo-organizationalunitdistinguishedname
	OrganizationalUnitDistinguishedName *Value `json:"OrganizationalUnitDistinguishedName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamImageBuilder_DomainJoinInfo) AWSCloudFormationType() string {
	return "AWS::AppStream::ImageBuilder.DomainJoinInfo"
}

func (r *AWSAppStreamImageBuilder_DomainJoinInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
