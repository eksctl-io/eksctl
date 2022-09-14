package outposts_test

import (
	"context"
	"sync"

	"github.com/weaveworks/eksctl/pkg/outposts/fakes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awsoutposts "github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/outposts"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Outposts Service", func() {
	Context("API calls count", func() {
		var (
			provider        *mockprovider.MockProvider
			outpostsService *outposts.Service
		)

		BeforeEach(func() {
			provider = mockprovider.NewMockProvider()
			mockOutpostInstanceTypes(provider)
			outpostsService = &outposts.Service{
				EC2API:      provider.EC2(),
				OutpostsAPI: provider.Outposts(),
				OutpostID:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		})

		runAssertion := func(count int, doAssertion func(*sync.WaitGroup)) chan struct{} {
			var wg sync.WaitGroup
			wg.Add(count)
			doAssertion(&wg)
			doneCh := make(chan struct{})
			go func() {
				wg.Wait()
				close(doneCh)
			}()
			return doneCh
		}

		It("should not make redundant API calls", func() {
			By("validating instance types")
			doneCh := runAssertion(len(instanceTypeInfoList), func(wg *sync.WaitGroup) {
				for _, it := range instanceTypeInfoList {
					go func(instanceType ec2types.InstanceType) {
						defer GinkgoRecover()
						defer wg.Done()
						Expect(outpostsService.ValidateInstanceType(context.Background(), string(instanceType))).To(Succeed())
					}(it.InstanceType)
				}
			})
			err := outpostsService.ValidateInstanceType(context.Background(), "t2.medium")
			Expect(err).To(MatchError(`instance type "t2.medium" does not exist in Outpost "arn:aws:outposts:us-west-2:1234:outpost/op-1234"`))

			Eventually(doneCh).Should(BeClosed())
			provider.MockOutposts().AssertNumberOfCalls(GinkgoT(), "GetOutpostInstanceTypes", 1)
			provider.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeInstanceTypes", 0)

			By("fetching the smallest available instance type")
			count := 5
			doneCh = runAssertion(count, func(wg *sync.WaitGroup) {
				for i := 0; i < count; i++ {
					go func() {
						defer GinkgoRecover()
						defer wg.Done()
						_, err := outpostsService.GetSmallestInstanceType(context.Background())
						Expect(err).NotTo(HaveOccurred())
					}()
				}
			})
			Eventually(doneCh).Should(BeClosed())
			provider.MockOutposts().AssertNumberOfCalls(GinkgoT(), "GetOutpostInstanceTypes", 1)
			provider.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeInstanceTypes", 1)
		})
	})

	type setOrValidateEntry struct {
		outpostInstance *fakes.FakeOutpostInstance

		expectedInstanceType string
	}

	DescribeTable("SetOrValidateInstanceType", func(e setOrValidateEntry) {
		provider := mockprovider.NewMockProvider()
		mockOutpostInstanceTypes(provider)
		outpostsService := &outposts.Service{
			OutpostsAPI: provider.Outposts(),
			EC2API:      provider.EC2(),
			OutpostID:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
		}

		Expect(outpostsService.SetOrValidateOutpostInstanceType(context.Background(), e.outpostInstance)).To(Succeed())
		Expect(e.outpostInstance.GetInstanceTypeCallCount()).To(Equal(1))
		if e.expectedInstanceType != "" {
			Expect(e.outpostInstance.SetInstanceTypeCallCount()).To(Equal(1))
			instanceType := e.outpostInstance.SetInstanceTypeArgsForCall(0)
			Expect(instanceType).To(Equal(e.expectedInstanceType))
		} else {
			Expect(e.outpostInstance.SetInstanceTypeCallCount()).To(Equal(0))
		}
	},
		Entry("instance type not set", setOrValidateEntry{
			outpostInstance: &fakes.FakeOutpostInstance{
				GetInstanceTypeStub: func() string { return "" },
				SetInstanceTypeStub: func(_ string) {},
			},

			expectedInstanceType: "m5a.large",
		}),

		Entry("instance type set", setOrValidateEntry{
			outpostInstance: &fakes.FakeOutpostInstance{
				GetInstanceTypeStub: func() string { return "m5.xlarge" },
				SetInstanceTypeStub: func(_ string) {},
			},
		}),
	)

	Context("GetOutpost", func() {
		var (
			provider        *mockprovider.MockProvider
			outpostsService *outposts.Service
		)

		BeforeEach(func() {
			provider = mockprovider.NewMockProvider()
			outpostsService = &outposts.Service{
				OutpostsAPI: provider.Outposts(),
				OutpostID:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		})

		It("should return the Outpost details", func() {
			outpostDetails := &outpoststypes.Outpost{
				OutpostId:        aws.String(outpostsService.OutpostID),
				Name:             aws.String("op-1234"),
				AvailabilityZone: aws.String("us-west-2a"),
				OutpostArn:       aws.String(outpostsService.OutpostID),
			}
			provider.MockOutposts().On("GetOutpost", mock.Anything, &awsoutposts.GetOutpostInput{
				OutpostId: aws.String(outpostsService.OutpostID),
			}).Return(&awsoutposts.GetOutpostOutput{
				Outpost: outpostDetails,
			}, nil)

			outpost, err := outpostsService.GetOutpost(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(outpost).To(Equal(outpostDetails))
		})

		It("should return an error if the Outpost does not exist", func() {
			provider.MockOutposts().On("GetOutpost", mock.Anything, &awsoutposts.GetOutpostInput{
				OutpostId: aws.String(outpostsService.OutpostID),
			}).Return(nil, &outpoststypes.NotFoundException{
				Message: aws.String("Outpost does not exist"),
			})

			_, err := outpostsService.GetOutpost(context.Background())
			Expect(err).To(MatchError(ContainSubstring("Outpost does not exist")))
		})
	})
})

var instanceTypeInfoList = []ec2types.InstanceTypeInfo{
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
}

func mockOutpostInstanceTypes(provider *mockprovider.MockProvider) {
	instanceTypeItems := make([]outpoststypes.InstanceTypeItem, len(instanceTypeInfoList))
	instanceTypes := make([]ec2types.InstanceType, len(instanceTypeInfoList))
	for i, it := range instanceTypeInfoList {
		instanceTypeItems[i] = outpoststypes.InstanceTypeItem{
			InstanceType: aws.String(string(it.InstanceType)),
		}
		instanceTypes[i] = it.InstanceType
	}
	provider.MockOutposts().On("GetOutpostInstanceTypes", mock.Anything, mock.Anything).Return(&awsoutposts.GetOutpostInstanceTypesOutput{
		InstanceTypes: instanceTypeItems,
	}, nil)

	provider.MockEC2().On("DescribeInstanceTypes", mock.Anything, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: instanceTypes,
	}).Return(&ec2.DescribeInstanceTypesOutput{
		InstanceTypes: instanceTypeInfoList,
	}, nil)
}
