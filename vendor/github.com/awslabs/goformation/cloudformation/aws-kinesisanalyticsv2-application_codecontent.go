package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_CodeContent AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.CodeContent)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-codecontent.html
type AWSKinesisAnalyticsV2Application_CodeContent struct {

	// S3ContentLocation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-codecontent.html#cfn-kinesisanalyticsv2-application-codecontent-s3contentlocation
	S3ContentLocation *AWSKinesisAnalyticsV2Application_S3ContentLocation `json:"S3ContentLocation,omitempty"`

	// TextContent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-codecontent.html#cfn-kinesisanalyticsv2-application-codecontent-textcontent
	TextContent *Value `json:"TextContent,omitempty"`

	// ZipFileContent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-codecontent.html#cfn-kinesisanalyticsv2-application-codecontent-zipfilecontent
	ZipFileContent *Value `json:"ZipFileContent,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_CodeContent) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.CodeContent"
}

func (r *AWSKinesisAnalyticsV2Application_CodeContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
