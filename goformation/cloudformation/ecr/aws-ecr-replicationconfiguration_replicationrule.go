package ecr

import (
	"goformation/v4/cloudformation/policies"
)

// ReplicationConfiguration_ReplicationRule AWS CloudFormation Resource (AWS::ECR::ReplicationConfiguration.ReplicationRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecr-replicationconfiguration-replicationrule.html
type ReplicationConfiguration_ReplicationRule struct {

	// Destinations AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecr-replicationconfiguration-replicationrule.html#cfn-ecr-replicationconfiguration-replicationrule-destinations
	Destinations []ReplicationConfiguration_ReplicationDestination `json:"Destinations,omitempty"`

	// RepositoryFilters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecr-replicationconfiguration-replicationrule.html#cfn-ecr-replicationconfiguration-replicationrule-repositoryfilters
	RepositoryFilters []ReplicationConfiguration_RepositoryFilter `json:"RepositoryFilters,omitempty"`

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
func (r *ReplicationConfiguration_ReplicationRule) AWSCloudFormationType() string {
	return "AWS::ECR::ReplicationConfiguration.ReplicationRule"
}
