package cloudformation

import (
	"encoding/json"
)

// AWSSecretsManagerRotationSchedule_RotationRules AWS CloudFormation Resource (AWS::SecretsManager::RotationSchedule.RotationRules)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-rotationrules.html
type AWSSecretsManagerRotationSchedule_RotationRules struct {

	// AutomaticallyAfterDays AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-rotationrules.html#cfn-secretsmanager-rotationschedule-rotationrules-automaticallyafterdays
	AutomaticallyAfterDays *Value `json:"AutomaticallyAfterDays,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSecretsManagerRotationSchedule_RotationRules) AWSCloudFormationType() string {
	return "AWS::SecretsManager::RotationSchedule.RotationRules"
}

func (r *AWSSecretsManagerRotationSchedule_RotationRules) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
