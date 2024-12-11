package ssm

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResourceDataSync_AwsOrganizationsSource AWS CloudFormation Resource (AWS::SSM::ResourceDataSync.AwsOrganizationsSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-awsorganizationssource.html
type ResourceDataSync_AwsOrganizationsSource struct {

	// OrganizationSourceType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-awsorganizationssource.html#cfn-ssm-resourcedatasync-awsorganizationssource-organizationsourcetype
	OrganizationSourceType *types.Value `json:"OrganizationSourceType,omitempty"`

	// OrganizationalUnits AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-resourcedatasync-awsorganizationssource.html#cfn-ssm-resourcedatasync-awsorganizationssource-organizationalunits
	OrganizationalUnits *types.Value `json:"OrganizationalUnits,omitempty"`

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
func (r *ResourceDataSync_AwsOrganizationsSource) AWSCloudFormationType() string {
	return "AWS::SSM::ResourceDataSync.AwsOrganizationsSource"
}
