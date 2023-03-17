package manager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("StackCollection NodeGroup", func() {
	Describe("GetNodeGroupType", func() {
		createTags := func(tags map[string]string) []types.Tag {
			cfnTags := make([]types.Tag, 0)
			for k, v := range tags {
				cfnTags = append(cfnTags, types.Tag{
					Key:   aws.String(k),
					Value: aws.String(v),
				})
			}
			return cfnTags
		}

		DescribeTable("with tag for the nodegroup type", func(inputTags map[string]string, expectedType api.NodeGroupType) {
			ngType, err := GetNodeGroupType(createTags(inputTags))

			if expectedType == "" {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(ngType).To(Equal(expectedType))
			}
		},

			Entry("finds the type of a managed nodegroup",
				map[string]string{
					api.NodeGroupNameTag: "mng-1",
					api.NodeGroupTypeTag: "managed",
				},
				api.NodeGroupTypeManaged),

			Entry("finds the type of an un-managed nodegroup",
				map[string]string{
					api.NodeGroupNameTag: "ng-1",
					api.NodeGroupTypeTag: "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("finds the type of a legacy un-managed nodegroup",
				map[string]string{
					api.OldNodeGroupNameTag: "ng-1",
					api.NodeGroupTypeTag:    "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("finds the type of a legacy un-managed nodegroup",
				map[string]string{
					api.OldNodeGroupIDTag: "ng-1",
					api.NodeGroupTypeTag:  "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("doesn't return the type if the stack tags don't contain any ng name tag",
				map[string]string{
					"some-other-tag": "ng-1",
					"name":           "ng-1",
				},
				api.NodeGroupType("")),
		)
		DescribeTable("for legacy ngs without tag for the type", func(inputTags map[string]string, expectedType api.NodeGroupType) {
			ngType, err := GetNodeGroupType(createTags(inputTags))

			if expectedType == "" {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(ngType).To(Equal(expectedType))
			}
		},

			Entry("legacy ngs with old name tags are un-managed by default",
				map[string]string{
					api.NodeGroupNameTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("legacy ngs with old name tags are un-managed by default",
				map[string]string{
					api.OldNodeGroupNameTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("legacy ngs with old name tag group-id are un-managed by default",
				map[string]string{
					api.OldNodeGroupIDTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("doesn't return the type if the stack tags don't contain any ng name tag",
				map[string]string{
					"some-other-tag": "ng-1",
					"name":           "ng-1",
				},
				api.NodeGroupType("")),
		)
	})

	Describe("GetUnmanagedNodeGroupAutoScalingGroupName", func() {

		stackName := "stack"
		logicalResourceID := "NodeGroup"
		physicalResourceID := "asg"

		It("returns the asg name", func() {
			p := mockprovider.NewMockProvider()
			p.MockCloudFormation().On("DescribeStackResource", mock.Anything, &cloudformation.DescribeStackResourceInput{
				LogicalResourceId: aws.String(logicalResourceID),
				StackName:         aws.String(stackName),
			}).Return(&cloudformation.DescribeStackResourceOutput{
				StackResourceDetail: &types.StackResourceDetail{
					LogicalResourceId:  aws.String(logicalResourceID),
					StackName:          aws.String(stackName),
					PhysicalResourceId: aws.String(physicalResourceID),
				},
			}, nil)

			sm := NewStackCollection(p, api.NewClusterConfig())
			name, err := sm.GetUnmanagedNodeGroupAutoScalingGroupName(context.Background(), &types.Stack{
				StackName: aws.String(stackName),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(name).To(Equal(physicalResourceID))
		})

		When("The asg resource has no physical ID", func() {
			It("returns an error", func() {
				p := mockprovider.NewMockProvider()
				p.MockCloudFormation().On("DescribeStackResource", mock.Anything, &cloudformation.DescribeStackResourceInput{
					LogicalResourceId: aws.String(logicalResourceID),
					StackName:         aws.String(stackName),
				}).Return(&cloudformation.DescribeStackResourceOutput{
					StackResourceDetail: &types.StackResourceDetail{
						LogicalResourceId:  aws.String(logicalResourceID),
						StackName:          aws.String(stackName),
						PhysicalResourceId: nil,
					},
				}, fmt.Errorf("no PhysicalResourceId"))

				sm := NewStackCollection(p, api.NewClusterConfig())
				name, err := sm.GetUnmanagedNodeGroupAutoScalingGroupName(context.Background(), &types.Stack{
					StackName: aws.String(stackName),
				})
				Expect(err).To(HaveOccurred())
				Expect(name).To(BeEmpty())
			})
		})
	})

	Describe("Propagate ManagedNodeGroupTags to ASG", func() {
		mng := api.NewManagedNodeGroup()
		mng.Name = "test-managed-nodegroup"
		mng.Tags = map[string]string{
			"tag-key": "tag-value",
		}
		sc := StackCollection{
			spec: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Name: "test-cluster",
				},
			},
		}
		p := mockprovider.NewMockProvider()

		BeforeEach(func() {
			p.MockEKS().On("DescribeNodegroup", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&eks.DescribeNodegroupInput{}))
			}).Return(&eks.DescribeNodegroupOutput{
				Nodegroup: &ekstypes.Nodegroup{
					Resources: &ekstypes.NodegroupResources{},
				},
			}, nil)
			sc.eksAPI = p.EKS()
		})

		When("there are unique labels and taints present", func() {
			It("converts them to tags and propagates to ASG", func() {
				mng.Labels = map[string]string{
					"label-key": "label-value",
				}
				mng.Taints = []api.NodeGroupTaint{
					{
						Key:   "taint-key",
						Value: "taint-value",
					},
				}
				propagatedTags := make(map[string]string)
				propagateFunc := func(ngName string, tags map[string]string, asgNames []string, errorCh chan error) error {
					propagatedTags = tags
					return nil
				}

				err := sc.propagateManagedNodeGroupTagsToASGTask(
					context.Background(),
					make(chan error),
					mng,
					propagateFunc,
				)

				Expect(err).To(BeNil())
				Expect(propagatedTags).To(Equal(map[string]string{
					"tag-key": "tag-value",
					"k8s.io/cluster-autoscaler/node-template/label/label-key": "label-value",
					"k8s.io/cluster-autoscaler/node-template/taint/taint-key": "taint-value",
				}))
			})
		})
	})
})
