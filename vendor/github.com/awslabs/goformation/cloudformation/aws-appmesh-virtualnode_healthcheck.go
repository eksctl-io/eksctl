package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_HealthCheck AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.HealthCheck)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html
type AWSAppMeshVirtualNode_HealthCheck struct {

	// HealthyThreshold AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html#cfn-appmesh-virtualnode-healthcheck-healthythreshold
	HealthyThreshold *Value `json:"HealthyThreshold,omitempty"`

	// IntervalMillis AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html#cfn-appmesh-virtualnode-healthcheck-intervalmillis
	IntervalMillis *Value `json:"IntervalMillis,omitempty"`

	// Path AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html#cfn-appmesh-virtualnode-healthcheck-path
	Path *Value `json:"Path,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html#cfn-appmesh-virtualnode-healthcheck-port
	Port *Value `json:"Port,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html#cfn-appmesh-virtualnode-healthcheck-protocol
	Protocol *Value `json:"Protocol,omitempty"`

	// TimeoutMillis AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html#cfn-appmesh-virtualnode-healthcheck-timeoutmillis
	TimeoutMillis *Value `json:"TimeoutMillis,omitempty"`

	// UnhealthyThreshold AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-healthcheck.html#cfn-appmesh-virtualnode-healthcheck-unhealthythreshold
	UnhealthyThreshold *Value `json:"UnhealthyThreshold,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_HealthCheck) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.HealthCheck"
}

func (r *AWSAppMeshVirtualNode_HealthCheck) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
