package addon_test

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Get", func() {
	var (
		manager            *addon.Manager
		mockProvider       *mockprovider.MockProvider
		describeAddonInput *awseks.DescribeAddonInput
	)
	BeforeEach(func() {
		var err error
		mockProvider = mockprovider.NewMockProvider()
		manager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}, mockProvider.EKS(), nil, false, nil, nil, 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Get", func() {
		It("returns an addon", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []*awseks.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []*awseks.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0.0"),
							},
							{
								//not sure if all versions come with v prefix or not, so test a mix
								AddonVersion: aws.String("v1.1.0"),
							},
							{
								AddonVersion: aws.String("1.2.0"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[0].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &awseks.Addon{
					AddonName:             aws.String("my-addon"),
					AddonVersion:          aws.String("v1.0.0"),
					ServiceAccountRoleArn: aws.String("foo"),
					Status:                aws.String("created"),
					Health: &awseks.AddonHealth{
						Issues: []*awseks.AddonIssue{
							{
								Code:        aws.String("1"),
								Message:     aws.String("foo"),
								ResourceIds: aws.StringSlice([]string{"id-1"}),
							},
						},
					},
				},
			}, nil)

			summary, err := manager.Get(&api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:         "my-addon",
				Version:      "v1.0.0",
				NewerVersion: "v1.1.0,1.2.0",
				IAMRole:      "foo",
				Status:       "created",
				Issues:       []string{"{\n  Code: \"1\",\n  Message: \"foo\",\n  ResourceIds: [\"id-1\"]\n}"},
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
		})

		When("it fails to get the addon", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DescribeAddon", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
					describeAddonInput = args[0].(*awseks.DescribeAddonInput)
				}).Return(nil, fmt.Errorf("foo"))

				_, err := manager.Get(&api.Addon{
					Name: "my-addon",
				})
				Expect(err).To(MatchError(`failed to get addon "my-addon": foo`))
				Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
			})
		})
	})

	Describe("GetAll", func() {
		var listAddonsInput *awseks.ListAddonsInput
		It("returns an addon", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []*awseks.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []*awseks.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0.0"),
							},
							{
								//not sure if all versions come with v prefix or not, so test a mix
								AddonVersion: aws.String("v1.1.0"),
							},
							{
								AddonVersion: aws.String("1.2.0"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("ListAddons", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.ListAddonsInput{}))
				listAddonsInput = args[0].(*awseks.ListAddonsInput)
			}).Return(&awseks.ListAddonsOutput{
				Addons: aws.StringSlice([]string{"my-addon"}),
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[0].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &awseks.Addon{
					AddonName:             aws.String("my-addon"),
					AddonVersion:          aws.String("1.0.0"),
					ServiceAccountRoleArn: aws.String("foo"),
					Status:                aws.String("created"),
				},
			}, nil)

			summary, err := manager.GetAll()
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal([]addon.Summary{
				{
					Name:         "my-addon",
					Version:      "1.0.0",
					NewerVersion: "v1.1.0,1.2.0",
					IAMRole:      "foo",
					Status:       "created",
				},
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*listAddonsInput.ClusterName).To(Equal("my-cluster"))
		})

		When("it fails to get the addon", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("ListAddons", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.ListAddonsInput{}))
					listAddonsInput = args[0].(*awseks.ListAddonsInput)
				}).Return(&awseks.ListAddonsOutput{
					Addons: aws.StringSlice([]string{"my-addon"}),
				}, nil)

				mockProvider.MockEKS().On("DescribeAddon", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
					describeAddonInput = args[0].(*awseks.DescribeAddonInput)
				}).Return(nil, fmt.Errorf("foo"))

				_, err := manager.GetAll()
				Expect(err).To(MatchError(`failed to get addon "my-addon": foo`))
				Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*listAddonsInput.ClusterName).To(Equal("my-cluster"))

			})
		})

		When("it fails to list addons", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("ListAddons", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.ListAddonsInput{}))
					listAddonsInput = args[0].(*awseks.ListAddonsInput)
				}).Return(&awseks.ListAddonsOutput{
					Addons: aws.StringSlice([]string{"my-addon"}),
				}, fmt.Errorf("foo"))

				_, err := manager.GetAll()
				Expect(err).To(MatchError(`failed to list addons: foo`))
				Expect(*listAddonsInput.ClusterName).To(Equal("my-cluster"))
			})
		})
	})
})
