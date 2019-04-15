package cloudformation

import (
	"encoding/json"
)

// AWSSageMakerEndpointConfig_ProductionVariant AWS CloudFormation Resource (AWS::SageMaker::EndpointConfig.ProductionVariant)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html
type AWSSageMakerEndpointConfig_ProductionVariant struct {

	// AcceleratorType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-acceleratortype
	AcceleratorType *Value `json:"AcceleratorType,omitempty"`

	// InitialInstanceCount AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-initialinstancecount
	InitialInstanceCount *Value `json:"InitialInstanceCount,omitempty"`

	// InitialVariantWeight AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-initialvariantweight
	InitialVariantWeight *Value `json:"InitialVariantWeight,omitempty"`

	// InstanceType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-instancetype
	InstanceType *Value `json:"InstanceType,omitempty"`

	// ModelName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-modelname
	ModelName *Value `json:"ModelName,omitempty"`

	// VariantName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant.html#cfn-sagemaker-endpointconfig-productionvariant-variantname
	VariantName *Value `json:"VariantName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSageMakerEndpointConfig_ProductionVariant) AWSCloudFormationType() string {
	return "AWS::SageMaker::EndpointConfig.ProductionVariant"
}

func (r *AWSSageMakerEndpointConfig_ProductionVariant) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
