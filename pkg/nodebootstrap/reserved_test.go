package nodebootstrap

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
		mib int64
		exp string
	}
	DescribeTable("calculate memory", func(c memoryCase) {
		Expect(InstanceTypeInfo{Memory: c.mib}.DefaultMemoryToReserve()).To(Equal(c.exp))
	},
		Entry("0MiB", memoryCase{0, "255Mi"}),
		Entry("512MiB", memoryCase{512, "255Mi"}),
		Entry("2GiB", memoryCase{2048, "512Mi"}),
		Entry("4GiB", memoryCase{4096, "1024Mi"}),
		Entry("6GiB", memoryCase{6144, "1433Mi"}),
		Entry("8GiB", memoryCase{8192, "1843Mi"}),
		Entry("12GiB", memoryCase{12288, "2252Mi"}),
		Entry("64GiB", memoryCase{1 << 16, "5611Mi"}),
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
		Entry("a1.2xlarge", instanceCase{"a1.2xlarge", "2662Mi", "90m", "1Gi"}),
		Entry("t3.nano", instanceCase{"t3.nano", "255Mi", "70m", "1Gi"}),
		Entry("t3a.micro", instanceCase{"t3a.micro", "256Mi", "70m", "1Gi"}),
		Entry("t2.small", instanceCase{"t2.small", "512Mi", "60m", "1Gi"}),
		Entry("t2.medium", instanceCase{"t2.medium", "1024Mi", "70m", "1Gi"}),
		Entry("m5ad.large", instanceCase{"m5ad.large", "1843Mi", "70m", "1Gi"}),
		Entry("m5ad", instanceCase{"m5ad.xlarge", "2662Mi", "80m", "1Gi"}),
		Entry("m5ad.2xlarge", instanceCase{"m5ad.2xlarge", "3645Mi", "90m", "1Gi"}),
		Entry("m5ad.8xlarge", instanceCase{"m5ad.8xlarge", "9543Mi", "150m", "1Gi"}),
		Entry("m5ad.12xlarge", instanceCase{"m5ad.12xlarge", "10853Mi", "190m", "1Gi"}),
		Entry("m5ad.16xlarge", instanceCase{"m5ad.16xlarge", "12164Mi", "230m", "1Gi"}),
		Entry("m5ad.24xlarge", instanceCase{"m5ad.24xlarge", "14785Mi", "310m", "1Gi"}),
	)
})
