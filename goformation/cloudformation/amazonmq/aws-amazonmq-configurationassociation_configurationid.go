package amazonmq

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConfigurationAssociation_ConfigurationId AWS CloudFormation Resource (AWS::AmazonMQ::ConfigurationAssociation.ConfigurationId)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configurationassociation-configurationid.html
type ConfigurationAssociation_ConfigurationId struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configurationassociation-configurationid.html#cfn-amazonmq-configurationassociation-configurationid-id
	Id *types.Value `json:"Id,omitempty"`

	// Revision AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configurationassociation-configurationid.html#cfn-amazonmq-configurationassociation-configurationid-revision
	Revision *types.Value `json:"Revision"`

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
func (r *ConfigurationAssociation_ConfigurationId) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::ConfigurationAssociation.ConfigurationId"
}
