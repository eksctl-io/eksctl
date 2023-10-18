package instance

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// IsARMInstanceType returns true if the instance type is ARM architecture
func IsARMInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "a1") ||
		strings.HasPrefix(instanceType, "t4g") ||
		strings.HasPrefix(instanceType, "m6g") ||
		strings.HasPrefix(instanceType, "m7g") ||
		strings.HasPrefix(instanceType, "c6g") ||
		strings.HasPrefix(instanceType, "c7g") ||
		strings.HasPrefix(instanceType, "r6g") ||
		strings.HasPrefix(instanceType, "r7g") ||
		strings.HasPrefix(instanceType, "im4g") ||
		strings.HasPrefix(instanceType, "is4g") ||
		strings.HasPrefix(instanceType, "g5g") ||
		strings.HasPrefix(instanceType, "x2g")
}

// IsGPUInstanceType returns true if the instance type is GPU optimised
func IsGPUInstanceType(instanceType string) bool {
	return IsNvidiaInstanceType(instanceType) ||
		IsInferentiaInstanceType(instanceType) ||
		IsTrainiumInstanceType(instanceType)
}

// IsNeuronInstanceType returns true if the instance type requires AWS Neuron
func IsNeuronInstanceType(instanceType string) bool {
	return IsInferentiaInstanceType(instanceType) ||
		IsTrainiumInstanceType(instanceType)
}

// IsARMGPUInstanceType returns true if the instance type is ARM-GPU architecture
func IsARMGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "g5g")
}

// IsNvidiaInstanceType returns true if the instance type has NVIDIA accelerated hardware
func IsNvidiaInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") ||
		strings.HasPrefix(instanceType, "p3") ||
		strings.HasPrefix(instanceType, "p4") ||
		strings.HasPrefix(instanceType, "p5") ||
		strings.HasPrefix(instanceType, "g3") ||
		strings.HasPrefix(instanceType, "g4") ||
		strings.HasPrefix(instanceType, "g5")
}

// IsInferentiaInstanceType returns true if the instance type requires AWS Neuron
func IsInferentiaInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "inf1")
}

// IsTrainiumnstanceType returns true if the instance type requires AWS Neuron
func IsTrainiumInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "trn1")
}

// GetSmallestInstanceType returns the smallest instance type in instanceTypes.
// Instance types that have a smaller vCPU are considered smaller.
// instanceTypes must be non-empty or it will panic.
func GetSmallestInstanceType(instanceTypes []ec2types.InstanceTypeInfo) string {
	smallestInstanceTypeInfo := instanceTypes[0]
	for _, it := range instanceTypes[1:] {
		switch vCPUs, smallestVCPUs := aws.ToInt32(it.VCpuInfo.DefaultVCpus), aws.ToInt32(smallestInstanceTypeInfo.VCpuInfo.DefaultVCpus); {
		case vCPUs < smallestVCPUs:
			smallestInstanceTypeInfo = it
		case vCPUs == smallestVCPUs:
			if aws.ToInt64(it.MemoryInfo.SizeInMiB) < aws.ToInt64(smallestInstanceTypeInfo.MemoryInfo.SizeInMiB) {
				smallestInstanceTypeInfo = it
			}
		}
	}

	return string(smallestInstanceTypeInfo.InstanceType)
}
