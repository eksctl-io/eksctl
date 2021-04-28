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
		}}, mockProvider.EKS(), nil, false, nil, nil, 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("DescribeVersions", func() {
		It("returns an addon", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
				describeAddonVersonsInput = args[0].(*awseks.DescribeAddonVersionsInput)
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []*awseks.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []*awseks.AddonVersionInfo{
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

			summary, err := manager.DescribeVersions(&api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(awseks.DescribeAddonVersionsOutput{
				Addons: []*awseks.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []*awseks.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
					},
				},
			}.String()))

			Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
			Expect(*describeAddonVersonsInput.AddonName).To(Equal("my-addon"))
		})

		When("it fails to describe addon versions", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
					describeAddonVersonsInput = args[0].(*awseks.DescribeAddonVersionsInput)
				}).Return(&awseks.DescribeAddonVersionsOutput{}, fmt.Errorf("foo"))

				_, err := manager.DescribeVersions(&api.Addon{
					Name: "my-addon",
				})
				Expect(err).To(MatchError(`failed to describe addon versions: foo`))
				Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
				Expect(*describeAddonVersonsInput.AddonName).To(Equal("my-addon"))
			})
		})
	})

	Describe("DescribeAllVersions", func() {
		It("returns an addon", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
				describeAddonVersonsInput = args[0].(*awseks.DescribeAddonVersionsInput)
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []*awseks.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []*awseks.AddonVersionInfo{
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

			summary, err := manager.DescribeAllVersions()
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(awseks.DescribeAddonVersionsOutput{
				Addons: []*awseks.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []*awseks.AddonVersionInfo{
							{
								AddonVersion: aws.String("1.0"),
							},
							{
								AddonVersion: aws.String("1.1"),
							},
						},
					},
				},
			}.String()))

			Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
			Expect(describeAddonVersonsInput.AddonName).To(BeNil())
		})

		When("it fails to describe addon versions", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
					describeAddonVersonsInput = args[0].(*awseks.DescribeAddonVersionsInput)
				}).Return(&awseks.DescribeAddonVersionsOutput{}, fmt.Errorf("foo"))

				_, err := manager.DescribeAllVersions()
				Expect(err).To(MatchError(`failed to describe addon versions: foo`))
				Expect(*describeAddonVersonsInput.KubernetesVersion).To(Equal("1.18"))
				Expect(describeAddonVersonsInput.AddonName).To(BeNil())
			})
		})
	})
})
