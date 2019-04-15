package cloudformation

import (
	"encoding/json"
)

// AWSDMSEndpoint_ElasticsearchSettings AWS CloudFormation Resource (AWS::DMS::Endpoint.ElasticsearchSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-elasticsearchsettings.html
type AWSDMSEndpoint_ElasticsearchSettings struct {

	// EndpointUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-elasticsearchsettings.html#cfn-dms-endpoint-elasticsearchsettings-endpointuri
	EndpointUri *Value `json:"EndpointUri,omitempty"`

	// ErrorRetryDuration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-elasticsearchsettings.html#cfn-dms-endpoint-elasticsearchsettings-errorretryduration
	ErrorRetryDuration *Value `json:"ErrorRetryDuration,omitempty"`

	// FullLoadErrorPercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-elasticsearchsettings.html#cfn-dms-endpoint-elasticsearchsettings-fullloaderrorpercentage
	FullLoadErrorPercentage *Value `json:"FullLoadErrorPercentage,omitempty"`

	// ServiceAccessRoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-elasticsearchsettings.html#cfn-dms-endpoint-elasticsearchsettings-serviceaccessrolearn
	ServiceAccessRoleArn *Value `json:"ServiceAccessRoleArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDMSEndpoint_ElasticsearchSettings) AWSCloudFormationType() string {
	return "AWS::DMS::Endpoint.ElasticsearchSettings"
}

func (r *AWSDMSEndpoint_ElasticsearchSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
