package ami_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

func TestUseAMI(t *testing.T) {
	amiTests := []struct {
		blockDeviceMappings []*ec2.BlockDeviceMapping
		rootDeviceName      string
		amiFamily           string
		volumeName          string
		description         string

		expectedVolumeName string
		expectedEncrypted  bool
	}{
		{
			description:    "Root device mapping not at index 0 (Windows AMIs in some regions)",
			rootDeviceName: "/dev/sda1",
			blockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName:  aws.String("xvdca"),
					VirtualName: aws.String("ephemeral0"),
				},
				{
					DeviceName:  aws.String("xvdcb"),
					VirtualName: aws.String("ephemeral1"),
				},
				{
					DeviceName: aws.String("/dev/sda1"),
					Ebs: &ec2.EbsBlockDevice{
						Encrypted: aws.Bool(true),
					},
				},
			},

			expectedEncrypted:  true,
			expectedVolumeName: "/dev/sda1",
		},
		{
			description:    "Only one device mapping (AL2 AMIs)",
			rootDeviceName: "/dev/sda1",
			blockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName: aws.String("/dev/sda1"),
					Ebs: &ec2.EbsBlockDevice{
						Encrypted: aws.Bool(true),
					},
				},
			},

			expectedEncrypted:  true,
			expectedVolumeName: "/dev/sda1",
		},
		{
			description:    "Different root device name",
			rootDeviceName: "/dev/xvda",
			blockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName:  aws.String("xvdca"),
					VirtualName: aws.String("ephemeral0"),
				},
				{
					DeviceName: aws.String("/dev/xvda"),
					Ebs: &ec2.EbsBlockDevice{
						Encrypted: aws.Bool(false),
					},
				},
				{
					DeviceName:  aws.String("xvdcb"),
					VirtualName: aws.String("ephemeral1"),
				},
			},

			expectedEncrypted:  false,
			expectedVolumeName: "/dev/xvda",
		},
		{
			description:    "volumeName for Bottlerocket is not modified",
			rootDeviceName: "/dev/xvda",
			amiFamily:      api.NodeImageFamilyBottlerocket,
			volumeName:     "/dev/xvdb",
			blockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName: aws.String("/dev/xvda"),
					Ebs: &ec2.EbsBlockDevice{
						Encrypted: aws.Bool(false),
					},
				},
			},

			expectedEncrypted:  false,
			expectedVolumeName: "/dev/xvdb",
		},
	}

	for i, tt := range amiTests {
		t.Run(fmt.Sprintf("%d: %s", i, tt.description), func(t *testing.T) {
			mockProvider := mockDescribeImages(tt.blockDeviceMappings, tt.rootDeviceName)
			ng := &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AMI:       "ami-0121d8347f8191f90",
					AMIFamily: tt.amiFamily,
				},
			}
			if tt.volumeName != "" {
				ng.VolumeName = aws.String(tt.volumeName)
			}

			err := ami.Use(context.Background(), mockProvider.MockEC2(), ng.NodeGroupBase)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if *ng.VolumeEncrypted != tt.expectedEncrypted {
				t.Errorf("expected VolumeEncrypted to be %v", tt.expectedEncrypted)
			}

			if *ng.VolumeName != tt.expectedVolumeName {
				t.Errorf("expected VolumeName to be %v", tt.expectedVolumeName)
			}
		})
	}

}

func mockDescribeImages(blockDeviceMappings []*ec2.BlockDeviceMapping, rootDeviceName string) *mockprovider.MockProvider {
	mockProvider := mockprovider.NewMockProvider()

	mockProvider.MockEC2().On("DescribeImagesWithContext", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
		return len(input.ImageIds) == 1 && strings.HasPrefix(*input.ImageIds[0], "ami-")
	})).Return(func(_ context.Context, input *ec2.DescribeImagesInput, _ ...request.Option) *ec2.DescribeImagesOutput {
		return &ec2.DescribeImagesOutput{
			Images: []*ec2.Image{
				{
					ImageId:             input.ImageIds[0],
					RootDeviceName:      aws.String(rootDeviceName),
					RootDeviceType:      aws.String("ebs"),
					BlockDeviceMappings: blockDeviceMappings,
				},
			},
		}
	}, func(_ context.Context, _ *ec2.DescribeImagesInput, _ ...request.Option) error {
		return nil
	})

	return mockProvider
}
