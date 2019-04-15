package cloudformation

import (
	"encoding/json"
)

// AWSRDSDBCluster_ScalingConfiguration AWS CloudFormation Resource (AWS::RDS::DBCluster.ScalingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbcluster-scalingconfiguration.html
type AWSRDSDBCluster_ScalingConfiguration struct {

	// AutoPause AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbcluster-scalingconfiguration.html#cfn-rds-dbcluster-scalingconfiguration-autopause
	AutoPause *Value `json:"AutoPause,omitempty"`

	// MaxCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbcluster-scalingconfiguration.html#cfn-rds-dbcluster-scalingconfiguration-maxcapacity
	MaxCapacity *Value `json:"MaxCapacity,omitempty"`

	// MinCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbcluster-scalingconfiguration.html#cfn-rds-dbcluster-scalingconfiguration-mincapacity
	MinCapacity *Value `json:"MinCapacity,omitempty"`

	// SecondsUntilAutoPause AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbcluster-scalingconfiguration.html#cfn-rds-dbcluster-scalingconfiguration-secondsuntilautopause
	SecondsUntilAutoPause *Value `json:"SecondsUntilAutoPause,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRDSDBCluster_ScalingConfiguration) AWSCloudFormationType() string {
	return "AWS::RDS::DBCluster.ScalingConfiguration"
}

func (r *AWSRDSDBCluster_ScalingConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
