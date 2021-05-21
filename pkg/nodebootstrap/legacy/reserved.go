package legacy

import (
	"fmt"
	"math"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	minimumMemoryToReserve = 255
)

// progression allows to define a progressive approach to take
// fractions of values.
type progression []struct {
	upper    int64
	fraction float64
}

var (
	cpuProgression = progression{
		{upper: 1000, fraction: 0.06},
		{upper: 2000, fraction: 0.01},
		{upper: 4000, fraction: 0.005},
		{upper: math.MaxInt64, fraction: 0.0025},
	}
)

func (p progression) calculate(value, min int64) int64 {
	var lower, reserve int64
	for _, f := range p {
		if value <= f.upper {
			reserve += int64(float64(value-lower) * f.fraction)
			break
		}
		reserve += int64(float64(f.upper-lower) * f.fraction)
		lower = f.upper
	}
	if reserve < min {
		return min
	}
	return reserve

}

// InstanceTypeInfo holds minimal instance info required to
// calculate resources to reserve.
type InstanceTypeInfo struct {
	// Storage (ephemeral) available (GiB).
	// Is 0 if not supported or none available.
	Storage int64
	// Max pods per node.
	MaxPodsPerNode int64
	// CPU count.
	CPU int64
}

// NewInstanceTypeInfo creates a simple version of ec2.InstanceTypeInfo
// that provides functions to calculate defaults.
func NewInstanceTypeInfo(ec2info *ec2.InstanceTypeInfo) InstanceTypeInfo {
	i := InstanceTypeInfo{}
	if ec2info == nil {
		return i
	}
	if aws.BoolValue(ec2info.InstanceStorageSupported) && ec2info.InstanceStorageInfo != nil {
		i.Storage = aws.Int64Value(ec2info.InstanceStorageInfo.TotalSizeInGB)
	}
	if maxPodsPerNode, exists := maxPodsPerNodeType[*ec2info.InstanceType]; exists {
		i.MaxPodsPerNode = int64(maxPodsPerNode)
	}
	if ec2info.VCpuInfo != nil {
		i.CPU = aws.Int64Value(ec2info.VCpuInfo.DefaultVCpus)
	}
	return i
}

// DefaultStorageToReserve returns how much storage to reserve.
//
// See https://github.com/awslabs/amazon-eks-ami/blob/ff690788dfaf399e6919eebb59371ee923617df4/files/bootstrap.sh#L306
// This is always 1GiB
func (i InstanceTypeInfo) DefaultStorageToReserve() string {
	return "1Gi"
}

// DefaultMemoryToReserve returns how much memory to reserve.
//
// See https://github.com/awslabs/amazon-eks-ami/blob/21426e27e3845319dbca92e7df32e5c4b984a1d1/files/bootstrap.sh#L154-L165
func (i InstanceTypeInfo) DefaultMemoryToReserve() string {
	mib := 11*i.MaxPodsPerNode + minimumMemoryToReserve
	return fmt.Sprintf("%dMi", mib)
}

// DefaultCPUToReserve returns the millicores to reserve.
//
// See https://github.com/awslabs/amazon-eks-ami/blob/ff690788dfaf399e6919eebb59371ee923617df4/files/bootstrap.sh#L183-L208
// which takes it form https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-architecture#node_allocatable
//
// 6% of the first core
// 1% of the next core (up to 2 cores)
// 0.5% of the next 2 cores (up to 4 cores)
// 0.25% of any cores above 4 cores
func (i InstanceTypeInfo) DefaultCPUToReserve() string {
	millicores := cpuProgression.calculate(i.CPU*1000, 0)
	return fmt.Sprintf("%dm", millicores)
}
