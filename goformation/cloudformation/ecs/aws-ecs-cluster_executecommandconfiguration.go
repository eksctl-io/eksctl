package ecs

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Cluster_ExecuteCommandConfiguration AWS CloudFormation Resource (AWS::ECS::Cluster.ExecuteCommandConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-executecommandconfiguration.html
type Cluster_ExecuteCommandConfiguration struct {

	// KmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-executecommandconfiguration.html#cfn-ecs-cluster-executecommandconfiguration-kmskeyid
	KmsKeyId *types.Value `json:"KmsKeyId,omitempty"`

	// LogConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-executecommandconfiguration.html#cfn-ecs-cluster-executecommandconfiguration-logconfiguration
	LogConfiguration *Cluster_ExecuteCommandLogConfiguration `json:"LogConfiguration,omitempty"`

	// Logging AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-executecommandconfiguration.html#cfn-ecs-cluster-executecommandconfiguration-logging
	Logging *types.Value `json:"Logging,omitempty"`

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
func (r *Cluster_ExecuteCommandConfiguration) AWSCloudFormationType() string {
	return "AWS::ECS::Cluster.ExecuteCommandConfiguration"
}
