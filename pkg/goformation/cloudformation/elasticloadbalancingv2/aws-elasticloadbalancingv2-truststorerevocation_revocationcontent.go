package elasticloadbalancingv2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// TrustStoreRevocation_RevocationContent AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::TrustStoreRevocation.RevocationContent)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-truststorerevocation-revocationcontent.html
type TrustStoreRevocation_RevocationContent struct {

	// RevocationType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-truststorerevocation-revocationcontent.html#cfn-elasticloadbalancingv2-truststorerevocation-revocationcontent-revocationtype
	RevocationType *types.Value `json:"RevocationType,omitempty"`

	// S3Bucket AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-truststorerevocation-revocationcontent.html#cfn-elasticloadbalancingv2-truststorerevocation-revocationcontent-s3bucket
	S3Bucket *types.Value `json:"S3Bucket,omitempty"`

	// S3Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-truststorerevocation-revocationcontent.html#cfn-elasticloadbalancingv2-truststorerevocation-revocationcontent-s3key
	S3Key *types.Value `json:"S3Key,omitempty"`

	// S3ObjectVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-truststorerevocation-revocationcontent.html#cfn-elasticloadbalancingv2-truststorerevocation-revocationcontent-s3objectversion
	S3ObjectVersion *types.Value `json:"S3ObjectVersion,omitempty"`

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
func (r *TrustStoreRevocation_RevocationContent) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::TrustStoreRevocation.RevocationContent"
}
