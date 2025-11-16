package manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("IAM Service Accounts", func() {
	var (
		p            *mockprovider.MockProvider
		cfg          *api.ClusterConfig
		stackManager *StackCollection
		ctx          context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		p = mockprovider.NewMockProvider()
		cfg = &api.ClusterConfig{
			Metadata: &api.ClusterMeta{
				Name: "test-cluster",
			},
		}
		stackManager = NewStackCollection(p, cfg).(*StackCollection)
	})

	Describe("GetIAMServiceAccounts", func() {
		It("returns service accounts when found", func() {
			// Setup mock response for DescribeIAMServiceAccountStacks
			testCases := []iamServiceAccountTestCase{
				{
					Name:      "app-service-account",
					Namespace: "default",
				},
				{
					Name:      "monitoring-service-account",
					Namespace: "monitoring",
				},
			}

			// Mock the ListStacks call that DescribeIAMServiceAccountStacks uses
			stacks := getStacks(testCases)
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(
				&cfn.ListStacksOutput{
					StackSummaries: getStackSummaries(stacks),
				}, nil)

			// Mock the DescribeStacks call that outputs.Collect uses
			for _, stack := range stacks {
				p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.MatchedBy(func(input interface{}) bool {
					if describeInput, ok := input.(*cfn.DescribeStacksInput); ok {
						return describeInput.StackName != nil && *describeInput.StackName == *stack.StackName
					}
					return false
				})).Return(&cfn.DescribeStacksOutput{
					Stacks: []types.Stack{
						{
							StackName:    stack.StackName,
							CreationTime: stack.CreationTime,
							StackStatus:  stack.StackStatus,
							Tags:         stack.Tags,
							Outputs: []types.Output{
								{
									OutputKey: aws.String(outputs.IAMServiceAccountRoleName),
									OutputValue: aws.String(fmt.Sprintf("arn:aws:iam::123456789012:role/eksctl-%s-%s-%s",
										cfg.Metadata.Name, *stack.Tags[0].Value, *stack.Tags[1].Value)),
								},
							},
						},
					},
				}, nil)
			}

			// Call the function with no filters - this should use the real implementation
			serviceAccounts, err := stackManager.GetIAMServiceAccounts(ctx, "", "")

			// Verify results
			Expect(err).NotTo(HaveOccurred())
			Expect(serviceAccounts).To(HaveLen(2))

			// Verify the service accounts have the expected values
			for _, sa := range serviceAccounts {
				Expect(sa.Status.RoleARN).NotTo(BeNil())
				switch sa.Name {
				case "app-service-account":
					Expect(sa.Namespace).To(Equal("default"))
				case "monitoring-service-account":
					Expect(sa.Namespace).To(Equal("monitoring"))
				default:
					Fail(fmt.Sprintf("Unexpected service account name: %s", sa.Name))
				}
			}
		})

		It("filters service accounts by name", func() {
			// Setup mock response
			testCases := []iamServiceAccountTestCase{
				{
					Name:      "monitoring-service-account",
					Namespace: "monitoring",
				},
				{
					Name:      "app-service-account",
					Namespace: "default",
				},
				{
					Name:      "another-app-service-account",
					Namespace: "default",
				},
			}

			// Mock the ListStacks call
			stacks := getStacks(testCases)
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(
				&cfn.ListStacksOutput{
					StackSummaries: getStackSummaries(stacks),
				}, nil)

			// Mock the DescribeStacks call for each stack
			for _, stack := range stacks {
				p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.MatchedBy(func(input interface{}) bool {
					if describeInput, ok := input.(*cfn.DescribeStacksInput); ok {
						return describeInput.StackName != nil && *describeInput.StackName == *stack.StackName
					}
					return false
				})).Return(&cfn.DescribeStacksOutput{
					Stacks: []types.Stack{
						{
							StackName:    stack.StackName,
							CreationTime: stack.CreationTime,
							StackStatus:  stack.StackStatus,
							Tags:         stack.Tags,
							Outputs: []types.Output{
								{
									OutputKey: aws.String(outputs.IAMServiceAccountRoleName),
									OutputValue: aws.String(fmt.Sprintf("arn:aws:iam::123456789012:role/eksctl-%s-%s-%s",
										cfg.Metadata.Name, *stack.Tags[0].Value, *stack.Tags[1].Value)),
								},
							},
						},
					},
				}, nil)
			}

			// Call the function with name filter
			serviceAccounts, err := stackManager.GetIAMServiceAccounts(ctx, "app-service-account", "")

			// Verify results
			Expect(err).NotTo(HaveOccurred())
			Expect(serviceAccounts).To(HaveLen(1))
			Expect(serviceAccounts[0].Name).To(Equal("app-service-account"))
			Expect(serviceAccounts[0].Namespace).To(Equal("default"))
		})

		It("filters service accounts by namespace", func() {
			// Setup mock response
			testCases := []iamServiceAccountTestCase{
				{
					Name:      "app-service-account",
					Namespace: "default",
				},
				{
					Name:      "monitoring-service-account",
					Namespace: "monitoring",
				},
				{
					Name:      "observability-account",
					Namespace: "monitoring",
				},
				{
					Name:      "another-app-service-account",
					Namespace: "default",
				},
			}

			// Mock the ListStacks call
			stacks := getStacks(testCases)
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(
				&cfn.ListStacksOutput{
					StackSummaries: getStackSummaries(stacks),
				}, nil)

			// Mock the DescribeStacks call for each stack
			for _, stack := range stacks {
				p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.MatchedBy(func(input interface{}) bool {
					if describeInput, ok := input.(*cfn.DescribeStacksInput); ok {
						return describeInput.StackName != nil && *describeInput.StackName == *stack.StackName
					}
					return false
				})).Return(&cfn.DescribeStacksOutput{
					Stacks: []types.Stack{
						{
							StackName:    stack.StackName,
							CreationTime: stack.CreationTime,
							StackStatus:  stack.StackStatus,
							Tags:         stack.Tags,
							Outputs: []types.Output{
								{
									OutputKey: aws.String(outputs.IAMServiceAccountRoleName),
									OutputValue: aws.String(fmt.Sprintf("arn:aws:iam::123456789012:role/eksctl-%s-%s-%s",
										cfg.Metadata.Name, *stack.Tags[0].Value, *stack.Tags[1].Value)),
								},
							},
						},
					},
				}, nil)
			}

			// Call the function with namespace filter
			serviceAccounts, err := stackManager.GetIAMServiceAccounts(ctx, "", "monitoring")

			// Verify results
			Expect(err).NotTo(HaveOccurred())
			Expect(serviceAccounts).To(HaveLen(2))
			for _, sa := range serviceAccounts {
				Expect(sa.Namespace).To(Equal("monitoring"))
			}
		})

		It("filters service accounts by both name and namespace", func() {
			// Setup mock response
			testCases := []iamServiceAccountTestCase{
				{
					Name:      "app-service-account",
					Namespace: "default",
				},
				{
					Name:      "app-service-account",
					Namespace: "monitoring",
				},
				{
					Name:      "another-app-service-account",
					Namespace: "monitoring",
				},
				{
					Name:      "another-app-service-account",
					Namespace: "default",
				},
			}

			// Mock the ListStacks call
			stacks := getStacks(testCases)
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(
				&cfn.ListStacksOutput{
					StackSummaries: getStackSummaries(stacks),
				}, nil)

			// Mock the DescribeStacks call for each stack
			for _, stack := range stacks {
				p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.MatchedBy(func(input interface{}) bool {
					if describeInput, ok := input.(*cfn.DescribeStacksInput); ok {
						return describeInput.StackName != nil && *describeInput.StackName == *stack.StackName
					}
					return false
				})).Return(&cfn.DescribeStacksOutput{
					Stacks: []types.Stack{
						{
							StackName:    stack.StackName,
							CreationTime: stack.CreationTime,
							StackStatus:  stack.StackStatus,
							Tags:         stack.Tags,
							Outputs: []types.Output{
								{
									OutputKey: aws.String(outputs.IAMServiceAccountRoleName),
									OutputValue: aws.String(fmt.Sprintf("arn:aws:iam::123456789012:role/eksctl-%s-%s-%s",
										cfg.Metadata.Name, *stack.Tags[0].Value, *stack.Tags[1].Value)),
								},
							},
						},
					},
				}, nil)
			}

			serviceAccounts, err := stackManager.GetIAMServiceAccounts(ctx, "app-service-account", "default")

			Expect(err).NotTo(HaveOccurred())
			Expect(serviceAccounts).To(HaveLen(1))
			Expect(serviceAccounts[0].Name).To(Equal("app-service-account"))
			Expect(serviceAccounts[0].Namespace).To(Equal("default"))
		})

		It("handles errors from the CloudFormation API", func() {
			// Setup mock error response
			expectedError := errors.New("failed to describe IAM service account stacks")
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

			// Call the function
			serviceAccounts, err := stackManager.GetIAMServiceAccounts(ctx, "", "")

			// Verify results
			Expect(err).To(MatchError(expectedError))
			Expect(serviceAccounts).To(BeNil())
		})

		It("returns empty slice when no service accounts match filters", func() {
			// Setup mock response with stacks that won't match our filter
			testCases := []iamServiceAccountTestCase{
				{
					Name:      "app-service-account",
					Namespace: "default",
				},
				{
					Name:      "monitoring-service-account",
					Namespace: "monitoring",
				},
			}

			// Mock the ListStacks call
			stacks := getStacks(testCases)
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(
				&cfn.ListStacksOutput{
					StackSummaries: getStackSummaries(stacks),
				}, nil)

			// Mock the DescribeStacks call for each stack
			for _, stack := range stacks {
				p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.MatchedBy(func(input interface{}) bool {
					if describeInput, ok := input.(*cfn.DescribeStacksInput); ok {
						return describeInput.StackName != nil && *describeInput.StackName == *stack.StackName
					}
					return false
				})).Return(&cfn.DescribeStacksOutput{
					Stacks: []types.Stack{
						{
							StackName:    stack.StackName,
							CreationTime: stack.CreationTime,
							StackStatus:  stack.StackStatus,
							Tags:         stack.Tags,
							Outputs: []types.Output{
								{
									OutputKey: aws.String(outputs.IAMServiceAccountRoleName),
									OutputValue: aws.String(fmt.Sprintf("arn:aws:iam::123456789012:role/eksctl-%s-%s-%s",
										cfg.Metadata.Name, *stack.Tags[0].Value, *stack.Tags[1].Value)),
								},
							},
						},
					},
				}, nil)
			}

			serviceAccounts, err := stackManager.GetIAMServiceAccounts(ctx, "non-existent", "non-existent")

			Expect(err).NotTo(HaveOccurred())
			Expect(serviceAccounts).To(BeEmpty())
		})
	})
})

