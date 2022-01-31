package fargate_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/fargate"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Fargate", func() {
	var (
		mockProvider     *mockprovider.MockProvider
		cfg              *api.ClusterConfig
		fargateManager   *fargate.Manager
		fakeStackManager *fakes.FakeStackManager
		clusterName      string
		fakeClientSet    *fake.Clientset
	)

	BeforeEach(func() {
		mockProvider = mockprovider.NewMockProvider()
		fakeStackManager = new(fakes.FakeStackManager)
		cfg = api.NewClusterConfig()
		cfg.FargateProfiles = []*api.FargateProfile{
			{
				Name: "fp-1",
				Selectors: []api.FargateProfileSelector{
					{
						Namespace: "default",
					},
				},
			},
		}

		clusterName = "my-cluster"
		cfg.Metadata.Name = clusterName

		fargateManager = fargate.New(cfg, &eks.ClusterProvider{Provider: mockProvider, Status: &eks.ProviderStatus{}}, fakeStackManager)
		fakeClientSet = fake.NewSimpleClientset()

		fargateManager.SetNewClientSet(func() (kubernetes.Interface, error) {
			return fakeClientSet, nil
		})

		mockProvider.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
			Expect(*input.Name).To(Equal(clusterName))
			return true
		})).Return(&awseks.DescribeClusterOutput{
			Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive),
		}, nil)
	})

	Context("owned cluster", func() {
		When("creating a farage role without specifying a role", func() {
			When("the fargate role doesn't exist", func() {
				BeforeEach(func() {
					fakeStackManager.ListStacksReturns(nil, nil)
					fakeStackManager.DescribeClusterStackReturns(&cloudformation.Stack{
						Outputs: []*cloudformation.Output{
							{
								OutputKey:   aws.String("VPC"),
								OutputValue: aws.String("vpc-123"),
							},
							{
								OutputKey:   aws.String("SecurityGroup"),
								OutputValue: aws.String("sg-123"),
							},
						},
					}, nil)
					fakeStackManager.CreateStackStub = func(_ string, _ builder.ResourceSet, _ map[string]string, _ map[string]string, errchan chan error) error {
						go func() {
							errchan <- nil
						}()
						return nil
					}

					fakeStackManager.RefreshFargatePodExecutionRoleARNStub = func() error {
						cfg.IAM.FargatePodExecutionRoleARN = aws.String("fargate-role-arn")
						return nil
					}

					mockProvider.MockEKS().On("CreateFargateProfile", &awseks.CreateFargateProfileInput{
						PodExecutionRoleArn: aws.String("fargate-role-arn"),
						ClusterName:         &clusterName,
						Selectors: []*awseks.FargateProfileSelector{
							{
								Namespace: aws.String("default"),
							},
						},
						FargateProfileName: aws.String("fp-1"),
					}).Return(nil, nil)

					mockProvider.MockEKS().On("DescribeFargateProfile", &awseks.DescribeFargateProfileInput{
						ClusterName:        &clusterName,
						FargateProfileName: aws.String("fp-1"),
					}).Return(&awseks.DescribeFargateProfileOutput{
						FargateProfile: &awseks.FargateProfile{
							Status: aws.String(awseks.FargateProfileStatusActive),
						},
					}, nil)
				})

				It("creates the fargateprofile using the newly created role", func() {
					err := fargateManager.Create()
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
					name, stack, _, _, _ := fakeStackManager.CreateStackArgsForCall(0)
					Expect(name).To(Equal("eksctl-my-cluster-fargate"))
					output, err := stack.RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("AWS::IAM::Role"))
					Expect(string(output)).To(ContainSubstring("FargatePodExecutionRole"))
					Expect(fakeStackManager.RefreshFargatePodExecutionRoleARNCallCount()).To(Equal(1))
				})
			})

			When("the fargate role exists in the cluster stack", func() {
				BeforeEach(func() {
					fakeStackManager.DescribeClusterStackReturns(&cloudformation.Stack{
						Outputs: []*cloudformation.Output{
							{
								OutputKey:   aws.String("VPC"),
								OutputValue: aws.String("vpc-123"),
							},
							{
								OutputKey:   aws.String("SecurityGroup"),
								OutputValue: aws.String("sg-123"),
							},
							{
								OutputKey:   aws.String("FargatePodExecutionRoleARN"),
								OutputValue: aws.String("fargate-existing-role-arn"),
							},
						},
					}, nil)

					fakeStackManager.RefreshFargatePodExecutionRoleARNStub = func() error {
						cfg.IAM.FargatePodExecutionRoleARN = aws.String("fargate-existing-role-arn")
						return nil
					}

					mockProvider.MockEKS().On("CreateFargateProfile", &awseks.CreateFargateProfileInput{
						PodExecutionRoleArn: aws.String("fargate-existing-role-arn"),
						ClusterName:         &clusterName,
						Selectors: []*awseks.FargateProfileSelector{
							{
								Namespace: aws.String("default"),
							},
						},
						FargateProfileName: aws.String("fp-1"),
					}).Return(nil, nil)

					mockProvider.MockEKS().On("DescribeFargateProfile", &awseks.DescribeFargateProfileInput{
						ClusterName:        &clusterName,
						FargateProfileName: aws.String("fp-1"),
					}).Return(&awseks.DescribeFargateProfileOutput{
						FargateProfile: &awseks.FargateProfile{
							Status: aws.String(awseks.FargateProfileStatusActive),
						},
					}, nil)
				})

				It("creates the fargate profile using the existing role", func() {
					err := fargateManager.Create()
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
					Expect(fakeStackManager.RefreshFargatePodExecutionRoleARNCallCount()).To(Equal(1))
				})
			})

			When("the fargate role exists in a separate stack", func() {
				BeforeEach(func() {
					fakeStackManager.GetFargateStackReturns(&manager.Stack{
						StackName: aws.String("fargate"),
					}, nil)

					fakeStackManager.DescribeClusterStackReturns(&cloudformation.Stack{
						Outputs: []*cloudformation.Output{
							{
								OutputKey:   aws.String("VPC"),
								OutputValue: aws.String("vpc-123"),
							},
							{
								OutputKey:   aws.String("SecurityGroup"),
								OutputValue: aws.String("sg-123"),
							},
						},
					}, nil)
					fakeStackManager.CreateStackStub = func(_ string, _ builder.ResourceSet, _ map[string]string, _ map[string]string, errchan chan error) error {
						go func() {
							errchan <- nil
						}()
						return nil
					}

					fakeStackManager.RefreshFargatePodExecutionRoleARNStub = func() error {
						cfg.IAM.FargatePodExecutionRoleARN = aws.String("fargate-role-arn")
						return nil
					}

					mockProvider.MockEKS().On("CreateFargateProfile", &awseks.CreateFargateProfileInput{
						PodExecutionRoleArn: aws.String("fargate-role-arn"),
						ClusterName:         &clusterName,
						Selectors: []*awseks.FargateProfileSelector{
							{
								Namespace: aws.String("default"),
							},
						},
						FargateProfileName: aws.String("fp-1"),
					}).Return(nil, nil)

					mockProvider.MockEKS().On("DescribeFargateProfile", &awseks.DescribeFargateProfileInput{
						ClusterName:        &clusterName,
						FargateProfileName: aws.String("fp-1"),
					}).Return(&awseks.DescribeFargateProfileOutput{
						FargateProfile: &awseks.FargateProfile{
							Status: aws.String(awseks.FargateProfileStatusActive),
						},
					}, nil)
				})

				It("creates the fargateprofile using the existing role", func() {
					err := fargateManager.Create()
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
					Expect(fakeStackManager.RefreshFargatePodExecutionRoleARNCallCount()).To(Equal(1))
				})
			})
		})
	})
})
