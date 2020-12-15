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

// Variations of image classes
const (
	ImageClassGeneral = iota
	ImageClassGPU
	ImageClassARM
)

// Bottlerocket disk used by kubelet
const (
	bottlerocketDataDisk = "/dev/xvdb"
	bottlerocketOSDisk   = "/dev/xvda"
)

// ImageClasses is a list of image class names
var ImageClasses = []string{
	"ImageClassGeneral",
	"ImageClassGPU",
	"ImageClassARM",
}

// Use checks if a given AMI ID is available in AWS EC2 as well as checking and populating RootDevice information
func Use(ec2api ec2iface.EC2API, ng *api.NodeGroupBase) error {
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

	image := output.Images[0]

	switch *image.RootDeviceType {
	// Instance-store AMIs cannot have their root volume size managed
	case "instance-store":
		return fmt.Errorf("%q is an instance-store AMI and EBS block device mappings are not supported for instance-store AMIs", ng.AMI)

	case "ebs":
		if !api.IsSetAndNonEmptyString(ng.VolumeName) {
			ng.VolumeName = image.RootDeviceName
			if ng.AMIFamily == api.NodeImageFamilyBottlerocket {
				ng.VolumeName = aws.String(bottlerocketDataDisk)
				ng.OSVolumeName = aws.String(bottlerocketOSDisk)
			}
		}
		rootDeviceMapping, err := findRootDeviceMapping(image)
		if err != nil {
			return err
		}

		amiEncrypted := rootDeviceMapping.Ebs.Encrypted
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

func findRootDeviceMapping(image *ec2.Image) (*ec2.BlockDeviceMapping, error) {
	for _, deviceMapping := range image.BlockDeviceMappings {
		if *deviceMapping.DeviceName == *image.RootDeviceName {
			return deviceMapping, nil
		}
	}
	return nil, errors.Errorf("failed to find root device mapping for AMI %q", *image.ImageId)
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
