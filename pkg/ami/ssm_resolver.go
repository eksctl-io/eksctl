package ami

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/weaveworks/eksctl/pkg/utils"
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
		return "", errors.Wrap(err, "error getting AMI from SSM Parameter Store")
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

	switch imageFamily {
	case api.NodeImageFamilyAmazonLinux2:
		return fmt.Sprintf("/aws/service/eks/optimized-ami/%s/%s/recommended/image_id", version, imageType(imageFamily, instanceType)), nil
	case api.NodeImageFamilyWindowsServer2019CoreContainer:
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-2019-English-Core-EKS_Optimized-%s/image_id", version), nil
	case api.NodeImageFamilyWindowsServer2019FullContainer:
		return fmt.Sprintf("/aws/service/ami-windows-latest/Windows_Server-2019-English-Full-EKS_Optimized-%s/image_id", version), nil
	case api.NodeImageFamilyBottlerocket:
		return fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/%s/latest/image_id", version, instanceEC2ArchName(instanceType)), nil
	case api.NodeImageFamilyUbuntu1804:
		return "", &UnsupportedQueryError{msg: fmt.Sprintf("SSM Parameter lookups for %s AMIs is not supported yet", imageFamily)}
	default:
		return "", fmt.Errorf("unknown image family %s", imageFamily)
	}
}

// instanceEC2ArchName returns the name of the architecture as used by EC2
// resources.
func instanceEC2ArchName(instanceType string) string {
	// eg: a1.large - an ARM instance type.
	if strings.HasPrefix(instanceType, "a") {
		return "arm64"
	}
	return "x86_64"
}

func imageType(imageFamily, instanceType string) string {
	family := utils.ToKebabCase(imageFamily)
	if utils.IsGPUInstanceType(instanceType) {
		return family + "-gpu"
	}
	return family
}
