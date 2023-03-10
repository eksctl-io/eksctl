package fargate_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/fargate"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Fargate", func() {
	var (
		mockProvider     *mockprovider.MockProvider
		cfg              *api.ClusterConfig
		fargateManager   *fargate.Manager
		fakeStackManager *fakes.FakeStackManager
		clusterName      string
		region           string
		accountID        string
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
		region = "eu-north-1"
		accountID = "111122223333"
		cfg.Metadata.Name = clusterName
		cfg.Metadata.Region = region
		cfg.Status = &api.ClusterStatus{
			ARN: fmt.Sprintf("arn:aws:eks:%s:%s:cluster/%s", region, accountID, clusterName),
		}
		ctl := &eks.ClusterProvider{AWSProvider: mockProvider, Status: &eks.ProviderStatus{
			ClusterInfo: &eks.ClusterInfo{
				Cluster: &ekstypes.Cluster{
					Status:  ekstypes.ClusterStatusActive,
					Version: aws.String("1.22"),
				},
			},
		}}
		fargateManager = fargate.New(cfg, ctl, fakeStackManager)
		fakeClientSet = fake.NewSimpleClientset()

		fargateManager.SetNewClientSet(func() (kubernetes.Interface, error) {
			return fakeClientSet, nil
		})

		mockProvider.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
			Expect(*input.Name).To(Equal(clusterName))
			return true
		})).Return(&awseks.DescribeClusterOutput{
			Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
		}, nil)
	})

	Context("owned cluster", func() {
		When("creating a farage role without specifying a role", func() {
			When("the fargate role doesn't exist", func() {
				BeforeEach(func() {
					fakeStackManager.ListStacksWithStatusesReturns(nil, nil)
					fakeStackManager.DescribeClusterStackIfExistsReturns(&types.Stack{
						Outputs: []types.Output{
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
					fakeStackManager.CreateStackStub = func(_ context.Context, _ string, _ builder.ResourceSetReader, _ map[string]string, _ map[string]string, errchan chan error) error {
						go func() {
							errchan <- nil
						}()
						return nil
					}

					fakeStackManager.RefreshFargatePodExecutionRoleARNStub = func(_ context.Context) error {
						cfg.IAM.FargatePodExecutionRoleARN = aws.String("fargate-role-arn")
						return nil
					}

					mockProvider.MockEKS().On("CreateFargateProfile", mock.Anything, &awseks.CreateFargateProfileInput{
						PodExecutionRoleArn: aws.String("fargate-role-arn"),
						ClusterName:         &clusterName,
						Selectors: []ekstypes.FargateProfileSelector{
							{
								Namespace: aws.String("default"),
							},
						},
						FargateProfileName: aws.String("fp-1"),
					}).Return(nil, nil)

					mockProvider.MockEKS().On("DescribeFargateProfile", mock.Anything, &awseks.DescribeFargateProfileInput{
						ClusterName:        &clusterName,
						FargateProfileName: aws.String("fp-1"),
					}).Return(&awseks.DescribeFargateProfileOutput{
						FargateProfile: &ekstypes.FargateProfile{
							Status: ekstypes.FargateProfileStatusActive,
						},
					}, nil)

					mockProvider.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
						ClusterName: &clusterName,
					}).Return(&awseks.ListFargateProfilesOutput{
						FargateProfileNames: []string{
							"fp-1",
						},
					}, nil)
				})

				It("creates the fargateprofile using the newly created role", func() {
					err := fargateManager.Create(context.Background())
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
					_, name, stack, _, _, _ := fakeStackManager.CreateStackArgsForCall(0)
					Expect(name).To(Equal("eksctl-my-cluster-fargate"))
					output, err := stack.RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("AWS::IAM::Role"))
					Expect(string(output)).To(ContainSubstring("FargatePodExecutionRole"))
					Expect(string(output)).To(ContainSubstring(fmt.Sprintf("\"aws:SourceArn\": \"arn:aws:eks:%s:%s:fargateprofile/%s/*\"", region, accountID, clusterName)))
					Expect(fakeStackManager.RefreshFargatePodExecutionRoleARNCallCount()).To(Equal(1))
				})
			})

			When("the fargate role exists in the cluster stack", func() {
				BeforeEach(func() {
					fakeStackManager.DescribeClusterStackIfExistsReturns(&types.Stack{
						Outputs: []types.Output{
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

					fakeStackManager.RefreshFargatePodExecutionRoleARNStub = func(_ context.Context) error {
						cfg.IAM.FargatePodExecutionRoleARN = aws.String("fargate-existing-role-arn")
						return nil
					}

					mockProvider.MockEKS().On("CreateFargateProfile", mock.Anything, &awseks.CreateFargateProfileInput{
						PodExecutionRoleArn: aws.String("fargate-existing-role-arn"),
						ClusterName:         &clusterName,
						Selectors: []ekstypes.FargateProfileSelector{
							{
								Namespace: aws.String("default"),
							},
						},
						FargateProfileName: aws.String("fp-1"),
					}).Return(nil, nil)

					mockProvider.MockEKS().On("DescribeFargateProfile", mock.Anything, &awseks.DescribeFargateProfileInput{
						ClusterName:        &clusterName,
						FargateProfileName: aws.String("fp-1"),
					}).Return(&awseks.DescribeFargateProfileOutput{
						FargateProfile: &ekstypes.FargateProfile{
							Status: ekstypes.FargateProfileStatusActive,
						},
					}, nil)

					mockProvider.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
						ClusterName: &clusterName,
					}).Return(&awseks.ListFargateProfilesOutput{
						FargateProfileNames: []string{
							"fp-1",
						},
					}, nil)
				})

				It("creates the fargate profile using the existing role", func() {
					err := fargateManager.Create(context.Background())
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

					fakeStackManager.DescribeClusterStackIfExistsReturns(&types.Stack{
						Outputs: []types.Output{
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
					fakeStackManager.CreateStackStub = func(_ context.Context, _ string, _ builder.ResourceSetReader, _ map[string]string, _ map[string]string, errchan chan error) error {
						go func() {
							errchan <- nil
						}()
						return nil
					}

					fakeStackManager.RefreshFargatePodExecutionRoleARNStub = func(_ context.Context) error {
						cfg.IAM.FargatePodExecutionRoleARN = aws.String("fargate-role-arn")
						return nil
					}

					mockProvider.MockEKS().On("CreateFargateProfile", mock.Anything, &awseks.CreateFargateProfileInput{
						PodExecutionRoleArn: aws.String("fargate-role-arn"),
						ClusterName:         &clusterName,
						Selectors: []ekstypes.FargateProfileSelector{
							{
								Namespace: aws.String("default"),
							},
						},
						FargateProfileName: aws.String("fp-1"),
					}).Return(nil, nil)

					mockProvider.MockEKS().On("DescribeFargateProfile", mock.Anything, &awseks.DescribeFargateProfileInput{
						ClusterName:        &clusterName,
						FargateProfileName: aws.String("fp-1"),
					}).Return(&awseks.DescribeFargateProfileOutput{
						FargateProfile: &ekstypes.FargateProfile{
							Status: ekstypes.FargateProfileStatusActive,
						},
					}, nil)

					mockProvider.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
						ClusterName: &clusterName,
					}).Return(&awseks.ListFargateProfilesOutput{
						FargateProfileNames: []string{
							"fp-1",
						},
					}, nil)
				})

				It("creates the fargateprofile using the existing role", func() {
					err := fargateManager.Create(context.Background())
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(0))
					Expect(fakeStackManager.RefreshFargatePodExecutionRoleARNCallCount()).To(Equal(1))
				})
			})
		})
	})
})
