package cloudformation

import (
	"encoding/json"
)

// AWSEMRCluster_PlacementType AWS CloudFormation Resource (AWS::EMR::Cluster.PlacementType)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-placementtype.html
type AWSEMRCluster_PlacementType struct {

	// AvailabilityZone AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-placementtype.html#cfn-elasticmapreduce-cluster-placementtype-availabilityzone
	AvailabilityZone *Value `json:"AvailabilityZone,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRCluster_PlacementType) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.PlacementType"
}

func (r *AWSEMRCluster_PlacementType) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
