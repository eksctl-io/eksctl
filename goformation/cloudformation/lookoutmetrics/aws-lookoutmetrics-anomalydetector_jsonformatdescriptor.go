package lookoutmetrics

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AnomalyDetector_JsonFormatDescriptor AWS CloudFormation Resource (AWS::LookoutMetrics::AnomalyDetector.JsonFormatDescriptor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-jsonformatdescriptor.html
type AnomalyDetector_JsonFormatDescriptor struct {

	// Charset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-jsonformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-jsonformatdescriptor-charset
	Charset *types.Value `json:"Charset,omitempty"`

	// FileCompression AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-jsonformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-jsonformatdescriptor-filecompression
	FileCompression *types.Value `json:"FileCompression,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AnomalyDetector_JsonFormatDescriptor) AWSCloudFormationType() string {
	return "AWS::LookoutMetrics::AnomalyDetector.JsonFormatDescriptor"
}
