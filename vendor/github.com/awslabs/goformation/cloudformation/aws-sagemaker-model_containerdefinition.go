package cloudformation

// AWSSageMakerModel_ContainerDefinition AWS CloudFormation Resource (AWS::SageMaker::Model.ContainerDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-model-containerdefinition.html
type AWSSageMakerModel_ContainerDefinition struct {

	// ContainerHostname AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-model-containerdefinition.html#cfn-sagemaker-model-containerdefinition-containerhostname
	ContainerHostname *StringIntrinsic `json:"ContainerHostname,omitempty"`

	// Environment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-model-containerdefinition.html#cfn-sagemaker-model-containerdefinition-environment
	Environment interface{} `json:"Environment,omitempty"`

	// Image AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-model-containerdefinition.html#cfn-sagemaker-model-containerdefinition-image
	Image *StringIntrinsic `json:"Image,omitempty"`

	// ModelDataUrl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-model-containerdefinition.html#cfn-sagemaker-model-containerdefinition-modeldataurl
	ModelDataUrl *StringIntrinsic `json:"ModelDataUrl,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSageMakerModel_ContainerDefinition) AWSCloudFormationType() string {
	return "AWS::SageMaker::Model.ContainerDefinition"
}
