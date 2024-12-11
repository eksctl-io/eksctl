package quicksight

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSet_GeoSpatialColumnGroup AWS CloudFormation Resource (AWS::QuickSight::DataSet.GeoSpatialColumnGroup)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-geospatialcolumngroup.html
type DataSet_GeoSpatialColumnGroup struct {

	// Columns AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-geospatialcolumngroup.html#cfn-quicksight-dataset-geospatialcolumngroup-columns
	Columns *types.Value `json:"Columns,omitempty"`

	// CountryCode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-geospatialcolumngroup.html#cfn-quicksight-dataset-geospatialcolumngroup-countrycode
	CountryCode *types.Value `json:"CountryCode,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-geospatialcolumngroup.html#cfn-quicksight-dataset-geospatialcolumngroup-name
	Name *types.Value `json:"Name,omitempty"`

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
func (r *DataSet_GeoSpatialColumnGroup) AWSCloudFormationType() string {
	return "AWS::QuickSight::DataSet.GeoSpatialColumnGroup"
}
