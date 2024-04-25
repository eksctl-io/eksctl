package ami

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
)

// SSMResolver resolves the AMI to the defaults for the region
// by querying AWS SSM get parameter API
type SSMResolver struct {
	ssmAPI awsapi.SSM
}

// Resolve will return an AMI to use based on the default AMI for
// each region
func (r *SSMResolver) Resolve(ctx context.Context, region, version, instanceType, imageFamily string) (string, error) {
	logger.Debug("resolving AMI using SSM Parameter resolver for region %s, instanceType %s and imageFamily %s", region, instanceType, imageFamily)

	parameterName, err := MakeSSMParameterName(version, instanceType, imageFamily)
	if err != nil {
		return "", err
	}

	output, err := r.ssmAPI.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(parameterName),
	})
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
	const fieldName = "image_id"

	switch imageFamily {
	case api.NodeImageFamilyAmazonLinux2023:
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/%s/standard/recommended/%s",
			version, utils.ToKebabCase(imageFamily), instanceEC2ArchName(instanceType), fieldName), nil
	case api.NodeImageFamilyAmazonLinux2:
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/recommended/%s", version, imageType(imageFamily, instanceType, version), fieldName), nil
	case api.NodeImageFamilyWindowsServer2019CoreContainer,
		api.NodeImageFamilyWindowsServer2019FullContainer:
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-2019-English-%s-EKS_Optimized-%s/%s", windowsAmiType(imageFamily), version, fieldName), nil
	case api.NodeImageFamilyWindowsServer2022CoreContainer,
		api.NodeImageFamilyWindowsServer2022FullContainer:
		const minVersion = api.Version1_23
		supportsWindows2022, err := utils.IsMinVersion(minVersion, version)
		if err != nil {
			return "", err
		}
		if !supportsWindows2022 {
			return "", errors.Errorf("Windows Server 2022 %s requires EKS version %s and above", windowsAmiType(imageFamily), minVersion)
		}
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-2022-English-%s-EKS_Optimized-%s/%s", windowsAmiType(imageFamily), version, fieldName), nil
	case api.NodeImageFamilyBottlerocket:
		return fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/%s/latest/%s", imageType(imageFamily, instanceType, version), instanceEC2ArchName(instanceType), fieldName), nil
	case api.NodeImageFamilyUbuntu2204, api.NodeImageFamilyUbuntu2004, api.NodeImageFamilyUbuntu1804:
		// FIXME: SSM lookup for Ubuntu EKS images is supported nowadays
		return "", &UnsupportedQueryError{msg: fmt.Sprintf("SSM Parameter lookups for %s AMIs is not supported yet", imageFamily)}
	default:
		return "", fmt.Errorf("unknown image family %s", imageFamily)
	}
}

// MakeManagedSSMParameterName creates an SSM parameter name for a managed nodegroup
func MakeManagedSSMParameterName(version string, amiType ekstypes.AMITypes) (string, error) {
	switch amiType {
	case ekstypes.AMITypesAl2023X8664Standard:
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/x86_64/standard/recommended/release_version", version, utils.ToKebabCase(api.NodeImageFamilyAmazonLinux2023)), nil
	case ekstypes.AMITypesAl2023Arm64Standard:
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/arm64/standard/recommended/release_version", version, utils.ToKebabCase(api.NodeImageFamilyAmazonLinux2023)), nil
	case ekstypes.AMITypesAl2X8664:
		imageType := utils.ToKebabCase(api.NodeImageFamilyAmazonLinux2)
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/recommended/release_version", version, imageType), nil
	case ekstypes.AMITypesAl2X8664Gpu:
		imageType := utils.ToKebabCase(api.NodeImageFamilyAmazonLinux2) + "-gpu"
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/recommended/release_version", version, imageType), nil
	case ekstypes.AMITypesBottlerocketArm64, ekstypes.AMITypesBottlerocketArm64Nvidia:
		return fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/arm64/latest/image_version", version), nil
	case ekstypes.AMITypesBottlerocketX8664, ekstypes.AMITypesBottlerocketX8664Nvidia:
		return fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/x86_64/latest/image_version", version), nil
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

func windowsAmiType(imageFamily string) string {
	if strings.Contains(imageFamily, "Core") {
		return "Core"
	}
	return "Full"
}
