package iotanalytics

import (
	"goformation/v4/cloudformation/policies"
)

// Datastore_DatastorePartitions AWS CloudFormation Resource (AWS::IoTAnalytics::Datastore.DatastorePartitions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-datastorepartitions.html
type Datastore_DatastorePartitions struct {

	// Partitions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-datastorepartitions.html#cfn-iotanalytics-datastore-datastorepartitions-partitions
	Partitions []Datastore_DatastorePartition `json:"Partitions,omitempty"`

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
func (r *Datastore_DatastorePartitions) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Datastore.DatastorePartitions"
}
