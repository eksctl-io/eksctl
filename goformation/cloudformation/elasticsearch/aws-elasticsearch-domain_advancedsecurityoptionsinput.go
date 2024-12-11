package elasticsearch

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Domain_AdvancedSecurityOptionsInput AWS CloudFormation Resource (AWS::Elasticsearch::Domain.AdvancedSecurityOptionsInput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-advancedsecurityoptionsinput.html
type Domain_AdvancedSecurityOptionsInput struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-advancedsecurityoptionsinput.html#cfn-elasticsearch-domain-advancedsecurityoptionsinput-enabled
	Enabled *types.Value `json:"Enabled,omitempty"`

	// InternalUserDatabaseEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-advancedsecurityoptionsinput.html#cfn-elasticsearch-domain-advancedsecurityoptionsinput-internaluserdatabaseenabled
	InternalUserDatabaseEnabled *types.Value `json:"InternalUserDatabaseEnabled,omitempty"`

	// MasterUserOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-advancedsecurityoptionsinput.html#cfn-elasticsearch-domain-advancedsecurityoptionsinput-masteruseroptions
	MasterUserOptions *Domain_MasterUserOptions `json:"MasterUserOptions,omitempty"`

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
func (r *Domain_AdvancedSecurityOptionsInput) AWSCloudFormationType() string {
	return "AWS::Elasticsearch::Domain.AdvancedSecurityOptionsInput"
}
