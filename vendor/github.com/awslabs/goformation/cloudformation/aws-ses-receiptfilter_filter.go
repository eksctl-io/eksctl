package cloudformation

// AWSSESReceiptFilter_Filter AWS CloudFormation Resource (AWS::SES::ReceiptFilter.Filter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptfilter-filter.html
type AWSSESReceiptFilter_Filter struct {

	// IpFilter AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptfilter-filter.html#cfn-ses-receiptfilter-filter-ipfilter
	IpFilter *AWSSESReceiptFilter_IpFilter `json:"IpFilter,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptfilter-filter.html#cfn-ses-receiptfilter-filter-name
	Name *StringIntrinsic `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptFilter_Filter) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptFilter.Filter"
}
