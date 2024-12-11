package apigatewayv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DomainName_DomainNameConfiguration AWS CloudFormation Resource (AWS::ApiGatewayV2::DomainName.DomainNameConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html
type DomainName_DomainNameConfiguration struct {

	// CertificateArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-certificatearn
	CertificateArn *types.Value `json:"CertificateArn,omitempty"`

	// CertificateName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-certificatename
	CertificateName *types.Value `json:"CertificateName,omitempty"`

	// EndpointType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-endpointtype
	EndpointType *types.Value `json:"EndpointType,omitempty"`

	// OwnershipVerificationCertificateArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-ownershipverificationcertificatearn
	OwnershipVerificationCertificateArn *types.Value `json:"OwnershipVerificationCertificateArn,omitempty"`

	// SecurityPolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-securitypolicy
	SecurityPolicy *types.Value `json:"SecurityPolicy,omitempty"`

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
func (r *DomainName_DomainNameConfiguration) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::DomainName.DomainNameConfiguration"
}
