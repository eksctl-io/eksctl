package dynamodb

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// GlobalTable_WriteProvisionedThroughputSettings AWS CloudFormation Resource (AWS::DynamoDB::GlobalTable.WriteProvisionedThroughputSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-writeprovisionedthroughputsettings.html
type GlobalTable_WriteProvisionedThroughputSettings struct {

	// WriteCapacityAutoScalingSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-writeprovisionedthroughputsettings.html#cfn-dynamodb-globaltable-writeprovisionedthroughputsettings-writecapacityautoscalingsettings
	WriteCapacityAutoScalingSettings *GlobalTable_CapacityAutoScalingSettings `json:"WriteCapacityAutoScalingSettings,omitempty"`

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
func (r *GlobalTable_WriteProvisionedThroughputSettings) AWSCloudFormationType() string {
	return "AWS::DynamoDB::GlobalTable.WriteProvisionedThroughputSettings"
}
