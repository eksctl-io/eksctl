package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/types"

	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_AthenaParameters AWS CloudFormation Resource (AWS::QuickSight::DataSource.AthenaParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-athenaparameters.html
type DataSource_AthenaParameters struct {

	// WorkGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-athenaparameters.html#cfn-quicksight-datasource-athenaparameters-workgroup
	WorkGroup *types.Value `json:"WorkGroup,omitempty"`

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
func (r *DataSource_AthenaParameters) AWSCloudFormationType() string {
	return "AWS::QuickSight::DataSource.AthenaParameters"
}
