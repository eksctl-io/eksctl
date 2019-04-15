package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_CpuOptions AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.CpuOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-cpuoptions.html
type AWSEC2LaunchTemplate_CpuOptions struct {

	// CoreCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-cpuoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-cpuoptions-corecount
	CoreCount *Value `json:"CoreCount,omitempty"`

	// ThreadsPerCore AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-cpuoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-cpuoptions-threadspercore
	ThreadsPerCore *Value `json:"ThreadsPerCore,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_CpuOptions) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.CpuOptions"
}

func (r *AWSEC2LaunchTemplate_CpuOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
