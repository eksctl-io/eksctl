package ami

import (
	"fmt"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/utils"
)

//go:generate go run ./static_resolver_ami_generate.go

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
