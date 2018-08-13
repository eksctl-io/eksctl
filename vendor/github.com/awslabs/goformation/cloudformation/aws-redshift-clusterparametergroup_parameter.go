package cloudformation

import (
	"encoding/json"
)

// AWSRedshiftClusterParameterGroup_Parameter AWS CloudFormation Resource (AWS::Redshift::ClusterParameterGroup.Parameter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-property-redshift-clusterparametergroup-parameter.html
type AWSRedshiftClusterParameterGroup_Parameter struct {

	// ParameterName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-property-redshift-clusterparametergroup-parameter.html#cfn-redshift-clusterparametergroup-parameter-parametername
	ParameterName *Value `json:"ParameterName,omitempty"`

	// ParameterValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-property-redshift-clusterparametergroup-parameter.html#cfn-redshift-clusterparametergroup-parameter-parametervalue
	ParameterValue *Value `json:"ParameterValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRedshiftClusterParameterGroup_Parameter) AWSCloudFormationType() string {
	return "AWS::Redshift::ClusterParameterGroup.Parameter"
}

func (r *AWSRedshiftClusterParameterGroup_Parameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
