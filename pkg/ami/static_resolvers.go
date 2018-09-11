package ami

import (
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/utils"
)

// StaticDefaultResolver resolves the AMi to the defaults for the region
type StaticDefaultResolver struct {
}

// Resolve will return an AMI to use based on the default AMI for each region
// TODO: https://github.com/weaveworks/eksctl/issues/49
// currently source of truth for these is here:
// https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
func (r *StaticDefaultResolver) Resolve(region string, instanceType string) (string, error) {
	logger.Debug("resolving AMI using StaticDefaultResolver for region %s and instanceType %s", region, instanceType)

	return regionalAMIs[region], nil
}

// StaticGPUResolver resolves the AMI for GPU instances types.
type StaticGPUResolver struct {
}

// Resolve will return an AMI based on the region for GPU instance types
func (r *StaticGPUResolver) Resolve(region string, instanceType string) (string, error) {
	logger.Debug("resolving AMI using StaticGPUResolver for region %s and instanceType %s", region, instanceType)

	if !utils.IsGPUInstanceType(instanceType) {
		logger.Debug("can't resolve AMI using StaticGPUResolver as instance type %s is GPU optimized", instanceType)
		return "", nil
	}

	return gpuRegionalAMIs[region], nil
}
