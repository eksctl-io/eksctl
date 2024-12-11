package iotanalytics

import (
	"goformation/v4/cloudformation/policies"
)

// Datastore_DatastorePartition AWS CloudFormation Resource (AWS::IoTAnalytics::Datastore.DatastorePartition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-datastorepartition.html
type Datastore_DatastorePartition struct {

	// Partition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-datastorepartition.html#cfn-iotanalytics-datastore-datastorepartition-partition
	Partition *Datastore_Partition `json:"Partition,omitempty"`

	// TimestampPartition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-datastorepartition.html#cfn-iotanalytics-datastore-datastorepartition-timestamppartition
	TimestampPartition *Datastore_TimestampPartition `json:"TimestampPartition,omitempty"`

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
func (r *Datastore_DatastorePartition) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Datastore.DatastorePartition"
}
