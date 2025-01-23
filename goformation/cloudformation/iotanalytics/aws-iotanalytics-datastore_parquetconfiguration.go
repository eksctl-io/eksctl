package iotanalytics

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Datastore_ParquetConfiguration AWS CloudFormation Resource (AWS::IoTAnalytics::Datastore.ParquetConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-parquetconfiguration.html
type Datastore_ParquetConfiguration struct {

	// SchemaDefinition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-parquetconfiguration.html#cfn-iotanalytics-datastore-parquetconfiguration-schemadefinition
	SchemaDefinition *Datastore_SchemaDefinition `json:"SchemaDefinition,omitempty"`

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
func (r *Datastore_ParquetConfiguration) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Datastore.ParquetConfiguration"
}
