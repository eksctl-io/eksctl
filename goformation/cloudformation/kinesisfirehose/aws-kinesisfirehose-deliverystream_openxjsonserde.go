package kinesisfirehose

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DeliveryStream_OpenXJsonSerDe AWS CloudFormation Resource (AWS::KinesisFirehose::DeliveryStream.OpenXJsonSerDe)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-openxjsonserde.html
type DeliveryStream_OpenXJsonSerDe struct {

	// CaseInsensitive AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-openxjsonserde.html#cfn-kinesisfirehose-deliverystream-openxjsonserde-caseinsensitive
	CaseInsensitive *types.Value `json:"CaseInsensitive,omitempty"`

	// ColumnToJsonKeyMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-openxjsonserde.html#cfn-kinesisfirehose-deliverystream-openxjsonserde-columntojsonkeymappings
	ColumnToJsonKeyMappings map[string]*types.Value `json:"ColumnToJsonKeyMappings,omitempty"`

	// ConvertDotsInJsonKeysToUnderscores AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-openxjsonserde.html#cfn-kinesisfirehose-deliverystream-openxjsonserde-convertdotsinjsonkeystounderscores
	ConvertDotsInJsonKeysToUnderscores *types.Value `json:"ConvertDotsInJsonKeysToUnderscores,omitempty"`

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
func (r *DeliveryStream_OpenXJsonSerDe) AWSCloudFormationType() string {
	return "AWS::KinesisFirehose::DeliveryStream.OpenXJsonSerDe"
}
