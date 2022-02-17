package ami

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
)

// SSMResolver resolves the AMI to the defaults for the region
// by querying AWS SSM get parameter API
type SSMResolver struct {
	ssmAPI ssmiface.SSMAPI
}

// Resolve will return an AMI to use based on the default AMI for
// each region
func (r *SSMResolver) Resolve(region, version, instanceType, imageFamily string) (string, error) {
	logger.Debug("resolving AMI using SSM Parameter resolver for region %s, instanceType %s and imageFamily %s", region, instanceType, imageFamily)

	parameterName, err := MakeSSMParameterName(version, instanceType, imageFamily)
	if err != nil {
		return "", err
	}
	input := ssm.GetParameterInput{
		Name: aws.String(parameterName),
	}
	output, err := r.ssmAPI.GetParameter(&input)
	if err != nil {
		return "", fmt.Errorf("error getting AMI from SSM Parameter Store: %w. please verify that AMI Family is supported", err)
	}

	if output == nil || output.Parameter == nil || *output.Parameter.Value == "" {
		return "", NewErrFailedResolution(region, version, instanceType, imageFamily)
	}

	return *output.Parameter.Value, nil
}

// MakeSSMParameterName creates an SSM parameter name
func MakeSSMParameterName(version, instanceType, imageFamily string) (string, error) {
	if api.IsWindowsImage(imageFamily) {
		if supportsWindows, err := utils.IsMinVersion(api.Version1_14, version); err != nil {
			return "", err
		} else if !supportsWindows {
			return "", fmt.Errorf("cannot find Windows AMI for Kubernetes version %s. Minimum version supported: %s", version, api.Version1_14)
		}
	}

	const fieldName = "image_id"

	switch imageFamily {
	case api.NodeImageFamilyAmazonLinux2:
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/recommended/%s", version, imageType(imageFamily, instanceType, version), fieldName), nil
	case api.NodeImageFamilyWindowsServer2019CoreContainer:
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-2019-English-Core-EKS_Optimized-%s/%s", version, fieldName), nil
	case api.NodeImageFamilyWindowsServer2019FullContainer:
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-2019-English-Full-EKS_Optimized-%s/%s", version, fieldName), nil
	case api.NodeImageFamilyWindowsServer2004CoreContainer:
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-2004-English-Core-EKS_Optimized-%s/%s", version, fieldName), nil
	case api.NodeImageFamilyWindowsServer20H2CoreContainer:
		const minVersion = api.Version1_21
		supportsWindows20H2, err := utils.IsMinVersion(minVersion, version)
		if err != nil {
			return "", err
		}
		if !supportsWindows20H2 {
			return "", errors.Errorf("Windows Server 20H2 Core requires EKS version %s and above", minVersion)
		}
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-20H2-English-Core-EKS_Optimized-%s/%s", version, fieldName), nil
	case api.NodeImageFamilyBottlerocket:
		return fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/%s/latest/%s", imageType(imageFamily, instanceType, version), instanceEC2ArchName(instanceType), fieldName), nil
	case api.NodeImageFamilyUbuntu2004, api.NodeImageFamilyUbuntu1804:
		return "", &UnsupportedQueryError{msg: fmt.Sprintf("SSM Parameter lookups for %s AMIs is not supported yet", imageFamily)}
	default:
		return "", fmt.Errorf("unknown image family %s", imageFamily)
	}
}

// MakeManagedSSMParameterName creates an SSM parameter name for a managed nodegroup
func MakeManagedSSMParameterName(version, amiType string) (string, error) {
	switch amiType {
	case eks.AMITypesAl2X8664:
		imageType := utils.ToKebabCase(api.NodeImageFamilyAmazonLinux2)
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/recommended/release_version", version, imageType), nil
	case eks.AMITypesAl2X8664Gpu:
		imageType := utils.ToKebabCase(api.NodeImageFamilyAmazonLinux2) + "-gpu"
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/recommended/release_version", version, imageType), nil
	}
	return "", nil
}

// instanceEC2ArchName returns the name of the architecture as used by EC2
// resources.
func instanceEC2ArchName(instanceType string) string {
	if instanceutils.IsARMInstanceType(instanceType) {
		return "arm64"
	}
	return "x86_64"
}

func imageType(imageFamily, instanceType, version string) string {
	family := utils.ToKebabCase(imageFamily)
	switch imageFamily {
	case api.NodeImageFamilyBottlerocket:
		if instanceutils.IsNvidiaInstanceType(instanceType) {
			return fmt.Sprintf("%s-%s", version, "nvidia")
		}
		return version
	default:
		if instanceutils.IsGPUInstanceType(instanceType) {
			return family + "-gpu"
		}
		if instanceutils.IsARMInstanceType(instanceType) {
			return family + "-arm64"
		}
		return family
	}
}
