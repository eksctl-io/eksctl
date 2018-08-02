package cloudformation

// AWSRedshiftClusterParameterGroup_Parameter AWS CloudFormation Resource (AWS::Redshift::ClusterParameterGroup.Parameter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-property-redshift-clusterparametergroup-parameter.html
type AWSRedshiftClusterParameterGroup_Parameter struct {

	// ParameterName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-property-redshift-clusterparametergroup-parameter.html#cfn-redshift-clusterparametergroup-parameter-parametername
	ParameterName *StringIntrinsic `json:"ParameterName,omitempty"`

	// ParameterValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-property-redshift-clusterparametergroup-parameter.html#cfn-redshift-clusterparametergroup-parameter-parametervalue
	ParameterValue *StringIntrinsic `json:"ParameterValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRedshiftClusterParameterGroup_Parameter) AWSCloudFormationType() string {
	return "AWS::Redshift::ClusterParameterGroup.Parameter"
}
