package addon_test

import (
	"fmt"
	"time"

	"github.com/weaveworks/eksctl/pkg/testutils"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/aws"
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
		clusterConfig          *api.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = &api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}
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
			Expect(err).NotTo(HaveOccurred())
			_, err = rc.CreateOrReplace(false)
			Expect(err).NotTo(HaveOccurred())
		}

		ct := rawClient.Collection

		Expect(ct.Updated()).To(BeEmpty())
		Expect(ct.Created()).NotTo(BeEmpty())
		Expect(ct.CreatedItems()).To(HaveLen(10))
	})

	JustBeforeEach(func() {
		var err error

		oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
		Expect(err).NotTo(HaveOccurred())
		oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

		mockProvider.MockEKS().On("CreateAddon", mock.Anything).Run(func(args mock.Arguments) {
			Expect(args).To(HaveLen(1))
			Expect(args[0]).To(BeAssignableToTypeOf(&awseks.CreateAddonInput{}))
			createAddonInput = args[0].(*awseks.CreateAddonInput)
		}).Return(nil, returnedErr)

		manager, err = addon.New(clusterConfig, mockProvider.EKS(), fakeStackManager, withOIDC, oidc, rawClient.ClientSet(), 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())
		manager.SetTimeout(time.Second)

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
							AddonVersion: aws.String("v1.0.0-eksbuild.1"),
						},
						{
							AddonVersion: aws.String("v1.7.5-eksbuild.1"),
						},
						{
							AddonVersion: aws.String("v1.7.5-eksbuild.2"),
						},
						{
							//not sure if all versions come with v prefix or not, so test a mix
							AddonVersion: aws.String("v1.7.7-eksbuild.2"),
						},
						{
							AddonVersion: aws.String("v1.7.6"),
						},
					},
				},
			},
		}, nil)
	})

	When("it fails to create addon", func() {
		BeforeEach(func() {
			returnedErr = fmt.Errorf("foo")
		})
		It("returns an error", func() {
			err := manager.Create(&api.Addon{
				Name:    "my-addon",
				Version: "v1.0.0-eksbuild.1",
			}, false)
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
				Version:          "v1.0.0-eksbuild.1",
				AttachPolicyARNs: []string{"arn-1"},
			}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
			Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
			Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
		})
	})

	When("version is specified", func() {
		When("the versions are valid", func() {
			BeforeEach(func() {
				withOIDC = false
			})

			When("version is set to a numeric value", func() {
				It("discovers and uses the latest available version", func() {
					err := manager.Create(&api.Addon{
						Name:             "my-addon",
						Version:          "1.7.5",
						AttachPolicyARNs: []string{"arn-1"},
					}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
					Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*createAddonInput.AddonVersion).To(Equal("v1.7.5-eksbuild.2"))
					Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
				})
			})

			When("version is set to an alphanumeric value", func() {
				It("discovers and uses the latest available version", func() {
					err := manager.Create(&api.Addon{
						Name:             "my-addon",
						Version:          "1.7.5-eksbuild",
						AttachPolicyARNs: []string{"arn-1"},
					}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
					Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*createAddonInput.AddonVersion).To(Equal("v1.7.5-eksbuild.2"))
					Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
				})
			})

			When("version is set to latest", func() {
				It("discovers and uses the latest available version", func() {
					err := manager.Create(&api.Addon{
						Name:             "my-addon",
						Version:          "latest",
						AttachPolicyARNs: []string{"arn-1"},
					}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
					Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*createAddonInput.AddonVersion).To(Equal("v1.7.7-eksbuild.2"))
					Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
				})
			})

			When("the version is set to a version that does not exist", func() {
				It("returns an error", func() {
					err := manager.Create(&api.Addon{
						Name:             "my-addon",
						Version:          "1.7.8",
						AttachPolicyARNs: []string{"arn-1"},
					}, false)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("no version(s) found matching \"1.7.8\" for \"my-addon\"")))
				})
			})
		})

		When("the versions are invalid", func() {
			BeforeEach(func() {
				withOIDC = false

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
									AddonVersion: aws.String("v1.7.5-eksbuild.1"),
								},
								{
									//not sure if all versions come with v prefix or not, so test a mix
									AddonVersion: aws.String("v1.7.7-eksbuild.1"),
								},
								{
									AddonVersion: aws.String("totally not semver"),
								},
							},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				err := manager.Create(&api.Addon{
					Name:             "my-addon",
					Version:          "latest",
					AttachPolicyARNs: []string{"arn-1"},
				}, false)
				Expect(err).To(MatchError(ContainSubstring("failed to parse version \"totally not semver\":")))
			})
		})

		When("there are no versions returned", func() {
			BeforeEach(func() {
				withOIDC = false

				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []*awseks.AddonInfo{
						{
							AddonName:     aws.String("my-addon"),
							Type:          aws.String("type"),
							AddonVersions: []*awseks.AddonVersionInfo{},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				err := manager.Create(&api.Addon{
					Name:             "my-addon",
					Version:          "latest",
					AttachPolicyARNs: []string{"arn-1"},
				}, false)
				Expect(err).To(MatchError(ContainSubstring("no versions available for \"my-addon\"")))
			})
		})
	})

	When("force is true", func() {
		BeforeEach(func() {
			withOIDC = false
		})

		It("creates the addons but not the policies", func() {
			err := manager.Create(&api.Addon{
				Name:             "my-addon",
				Version:          "v1.0.0-eksbuild.1",
				AttachPolicyARNs: []string{"arn-1"},
				Force:            true,
			}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
			Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
			Expect(*createAddonInput.ResolveConflicts).To(Equal("overwrite"))
			Expect(createAddonInput.ServiceAccountRoleArn).To(BeNil())
		})
	})

	When("wait is true", func() {
		When("the addon creation succeeds", func() {
			BeforeEach(func() {
				withOIDC = false
				mockProvider.MockEKS().On("DescribeAddon", mock.Anything).
					Return(&awseks.DescribeAddonOutput{
						Addon: &awseks.Addon{
							AddonName: aws.String("my-addon"),
							Status:    aws.String("ACTIVE"),
						},
					}, nil)
			})

			It("creates the addon and waits for it to be active", func() {
				err := manager.Create(&api.Addon{
					Name:    "my-addon",
					Version: "v1.0.0-eksbuild.1",
				}, true)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
				Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
			})
		})

		When("the addon creation fails", func() {
			BeforeEach(func() {
				withOIDC = false
				mockProvider.MockEKS().On("DescribeAddon", mock.Anything).
					Return(&awseks.DescribeAddonOutput{
						Addon: &awseks.Addon{
							AddonName: aws.String("my-addon"),
							Status:    aws.String("DEGRADED"),
						},
					}, nil)
			})

			It("returns an error", func() {
				err := manager.Create(&api.Addon{
					Name:    "my-addon",
					Version: "v1.0.0-eksbuild.1",
				}, true)
				Expect(err).To(MatchError("timed out waiting for addon \"my-addon\" to become active, status: \"DEGRADED\""))
			})
		})
	})

	When("No policy/role is specified", func() {
		When("we don't know the recommended policies for the specified addon", func() {
			It("does not provide a role", func() {
				err := manager.Create(&api.Addon{
					Name:    "my-addon",
					Version: "v1.0.0-eksbuild.1",
				}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
				Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
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

			When("it's the vpc-cni addon", func() {
				Context("ipv4", func() {
					It("creates a role with the recommended policies and attaches it to the addon", func() {
						err := manager.Create(&api.Addon{
							Name:    "vpc-cni",
							Version: "v1.0.0-eksbuild.1",
						}, false)
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
						Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
						Expect(*createAddonInput.ServiceAccountRoleArn).To(Equal("role-arn"))
					})
				})

				Context("ipv6", func() {
					BeforeEach(func() {
						clusterConfig.VPC = api.NewClusterVPC(false)
						clusterConfig.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
							IPFamily: api.IPV6Family,
						}
					})

					It("creates a role with the recommended policies and attaches it to the addon", func() {
						err := manager.Create(&api.Addon{
							Name:    "vpc-cni",
							Version: "v1.0.0-eksbuild.1",
						}, false)
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
						Expect(string(output)).To(ContainSubstring("AssignIpv6Addresses"))
						Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
						Expect(*createAddonInput.AddonName).To(Equal("vpc-cni"))
						Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
						Expect(*createAddonInput.ServiceAccountRoleArn).To(Equal("role-arn"))
					})
				})
			})

			When("it's the aws-ebs-csi-driver addon", func() {
				It("creates a role with the recommended policies and attaches it to the addon", func() {
					err := manager.Create(&api.Addon{
						Name:    "aws-ebs-csi-driver",
						Version: "v1.0.0-eksbuild.1",
					}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
					name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
					Expect(name).To(Equal("eksctl-my-cluster-addon-aws-ebs-csi-driver"))
					Expect(resourceSet).NotTo(BeNil())
					Expect(tags).To(Equal(map[string]string{
						api.AddonNameTag: "aws-ebs-csi-driver",
					}))
					output, err := resourceSet.RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("PolicyEBSCSIController"))
					Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*createAddonInput.AddonName).To(Equal("aws-ebs-csi-driver"))
					Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
					Expect(*createAddonInput.ServiceAccountRoleArn).To(Equal("role-arn"))
				})
			})
		})
	})

	When("attachPolicyARNs is configured", func() {
		It("uses AttachPolicyARNS to create a role to attach to the addon", func() {
			err := manager.Create(&api.Addon{
				Name:             "my-addon",
				Version:          "v1.0.0-eksbuild.1",
				AttachPolicyARNs: []string{"arn-1"},
			}, false)
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

	When("wellKnownPolicies is configured", func() {
		It("uses wellKnownPolicies to create a role to attach to the addon", func() {
			err := manager.Create(&api.Addon{
				Name:    "my-addon",
				Version: "v1.0.0-eksbuild.1",
				WellKnownPolicies: api.WellKnownPolicies{
					AutoScaler: true,
				},
			}, false)
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
			Expect(string(output)).To(ContainSubstring("autoscaling:SetDesiredCapacity"))
		})
	})

	When("AttachPolicy is configured", func() {
		It("uses AttachPolicy to create a role to attach to the addon", func() {
			err := manager.Create(&api.Addon{
				Name:    "my-addon",
				Version: "v1.0.0-eksbuild.1",
				AttachPolicy: api.InlineDocument{
					"foo": "policy-bar",
				},
			}, false)
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
				Version:               "v1.0.0-eksbuild.1",
				ServiceAccountRoleARN: "foo",
			}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
			Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
			Expect(*createAddonInput.ServiceAccountRoleArn).To(Equal("foo"))
		})
	})

	When("tags are configured", func() {
		It("uses the Tags to create the addon", func() {
			err := manager.Create(&api.Addon{
				Name:    "my-addon",
				Version: "v1.0.0-eksbuild.1",
				Tags:    map[string]string{"foo": "bar", "fox": "brown"},
			}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
			Expect(*createAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*createAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*createAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.1"))
			Expect(*createAddonInput.Tags["foo"]).To(Equal("bar"))
			Expect(*createAddonInput.Tags["fox"]).To(Equal("brown"))
		})
	})
})
