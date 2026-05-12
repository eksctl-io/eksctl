package addons

import (
	"fmt"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	corev1 "k8s.io/api/core/v1"
)

// awsDNSSuffixForRegion returns the AWS DNS suffix (amazonaws.com or amazonaws.com.cn) for the specified region
func awsDNSSuffixForRegion(region string) (string, error) {
	return api.Partitions.V1SDKDNSPrefixForRegion(region)
}

// UseRegionalImage sets the region and AWS DNS suffix for all container images
// in format '%s.dkr.ecr.%s.%s/image:tag'
func UseRegionalImage(spec *corev1.PodTemplateSpec, region string) error {
	dnsSuffix, err := awsDNSSuffixForRegion(region)
	if err != nil {
		return err
	}

	for i := range spec.Spec.Containers {
		imageFormat := spec.Spec.Containers[i].Image
		if isRegionalImageFormat(imageFormat) {
			regionalImage := fmt.Sprintf(imageFormat, api.EKSResourceAccountID(region), region, dnsSuffix)
			spec.Spec.Containers[i].Image = regionalImage
		}
	}

	for i := range spec.Spec.InitContainers {
		imageFormat := spec.Spec.InitContainers[i].Image
		if isRegionalImageFormat(imageFormat) {
			regionalImage := fmt.Sprintf(imageFormat, api.EKSResourceAccountID(region), region, dnsSuffix)
			spec.Spec.InitContainers[i].Image = regionalImage
		}
	}

	return nil
}

// isRegionalImageFormat checks whether an image string contains format verbs
// (i.e., it's a template like "%s.dkr.ecr.%s.%s/image:tag" rather than a
// fully-resolved image URI).
func isRegionalImageFormat(image string) bool {
	return strings.Contains(image, "%s")
}

// ImageTag extracts the container image's tag.
func ImageTag(image string) (string, error) {
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected image format %q", image)
	}

	return parts[1], nil
}

// ImageTagsDiffer returns true if the image tags are not the same
// while ignoring the image name.
func ImageTagsDiffer(image1, image2 string) (bool, error) {
	tag1, err := ImageTag(image1)
	if err != nil {
		return false, err
	}
	tag2, err := ImageTag(image2)
	if err != nil {
		return false, err
	}
	return tag1 != tag2, nil
}
