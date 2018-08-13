package cloudformation

import (
	"encoding/json"
)

// AWSGameLiftAlias_RoutingStrategy AWS CloudFormation Resource (AWS::GameLift::Alias.RoutingStrategy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-alias-routingstrategy.html
type AWSGameLiftAlias_RoutingStrategy struct {

	// FleetId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-alias-routingstrategy.html#cfn-gamelift-alias-routingstrategy-fleetid
	FleetId *Value `json:"FleetId,omitempty"`

	// Message AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-alias-routingstrategy.html#cfn-gamelift-alias-routingstrategy-message
	Message *Value `json:"Message,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-alias-routingstrategy.html#cfn-gamelift-alias-routingstrategy-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGameLiftAlias_RoutingStrategy) AWSCloudFormationType() string {
	return "AWS::GameLift::Alias.RoutingStrategy"
}

func (r *AWSGameLiftAlias_RoutingStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
