package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_RepositoryCredentials AWS CloudFormation Resource (AWS::ECS::TaskDefinition.RepositoryCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-repositorycredentials.html
type AWSECSTaskDefinition_RepositoryCredentials struct {

	// CredentialsParameter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-repositorycredentials.html#cfn-ecs-taskdefinition-repositorycredentials-credentialsparameter
	CredentialsParameter *Value `json:"CredentialsParameter,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_RepositoryCredentials) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.RepositoryCredentials"
}

func (r *AWSECSTaskDefinition_RepositoryCredentials) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
