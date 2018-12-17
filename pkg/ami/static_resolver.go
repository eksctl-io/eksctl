package ami

import (
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/utils"
)

//go:generate go run ./static_resolver_ami_generate.go

// StaticDefaultResolver resolves the AMI to the defaults for the region
type StaticDefaultResolver struct {
}

// Resolve will return an AMI to use based on the default AMI for each region
// currently source of truth for these is here
func (r *StaticDefaultResolver) Resolve(region, version, instanceType, imageFamily string) (string, error) {
	logger.Debug("resolving AMI using StaticDefaultResolver for region %s, version %s, instanceType %s and imageFamily %s", region, instanceType, imageFamily)

	regionalAMIs := StaticImages[version][imageFamily][ImageClassGeneral]
	return regionalAMIs[region], nil
}

// StaticGPUResolver resolves the AMI for GPU instances types.
type StaticGPUResolver struct {
}

// Resolve will return an AMI based on the region for GPU instance types
func (r *StaticGPUResolver) Resolve(region, version, instanceType, imageFamily string) (string, error) {
	logger.Debug("resolving AMI using StaticGPUResolver for region %s, instanceType %s and imageFamily %s", region, instanceType, imageFamily)

	if !utils.IsGPUInstanceType(instanceType) {
		logger.Debug("can't resolve AMI using StaticGPUResolver as instance type %s is non-GPU", instanceType)
		return "", nil
	}

	regionalAMIs, ok := StaticImages[version][imageFamily][ImageClassGPU]
	if !ok {
		logger.Critical("image family %s doesn't support GPU image class", imageFamily)
		return "", NewErrFailedResolution(region, version, instanceType, imageFamily)
	}

	return regionalAMIs[region], nil
}
