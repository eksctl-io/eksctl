package cloudformation

import (
	"encoding/json"
)

// AWSGlueTable_StorageDescriptor AWS CloudFormation Resource (AWS::Glue::Table.StorageDescriptor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html
type AWSGlueTable_StorageDescriptor struct {

	// BucketColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-bucketcolumns
	BucketColumns []*Value `json:"BucketColumns,omitempty"`

	// Columns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-columns
	Columns []AWSGlueTable_Column `json:"Columns,omitempty"`

	// Compressed AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-compressed
	Compressed *Value `json:"Compressed,omitempty"`

	// InputFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-inputformat
	InputFormat *Value `json:"InputFormat,omitempty"`

	// Location AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-location
	Location *Value `json:"Location,omitempty"`

	// NumberOfBuckets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-numberofbuckets
	NumberOfBuckets *Value `json:"NumberOfBuckets,omitempty"`

	// OutputFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-outputformat
	OutputFormat *Value `json:"OutputFormat,omitempty"`

	// Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-parameters
	Parameters interface{} `json:"Parameters,omitempty"`

	// SerdeInfo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-serdeinfo
	SerdeInfo *AWSGlueTable_SerdeInfo `json:"SerdeInfo,omitempty"`

	// SkewedInfo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-skewedinfo
	SkewedInfo *AWSGlueTable_SkewedInfo `json:"SkewedInfo,omitempty"`

	// SortColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-sortcolumns
	SortColumns []AWSGlueTable_Order `json:"SortColumns,omitempty"`

	// StoredAsSubDirectories AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-storagedescriptor.html#cfn-glue-table-storagedescriptor-storedassubdirectories
	StoredAsSubDirectories *Value `json:"StoredAsSubDirectories,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueTable_StorageDescriptor) AWSCloudFormationType() string {
	return "AWS::Glue::Table.StorageDescriptor"
}

func (r *AWSGlueTable_StorageDescriptor) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
