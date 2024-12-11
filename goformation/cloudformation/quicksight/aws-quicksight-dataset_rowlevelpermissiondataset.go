package quicksight

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSet_RowLevelPermissionDataSet AWS CloudFormation Resource (AWS::QuickSight::DataSet.RowLevelPermissionDataSet)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-rowlevelpermissiondataset.html
type DataSet_RowLevelPermissionDataSet struct {

	// Arn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-rowlevelpermissiondataset.html#cfn-quicksight-dataset-rowlevelpermissiondataset-arn
	Arn *types.Value `json:"Arn,omitempty"`

	// FormatVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-rowlevelpermissiondataset.html#cfn-quicksight-dataset-rowlevelpermissiondataset-formatversion
	FormatVersion *types.Value `json:"FormatVersion,omitempty"`

	// Namespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-rowlevelpermissiondataset.html#cfn-quicksight-dataset-rowlevelpermissiondataset-namespace
	Namespace *types.Value `json:"Namespace,omitempty"`

	// PermissionPolicy AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-rowlevelpermissiondataset.html#cfn-quicksight-dataset-rowlevelpermissiondataset-permissionpolicy
	PermissionPolicy *types.Value `json:"PermissionPolicy,omitempty"`

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
func (r *DataSet_RowLevelPermissionDataSet) AWSCloudFormationType() string {
	return "AWS::QuickSight::DataSet.RowLevelPermissionDataSet"
}
