package ami

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils"
)

// MakeImageSearchPatterns creates a map of image search patterns by image OS family and class
func MakeImageSearchPatterns(version string) map[string]map[int]string {
	return map[string]map[int]string{
		ImageFamilyAmazonLinux2: {
			ImageClassGeneral: fmt.Sprintf("amazon-eks-node-%s-v*", version),
			ImageClassGPU:     fmt.Sprintf("amazon-eks-gpu-node-%s-*", version),
		},
		ImageFamilyUbuntu1804: {
			ImageClassGeneral: fmt.Sprintf("ubuntu-eks/k8s_%s/images/*", version),
		},
	}
}

// OwnerAccountID returns the AWS account ID that owns worker AMI.
func OwnerAccountID(imageFamily, region string) (string, error) {
	switch imageFamily {
	case ImageFamilyUbuntu1804:
		return "099720109477", nil
	case ImageFamilyAmazonLinux2:
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

// NewAutoResolver creates a new AutoResolver
func NewAutoResolver(api ec2iface.EC2API) *AutoResolver {
	return &AutoResolver{api: api}
}