type iamServiceAccountTestCase struct {
	Name      string
	Namespace string
}

func getStacks(testCases []iamServiceAccountTestCase) []types.Stack {
	stacks := make([]types.Stack, 0)

	for _, testCase := range testCases {
		stackName := fmt.Sprintf("eksctl-test-cluster-addon-iamserviceaccount-%s-%s", testCase.Namespace, testCase.Name)
		stack := types.Stack{
			StackName:    aws.String(stackName),
			CreationTime: aws.Time(time.Now()),
			StackStatus:  types.StackStatusCreateComplete,
			Tags: []types.Tag{
				{
					Key:   aws.String(api.IAMServiceAccountNameTag),
					Value: aws.String(fmt.Sprintf("%s/%s", testCase.Namespace, testCase.Name)),
				},
				{
					Key:   aws.String("Namespace"),
					Value: aws.String(testCase.Namespace),
				},
				{
					Key:   aws.String("ServiceAccount"),
					Value: aws.String(testCase.Name),
				},
			},
		}
		stacks = append(stacks, stack)
	}

	return stacks
}

func getStackSummaries(stacks []types.Stack) []types.StackSummary {
	summaries := make([]types.StackSummary, 0, len(stacks))

	for _, stack := range stacks {
		summary := types.StackSummary{
			StackName:    stack.StackName,
			CreationTime: stack.CreationTime,
			StackStatus:  stack.StackStatus,
		}
		summaries = append(summaries, summary)
	}

	return summaries
}
