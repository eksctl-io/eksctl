package cloudformation

import (
	"encoding/json"
)

// AWSGlueCrawler_JdbcTarget AWS CloudFormation Resource (AWS::Glue::Crawler.JdbcTarget)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-crawler-jdbctarget.html
type AWSGlueCrawler_JdbcTarget struct {

	// ConnectionName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-crawler-jdbctarget.html#cfn-glue-crawler-jdbctarget-connectionname
	ConnectionName *Value `json:"ConnectionName,omitempty"`

	// Exclusions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-crawler-jdbctarget.html#cfn-glue-crawler-jdbctarget-exclusions
	Exclusions []*Value `json:"Exclusions,omitempty"`

	// Path AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-crawler-jdbctarget.html#cfn-glue-crawler-jdbctarget-path
	Path *Value `json:"Path,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueCrawler_JdbcTarget) AWSCloudFormationType() string {
	return "AWS::Glue::Crawler.JdbcTarget"
}

func (r *AWSGlueCrawler_JdbcTarget) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
