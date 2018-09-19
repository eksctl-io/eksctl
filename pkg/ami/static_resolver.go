package ami

import (
	"fmt"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils"
)

// TODO: This file will be go generated in the future
// https://github.com/weaveworks/eksctl/issues/49

// StaticImages is a map of images in each region by
// image OS family and by class
var StaticImages = map[string]map[int]map[string]string{
	ImageFamilyAmazonLinux2: {
		ImageClassGeneral: {
			api.EKS_REGION_US_WEST_2: "ami-0a54c984b9f908c81",
			api.EKS_REGION_US_EAST_1: "ami-0440e4f6b9713faf6",
			api.EKS_REGION_EU_WEST_1: "ami-0c7a4976cb6fafd3a",
		},
		ImageClassGPU: {
			api.EKS_REGION_US_WEST_2: "ami-0731694d53ef9604b",
			api.EKS_REGION_US_EAST_1: "ami-058bfb8c236caae89",
			api.EKS_REGION_EU_WEST_1: "ami-0706dc8a5eed2eed9",
		},
	},
}

// StaticDefaultResolver resolves the AMI to the defaults for the region
type StaticDefaultResolver struct {
}

// Resolve will return an AMI to use based on the default AMI for each region
// currently source of truth for these is here
func (r *StaticDefaultResolver) Resolve(region string, instanceType string) (string, error) {
	logger.Debug("resolving AMI using StaticDefaultResolver for region %s and instanceType %s", region, instanceType)

	regionalAMIs := StaticImages[DefaultImageFamily][ImageClassGeneral]
	return regionalAMIs[region], nil
}

// StaticGPUResolver resolves the AMI for GPU instances types.
type StaticGPUResolver struct {
}

// Resolve will return an AMI based on the region for GPU instance types
func (r *StaticGPUResolver) Resolve(region string, instanceType string) (string, error) {
	logger.Debug("resolving AMI using StaticGPUResolver for region %s and instanceType %s", region, instanceType)

	regionalAMIs, ok := StaticImages[DefaultImageFamily][ImageClassGPU]
	if !ok {
		return "", fmt.Errorf("image family %s doesn't support GPU image class", DefaultImageFamily)
	}
	if !utils.IsGPUInstanceType(instanceType) {
		logger.Debug("can't resolve AMI using StaticGPUResolver as instance type %s is non-GPU", instanceType)
		return "", nil
	}

	return regionalAMIs[region], nil
}
