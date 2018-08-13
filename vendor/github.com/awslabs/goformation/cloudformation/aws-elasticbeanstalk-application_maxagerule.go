package cloudformation

import (
	"encoding/json"
)

// AWSElasticBeanstalkApplication_MaxAgeRule AWS CloudFormation Resource (AWS::ElasticBeanstalk::Application.MaxAgeRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxagerule.html
type AWSElasticBeanstalkApplication_MaxAgeRule struct {

	// DeleteSourceFromS3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxagerule.html#cfn-elasticbeanstalk-application-maxagerule-deletesourcefroms3
	DeleteSourceFromS3 *Value `json:"DeleteSourceFromS3,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxagerule.html#cfn-elasticbeanstalk-application-maxagerule-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// MaxAgeInDays AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxagerule.html#cfn-elasticbeanstalk-application-maxagerule-maxageindays
	MaxAgeInDays *Value `json:"MaxAgeInDays,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticBeanstalkApplication_MaxAgeRule) AWSCloudFormationType() string {
	return "AWS::ElasticBeanstalk::Application.MaxAgeRule"
}

func (r *AWSElasticBeanstalkApplication_MaxAgeRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
