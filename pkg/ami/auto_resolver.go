package ami

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils"
)

const (
	// ownerIDUbuntu1804Family is the owner ID used for Ubuntu AMIs
	ownerIDUbuntu1804Family = "099720109477"

	// ownerIDWindowsFamily is the owner ID used for Ubuntu AMIs
	ownerIDWindowsFamily = "801119661308"
)

// MakeImageSearchPatterns creates a map of image search patterns by image OS family and class
func MakeImageSearchPatterns(version string) map[string]map[int]string {
	return map[string]map[int]string{
		api.NodeImageFamilyAmazonLinux2: {
			ImageClassGeneral: fmt.Sprintf("amazon-eks-node-%s-v*", version),
			ImageClassGPU:     fmt.Sprintf("amazon-eks-gpu-node-%s-*", version),
		},
		api.NodeImageFamilyUbuntu1804: {
			ImageClassGeneral: fmt.Sprintf("ubuntu-eks/k8s_%s/images/*", version),
		},
		api.NodeImageFamilyWindowsServer2019CoreContainer: {
			ImageClassGeneral: fmt.Sprintf("Windows_Server-2019-English-Core-EKS_Optimized-%v-*", version),
		},
		api.NodeImageFamilyWindowsServer2019FullContainer: {
			ImageClassGeneral: fmt.Sprintf("Windows_Server-2019-English-Full-EKS_Optimized-%v-*", version),
		},
	}
}

// OwnerAccountID returns the AWS account ID that owns worker AMI.
func OwnerAccountID(imageFamily, region string) (string, error) {
	switch imageFamily {
	case api.NodeImageFamilyUbuntu1804:
		return ownerIDUbuntu1804Family, nil
	case api.NodeImageFamilyWindowsServer2019CoreContainer, api.NodeImageFamilyWindowsServer2019FullContainer:
		return ownerIDWindowsFamily, nil
	case api.NodeImageFamilyAmazonLinux2:
		return api.EKSResourceAccountID(region), nil
	default:
		return "", fmt.Errorf("unable to determine the account owner for image family %s", imageFamily)
	}
}

// AutoResolver resolves the AMi to the defaults for the region
// by querying AWS EC2 API for the AMI to use
type AutoResolver struct {
	api ec2iface.EC2API
}

// Resolve will return an AMI to use based on the default AMI for
// each region
func (r *AutoResolver) Resolve(region, version, instanceType, imageFamily string) (string, error) {
	logger.Debug("resolving AMI using AutoResolver for region %s, instanceType %s and imageFamily %s", region, instanceType, imageFamily)

	imageClasses := MakeImageSearchPatterns(version)[imageFamily]
	namePattern := imageClasses[ImageClassGeneral]
	if utils.IsGPUInstanceType(instanceType) {
		var ok bool
		namePattern, ok = imageClasses[ImageClassGPU]
		if !ok {
			logger.Critical("image family %s doesn't support GPU image class", imageFamily)
			return "", NewErrFailedResolution(region, version, instanceType, imageFamily)
		}
	}

	ownerAccount, err := OwnerAccountID(imageFamily, region)
	if err != nil {
		logger.Critical("%v", err)
		return "", NewErrFailedResolution(region, version, instanceType, imageFamily)
	}

	id, err := FindImage(r.api, ownerAccount, namePattern)
	if err != nil {
		return "", errors.Wrap(err, "error getting AMI")
	}

	return id, nil
}
