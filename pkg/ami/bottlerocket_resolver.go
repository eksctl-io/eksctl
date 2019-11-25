package ami

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// StaticBottlerocketResolver retrieves Bottlerocket AMI IDs from a
// static mapping.
type StaticBottlerocketResolver struct{}

// Resolve will resolve a Bottlerocket AMI ID for the supplied region
// and instance type from a static mapping.
func (*StaticBottlerocketResolver) Resolve(region, kubeVersion, instanceType, imageFamily string) (string, error) {
	if imageFamily != api.NodeImageFamilyBottlerocket {
		return "", nil
	}

	if inRegion, ok := bottlerocketByKubeVersion[kubeVersion]; ok {
		if image, ok := inRegion[region]; ok {
			return image, nil
		}
	}

	return "", NewErrFailedResolution(region, kubeVersion, instanceType, imageFamily)
}

var bottlerocketByKubeVersion = map[string]map[string]string{
	"1.15": {
		// v0.3.1
		"ap-northeast-1": "ami-0f96126d91bdfcb26",
		"eu-central-1":   "ami-0d0d1fb5b0eccae93",
		"us-east-1":      "ami-07fc2dd9e0d7741ed",
		"us-west-2":      "ami-04e994d5dc1edce6b",
	},
}
