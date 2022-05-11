package ami

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
)

// Variations of image classes
const (
	ImageClassGeneral = iota
	ImageClassGPU
	ImageClassARM
)

// ImageClasses is a list of image class names
var ImageClasses = []string{
	"ImageClassGeneral",
	"ImageClassGPU",
	"ImageClassARM",
}

// Use checks if a given AMI ID is available in AWS EC2 as well as checking and populating RootDevice information
func Use(ctx context.Context, ec2API awsapi.EC2, ng *api.NodeGroupBase) error {
	output, err := ec2API.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{ng.AMI},
	})
	if err != nil {
		return errors.Wrapf(err, "unable to find image %q", ng.AMI)
	}

	// This will never return more than one as we are looking up a single ami id
	if len(output.Images) < 1 {
		return NewErrNotFound(ng.AMI)
	}

	image := output.Images[0]

	switch image.RootDeviceType {
	// Instance-store AMIs cannot have their root volume size managed
	case ec2types.DeviceTypeInstanceStore:
		return fmt.Errorf("%q is an instance-store AMI and EBS block device mappings are not supported for instance-store AMIs", ng.AMI)

	case ec2types.DeviceTypeEbs:
		if ng.AMIFamily != api.NodeImageFamilyBottlerocket && !api.IsSetAndNonEmptyString(ng.VolumeName) {
			// Volume name is preset for Bottlerocket.
			ng.VolumeName = image.RootDeviceName
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

func findRootDeviceMapping(image ec2types.Image) (ec2types.BlockDeviceMapping, error) {
	for _, deviceMapping := range image.BlockDeviceMappings {
		if *deviceMapping.DeviceName == *image.RootDeviceName {
			return deviceMapping, nil
		}
	}
	return ec2types.BlockDeviceMapping{}, errors.Errorf("failed to find root device mapping for AMI %q", *image.ImageId)
}

// FindImage will get the AMI to use for the EKS nodes by querying AWS EC2 API.
// It will only look for images with a status of available and it will pick the
// image with the newest creation date.
func FindImage(ctx context.Context, ec2API awsapi.EC2, ownerAccount, namePattern string) (string, error) {
	input := &ec2.DescribeImagesInput{
		Owners: []string{ownerAccount},
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{namePattern},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
			{
				Name:   aws.String("root-device-type"),
				Values: []string{"ebs"},
			},
			{
				Name:   aws.String("is-public"),
				Values: []string{"true"},
			},
			{
				Name:   aws.String("state"),
				Values: []string{"available"},
			},
		},
	}

	output, err := ec2API.DescribeImages(ctx, input)
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
