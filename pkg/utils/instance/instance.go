package instance

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// IsARMInstanceType returns true if the instance type is ARM architecture
func IsARMInstanceType(instanceType string) bool {
	return InstanceTypesMap[instanceType].CPUArch == "arm64"
}

// IsGPUInstanceType returns true if the instance type is GPU optimised
func IsGPUInstanceType(instanceType string) bool {
	itype := InstanceTypesMap[instanceType]
	return itype.NvidiaGPUSupported || itype.NeuronSupported
}

// IsNeuronInstanceType returns true if the instance type requires AWS Neuron
func IsNeuronInstanceType(instanceType string) bool {
	return InstanceTypesMap[instanceType].NeuronSupported
}

// IsNvidiaInstanceType returns true if the instance type has NVIDIA accelerated hardware
func IsNvidiaInstanceType(instanceType string) bool {
	return InstanceTypesMap[instanceType].NvidiaGPUSupported
}

// IsInferentiaInstanceType returns true if the instance type requires AWS Neuron Inferentia/Inferentia2
func IsInferentiaInstanceType(instanceType string) bool {
	itype := InstanceTypesMap[instanceType]
	return itype.NeuronSupported &&
		(itype.NeuronDeviceType == "Inferentia" || itype.NeuronDeviceType == "Inferentia2")
}

// IsTrainiumnstanceType returns true if the instance type requires AWS Neuron Trainium/Trainium2
func IsTrainiumInstanceType(instanceType string) bool {
	itype := InstanceTypesMap[instanceType]
	return itype.NeuronSupported &&
		(itype.NeuronDeviceType == "Trainium" || itype.NeuronDeviceType == "Trainium2")
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
