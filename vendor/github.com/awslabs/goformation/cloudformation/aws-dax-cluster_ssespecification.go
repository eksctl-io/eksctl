package cloudformation

// AWSDAXCluster_SSESpecification AWS CloudFormation Resource (AWS::DAX::Cluster.SSESpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dax-cluster-ssespecification.html
type AWSDAXCluster_SSESpecification struct {

	// SSEEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dax-cluster-ssespecification.html#cfn-dax-cluster-ssespecification-sseenabled
	SSEEnabled bool `json:"SSEEnabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDAXCluster_SSESpecification) AWSCloudFormationType() string {
	return "AWS::DAX::Cluster.SSESpecification"
}
