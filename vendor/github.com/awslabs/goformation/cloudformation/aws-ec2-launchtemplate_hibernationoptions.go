package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_HibernationOptions AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.HibernationOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-hibernationoptions.html
type AWSEC2LaunchTemplate_HibernationOptions struct {

	// Configured AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-hibernationoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-hibernationoptions-configured
	Configured *Value `json:"Configured,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_HibernationOptions) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.HibernationOptions"
}

func (r *AWSEC2LaunchTemplate_HibernationOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
