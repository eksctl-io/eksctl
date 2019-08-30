package ami

import (
	"fmt"
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
	ImageFamilyAmazonLinux2 = api.NodeImageFamilyAmazonLinux2 // Owner by EKS(depends on aws partition and opt-in region)

	// ImageFamilyUbuntu1804 represents Ubuntu 18.04 family
	ImageFamilyUbuntu1804 = api.NodeImageFamilyUbuntu1804 // Owner 099720109477

	// ResolverStatic is used to indicate that the static (i.e. compiled into eksctl) AMIs should be used
	ResolverStatic = api.NodeImageResolverStatic
	// ResolverAuto is used to indicate that the latest EKS AMIs should be used for the nodes. This implies
	// that automatic resolution of AMI will occur.
	ResolverAuto = api.NodeImageResolverAuto
)

// Variations of image classes
const (
	ImageClassGeneral = iota
	ImageClassGPU
)

// ImageClasses is a list of image class names
var ImageClasses = []string{
	"ImageClassGeneral",
	"ImageClassGPU",
}

// Use checks if a given AMI ID is available in AWS EC2 as well as checking and populating RootDevice information
func Use(ec2api ec2iface.EC2API, ng *api.NodeGroup) error {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{&ng.AMI},
	}

	output, err := ec2api.DescribeImages(input)
	if err != nil {
		return errors.Wrapf(err, "unable to find image %q", ng.AMI)
	}

	// This will never return more than one as we are looking up a single ami id
	if len(output.Images) < 1 {
		return NewErrNotFound(ng.AMI)
	}

	// Instance-store AMIs cannot have their root volume size managed
	if *output.Images[0].RootDeviceType == "instance-store" {
		return fmt.Errorf("%q is an instance-store AMI and EBS block device mappings not supported for instance-store AMIs", ng.AMI)
	}

	if *output.Images[0].RootDeviceType == "ebs" {
		if !api.IsSetAndNonEmptyString(ng.VolumeName) {
			ng.VolumeName = output.Images[0].RootDeviceName
		}

		amiEncrypted := output.Images[0].BlockDeviceMappings[0].Ebs.Encrypted
		if ng.VolumeEncrypted == nil {
			ng.VolumeEncrypted = amiEncrypted
		} else {
			// VolumeEncrypted cannot be false if the AMI being used is already encrypted.
			if api.IsDisabled(ng.VolumeEncrypted) && api.IsEnabled(amiEncrypted) {
				return fmt.Errorf("%q is an encrypted AMI and volumeEncrypted has been set to false", ng.AMI)
			}
		}
	}

	return nil
}

// FindImage will get the AMI to use for the EKS nodes by querying AWS EC2 API.
// It will only look for images with a status of available and it will pick the
// image with the newest creation date.
func FindImage(ec2api ec2iface.EC2API, ownerAccount, namePattern string) (string, error) {
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

	output, err := ec2api.DescribeImages(input)
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
