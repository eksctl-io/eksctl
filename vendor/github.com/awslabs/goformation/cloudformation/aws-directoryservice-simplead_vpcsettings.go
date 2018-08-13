package cloudformation

import (
	"encoding/json"
)

// AWSDirectoryServiceSimpleAD_VpcSettings AWS CloudFormation Resource (AWS::DirectoryService::SimpleAD.VpcSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-directoryservice-simplead-vpcsettings.html
type AWSDirectoryServiceSimpleAD_VpcSettings struct {

	// SubnetIds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-directoryservice-simplead-vpcsettings.html#cfn-directoryservice-simplead-vpcsettings-subnetids
	SubnetIds []*Value `json:"SubnetIds,omitempty"`

	// VpcId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-directoryservice-simplead-vpcsettings.html#cfn-directoryservice-simplead-vpcsettings-vpcid
	VpcId *Value `json:"VpcId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDirectoryServiceSimpleAD_VpcSettings) AWSCloudFormationType() string {
	return "AWS::DirectoryService::SimpleAD.VpcSettings"
}

func (r *AWSDirectoryServiceSimpleAD_VpcSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
