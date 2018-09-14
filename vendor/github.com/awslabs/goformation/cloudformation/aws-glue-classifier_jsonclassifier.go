package cloudformation

// AWSGlueClassifier_JsonClassifier AWS CloudFormation Resource (AWS::Glue::Classifier.JsonClassifier)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-classifier-jsonclassifier.html
type AWSGlueClassifier_JsonClassifier struct {

	// JsonPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-classifier-jsonclassifier.html#cfn-glue-classifier-jsonclassifier-jsonpath
	JsonPath string `json:"JsonPath,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-classifier-jsonclassifier.html#cfn-glue-classifier-jsonclassifier-name
	Name string `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueClassifier_JsonClassifier) AWSCloudFormationType() string {
	return "AWS::Glue::Classifier.JsonClassifier"
}
