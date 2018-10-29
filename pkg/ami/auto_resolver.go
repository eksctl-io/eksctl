package ami

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/utils"
)

// ImageSearchPatterns is a map of image search patterns by
// image OS family and by class
var ImageSearchPatterns = map[string]map[int]string{
	ImageFamilyAmazonLinux2: {
		ImageClassGeneral: "amazon-eks-node-*",
		ImageClassGPU:     "amazon-eks-gpu-node-*",
	},
	ImageFamilyUbuntu1804: {
		ImageClassGeneral: "ubuntu-eks/1.10.3/*",
	},
}

// AutoResolver resolves the AMi to the defaults for the region
// by querying AWS EC2 API for the AMI to use
type AutoResolver struct {
	api ec2iface.EC2API
}

// Resolve will return an AMI to use based on the default AMI for
// each region
func (r *AutoResolver) Resolve(region string, instanceType string, imageFamily string) (string, error) {
	logger.Debug("resolving AMI using AutoResolver for region %s, instanceType %s and imageFamily %s", region, instanceType, imageFamily)

	p := ImageSearchPatterns[imageFamily][ImageClassGeneral]
	if utils.IsGPUInstanceType(instanceType) {
		var ok bool
		p, ok = ImageSearchPatterns[imageFamily][ImageClassGPU]
		if !ok {
			logger.Critical("image family %s doesn't support GPU image class", imageFamily)
			return "", NewErrFailedResolution(region, instanceType, imageFamily)
		}
	}

	id, err := FindImage(r.api, p)
	if err != nil {
		return "", errors.Wrap(err, "error getting AMI")
	}

	return id, nil
}

// NewAutoResolver creates a new AutoResolver
func NewAutoResolver(api ec2iface.EC2API) *AutoResolver {
	return &AutoResolver{api: api}
}
