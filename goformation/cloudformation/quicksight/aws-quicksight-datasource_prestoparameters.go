package quicksight

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_PrestoParameters AWS CloudFormation Resource (AWS::QuickSight::DataSource.PrestoParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-prestoparameters.html
type DataSource_PrestoParameters struct {

	// Catalog AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-prestoparameters.html#cfn-quicksight-datasource-prestoparameters-catalog
	Catalog *types.Value `json:"Catalog,omitempty"`

	// Host AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-prestoparameters.html#cfn-quicksight-datasource-prestoparameters-host
	Host *types.Value `json:"Host,omitempty"`

	// Port AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-prestoparameters.html#cfn-quicksight-datasource-prestoparameters-port
	Port *types.Value `json:"Port"`

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
func (r *DataSource_PrestoParameters) AWSCloudFormationType() string {
	return "AWS::QuickSight::DataSource.PrestoParameters"
}
