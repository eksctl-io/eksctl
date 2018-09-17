package cloudformation

import (
	"encoding/json"
)

// AWSElasticBeanstalkApplication_MaxCountRule AWS CloudFormation Resource (AWS::ElasticBeanstalk::Application.MaxCountRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxcountrule.html
type AWSElasticBeanstalkApplication_MaxCountRule struct {

	// DeleteSourceFromS3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxcountrule.html#cfn-elasticbeanstalk-application-maxcountrule-deletesourcefroms3
	DeleteSourceFromS3 *Value `json:"DeleteSourceFromS3,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxcountrule.html#cfn-elasticbeanstalk-application-maxcountrule-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// MaxCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticbeanstalk-application-maxcountrule.html#cfn-elasticbeanstalk-application-maxcountrule-maxcount
	MaxCount *Value `json:"MaxCount,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticBeanstalkApplication_MaxCountRule) AWSCloudFormationType() string {
	return "AWS::ElasticBeanstalk::Application.MaxCountRule"
}

func (r *AWSElasticBeanstalkApplication_MaxCountRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
