package efa

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/version"
)

// IsBuiltInSupported returns true if the Kubernetes version supports built-in EFA in the default security group
func IsBuiltInSupported(kubernetesVersion string) (bool, error) {
	supported, err := version.IsMinVersion(api.EFABuiltInSupportVersion, kubernetesVersion)
	if err != nil {
		return false, fmt.Errorf("failed to determine EFA built-in support for Kubernetes version %q (minimum required: %s): %w",
			kubernetesVersion, api.EFABuiltInSupportVersion, err)
	}
	return supported, nil
}
