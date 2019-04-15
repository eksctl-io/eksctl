package cloudformation

import (
	"encoding/json"
)

// AWSEC2SpotFleet_TargetGroupsConfig AWS CloudFormation Resource (AWS::EC2::SpotFleet.TargetGroupsConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-targetgroupsconfig.html
type AWSEC2SpotFleet_TargetGroupsConfig struct {

	// TargetGroups AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-targetgroupsconfig.html#cfn-ec2-spotfleet-targetgroupsconfig-targetgroups
	TargetGroups []AWSEC2SpotFleet_TargetGroup `json:"TargetGroups,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_TargetGroupsConfig) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.TargetGroupsConfig"
}

func (r *AWSEC2SpotFleet_TargetGroupsConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
