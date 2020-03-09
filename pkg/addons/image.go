package addons

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	corev1 "k8s.io/api/core/v1"
)

// awsDNSSuffixForRegion returns the AWS DNS suffix (amazonaws.com or amazonaws.com.cn) for the specified region
func awsDNSSuffixForRegion(region string) (string, error) {
	for _, p := range endpoints.DefaultPartitions() {
		if _, ok := p.Regions()[region]; ok {
			return p.DNSSuffix(), nil
		}
	}
	return "", fmt.Errorf("failed to find DNS suffix for region %s", region)
}

// UseRegionalImage sets the region and AWS DNS suffix for a container image
// in format '%s.dkr.ecr.%s.%s/image:tag'
func UseRegionalImage(spec *corev1.PodTemplateSpec, region string) error {
	imageFormat := spec.Spec.Containers[0].Image
	dnsSuffix, err := awsDNSSuffixForRegion(region)
	if err != nil {
		return err
	}
	regionalImage := fmt.Sprintf(imageFormat, api.EKSResourceAccountID(region), region, dnsSuffix)
	spec.Spec.Containers[0].Image = regionalImage
	return nil
}

// imageTag extracts the container image's tag.
func imageTag(image string) (string, error) {
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected image format %q", image)
	}

	return parts[1], nil
}

// ImageTagsDiffer returns true if the image tags are not the same
// while ignoring the image name.
func ImageTagsDiffer(image1, image2 string) (bool, error) {
	tag1, err := imageTag(image1)
	if err != nil {
		return false, err
	}
	tag2, err := imageTag(image2)
	if err != nil {
		return false, err
	}
	return tag1 != tag2, nil
}
