package fargate

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const (
	// MinPlatformVersion is the minimum platform version which supports
	// Fargate, i.e. represents "eks.5":
	MinPlatformVersion = 5

	// MinKubernetesVersion is the minimum Kubernetes version which supports
	// Fargate.
	MinKubernetesVersion = api.Version1_14
)
