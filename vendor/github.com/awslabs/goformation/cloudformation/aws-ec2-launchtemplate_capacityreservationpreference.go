package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_CapacityReservationPreference AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.CapacityReservationPreference)
// See:
type AWSEC2LaunchTemplate_CapacityReservationPreference struct {
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_CapacityReservationPreference) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.CapacityReservationPreference"
}

func (r *AWSEC2LaunchTemplate_CapacityReservationPreference) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
