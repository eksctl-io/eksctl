package fargate

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils"
)

const (
	// MinPlatformVersion is the minimum platform version which supports
	// Fargate, i.e. represents "eks.5":
	MinPlatformVersion = 5

	// MinKubernetesVersion is the minimum Kubernetes version which supports
	// Fargate.
	MinKubernetesVersion = api.Version1_14
)

// IsSupportedBy reports whether the control plane version can support Fargate.
func IsSupportedBy(controlPlaneVersion string) (bool, error) {
	supportsFargate, err := utils.IsMinVersion(MinKubernetesVersion, controlPlaneVersion)
	if err != nil {
		return false, err
	}
	return supportsFargate, nil
}
