package athena

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WorkGroup_WorkGroupConfigurationUpdates AWS CloudFormation Resource (AWS::Athena::WorkGroup.WorkGroupConfigurationUpdates)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html
type WorkGroup_WorkGroupConfigurationUpdates struct {

	// BytesScannedCutoffPerQuery AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html#cfn-athena-workgroup-workgroupconfigurationupdates-bytesscannedcutoffperquery
	BytesScannedCutoffPerQuery *types.Value `json:"BytesScannedCutoffPerQuery,omitempty"`

	// EnforceWorkGroupConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html#cfn-athena-workgroup-workgroupconfigurationupdates-enforceworkgroupconfiguration
	EnforceWorkGroupConfiguration *types.Value `json:"EnforceWorkGroupConfiguration,omitempty"`

	// EngineVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html#cfn-athena-workgroup-workgroupconfigurationupdates-engineversion
	EngineVersion *WorkGroup_EngineVersion `json:"EngineVersion,omitempty"`

	// PublishCloudWatchMetricsEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html#cfn-athena-workgroup-workgroupconfigurationupdates-publishcloudwatchmetricsenabled
	PublishCloudWatchMetricsEnabled *types.Value `json:"PublishCloudWatchMetricsEnabled,omitempty"`

	// RemoveBytesScannedCutoffPerQuery AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html#cfn-athena-workgroup-workgroupconfigurationupdates-removebytesscannedcutoffperquery
	RemoveBytesScannedCutoffPerQuery *types.Value `json:"RemoveBytesScannedCutoffPerQuery,omitempty"`

	// RequesterPaysEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html#cfn-athena-workgroup-workgroupconfigurationupdates-requesterpaysenabled
	RequesterPaysEnabled *types.Value `json:"RequesterPaysEnabled,omitempty"`

	// ResultConfigurationUpdates AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-workgroup-workgroupconfigurationupdates.html#cfn-athena-workgroup-workgroupconfigurationupdates-resultconfigurationupdates
	ResultConfigurationUpdates *WorkGroup_ResultConfigurationUpdates `json:"ResultConfigurationUpdates,omitempty"`

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
func (r *WorkGroup_WorkGroupConfigurationUpdates) AWSCloudFormationType() string {
	return "AWS::Athena::WorkGroup.WorkGroupConfigurationUpdates"
}
