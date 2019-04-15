package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamFleet_DomainJoinInfo AWS CloudFormation Resource (AWS::AppStream::Fleet.DomainJoinInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-domainjoininfo.html
type AWSAppStreamFleet_DomainJoinInfo struct {

	// DirectoryName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-domainjoininfo.html#cfn-appstream-fleet-domainjoininfo-directoryname
	DirectoryName *Value `json:"DirectoryName,omitempty"`

	// OrganizationalUnitDistinguishedName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-domainjoininfo.html#cfn-appstream-fleet-domainjoininfo-organizationalunitdistinguishedname
	OrganizationalUnitDistinguishedName *Value `json:"OrganizationalUnitDistinguishedName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamFleet_DomainJoinInfo) AWSCloudFormationType() string {
	return "AWS::AppStream::Fleet.DomainJoinInfo"
}

func (r *AWSAppStreamFleet_DomainJoinInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
