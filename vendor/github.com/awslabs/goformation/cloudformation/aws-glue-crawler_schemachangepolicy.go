package cloudformation

import (
	"encoding/json"
)

// AWSGlueCrawler_SchemaChangePolicy AWS CloudFormation Resource (AWS::Glue::Crawler.SchemaChangePolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-crawler-schemachangepolicy.html
type AWSGlueCrawler_SchemaChangePolicy struct {

	// DeleteBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-crawler-schemachangepolicy.html#cfn-glue-crawler-schemachangepolicy-deletebehavior
	DeleteBehavior *Value `json:"DeleteBehavior,omitempty"`

	// UpdateBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-crawler-schemachangepolicy.html#cfn-glue-crawler-schemachangepolicy-updatebehavior
	UpdateBehavior *Value `json:"UpdateBehavior,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueCrawler_SchemaChangePolicy) AWSCloudFormationType() string {
	return "AWS::Glue::Crawler.SchemaChangePolicy"
}

func (r *AWSGlueCrawler_SchemaChangePolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
