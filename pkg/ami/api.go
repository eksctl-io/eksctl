package ami

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
)

const (
	DefaultImageFamily = ImageFamilyAmazonLinux2

	ImageFamilyAmazonLinux2 = "AmazonLinux2"

	// ResolverStatic is used to indicate that the stqtic (i.e. compiled into eksctl) AMIs should be used
	ResolverStatic = "static"
	// ResolverAuto is used to indicate that the latest EKS AMIs should be used for the nodes. This implies
	// that automatic resolution of AMI will occur.
	ResolverAuto = "auto"
)

const (
	ImageClassGeneral int = iota
	ImageClassGPU
)

// IsAvailable checks if a given AMI ID is available in AWS EC2
func IsAvailable(api ec2iface.EC2API, id string) (bool, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(id)},
	}

	output, err := api.DescribeImages(input)
	if err != nil {
		return false, errors.Wrapf(err, "unable to find %q", id)
	}

	if len(output.Images) < 1 {
		return false, nil
	}

	return *output.Images[0].State == "available", nil
}

// FindImage will get the AMI to use for the EKS nodes by querying AWS EC2 API
func FindImage(api ec2iface.EC2API, namePattern string) (string, error) {
	input := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("name"),
				Values: []*string{aws.String(namePattern)},
			},
			&ec2.Filter{
				Name:   aws.String("virtualization-type"),
				Values: []*string{aws.String("hvm")},
			},
			&ec2.Filter{
				Name:   aws.String("root-device-type"),
				Values: []*string{aws.String("ebs")},
			},
			&ec2.Filter{
				Name:   aws.String("is-public"),
				Values: []*string{aws.String("true")},
			},
		},
	}

	output, err := api.DescribeImages(input)
	if err != nil {
		return "", errors.Wrapf(err, "cannot find image")
	}

	if len(output.Images) < 1 {
		return "", nil
	}

	if *output.Images[0].State == "available" {
		return *output.Images[0].ImageId, nil
	}

	return "", nil
}
