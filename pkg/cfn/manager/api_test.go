package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	asTypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("StackCollection", func() {
	Context("PropagateManagedNodeGroupTagsToASG", func() {
		var (
			asgName string
			ngName  string
			ngTags  map[string]string
			errCh   chan error
			p       *mockprovider.MockProvider
		)
		BeforeEach(func() {
			asgName = "asg-test-name"
			ngName = "ng-test-name"
			ngTags = map[string]string{
				"tag_key_1": "tag_value_1",
			}
			errCh = make(chan error)
			p = mockprovider.NewMockProvider()
		})

		It("can create propagate tag", func() {
			// DescribeTags classic mock
			describeTagsInput := &autoscaling.DescribeTagsInput{
				Filters: []asTypes.Filter{{Name: aws.String(resourceTypeAutoScalingGroup), Values: []string{asgName}}},
			}
			p.MockASG().On("DescribeTags", mock.Anything, describeTagsInput).Return(&autoscaling.DescribeTagsOutput{}, nil)

			// CreateOrUpdateTags classic mock
			createOrUpdateTagsInput := &autoscaling.CreateOrUpdateTagsInput{
				Tags: []asTypes.Tag{
					{
						ResourceId:        aws.String(asgName),
						ResourceType:      aws.String(resourceTypeAutoScalingGroup),
						Key:               aws.String("tag_key_1"),
						Value:             aws.String("tag_value_1"),
						PropagateAtLaunch: aws.Bool(false),
					},
				},
			}
			p.MockASG().On("CreateOrUpdateTags", mock.Anything, createOrUpdateTagsInput).Return(&autoscaling.CreateOrUpdateTagsOutput{}, nil)

			sm := NewStackCollection(p, api.NewClusterConfig())
			err := sm.PropagateManagedNodeGroupTagsToASG(ngName, ngTags, []string{asgName}, errCh)
			Expect(err).NotTo(HaveOccurred())
			err = <-errCh
			Expect(err).NotTo(HaveOccurred())
		})
		It("cannot propagate tags in chunks of 25", func() {
			// populate the createOrUpdateTagsSliceInput for easier generation of chunks
			createOrUpdateTagsSliceInput := []asTypes.Tag{}
			for i := 0; i < 30; i++ {
				tagKey, tagValue := fmt.Sprintf("tag_key_%d", i), fmt.Sprintf("tag_value_%d", i)
				ngTags[tagKey] = tagValue
				createOrUpdateTagsSliceInput = append(createOrUpdateTagsSliceInput, asTypes.Tag{
					ResourceId:        aws.String(asgName),
					ResourceType:      aws.String(resourceTypeAutoScalingGroup),
					Key:               aws.String(tagKey),
					Value:             aws.String(tagValue),
					PropagateAtLaunch: aws.Bool(false),
				})
			}

			// DescribeTags classic mock
			describeTagsInput := &autoscaling.DescribeTagsInput{
				Filters: []asTypes.Filter{{Name: aws.String(resourceTypeAutoScalingGroup), Values: []string{asgName}}},
			}
			p.MockASG().On("DescribeTags", mock.Anything, describeTagsInput).Return(&autoscaling.DescribeTagsOutput{}, nil)

			// CreateOrUpdateTags chunked mock
			// generate the expected chunk of tags
			chunkSize := builder.MaximumCreatedTagNumberPerCall
			firstchunkLenMatcher := func(input *autoscaling.CreateOrUpdateTagsInput) bool {
				return len(input.Tags) == len(createOrUpdateTagsSliceInput[:chunkSize])
			}
			secondChunkLenMatcher := func(input *autoscaling.CreateOrUpdateTagsInput) bool {
				return len(input.Tags) == len(createOrUpdateTagsSliceInput[chunkSize:])
			}

			// setup the call verification of the two chunks
			// NOTE: because of the use of map (unordered processing), we just verify size of chunk
			p.MockASG().On("CreateOrUpdateTags", mock.Anything, mock.MatchedBy(firstchunkLenMatcher)).Return(&autoscaling.CreateOrUpdateTagsOutput{}, nil)
			p.MockASG().On("CreateOrUpdateTags", mock.Anything, mock.MatchedBy(secondChunkLenMatcher)).Return(&autoscaling.CreateOrUpdateTagsOutput{}, nil)

			sm := NewStackCollection(p, api.NewClusterConfig())
			err := sm.PropagateManagedNodeGroupTagsToASG(ngName, ngTags, []string{asgName}, errCh)
			Expect(err).NotTo(HaveOccurred())
			err = <-errCh
			Expect(err).NotTo(HaveOccurred())
		})
		It("cannot propagate if too many tags", func() {
			// fill parameters
			for i := 0; i < builder.MaximumTagNumber+1; i++ {
				ngTags[fmt.Sprintf("tag_key_%d", i)] = fmt.Sprintf("tag_value_%d", i)
			}

			// DescribeTags classic mock
			describeTagsInput := &autoscaling.DescribeTagsInput{
				Filters: []asTypes.Filter{{Name: aws.String(resourceTypeAutoScalingGroup), Values: []string{asgName}}},
			}
			p.MockASG().On("DescribeTags", mock.Anything, describeTagsInput).Return(&autoscaling.DescribeTagsOutput{}, nil)

			sm := NewStackCollection(p, api.NewClusterConfig())
			err := sm.PropagateManagedNodeGroupTagsToASG(ngName, ngTags, []string{asgName}, errCh)
			Expect(err).NotTo(HaveOccurred())
			err = <-errCh
			Expect(err).To(MatchError(ContainSubstring("maximum amount for asg")))
		})
	})

	Context("UpdateStack", func() {
		It("succeeds if no changes required", func() {
			// Order of AWS SDK invocation
			// 1) DescribeStacks
			// 2) CreateChangeSet
			// 3) DescribeChangeSetRequest (FAILED to abort early)
			// 4) DescribeChangeSet (StatusReason contains "The submitted information didn't contain changes" to exit with noChangeError)

			stackName := "eksctl-stack"
			changeSetName := "eksctl-changeset"
			describeInput := &cfn.DescribeStacksInput{StackName: &stackName}
			describeOutput := &cfn.DescribeStacksOutput{Stacks: []types.Stack{{
				StackName:   &stackName,
				StackStatus: types.StackStatusCreateComplete,
			}}}
			describeChangeSetNoChange := &cfn.DescribeChangeSetOutput{
				StackName:    &stackName,
				StatusReason: aws.String("The submitted information didn't contain changes"),
			}
			p := mockprovider.NewMockProvider()
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, describeInput).Return(describeOutput, nil)
			p.MockCloudFormation().On("CreateChangeSet", mock.Anything, mock.Anything).Return(nil, nil)
			p.MockCloudFormation().On("DescribeChangeSet", mock.Anything, mock.Anything, mock.Anything).Return(describeChangeSetNoChange, nil)

			sm := NewStackCollection(p, api.NewClusterConfig())
			err := sm.UpdateStack(context.TODO(), UpdateStackOptions{
				StackName:     stackName,
				ChangeSetName: changeSetName,
				Description:   "description",
				TemplateData:  TemplateBody(""),
				Wait:          false,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("can update when only the stack is provided", func() {
			// Order of AWS SDK invocation
			// 1) DescribeStacks
			// 2) CreateChangeSet
			// 3) DescribeChangeSet (StatusReason contains "The submitted information didn't contain changes" to exit with noChangeError)

			stackName := "eksctl-stack"
			changeSetName := "eksctl-changeset"
			describeInput := &cfn.DescribeStacksInput{StackName: &stackName}
			describeOutput := &cfn.DescribeStacksOutput{Stacks: []types.Stack{{
				StackName:   &stackName,
				StackStatus: types.StackStatusCreateComplete,
			}}}
			describeChangeSetNoChange := &cfn.DescribeChangeSetOutput{
				StackName:    &stackName,
				StatusReason: aws.String("The submitted information didn't contain changes"),
			}
			p := mockprovider.NewMockProvider()
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, describeInput).Return(describeOutput, nil)
			p.MockCloudFormation().On("CreateChangeSet", mock.Anything, mock.Anything).Return(nil, nil)
			p.MockCloudFormation().On("DescribeChangeSet", mock.Anything, mock.Anything, mock.Anything).Return(describeChangeSetNoChange, nil)

			sm := NewStackCollection(p, api.NewClusterConfig())
			err := sm.UpdateStack(context.TODO(), UpdateStackOptions{
				Stack: &Stack{
					StackName: &stackName,
				},
				ChangeSetName: changeSetName,
				Description:   "description",
				TemplateData:  TemplateBody(""),
				Wait:          false,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("updates tags (existing + metadata + auto)", func() {
		// Order of execution
		// 1) DescribeStacks
		// 2) CreateChangeSet
		// 3) DescribeChangeSet
		// 4) ExecuteChangeSet

		clusterName := "clusteur"
		stackName := "eksctl-stack"
		changeSetName := "eksctl-changeset"
		describeInput := &cfn.DescribeStacksInput{StackName: &stackName}
		existingTag := types.Tag{
			Key:   aws.String("existing"),
			Value: aws.String("tag"),
		}
		describeOutput := &cfn.DescribeStacksOutput{Stacks: []types.Stack{{
			StackName:   &stackName,
			StackStatus: types.StackStatusCreateComplete,
			Tags:        []types.Tag{existingTag},
		}}}
		describeChangeSetCreateCompleteOutput := &cfn.DescribeChangeSetOutput{
			StackName:     &stackName,
			ChangeSetName: &changeSetName,
			Status:        types.ChangeSetStatusCreateComplete,
		}
		executeChangeSetInput := &cfn.ExecuteChangeSetInput{
			ChangeSetName: &changeSetName,
			StackName:     &stackName,
		}

		p := mockprovider.NewMockProvider()
		p.MockCloudFormation().On("DescribeStacks", mock.Anything, describeInput, mock.Anything).Return(describeOutput, nil)
		p.MockCloudFormation().On("CreateChangeSet", mock.Anything, mock.Anything).Return(nil, nil)
		p.MockCloudFormation().On("DescribeChangeSet", mock.Anything, mock.Anything, mock.Anything).Return(describeChangeSetCreateCompleteOutput, nil)
		p.MockCloudFormation().On("ExecuteChangeSet", mock.Anything, executeChangeSetInput).Return(nil, nil)

		spec := api.NewClusterConfig()
		spec.Metadata.Name = clusterName
		spec.Metadata.Tags = map[string]string{"meta": "data"}
		sm := NewStackCollection(p, spec)
		err := sm.UpdateStack(context.TODO(), UpdateStackOptions{
			StackName:     stackName,
			ChangeSetName: changeSetName,
			Description:   "description",
			TemplateData:  TemplateBody(""),
			Wait:          false,
		})
		Expect(err).NotTo(HaveOccurred())

		// Second is CreateChangeSet() call which we are interested in
		args := p.MockCloudFormation().Calls[1].Arguments.Get(1)
		createChangeSetInput := args.(*cfn.CreateChangeSetInput)
		// Existing tag
		Expect(createChangeSetInput.Tags).To(ContainElement(existingTag))
		// Auto-populated tag
		Expect(createChangeSetInput.Tags).To(ContainElement(types.Tag{Key: aws.String(api.ClusterNameTag), Value: &clusterName}))
		// Metadata tag
		Expect(createChangeSetInput.Tags).To(ContainElement(types.Tag{Key: aws.String("meta"), Value: aws.String("data")}))
	})
	When("wait is set to false", func() {
		It("will skip the last wait sequence", func() {
			clusterName := "cluster"
			stackName := "eksctl-stack"
			changeSetName := "eksctl-changeset"
			describeInput := &cfn.DescribeStacksInput{StackName: &stackName}
			existingTag := types.Tag{
				Key:   aws.String("existing"),
				Value: aws.String("tag"),
			}
			describeOutput := &cfn.DescribeStacksOutput{Stacks: []types.Stack{{
				StackName:   &stackName,
				StackStatus: types.StackStatusCreateComplete,
				Tags:        []types.Tag{existingTag},
			}}}
			describeChangeSetCreateCompleteOutput := &cfn.DescribeChangeSetOutput{
				StackName:     &stackName,
				ChangeSetName: &changeSetName,
				Status:        types.ChangeSetStatusCreateComplete,
			}
			executeChangeSetInput := &cfn.ExecuteChangeSetInput{
				ChangeSetName: &changeSetName,
				StackName:     &stackName,
			}

			p := mockprovider.NewMockProvider()
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, describeInput).Return(describeOutput, nil)
			p.MockCloudFormation().On("CreateChangeSet", mock.Anything, mock.Anything).Return(nil, nil)
			req := awstesting.NewClient(nil).NewRequest(&request.Operation{Name: "Operation"}, nil, describeChangeSetCreateCompleteOutput)
			p.MockCloudFormation().On("DescribeChangeSetRequest", mock.Anything, mock.Anything).Return(req, describeChangeSetCreateCompleteOutput)
			p.MockCloudFormation().On("DescribeChangeSet", mock.Anything, mock.Anything, mock.Anything).Return(describeChangeSetCreateCompleteOutput, nil)
			p.MockCloudFormation().On("ExecuteChangeSet", mock.Anything, executeChangeSetInput).Return(nil, nil)
			// For the future, this is the call we do not expect to happen, and this is the difference compared to the
			// above test case.
			// p.MockCloudFormation().On("DescribeStacksRequest", mock.Anything).Return(req, describeStacksUpdateCompleteOutput)

			spec := api.NewClusterConfig()
			spec.Metadata.Name = clusterName
			spec.Metadata.Tags = map[string]string{"meta": "data"}
			sm := NewStackCollection(p, spec)
			err := sm.UpdateStack(context.TODO(), UpdateStackOptions{
				StackName:     stackName,
				ChangeSetName: changeSetName,
				Description:   "description",
				TemplateData:  TemplateBody(""),
				Wait:          false,
			})
			Expect(err).NotTo(HaveOccurred())

			// Second is CreateChangeSet() call which we are interested in
			args := p.MockCloudFormation().Calls[1].Arguments.Get(1)
			createChangeSetInput := args.(*cfn.CreateChangeSetInput)
			// Existing tag
			Expect(createChangeSetInput.Tags).To(ContainElement(existingTag))
			// Auto-populated tag
			Expect(createChangeSetInput.Tags).To(ContainElement(types.Tag{Key: aws.String(api.ClusterNameTag), Value: &clusterName}))
			// Metadata tag
			Expect(createChangeSetInput.Tags).To(ContainElement(types.Tag{Key: aws.String("meta"), Value: aws.String("data")}))
		})
	})

	Context("HasClusterStackFromList", func() {
		type clusterInput struct {
			clusterName   string
			eksctlCreated bool
		}

		DescribeTable("should work for eksctl-created clusters", func(ci clusterInput) {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = ci.clusterName
			stackName := aws.String(fmt.Sprintf("eksctl-%s-cluster", clusterConfig.Metadata.Name))

			var out *cfn.DescribeStacksOutput
			if ci.eksctlCreated {
				out = &cfn.DescribeStacksOutput{
					Stacks: []types.Stack{
						{
							StackName: stackName,
							Tags: []types.Tag{
								{
									Key:   aws.String("alpha.eksctl.io/cluster-name"),
									Value: aws.String(clusterConfig.Metadata.Name),
								},
							},
						},
					},
				}
			} else {
				out = &cfn.DescribeStacksOutput{}
			}

			p := mockprovider.NewMockProvider()
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, &cfn.DescribeStacksInput{StackName: stackName}).Return(out, nil)

			s := NewStackCollection(p, clusterConfig)
			hasClusterStack, err := s.HasClusterStackFromList(context.TODO(), []string{
				"eksctl-test-cluster",
				*stackName,
			}, clusterConfig.Metadata.Name)

			if ci.eksctlCreated {
				Expect(err).NotTo(HaveOccurred())
				Expect(hasClusterStack).To(Equal(true))
			} else {
				Expect(err).To(MatchError(fmt.Sprintf("no CloudFormation stack found for %s", *stackName)))
			}
		},
			Entry("cluster stack exists", clusterInput{
				clusterName:   "web",
				eksctlCreated: true,
			}),
			Entry("cluster stack does not exist", clusterInput{
				clusterName:   "unowned",
				eksctlCreated: false,
			}),
		)
	})

	Context("GetClusterStackIfExists", func() {
		var (
			cfg                 *api.ClusterConfig
			p                   *mockprovider.MockProvider
			stackNameWithEksctl string
		)
		BeforeEach(func() {
			stackName := "confirm-this"
			stackNameWithEksctl = "eksctl-" + stackName + "-cluster"
			describeInput := &cfn.DescribeStacksInput{StackName: &stackNameWithEksctl}
			describeOutput := &cfn.DescribeStacksOutput{Stacks: []types.Stack{{
				StackName:   &stackName,
				StackStatus: types.StackStatusCreateComplete,
				Tags: []types.Tag{
					{
						Key:   aws.String(api.ClusterNameTag),
						Value: &stackName,
					},
				},
			}}}
			p = mockprovider.NewMockProvider()
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, describeInput).Return(describeOutput, nil)

			cfg = api.NewClusterConfig()
			cfg.Metadata.Name = stackName
		})

		It("can retrieve stacks", func() {
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(&cfn.ListStacksOutput{
				StackSummaries: []types.StackSummary{
					{
						StackName: &stackNameWithEksctl,
					},
				},
			}, nil)
			sm := NewStackCollection(p, cfg)
			stack, err := sm.GetClusterStackIfExists(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(stack).NotTo(BeNil())
		})

		When("the config stack doesn't match", func() {
			It("returns no stack", func() {
				p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(&cfn.ListStacksOutput{}, nil)
				cfg.Metadata.Name = "not-this"
				sm := NewStackCollection(p, cfg)
				stack, err := sm.GetClusterStackIfExists(context.TODO())
				Expect(err).NotTo(HaveOccurred())
				Expect(stack).To(BeNil())
			})
		})

		When("ListStacks errors", func() {
			It("errors", func() {
				p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(nil, errors.New("nope"))
				sm := NewStackCollection(p, cfg)
				_, err := sm.GetClusterStackIfExists(context.TODO())
				Expect(err).To(MatchError(ContainSubstring("nope")))
			})
		})
	})
})
