package ssm

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResourceDataSync_SyncSource AWS CloudFormation Resource (AWS::SSM::ResourceDataSync.SyncSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-syncsource.html
type ResourceDataSync_SyncSource struct {

	// AwsOrganizationsSource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-syncsource.html#cfn-ssm-resourcedatasync-syncsource-awsorganizationssource
	AwsOrganizationsSource *ResourceDataSync_AwsOrganizationsSource `json:"AwsOrganizationsSource,omitempty"`

	// IncludeFutureRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-syncsource.html#cfn-ssm-resourcedatasync-syncsource-includefutureregions
	IncludeFutureRegions *types.Value `json:"IncludeFutureRegions,omitempty"`

	// SourceRegions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-syncsource.html#cfn-ssm-resourcedatasync-syncsource-sourceregions
	SourceRegions *types.Value `json:"SourceRegions,omitempty"`

	// SourceType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-syncsource.html#cfn-ssm-resourcedatasync-syncsource-sourcetype
	SourceType *types.Value `json:"SourceType,omitempty"`

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
func (r *ResourceDataSync_SyncSource) AWSCloudFormationType() string {
	return "AWS::SSM::ResourceDataSync.SyncSource"
}
