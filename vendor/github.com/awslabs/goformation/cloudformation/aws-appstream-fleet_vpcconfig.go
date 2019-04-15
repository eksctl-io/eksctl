package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamFleet_VpcConfig AWS CloudFormation Resource (AWS::AppStream::Fleet.VpcConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-vpcconfig.html
type AWSAppStreamFleet_VpcConfig struct {

	// SecurityGroupIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-vpcconfig.html#cfn-appstream-fleet-vpcconfig-securitygroupids
	SecurityGroupIds []*Value `json:"SecurityGroupIds,omitempty"`

	// SubnetIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-vpcconfig.html#cfn-appstream-fleet-vpcconfig-subnetids
	SubnetIds []*Value `json:"SubnetIds,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamFleet_VpcConfig) AWSCloudFormationType() string {
	return "AWS::AppStream::Fleet.VpcConfig"
}

func (r *AWSAppStreamFleet_VpcConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
