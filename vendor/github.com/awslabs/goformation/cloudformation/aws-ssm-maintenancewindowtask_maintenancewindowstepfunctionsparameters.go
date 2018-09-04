package cloudformation

// AWSSSMMaintenanceWindowTask_MaintenanceWindowStepFunctionsParameters AWS CloudFormation Resource (AWS::SSM::MaintenanceWindowTask.MaintenanceWindowStepFunctionsParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-maintenancewindowstepfunctionsparameters.html
type AWSSSMMaintenanceWindowTask_MaintenanceWindowStepFunctionsParameters struct {

	// Input AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-maintenancewindowstepfunctionsparameters.html#cfn-ssm-maintenancewindowtask-maintenancewindowstepfunctionsparameters-input
	Input *StringIntrinsic `json:"Input,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-maintenancewindowstepfunctionsparameters.html#cfn-ssm-maintenancewindowtask-maintenancewindowstepfunctionsparameters-name
	Name *StringIntrinsic `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMMaintenanceWindowTask_MaintenanceWindowStepFunctionsParameters) AWSCloudFormationType() string {
	return "AWS::SSM::MaintenanceWindowTask.MaintenanceWindowStepFunctionsParameters"
}
