package ami

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/utils"
)

const (
	// NodePatternNonGPU is the pattern used to search for images
	// for use with nodes without GPU shupport
	NodePatternNonGPU = "amazon-eks-node-*"

	// NodePatternGpu is the pattern used to search for images
	// for use with nodes with GPU shupport
	NodePatternGpu = "amazon-eks-gpu-node-*"
)

// AutoResolver resolves the AMi to the defaults for the region
// by querying AWS for the AMI to use
type AutoResolver struct {
	api ec2iface.EC2API
}

// Resolve will return an AMI to use based on the default AMI for each region
func (r *AutoResolver) Resolve(region string, instanceType string) (string, error) {
	logger.Debug("resolving AMI using AutoResolver for region %s and instanceType %s", region, instanceType)

	nodePattern := NodePatternNonGPU
	if utils.IsGPUInstanceType(instanceType) {
		nodePattern = NodePatternGpu
	}

	imageID, err := FindImageForEKS(r.api, nodePattern)
	if err != nil {
		return "", errors.Wrap(err, "error getting ami to use for region")
	}

	return imageID, nil
}

// NewAutoResolver creates a new AutoResolver
func NewAutoResolver(api ec2iface.EC2API) *AutoResolver {
	return &AutoResolver{api: api}
}
