package eks

import (
	"goformation/v4/cloudformation/types"
)

type AccessEntry_AccessScope struct {

	// Effect AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-taint.html#cfn-eks-nodegroup-taint-effect
	Type *types.Value `json:"Type,omitempty"`

	Namespaces *types.Value `json:"Namespaces,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AccessEntry_AccessScope) AWSCloudFormationType() string {
	return "AWS::EKS::AccessEntry.AccessScope"
}
