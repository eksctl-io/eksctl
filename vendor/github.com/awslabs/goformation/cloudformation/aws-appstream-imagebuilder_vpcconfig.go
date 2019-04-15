package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamImageBuilder_VpcConfig AWS CloudFormation Resource (AWS::AppStream::ImageBuilder.VpcConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-vpcconfig.html
type AWSAppStreamImageBuilder_VpcConfig struct {

	// SecurityGroupIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-vpcconfig.html#cfn-appstream-imagebuilder-vpcconfig-securitygroupids
	SecurityGroupIds []*Value `json:"SecurityGroupIds,omitempty"`

	// SubnetIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-imagebuilder-vpcconfig.html#cfn-appstream-imagebuilder-vpcconfig-subnetids
	SubnetIds []*Value `json:"SubnetIds,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamImageBuilder_VpcConfig) AWSCloudFormationType() string {
	return "AWS::AppStream::ImageBuilder.VpcConfig"
}

func (r *AWSAppStreamImageBuilder_VpcConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
