package cloudformation

// AWSAmazonMQBroker_MaintenanceWindow AWS CloudFormation Resource (AWS::AmazonMQ::Broker.MaintenanceWindow)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-maintenancewindow.html
type AWSAmazonMQBroker_MaintenanceWindow struct {

	// DayOfWeek AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-maintenancewindow.html#cfn-amazonmq-broker-maintenancewindow-dayofweek
	DayOfWeek *StringIntrinsic `json:"DayOfWeek,omitempty"`

	// TimeOfDay AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-maintenancewindow.html#cfn-amazonmq-broker-maintenancewindow-timeofday
	TimeOfDay *StringIntrinsic `json:"TimeOfDay,omitempty"`

	// TimeZone AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-maintenancewindow.html#cfn-amazonmq-broker-maintenancewindow-timezone
	TimeZone *StringIntrinsic `json:"TimeZone,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAmazonMQBroker_MaintenanceWindow) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Broker.MaintenanceWindow"
}
