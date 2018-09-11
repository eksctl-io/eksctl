package ami

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
)

// IsAmiAvailableInAWS checks if a given ami is available in AWS
func IsAmiAvailableInAWS(api ec2iface.EC2API, ami string) (bool, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(ami)},
	}

	output, err := api.DescribeImages(input)
	if err != nil {
		errors.Wrapf(err, "Unable to find AMI with id %s", ami)
	}

	if len(output.Images) < 1 {
		return false, nil
	}

	return *output.Images[0].State == "available", nil
}

// QueryAWSForEksAmi will get the AMI is for the EKS nodesby querying
// AWS for theimage to use. You need to supply a name pattern to use.
func QueryAWSForEksAmi(api ec2iface.EC2API, namePattern string) (string, error) {
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
		errors.Wrap(err, "Error querying AWS for AMI")
	}

	if len(output.Images) < 1 {
		return "", nil
	}

	if *output.Images[0].State == "available" {
		return *output.Images[0].ImageId, nil
	}

	return "", nil
}
