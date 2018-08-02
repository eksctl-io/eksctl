package cloudformation

// AWSSageMakerEndpointConfig_ProductionVariant AWS CloudFormation Resource (AWS::SageMaker::EndpointConfig.ProductionVariant)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html
type AWSSageMakerEndpointConfig_ProductionVariant struct {

	// InitialInstanceCount AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-initialinstancecount
	InitialInstanceCount int `json:"InitialInstanceCount,omitempty"`

	// InitialVariantWeight AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-initialvariantweight
	InitialVariantWeight float64 `json:"InitialVariantWeight,omitempty"`

	// InstanceType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-instancetype
	InstanceType *StringIntrinsic `json:"InstanceType,omitempty"`

	// ModelName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-modelname
	ModelName *StringIntrinsic `json:"ModelName,omitempty"`

	// VariantName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-variantname
	VariantName *StringIntrinsic `json:"VariantName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSageMakerEndpointConfig_ProductionVariant) AWSCloudFormationType() string {
	return "AWS::SageMaker::EndpointConfig.ProductionVariant"
}
