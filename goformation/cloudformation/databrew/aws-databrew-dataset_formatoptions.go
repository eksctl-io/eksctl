package databrew

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Dataset_FormatOptions AWS CloudFormation Resource (AWS::DataBrew::Dataset.FormatOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-formatoptions.html
type Dataset_FormatOptions struct {

	// Csv AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-formatoptions.html#cfn-databrew-dataset-formatoptions-csv
	Csv *Dataset_CsvOptions `json:"Csv,omitempty"`

	// Excel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-formatoptions.html#cfn-databrew-dataset-formatoptions-excel
	Excel *Dataset_ExcelOptions `json:"Excel,omitempty"`

	// Json AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-formatoptions.html#cfn-databrew-dataset-formatoptions-json
	Json *Dataset_JsonOptions `json:"Json,omitempty"`

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
func (r *Dataset_FormatOptions) AWSCloudFormationType() string {
	return "AWS::DataBrew::Dataset.FormatOptions"
}
