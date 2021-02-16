package addon_test

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/testutils"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/stretchr/testify/mock"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Create", func() {
	var (
		manager                *addon.Manager
		withOIDC               bool
		oidc                   *iamoidc.OpenIDConnectManager
		fakeStackManager       *fakes.FakeStackManager
		mockProvider           *mockprovider.MockProvider
		createAddonInput       *awseks.CreateAddonInput
		returnedErr            error
		createStackReturnValue error
		rawClient              *testutils.FakeRawClient
	)

	BeforeEach(func() {
		withOIDC = true
		returnedErr = nil
		fakeStackManager = new(fakes.FakeStackManager)
		mockProvider = mockprovider.NewMockProvider()
		createStackReturnValue = nil

		fakeStackManager.CreateStackStub = func(_ string, rs builder.ResourceSet, _ map[string]string, _ map[string]string, errs chan error) error {
			go func() {
				errs <- nil
			}()
			return createStackReturnValue
		}

		sampleAddons := testutils.LoadSamples("testdata/aws-node.json")

		rawClient = testutils.NewFakeRawClient()

		rawClient.AssumeObjectsMissing = true

		for _, item := range sampleAddons {
			rc, err := rawClient.NewRawResource(item)
			Expect(err).ToNot(HaveOccurred())
			_, err = rc.CreateOrReplace(false)
			Expect(err).ToNot(HaveOccurred())
		}

		ct := rawClient.Collection

		Expect(ct.Updated()).To(BeEmpty())
		Expect(ct.Created()).ToNot(BeEmpty())
		Expect(ct.CreatedItems()).To(HaveLen(10))
	})

	JustBeforeEach(func() {
		var err error

		oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws")
		Expect(err).ToNot(HaveOccurred())
		oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

		mockProvider.MockEKS().On("CreateAddon", mock.Anything).Run(func(args mock.Arguments) {
			Expect(args).To(HaveLen(1))
			Expect(args[0]).To(BeAssignableToTypeOf(&awseks.CreateAddonInput{}))
			createAddonInput = args[0].(*awseks.CreateAddonInput)
		}).Return(nil, returnedErr)

		manager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}, mockProvider.EKS(), fakeStackManager, withOIDC, oidc, rawClient.ClientSet())
		Expect(err).NotTo(HaveOccurred())
	})

	When("it fails to create addon", func() {
		BeforeEach(func() {
			returnedErr = fmt.Errorf("foo")
		})
		It("returns an error", func() {
			err := manager.Create(&api.Addon{
				Name:    "my-addon",
				Version: "1.0",
			})
			Expect(err).To(MatchError(`failed to create addon "my-addon": foo`))

		})
	})

	When("OIDC is disabled", func() {
		BeforeEach(func() {
			withOIDC = false
		})
		It("creates the addons but not the policies", func() {
			err := manager.Create(&api.Addon{
				Name:             "my-addon",
				Version:          "1.0",
				AttachPolicyARNs: []string{"arn-1"},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
			Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*createAddonInput.AddonVersion).To(Equal("1.0"))
			Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
		})
	})

	When("force is true", func() {
		BeforeEach(func() {
			withOIDC = false
		})

		It("creates the addons but not the policies", func() {
			err := manager.Create(&api.Addon{
				Name:             "my-addon",
				Version:          "1.0",
				AttachPolicyARNs: []string{"arn-1"},
				Force:            true,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
			Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*createAddonInput.AddonVersion).To(Equal("1.0"))
			Expect(*createAddonInput.ResolveConflicts).To(Equal("overwrite"))
			Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
		})
	})

	When("No policy/role is specified", func() {
		When("we don't know the recommended policies for the specified addon", func() {
			It("does not provide a role", func() {
				err := manager.Create(&api.Addon{
					Name:    "my-addon",
					Version: "1.0",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
				Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*createAddonInput.AddonVersion).To(Equal("1.0"))
				Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
			})
		})

		When("we know the recommended policies for the specified addon", func() {
			BeforeEach(func() {
				fakeStackManager.CreateStackStub = func(_ string, rs builder.ResourceSet, _ map[string]string, _ map[string]string, errs chan error) error {
					go func() {
						errs <- nil
					}()
					Expect(rs).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
					rs.(*builder.IAMRoleResourceSet).OutputRole = "role-arn"
					return createStackReturnValue
				}

			})
			It("creates a role with the recommended policies and attaches it to the addon", func() {
				err := manager.Create(&api.Addon{
					Name:    "vpc-cni",
					Version: "1.0",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
				name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
				Expect(name).To(Equal("eksctl-my-cluster-addon-vpc-cni"))
				Expect(resourceSet).NotTo(BeNil())
				Expect(tags).To(Equal(map[string]string{
					api.AddonNameTag: "vpc-cni",
				}))
				output, err := resourceSet.RenderJSON()
				Expect(err).NotTo(HaveOccurred())
				Expect(string(output)).To(ContainSubstring("arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"))
				Expect(string(output)).To(ContainSubstring(":sub\":\"system:serviceaccount:kube-system:aws-node"))
				Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*createAddonInput.AddonName).To(Equal("vpc-cni"))
				Expect(*createAddonInput.AddonVersion).To(Equal("1.0"))
				Expect(*createAddonInput.ServiceAccountRoleArn).To(Equal("role-arn"))
			})
		})
	})

	When("attachPolicyARNs is configured", func() {
		It("uses AttachPolicyARNS to create a role to attach to the addon", func() {
			err := manager.Create(&api.Addon{
				Name:             "my-addon",
				Version:          "1.0",
				AttachPolicyARNs: []string{"arn-1"},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
			name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
			Expect(name).To(Equal("eksctl-my-cluster-addon-my-addon"))
			Expect(resourceSet).NotTo(BeNil())
			Expect(tags).To(Equal(map[string]string{
				api.AddonNameTag: "my-addon",
			}))
			output, err := resourceSet.RenderJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("arn-1"))
		})
	})

	When("AttachPolicy is configured", func() {
		It("uses AttachPolicy to create a role to attach to the addon", func() {
			err := manager.Create(&api.Addon{
				Name:    "my-addon",
				Version: "1.0",
				AttachPolicy: api.InlineDocument{
					"foo": "policy-bar",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
			name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
			Expect(name).To(Equal("eksctl-my-cluster-addon-my-addon"))
			Expect(resourceSet).NotTo(BeNil())
			Expect(tags).To(Equal(map[string]string{
				api.AddonNameTag: "my-addon",
			}))
			output, err := resourceSet.RenderJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("policy-bar"))
		})
	})

	When("serviceAccountRoleARN is configured", func() {
		It("uses the serviceAccountRoleARN to create the addon", func() {
			err := manager.Create(&api.Addon{
				Name:                  "my-addon",
				Version:               "1.0",
				ServiceAccountRoleARN: "foo",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
			Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*createAddonInput.AddonVersion).To(Equal("1.0"))
			Expect(*createAddonInput.ServiceAccountRoleArn).To(Equal("foo"))
		})
	})
})
