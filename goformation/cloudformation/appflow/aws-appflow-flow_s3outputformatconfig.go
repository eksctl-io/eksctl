package appflow

import (
	"github.com/awslabs/goformation/v4/cloudformation/types"

	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Flow_S3OutputFormatConfig AWS CloudFormation Resource (AWS::AppFlow::Flow.S3OutputFormatConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3outputformatconfig.html
type Flow_S3OutputFormatConfig struct {

	// AggregationConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3outputformatconfig.html#cfn-appflow-flow-s3outputformatconfig-aggregationconfig
	AggregationConfig *Flow_AggregationConfig `json:"AggregationConfig,omitempty"`

	// FileType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3outputformatconfig.html#cfn-appflow-flow-s3outputformatconfig-filetype
	FileType *types.Value `json:"FileType,omitempty"`

	// PrefixConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3outputformatconfig.html#cfn-appflow-flow-s3outputformatconfig-prefixconfig
	PrefixConfig *Flow_PrefixConfig `json:"PrefixConfig,omitempty"`

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
func (r *Flow_S3OutputFormatConfig) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.S3OutputFormatConfig"
}
