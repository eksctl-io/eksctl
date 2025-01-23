package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ModelQualityJobDefinition_ModelQualityJobInput AWS CloudFormation Resource (AWS::SageMaker::ModelQualityJobDefinition.ModelQualityJobInput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityjobinput.html
type ModelQualityJobDefinition_ModelQualityJobInput struct {

	// EndpointInput AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityjobinput.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityjobinput-endpointinput
	EndpointInput *ModelQualityJobDefinition_EndpointInput `json:"EndpointInput,omitempty"`

	// GroundTruthS3Input AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityjobinput.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityjobinput-groundtruths3input
	GroundTruthS3Input *ModelQualityJobDefinition_MonitoringGroundTruthS3Input `json:"GroundTruthS3Input,omitempty"`

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
func (r *ModelQualityJobDefinition_ModelQualityJobInput) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelQualityJobDefinition.ModelQualityJobInput"
}
