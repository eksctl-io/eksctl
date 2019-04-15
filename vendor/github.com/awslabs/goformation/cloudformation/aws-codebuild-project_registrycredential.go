package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_RegistryCredential AWS CloudFormation Resource (AWS::CodeBuild::Project.RegistryCredential)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-registrycredential.html
type AWSCodeBuildProject_RegistryCredential struct {

	// Credential AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-registrycredential.html#cfn-codebuild-project-registrycredential-credential
	Credential *Value `json:"Credential,omitempty"`

	// CredentialProvider AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-registrycredential.html#cfn-codebuild-project-registrycredential-credentialprovider
	CredentialProvider *Value `json:"CredentialProvider,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_RegistryCredential) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.RegistryCredential"
}

func (r *AWSCodeBuildProject_RegistryCredential) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
