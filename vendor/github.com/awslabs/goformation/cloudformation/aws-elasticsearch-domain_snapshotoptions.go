package cloudformation

import (
	"encoding/json"
)

// AWSElasticsearchDomain_SnapshotOptions AWS CloudFormation Resource (AWS::Elasticsearch::Domain.SnapshotOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-snapshotoptions.html
type AWSElasticsearchDomain_SnapshotOptions struct {

	// AutomatedSnapshotStartHour AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-snapshotoptions.html#cfn-elasticsearch-domain-snapshotoptions-automatedsnapshotstarthour
	AutomatedSnapshotStartHour *Value `json:"AutomatedSnapshotStartHour,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticsearchDomain_SnapshotOptions) AWSCloudFormationType() string {
	return "AWS::Elasticsearch::Domain.SnapshotOptions"
}

func (r *AWSElasticsearchDomain_SnapshotOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
