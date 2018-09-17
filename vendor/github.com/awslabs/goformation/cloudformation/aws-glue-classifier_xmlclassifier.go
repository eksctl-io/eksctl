package cloudformation

import (
	"encoding/json"
)

// AWSGlueClassifier_XMLClassifier AWS CloudFormation Resource (AWS::Glue::Classifier.XMLClassifier)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-classifier-xmlclassifier.html
type AWSGlueClassifier_XMLClassifier struct {

	// Classification AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-classifier-xmlclassifier.html#cfn-glue-classifier-xmlclassifier-classification
	Classification *Value `json:"Classification,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-classifier-xmlclassifier.html#cfn-glue-classifier-xmlclassifier-name
	Name *Value `json:"Name,omitempty"`

	// RowTag AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-classifier-xmlclassifier.html#cfn-glue-classifier-xmlclassifier-rowtag
	RowTag *Value `json:"RowTag,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueClassifier_XMLClassifier) AWSCloudFormationType() string {
	return "AWS::Glue::Classifier.XMLClassifier"
}

func (r *AWSGlueClassifier_XMLClassifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
