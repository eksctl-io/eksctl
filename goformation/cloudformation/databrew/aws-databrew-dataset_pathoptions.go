package databrew

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Dataset_PathOptions AWS CloudFormation Resource (AWS::DataBrew::Dataset.PathOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-pathoptions.html
type Dataset_PathOptions struct {

	// FilesLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-pathoptions.html#cfn-databrew-dataset-pathoptions-fileslimit
	FilesLimit *Dataset_FilesLimit `json:"FilesLimit,omitempty"`

	// LastModifiedDateCondition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-pathoptions.html#cfn-databrew-dataset-pathoptions-lastmodifieddatecondition
	LastModifiedDateCondition *Dataset_FilterExpression `json:"LastModifiedDateCondition,omitempty"`

	// Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-pathoptions.html#cfn-databrew-dataset-pathoptions-parameters
	Parameters []Dataset_PathParameter `json:"Parameters,omitempty"`

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
func (r *Dataset_PathOptions) AWSCloudFormationType() string {
	return "AWS::DataBrew::Dataset.PathOptions"
}
