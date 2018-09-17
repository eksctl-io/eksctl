package cloudformation

import (
	"encoding/json"
)

// AWSCodeCommitRepository_RepositoryTrigger AWS CloudFormation Resource (AWS::CodeCommit::Repository.RepositoryTrigger)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codecommit-repository-repositorytrigger.html
type AWSCodeCommitRepository_RepositoryTrigger struct {

	// Branches AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codecommit-repository-repositorytrigger.html#cfn-codecommit-repository-repositorytrigger-branches
	Branches []*Value `json:"Branches,omitempty"`

	// CustomData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codecommit-repository-repositorytrigger.html#cfn-codecommit-repository-repositorytrigger-customdata
	CustomData *Value `json:"CustomData,omitempty"`

	// DestinationArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codecommit-repository-repositorytrigger.html#cfn-codecommit-repository-repositorytrigger-destinationarn
	DestinationArn *Value `json:"DestinationArn,omitempty"`

	// Events AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codecommit-repository-repositorytrigger.html#cfn-codecommit-repository-repositorytrigger-events
	Events []*Value `json:"Events,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codecommit-repository-repositorytrigger.html#cfn-codecommit-repository-repositorytrigger-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeCommitRepository_RepositoryTrigger) AWSCloudFormationType() string {
	return "AWS::CodeCommit::Repository.RepositoryTrigger"
}

func (r *AWSCodeCommitRepository_RepositoryTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
