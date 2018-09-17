package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_IamInstanceProfile AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.IamInstanceProfile)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-iaminstanceprofile.html
type AWSEC2LaunchTemplate_IamInstanceProfile struct {

	// Arn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-iaminstanceprofile.html#cfn-ec2-launchtemplate-launchtemplatedata-iaminstanceprofile-arn
	Arn *Value `json:"Arn,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-iaminstanceprofile.html#cfn-ec2-launchtemplate-launchtemplatedata-iaminstanceprofile-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_IamInstanceProfile) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.IamInstanceProfile"
}

func (r *AWSEC2LaunchTemplate_IamInstanceProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
