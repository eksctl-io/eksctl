package instance_test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/utils/instance"
)

type instanceTypeEntry struct {
	instanceTypes        []ec2types.InstanceTypeInfo
	expectedInstanceType string
}

var _ = Describe("Instance utils", func() {
	DescribeTable("GetSmallestInstanceType", func(e instanceTypeEntry) {
		instanceType := instance.GetSmallestInstanceType(e.instanceTypes)
		Expect(instanceType).To(Equal(e.expectedInstanceType))
	},
		Entry("instance types with a distinct vCPUs count", instanceTypeEntry{
			instanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: "m5a.12xlarge",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(48),
						DefaultCores:          aws.Int32(24),
						DefaultThreadsPerCore: aws.Int32(2),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(196608),
					},
				},
				{
					InstanceType: "m5a.large",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(2),
						DefaultCores:          aws.Int32(1),
						DefaultThreadsPerCore: aws.Int32(2),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(196608),
					},
				},
				{
					InstanceType: "m5.xlarge",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(4),
						DefaultCores:          aws.Int32(2),
						DefaultThreadsPerCore: aws.Int32(2),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(16384),
					},
				},
				{
					InstanceType: "m5a.16xlarge",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(64),
						DefaultCores:          aws.Int32(32),
						DefaultThreadsPerCore: aws.Int32(2),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(262144),
					},
				},
			},
			expectedInstanceType: "m5a.large",
		}),
		Entry("instance types with the same vCPUs count", instanceTypeEntry{
			instanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: "t2.large",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(2),
						DefaultCores:          aws.Int32(2),
						DefaultThreadsPerCore: aws.Int32(1),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(8192),
					},
				},
				{
					InstanceType: "t2.medium",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(2),
						DefaultCores:          aws.Int32(2),
						DefaultThreadsPerCore: aws.Int32(1),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(4096),
					},
				},
				{
					InstanceType: "t3.large",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(2),
						DefaultCores:          aws.Int32(2),
						DefaultThreadsPerCore: aws.Int32(1),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(8192),
					},
				},
				{
					InstanceType: "t3a.nano",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(2),
						DefaultCores:          aws.Int32(1),
						DefaultThreadsPerCore: aws.Int32(2),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(512),
					},
				},
				{
					InstanceType: "t3.nano",
					VCpuInfo: &ec2types.VCpuInfo{
						DefaultVCpus:          aws.Int32(2),
						DefaultCores:          aws.Int32(1),
						DefaultThreadsPerCore: aws.Int32(2),
					},
					MemoryInfo: &ec2types.MemoryInfo{
						SizeInMiB: aws.Int64(512),
					},
				},
			},

			expectedInstanceType: "t3a.nano",
		}),
	)
})
