package manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("StackCollection", func() {
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
			describeOutput := &cfn.DescribeStacksOutput{Stacks: []*cfn.Stack{{
				StackName:   &stackName,
				StackStatus: aws.String(cfn.StackStatusCreateComplete),
			}}}
			describeChangeSetFailed := &cfn.DescribeChangeSetOutput{
				StackName:     &stackName,
				ChangeSetName: &changeSetName,
				Status:        aws.String(cfn.ChangeSetStatusFailed),
			}
			describeChangeSetNoChange := &cfn.DescribeChangeSetOutput{
				StackName:    &stackName,
				StatusReason: aws.String("The submitted information didn't contain changes"),
			}
			p := mockprovider.NewMockProvider()
			p.MockCloudFormation().On("DescribeStacks", describeInput).Return(describeOutput, nil)
			p.MockCloudFormation().On("CreateChangeSet", mock.Anything).Return(nil, nil)
			req := awstesting.NewClient(nil).NewRequest(&request.Operation{Name: "Operation"}, nil, describeChangeSetFailed)
			p.MockCloudFormation().On("DescribeChangeSetRequest", mock.Anything).Return(req, describeChangeSetFailed)
			p.MockCloudFormation().On("DescribeChangeSet", mock.Anything).Return(describeChangeSetNoChange, nil)

			sm := NewStackCollection(p, api.NewClusterConfig())
			err := sm.UpdateStack(stackName, changeSetName, "description", TemplateBody(""), nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("updates tags (existing + metadata + auto)", func() {
		// Order of execution
		// 1) DescribeStacks
		// 2) CreateChangeSet
		// 3) DescribeChangeSetRequest (until CREATE_COMPLETE)
		// 4) DescribeChangeSet
		// 5) ExecuteChangeSet
		// 6) DescribeStacksRequest (until UPDATE_COMPLETE)

		clusterName := "clusteur"
		stackName := "eksctl-stack"
		changeSetName := "eksctl-changeset"
		describeInput := &cfn.DescribeStacksInput{StackName: &stackName}
		existingTag := &cfn.Tag{
			Key:   aws.String("existing"),
			Value: aws.String("tag"),
		}
		describeOutput := &cfn.DescribeStacksOutput{Stacks: []*cfn.Stack{{
			StackName:   &stackName,
			StackStatus: aws.String(cfn.StackStatusCreateComplete),
			Tags:        []*cfn.Tag{existingTag},
		}}}
		describeChangeSetCreateCompleteOutput := &cfn.DescribeChangeSetOutput{
			StackName:     &stackName,
			ChangeSetName: &changeSetName,
			Status:        aws.String(cfn.ChangeSetStatusCreateComplete),
		}
		describeStacksUpdateCompleteOutput := &cfn.DescribeStacksOutput{
			Stacks: []*cfn.Stack{
				{
					StackName:   &stackName,
					StackStatus: aws.String(cfn.StackStatusUpdateComplete),
				},
			},
		}
		executeChangeSetInput := &cfn.ExecuteChangeSetInput{
			ChangeSetName: &changeSetName,
			StackName:     &stackName,
		}

		p := mockprovider.NewMockProvider()
		p.MockCloudFormation().On("DescribeStacks", describeInput).Return(describeOutput, nil)
		p.MockCloudFormation().On("CreateChangeSet", mock.Anything).Return(nil, nil)
		req := awstesting.NewClient(nil).NewRequest(&request.Operation{Name: "Operation"}, nil, describeChangeSetCreateCompleteOutput)
		p.MockCloudFormation().On("DescribeChangeSetRequest", mock.Anything).Return(req, describeChangeSetCreateCompleteOutput)
		p.MockCloudFormation().On("DescribeChangeSet", mock.Anything).Return(describeChangeSetCreateCompleteOutput, nil)
		p.MockCloudFormation().On("ExecuteChangeSet", executeChangeSetInput).Return(nil, nil)
		req = awstesting.NewClient(nil).NewRequest(&request.Operation{Name: "Operation"}, nil, describeStacksUpdateCompleteOutput)
		p.MockCloudFormation().On("DescribeStacksRequest", mock.Anything).Return(req, describeStacksUpdateCompleteOutput)

		spec := api.NewClusterConfig()
		spec.Metadata.Name = clusterName
		spec.Metadata.Tags = map[string]string{"meta": "data"}
		sm := NewStackCollection(p, spec)
		err := sm.UpdateStack(stackName, changeSetName, "description", TemplateBody(""), nil)
		Expect(err).ToNot(HaveOccurred())

		// Second is CreateChangeSet() call which we are interested in
		args := p.MockCloudFormation().Calls[1].Arguments.Get(0)
		createChangeSetInput := args.(*cfn.CreateChangeSetInput)
		// Existing tag
		Expect(createChangeSetInput.Tags).To(ContainElement(existingTag))
		// Auto-populated tag
		Expect(createChangeSetInput.Tags).To(ContainElement(&cfn.Tag{Key: aws.String(api.ClusterNameTag), Value: &clusterName}))
		// Metadata tag
		Expect(createChangeSetInput.Tags).To(ContainElement(&cfn.Tag{Key: aws.String("meta"), Value: aws.String("data")}))
	})
})
