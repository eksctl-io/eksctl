package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamFleet_ComputeCapacity AWS CloudFormation Resource (AWS::AppStream::Fleet.ComputeCapacity)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-computecapacity.html
type AWSAppStreamFleet_ComputeCapacity struct {

	// DesiredInstances AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-fleet-computecapacity.html#cfn-appstream-fleet-computecapacity-desiredinstances
	DesiredInstances *Value `json:"DesiredInstances,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamFleet_ComputeCapacity) AWSCloudFormationType() string {
	return "AWS::AppStream::Fleet.ComputeCapacity"
}

func (r *AWSAppStreamFleet_ComputeCapacity) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
