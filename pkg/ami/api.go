package ami

import (
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const (
	// ImageFamilyAmazonLinux2 represents Amazon Linux 2 family
	ImageFamilyAmazonLinux2 = api.NodeImageFamilyAmazonLinux2 // Owner 602401143452

	// ImageFamilyUbuntu1804 represents Ubuntu 18.04 family
	ImageFamilyUbuntu1804 = api.NodeImageFamilyUbuntu1804 // Owner 099720109477

	// ResolverStatic is used to indicate that the static (i.e. compiled into eksctl) AMIs should be used
	ResolverStatic = api.NodeImageResolverStatic
	// ResolverAuto is used to indicate that the latest EKS AMIs should be used for the nodes. This implies
	// that automatic resolution of AMI will occur.
	ResolverAuto = api.NodeImageResolverAuto
)

// Variations of iamge classes
const (
	ImageClassGeneral int = iota
	ImageClassGPU
)

// ImageClasses is a list of image class names
var ImageClasses = []string{
	"ImageClassGeneral",
	"ImageClassGPU",
}

// IsAvailable checks if a given AMI ID is available in AWS EC2
func IsAvailable(api ec2iface.EC2API, id string) (bool, string, string, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{&id},
	}

	output, err := api.DescribeImages(input)
	if err != nil {
		return false, "", "", errors.Wrapf(err, "unable to find %q", id)
	}

	// This will never return more than one as we are looking up a single ami id
	if len(output.Images) < 1 {
		return false, "", "", nil
	}

	return *output.Images[0].State == "available", *output.Images[0].RootDeviceName, *output.Images[0].RootDeviceType, nil
}

// FindImage will get the AMI to use for the EKS nodes by querying AWS EC2 API.
// It will only look for images with a status of available and it will pick the
// image with the newest creation date.
func FindImage(api ec2iface.EC2API, ownerAccount, namePattern string) (string, error) {
	input := &ec2.DescribeImagesInput{
		Owners: []*string{&ownerAccount},
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: []*string{&namePattern},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []*string{aws.String("hvm")},
			},
			{
				Name:   aws.String("root-device-type"),
				Values: []*string{aws.String("ebs")},
			},
			{
				Name:   aws.String("is-public"),
				Values: []*string{aws.String("true")},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
		},
	}

	output, err := api.DescribeImages(input)
	if err != nil {
		return "", errors.Wrapf(err, "error querying AWS for images")
	}

	if len(output.Images) < 1 {
		return "", nil
	}

	if len(output.Images) == 1 {
		return *output.Images[0].ImageId, nil
	}

	// Sort images so newest is first
	sort.Slice(output.Images, func(i, j int) bool {
		//nolint:gosec
		creationLeft, _ := time.Parse(time.RFC3339, *output.Images[i].CreationDate)
		//nolint:gosec
		creationRight, _ := time.Parse(time.RFC3339, *output.Images[j].CreationDate)
		return creationLeft.After(creationRight)
	})

	return *output.Images[0].ImageId, nil
}
