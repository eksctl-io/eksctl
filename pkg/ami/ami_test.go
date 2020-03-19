package ami

import (
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/mock"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

func TestUseAMI(t *testing.T) {
	amiTests := []struct {
		blockDeviceMappings []*ec2.BlockDeviceMapping
		rootDeviceName      string
	}{
		{
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
		},
	}

	for i, tt := range amiTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			mockProvider := mockDescribeImages(tt.blockDeviceMappings, tt.rootDeviceName)
			err := Use(mockProvider.MockEC2(), &api.NodeGroup{
				AMI: "ami-0121d8347f8191f90",
			})

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

}

func mockDescribeImages(blockDeviceMappings []*ec2.BlockDeviceMapping, rootDeviceName string) *mockprovider.MockProvider {
	mockProvider := mockprovider.NewMockProvider()

	mockProvider.MockEC2().On("DescribeImages", mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
		return len(input.ImageIds) == 1 && strings.HasPrefix(*input.ImageIds[0], "ami-")
	})).Return(func(input *ec2.DescribeImagesInput) *ec2.DescribeImagesOutput {
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
	}, func(*ec2.DescribeImagesInput) error {
		return nil
	})

	return mockProvider
}
