package addon_test

import (
	"context"
	"encoding/json"
	"fmt"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("DescribeVersions", func() {
	var (
		manager                   *addon.Manager
		mockProvider              *mockprovider.MockProvider
		describeAddonVersonsInput *awseks.DescribeAddonVersionsInput
	)
	BeforeEach(func() {
		var err error
		mockProvider = mockprovider.NewMockProvider()
		manager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}, mockProvider.EKS(), nil, false, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("DescribeVersions", func() {
		It("returns an addon", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
				describeAddonVersonsInput = args[1].(*awseks.DescribeAddonVersionsInput)
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
					},
				},
			}, nil)

			summary, err := manager.DescribeVersions(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.MarshalIndent(struct {
				Addons []ekstypes.AddonInfo
			}{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
					},
				},
			}, "", "\t")
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(string(expected)))
			Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
			Expect(*describeAddonVersonsInput.AddonName).To(Equal("my-addon"))
		})

		When("it fails to describe addon versions", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
					describeAddonVersonsInput = args[1].(*awseks.DescribeAddonVersionsInput)
				}).Return(&awseks.DescribeAddonVersionsOutput{}, fmt.Errorf("foo"))

				_, err := manager.DescribeVersions(context.Background(), &api.Addon{
					Name: "my-addon",
				})
				Expect(err).To(MatchError(`failed to describe addon versions: foo`))
				Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
				Expect(*describeAddonVersonsInput.AddonName).To(Equal("my-addon"))
			})
		})

		It("returns an addon with a publisher, type and owner", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
				describeAddonVersonsInput = args[1].(*awseks.DescribeAddonVersionsInput)
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("upbound_universal-crossplane"),
						Type:      aws.String("infra-management"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
						Publisher: aws.String("upbound"),
						Owner:     aws.String("aws-marketplace"),
					},
				},
			}, nil)

			summary, err := manager.DescribeVersions(context.Background(), &api.Addon{
				Name:       "upbound_universal-crossplane",
				Types:      []string{"infra-management"},
				Owners:     []string{"aws-marketplace"},
				Publishers: []string{"upbound"},
			})
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.MarshalIndent(struct {
				Addons []ekstypes.AddonInfo
			}{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("upbound_universal-crossplane"),
						Type:      aws.String("infra-management"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
						Publisher: aws.String("upbound"),
						Owner:     aws.String("aws-marketplace"),
					},
				},
			}, "", "\t")
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(string(expected)))
			Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
			Expect(*describeAddonVersonsInput.AddonName).To(Equal("upbound_universal-crossplane"))
			Expect(describeAddonVersonsInput.Types).To(Equal([]string{"infra-management"}))
			Expect(describeAddonVersonsInput.Publishers).To(Equal([]string{"upbound"}))
			Expect(describeAddonVersonsInput.Owners).To(Equal([]string{"aws-marketplace"}))
		})
	})

	Describe("DescribeAllVersions", func() {
		It("returns an addon", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
				describeAddonVersonsInput = args[1].(*awseks.DescribeAddonVersionsInput)
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
					},
				},
			}, nil)

			summary, err := manager.DescribeAllVersions(context.Background(), &api.Addon{})
			Expect(err).NotTo(HaveOccurred())

			expected, err := json.MarshalIndent(struct {
				Addons []ekstypes.AddonInfo
			}{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
					},
				},
			}, "", "\t")
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(string(expected)))

			Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
			Expect(describeAddonVersonsInput.AddonName).To(BeNil())
		})

		When("it fails to describe addon versions", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
					describeAddonVersonsInput = args[1].(*awseks.DescribeAddonVersionsInput)
				}).Return(&awseks.DescribeAddonVersionsOutput{}, fmt.Errorf("foo"))

				_, err := manager.DescribeAllVersions(context.Background(), &api.Addon{})
				Expect(err).To(MatchError(`failed to describe addon versions: foo`))
				Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
				Expect(describeAddonVersonsInput.AddonName).To(BeNil())
			})
		})

		It("returns an addon with a publisher, type and owner", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
				describeAddonVersonsInput = args[1].(*awseks.DescribeAddonVersionsInput)
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("upbound_universal-crossplane"),
						Type:      aws.String("infra-management"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
						Publisher: aws.String("upbound"),
						Owner:     aws.String("aws-marketplace"),
					},
				},
			}, nil)

			summary, err := manager.DescribeAllVersions(context.Background(), &api.Addon{
				Types:      []string{"infra-management"},
				Owners:     []string{"aws-marketplace"},
				Publishers: []string{"upbound"},
			})
			Expect(err).NotTo(HaveOccurred())

			expected, err := json.MarshalIndent(struct {
				Addons []ekstypes.AddonInfo
			}{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("upbound_universal-crossplane"),
						Type:      aws.String("infra-management"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
						Publisher: aws.String("upbound"),
						Owner:     aws.String("aws-marketplace"),
					},
				},
			}, "", "\t")
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(string(expected)))

			Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
			Expect(describeAddonVersonsInput.AddonName).To(BeNil())
			Expect(describeAddonVersonsInput.Types).To(Equal([]string{"infra-management"}))
			Expect(describeAddonVersonsInput.Publishers).To(Equal([]string{"upbound"}))
			Expect(describeAddonVersonsInput.Owners).To(Equal([]string{"aws-marketplace"}))
		})
	})
})
