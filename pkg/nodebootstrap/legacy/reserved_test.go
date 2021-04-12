package legacy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reservations defaults", func() {
	type storageCase struct {
		gib int64
		exp string
	}
	DescribeTable("calculate storage", func(c storageCase) {
		Expect(InstanceTypeInfo{Storage: c.gib}.DefaultStorageToReserve()).To(Equal(c.exp))
	},
		Entry("0GiB", storageCase{0, "1Gi"}),
		Entry("1GiB", storageCase{1, "1Gi"}),
		Entry("100GiB", storageCase{100, "1Gi"}),
		Entry("1000GiB", storageCase{1000, "1Gi"}),
	)

	type memoryCase struct {
		maxPodsPerNode int64
		exp            string
	}
	DescribeTable("calculate memory", func(c memoryCase) {
		Expect(InstanceTypeInfo{MaxPodsPerNode: c.maxPodsPerNode}.DefaultMemoryToReserve()).To(Equal(c.exp))
	},
		Entry("4 pods per node", memoryCase{4, "299Mi"}),
		Entry("8 pods per node", memoryCase{8, "343Mi"}),
		Entry("29 pods per node", memoryCase{29, "574Mi"}),
		Entry("58 pods per node", memoryCase{58, "893Mi"}),
		Entry("234 pods per node", memoryCase{234, "2829Mi"}),
		Entry("452 pods per node", memoryCase{452, "5227Mi"}),
		Entry("737 pods per node", memoryCase{737, "8362Mi"}),
	)

	type cpuCase struct {
		cpu int64
		exp string
	}
	DescribeTable("calculate cpu", func(c cpuCase) {
		Expect(InstanceTypeInfo{CPU: c.cpu}.DefaultCPUToReserve()).To(Equal(c.exp))
	},
		Entry("0 cpu", cpuCase{0, "0m"}),
		Entry("1 cpu", cpuCase{1, "60m"}),
		Entry("2 cpu", cpuCase{2, "70m"}),
		Entry("4 cpu", cpuCase{4, "80m"}),
		Entry("8 cpu", cpuCase{8, "90m"}),
		Entry("16 cpu", cpuCase{16, "110m"}),
		Entry("32 cpu", cpuCase{32, "150m"}),
		Entry("48 cpu", cpuCase{48, "190m"}),
		Entry("64 cpu", cpuCase{64, "230m"}),
		Entry("80 cpu", cpuCase{80, "270m"}),
		Entry("96 cpu", cpuCase{96, "310m"}),
	)

	type instanceCase struct {
		instance   string
		expMemory  string
		expCPU     string
		expStorage string
	}
	DescribeTable("calculate by instance type", func(c instanceCase) {
		Expect(instanceTypeInfos[c.instance].DefaultStorageToReserve()).To(Equal(c.expStorage))
		Expect(instanceTypeInfos[c.instance].DefaultMemoryToReserve()).To(Equal(c.expMemory))
		Expect(instanceTypeInfos[c.instance].DefaultCPUToReserve()).To(Equal(c.expCPU))
	},
		Entry("a1.2xlarge", instanceCase{"a1.2xlarge", "893Mi", "90m", "1Gi"}),
		Entry("t3.nano", instanceCase{"t3.nano", "299Mi", "70m", "1Gi"}),
		Entry("t3a.micro", instanceCase{"t3a.micro", "299Mi", "70m", "1Gi"}),
		Entry("t2.small", instanceCase{"t2.small", "376Mi", "60m", "1Gi"}),
		Entry("t2.medium", instanceCase{"t2.medium", "442Mi", "70m", "1Gi"}),
		Entry("m5ad.large", instanceCase{"m5ad.large", "574Mi", "70m", "1Gi"}),
		Entry("m5ad.xlarge", instanceCase{"m5ad.xlarge", "893Mi", "80m", "1Gi"}),
		Entry("m5ad.2xlarge", instanceCase{"m5ad.2xlarge", "893Mi", "90m", "1Gi"}),
		Entry("m5ad.12xlarge", instanceCase{"m5ad.12xlarge", "2829Mi", "190m", "1Gi"}),
		Entry("m5ad.24xlarge", instanceCase{"m5ad.24xlarge", "8362Mi", "310m", "1Gi"}),
	)
})
