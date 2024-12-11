package lookoutmetrics

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AnomalyDetector_CsvFormatDescriptor AWS CloudFormation Resource (AWS::LookoutMetrics::AnomalyDetector.CsvFormatDescriptor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-csvformatdescriptor.html
type AnomalyDetector_CsvFormatDescriptor struct {

	// Charset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-csvformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-csvformatdescriptor-charset
	Charset *types.Value `json:"Charset,omitempty"`

	// ContainsHeader AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-csvformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-csvformatdescriptor-containsheader
	ContainsHeader *types.Value `json:"ContainsHeader,omitempty"`

	// Delimiter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-csvformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-csvformatdescriptor-delimiter
	Delimiter *types.Value `json:"Delimiter,omitempty"`

	// FileCompression AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-csvformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-csvformatdescriptor-filecompression
	FileCompression *types.Value `json:"FileCompression,omitempty"`

	// HeaderList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-csvformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-csvformatdescriptor-headerlist
	HeaderList *types.Value `json:"HeaderList,omitempty"`

	// QuoteSymbol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-csvformatdescriptor.html#cfn-lookoutmetrics-anomalydetector-csvformatdescriptor-quotesymbol
	QuoteSymbol *types.Value `json:"QuoteSymbol,omitempty"`

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
func (r *AnomalyDetector_CsvFormatDescriptor) AWSCloudFormationType() string {
	return "AWS::LookoutMetrics::AnomalyDetector.CsvFormatDescriptor"
}
