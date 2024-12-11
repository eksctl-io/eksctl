package dynamodb

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GlobalTable_ReadProvisionedThroughputSettings AWS CloudFormation Resource (AWS::DynamoDB::GlobalTable.ReadProvisionedThroughputSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-readprovisionedthroughputsettings.html
type GlobalTable_ReadProvisionedThroughputSettings struct {

	// ReadCapacityAutoScalingSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-readprovisionedthroughputsettings.html#cfn-dynamodb-globaltable-readprovisionedthroughputsettings-readcapacityautoscalingsettings
	ReadCapacityAutoScalingSettings *GlobalTable_CapacityAutoScalingSettings `json:"ReadCapacityAutoScalingSettings,omitempty"`

	// ReadCapacityUnits AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-readprovisionedthroughputsettings.html#cfn-dynamodb-globaltable-readprovisionedthroughputsettings-readcapacityunits
	ReadCapacityUnits *types.Value `json:"ReadCapacityUnits,omitempty"`

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
func (r *GlobalTable_ReadProvisionedThroughputSettings) AWSCloudFormationType() string {
	return "AWS::DynamoDB::GlobalTable.ReadProvisionedThroughputSettings"
}
