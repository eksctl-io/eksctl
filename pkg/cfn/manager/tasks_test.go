package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/fake"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"
)

type task struct{ id int }

func (t *task) Describe() string {
	return fmt.Sprintf("task %d", t.id)
}

func (t *task) Do(chan error) error {
	return nil
}

var _ = Describe("StackCollection Tasks", func() {
	var (
		p   *mockprovider.MockProvider
		cfg *api.ClusterConfig

		stackManager StackManager
	)

	testAZs := []string{"us-west-2b", "us-west-2a", "us-west-2c"}

	newClusterConfig := func(clusterName string) *api.ClusterConfig {
		cfg := api.NewClusterConfig()
		*cfg.VPC.CIDR = api.DefaultCIDR()

		ng1 := cfg.NewNodeGroup()
		ng2 := cfg.NewNodeGroup()

		cfg.Metadata.Region = "us-west-2"
		cfg.Metadata.Name = clusterName
		cfg.AvailabilityZones = testAZs

		ng1.Name = "bar"
		ng1.InstanceType = "t2.medium"
		ng1.AMIFamily = "AmazonLinux2"
		ng2.Labels = map[string]string{"bar": "bar"}

		ng2.Name = "foo"
		ng2.InstanceType = "t2.medium"
		ng2.AMIFamily = "AmazonLinux2"
		ng2.Labels = map[string]string{"foo": "foo"}

		return cfg
	}

	Describe("TaskTree", func() {
		BeforeEach(func() {
			p = mockprovider.NewMockProvider()
			cfg = newClusterConfig("test-cluster")
			stackManager = NewStackCollection(p, cfg)
		})

		It("should have nice description", func() {
			fakeVPCImporter := new(vpcfakes.FakeImporter)
			// TODO use DescribeTable

			// The supportsManagedNodes argument has no effect on the Describe call, so the values are alternated
			// in these tests
			{
				tasks := stackManager.NewUnmanagedNodeGroupTask(makeNodeGroups("bar", "foo"), false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`
2 parallel tasks: { create nodegroup "bar", create nodegroup "foo" 
}
`))
			}
			{
				tasks := stackManager.NewUnmanagedNodeGroupTask(makeNodeGroups("bar"), false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`1 task: { create nodegroup "bar" }`))
			}
			{
				tasks := stackManager.NewUnmanagedNodeGroupTask(makeNodeGroups("foo"), false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`1 task: { create nodegroup "foo" }`))
			}
			{
				tasks := stackManager.NewUnmanagedNodeGroupTask(nil, false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`no tasks`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar", "foo"), nil)
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", 
    2 parallel sub-tasks: { 
        create nodegroup "bar",
        create nodegroup "foo",
    } 
}
`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar"), nil)
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", create nodegroup "bar" 
}
`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(nil, nil)
				Expect(tasks.Describe()).To(Equal(`1 task: { create cluster control plane "test-cluster" }`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar", "foo"), makeManagedNodeGroups("m1", "m2"))
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", 
    4 parallel sub-tasks: { 
        create nodegroup "bar",
        create nodegroup "foo",
        create managed nodegroup "m1",
        create managed nodegroup "m2",
    } 
}
`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("foo"), makeManagedNodeGroups("m1"))
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", 
    2 parallel sub-tasks: { 
        create nodegroup "foo",
        create managed nodegroup "m1",
    } 
}
`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar"), nil, &task{id: 1})
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", 
    2 sequential sub-tasks: { 
        task 1,
        create nodegroup "bar",
    } 
}
`))
			}
		})

		When("IPFamily is set to ipv6", func() {
			BeforeEach(func() {
				cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
			})
			It("appends the AssignIpv6AddressOnCreation task to occur after the cluster creation", func() {
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar", "foo"), nil)
				Expect(tasks.Describe()).To(Equal(`
3 sequential tasks: { create cluster control plane "test-cluster", set AssignIpv6AddressOnCreation to true for public subnets, 
    2 parallel sub-tasks: { 
        create nodegroup "bar",
        create nodegroup "foo",
    } 
}
`))
			})
		})
	})

	Describe("ManagedNodeGroupTask", func() {
		When("creating managed nodegroups on a ipv6 cluster", func() {
			var (
				p            *mockprovider.MockProvider
				cfg          *api.ClusterConfig
				stackManager StackManager
			)
			BeforeEach(func() {
				p = mockprovider.NewMockProvider()
				cfg = newClusterConfig("test-ipv6-cluster")
				cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
				stackManager = NewStackCollection(p, cfg)
			})
			It("returns an error", func() {
				p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil)
				ng := api.NewManagedNodeGroup()
				fakeVPCImporter := new(vpcfakes.FakeImporter)
				tasks := stackManager.NewManagedNodeGroupTask([]*api.ManagedNodeGroup{ng}, false, fakeVPCImporter)
				errs := tasks.DoAllSync()
				Expect(errs).To(HaveLen(1))
				Expect(errs[0]).To(MatchError(ContainSubstring("managed nodegroups cannot be created on IPv6 unowned clusters")))
			})
			When("finding the stack fails", func() {
				It("returns the stack error", func() {
					p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(errors.New("not found"))
					ng := api.NewManagedNodeGroup()
					fakeVPCImporter := new(vpcfakes.FakeImporter)
					tasks := stackManager.NewManagedNodeGroupTask([]*api.ManagedNodeGroup{ng}, false, fakeVPCImporter)
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(1))
					Expect(errs[0]).To(MatchError(ContainSubstring("not found")))
				})
			})
		})
	})

	Describe("IAMServiceAccountTask", func() {
		var (
			clusterName        string
			iamServiceAccounts []*api.ClusterIAMServiceAccount
			name               string
			namespace          string
			oidc               *iamoidc.OpenIDConnectManager
			stackArn           string
			stackArnPrefix     string
			stackName          string
		)
		// service account creation should succeed
		When("service account doesn't exist yet", func() {
			// setup
			BeforeEach(func() {
				clusterName = "cluster1"
				namespace = "default"
				name = "new"
				cfg = newClusterConfig(clusterName)
				iamServiceAccounts = []*api.ClusterIAMServiceAccount{
					{
						ClusterIAMMeta: api.ClusterIAMMeta{
							Name:      name,
							Namespace: namespace,
						},
						AttachPolicyARNs: []string{"arn-123"},
					},
				}
				cfg.IAM.ServiceAccounts = iamServiceAccounts
				var err error
				oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
				Expect(err).NotTo(HaveOccurred())
				oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
				p = mockprovider.NewMockProvider()
				stackManager = NewStackCollection(p, cfg)
				stackName = fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
				stackArn = fmt.Sprintf("arn:aws:cloudformation:eu-west-1:456123987123:stack/%s/01", stackName)
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}).Return(&cloudformation.DescribeStacksOutput{}, nil)
				p.MockCloudFormation().On("CreateStack", mock.Anything).Return(&cloudformation.CreateStackOutput{StackId: aws.String(stackArn)}, nil)
				describeStacksOutput := &cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackName:   aws.String(stackName),
							StackId:     aws.String(stackArn),
							StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
							Tags: []*cloudformation.Tag{
								{
									Key:   aws.String("alpha.eksctl.io/cluster-name"),
									Value: aws.String(clusterName),
								},
							},
							Outputs: []*cloudformation.Output{
								{
									Description: aws.String("IAM Role"),
									OutputKey:   aws.String("Role1"),
									OutputValue: aws.String("role-arn-123"),
									ExportName:  aws.String("role-arn-123"),
								},
							},
						},
					},
				}
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{StackName: aws.String(stackArn)}).Return(describeStacksOutput, nil)
				req := awstesting.NewClient(aws.NewConfig().WithRegion("us-west-2")).NewRequest(&request.Operation{Name: "Operation"}, nil, describeStacksOutput)
				p.MockCloudFormation().On("DescribeStacksRequest", &cloudformation.DescribeStacksInput{StackName: aws.String(stackArn)}).Return(req, describeStacksOutput)
			})
			It("creates a new stack without errors", func() {
				tasks := stackManager.NewTasksToCreateIAMServiceAccounts(iamServiceAccounts, oidc, kubernetes.NewCachedClientSet(fake.NewSimpleClientset()))
				Expect(tasks.Len()).To(Equal(1))
				Expect(tasks.Describe()).To(Equal(`1 task: { 
    2 sequential sub-tasks: { 
        create IAM role for serviceaccount "default/new",
        create serviceaccount "default/new",
    } }`))
				errs := tasks.DoAllSync()
				Expect(errs).To(HaveLen(0))
				// all calls should happen
				Expect(p.MockCloudFormation().AssertExpectations(GinkgoT()))
			})
		})
		// existing stack in CREATE_COMPLETE status, shouldn't fail
		When("service account already exists", func() {
			BeforeEach(func() {
				clusterName = "cluster1"
				namespace = "default"
				name = "existing"
				cfg = newClusterConfig(clusterName)
				iamServiceAccounts = []*api.ClusterIAMServiceAccount{
					{
						ClusterIAMMeta: api.ClusterIAMMeta{
							Name:      name,
							Namespace: namespace,
						},
						AttachPolicyARNs: []string{"arn-123"},
					},
				}
				cfg.IAM.ServiceAccounts = iamServiceAccounts
				var err error
				oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
				Expect(err).NotTo(HaveOccurred())
				oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
				p = mockprovider.NewMockProvider()
				stackManager = NewStackCollection(p, cfg)
				stackName = fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
				stackArnPrefix = fmt.Sprintf("arn:aws:cloudformation:eu-west-1:456123987123:stack/%s", stackName)
				describeStacksOutput01 := &cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackName:   aws.String(stackName),
							StackId:     aws.String(fmt.Sprintf("%s/01", stackArnPrefix)),
							StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
							Tags: []*cloudformation.Tag{
								{
									Key:   aws.String("alpha.eksctl.io/cluster-name"),
									Value: aws.String(clusterName),
								},
							},
							Outputs: []*cloudformation.Output{
								{
									Description: aws.String("IAM Role"),
									OutputKey:   aws.String("Role1"),
									OutputValue: aws.String("role-arn-123"),
									ExportName:  aws.String("role-arn-123"),
								},
							},
						},
					},
				}
				describeStacksOutput02 := &cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackName:   aws.String(stackName),
							StackId:     aws.String(fmt.Sprintf("%s/02", stackArnPrefix)),
							StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
							Tags: []*cloudformation.Tag{
								{
									Key:   aws.String("alpha.eksctl.io/cluster-name"),
									Value: aws.String(clusterName),
								},
							},
							Outputs: []*cloudformation.Output{
								{
									Description: aws.String("IAM Role"),
									OutputKey:   aws.String("Role1"),
									OutputValue: aws.String("role-arn-123"),
									ExportName:  aws.String("role-arn-123"),
								},
							},
						},
					},
				}
				// filter on stack name when looking for a unknown stack
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}).Return(describeStacksOutput01, nil)
				// filter on stack id when looking for a known stack
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{StackName: aws.String(fmt.Sprintf("%s/02", stackArnPrefix))}).Return(describeStacksOutput02, nil)
				p.MockCloudFormation().On("CreateStack", mock.Anything).Return(&cloudformation.CreateStackOutput{StackId: aws.String(fmt.Sprintf("%s/02", stackArnPrefix))}, nil)
				req := awstesting.NewClient(aws.NewConfig().WithRegion("us-west-2")).NewRequest(&request.Operation{Name: "Operation"}, nil, describeStacksOutput02)
				p.MockCloudFormation().On("DescribeStacksRequest", &cloudformation.DescribeStacksInput{StackName: aws.String(fmt.Sprintf("%s/02", stackArnPrefix))}).Return(req, describeStacksOutput02)
			})
			It("doesn't fail", func() {
				tasks := stackManager.NewTasksToCreateIAMServiceAccounts(iamServiceAccounts, oidc, kubernetes.NewCachedClientSet(fake.NewSimpleClientset()))
				Expect(tasks.Len()).To(Equal(1))
				errs := tasks.DoAllSync()
				Expect(errs).To(HaveLen(0))
				// all calls should happen
				Expect(p.MockCloudFormation().AssertExpectations(GinkgoT()))
			})
		})
		// should delete and attempt to recreate the stack
		When("service account exists and is in ROLLBACK_COMPLETE status", func() {
			BeforeEach(func() {
				clusterName = "cluster1"
				namespace = "default"
				name = "rollback"
				cfg = newClusterConfig(clusterName)
				iamServiceAccounts = []*api.ClusterIAMServiceAccount{
					{
						ClusterIAMMeta: api.ClusterIAMMeta{
							Name:      name,
							Namespace: namespace,
						},
						AttachPolicyARNs: []string{"arn-123"},
					},
				}
				cfg.IAM.ServiceAccounts = iamServiceAccounts
				var err error
				oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
				Expect(err).NotTo(HaveOccurred())
				oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
				p = mockprovider.NewMockProvider()
				p.SetRegion("us-west-2")
				stackManager = NewStackCollection(p, cfg)
				stackName = fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
				stackArnPrefix = fmt.Sprintf("arn:aws:cloudformation:eu-west-1:456123987123:stack/%s", stackName)
				describeStacksOutput01 := &cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackName:   aws.String(stackName),
							StackId:     aws.String(fmt.Sprintf("%s/01", stackArnPrefix)),
							StackStatus: aws.String(cloudformation.StackStatusRollbackComplete),
							Tags: []*cloudformation.Tag{
								{
									Key:   aws.String("alpha.eksctl.io/cluster-name"),
									Value: aws.String(clusterName),
								},
							},
						},
					},
				}
				describeStacksOutput01Deleted := &cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackName:   aws.String(stackName),
							StackId:     aws.String(fmt.Sprintf("%s/01", stackArnPrefix)),
							StackStatus: aws.String(cloudformation.StackStatusDeleteComplete),
							Tags: []*cloudformation.Tag{
								{
									Key:   aws.String("alpha.eksctl.io/cluster-name"),
									Value: aws.String(clusterName),
								},
							},
						},
					},
				}
				describeStacksOutput02 := &cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackName:   aws.String(stackName),
							StackId:     aws.String(fmt.Sprintf("%s/02", stackArnPrefix)),
							StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
							Tags: []*cloudformation.Tag{
								{
									Key:   aws.String("alpha.eksctl.io/cluster-name"),
									Value: aws.String(clusterName),
								},
							},
							Outputs: []*cloudformation.Output{
								{
									Description: aws.String("IAM Role"),
									OutputKey:   aws.String("Role1"),
									OutputValue: aws.String("role-arn-123"),
									ExportName:  aws.String("role-arn-123"),
								},
							},
						},
					},
				}
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}).Return(describeStacksOutput01, nil)
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{StackName: aws.String(fmt.Sprintf("%s/02", stackArnPrefix))}).Return(describeStacksOutput02, nil)
				p.MockCloudFormation().On("CreateStack", mock.Anything).Return(&cloudformation.CreateStackOutput{StackId: aws.String(fmt.Sprintf("%s/02", stackArnPrefix))}, nil)
				p.MockCloudFormation().On("DeleteStack", mock.Anything).Return(&cloudformation.DeleteStackOutput{}, nil)
				// for deletion request
				req01 := awstesting.NewClient(aws.NewConfig().WithRegion("us-west-2")).NewRequest(&request.Operation{Name: "Operation"}, nil, describeStacksOutput01Deleted)
				p.MockCloudFormation().On("DescribeStacksRequest", &cloudformation.DescribeStacksInput{StackName: aws.String(fmt.Sprintf("%s/01", stackArnPrefix))}).Return(req01, describeStacksOutput01Deleted)
				// for creation request
				req02 := awstesting.NewClient(aws.NewConfig().WithRegion("us-west-2")).NewRequest(&request.Operation{Name: "Operation"}, nil, describeStacksOutput02)
				p.MockCloudFormation().On("DescribeStacksRequest", &cloudformation.DescribeStacksInput{StackName: aws.String(fmt.Sprintf("%s/02", stackArnPrefix))}).Return(req02, describeStacksOutput02)
			})
			It("doesn't fail", func() {
				tasks := stackManager.NewTasksToCreateIAMServiceAccounts(iamServiceAccounts, oidc, kubernetes.NewCachedClientSet(fake.NewSimpleClientset()))
				Expect(tasks.Describe()).To(Equal(`1 task: { 
    2 sequential sub-tasks: { 
        create IAM role for serviceaccount "default/rollback",
        create serviceaccount "default/rollback",
    } }`))
				Expect(tasks.Len()).To(Equal(1))
				errs := tasks.DoAllSync()
				Expect(errs).To(HaveLen(0))
				// all calls should happen
				Expect(p.MockCloudFormation().AssertExpectations(GinkgoT()))
				// all calls should happen
				Expect(p.MockCloudFormation().AssertExpectations(GinkgoT()))
			})
		})
	})
})

func makeNodeGroups(names ...string) []*api.NodeGroup {
	var nodeGroups []*api.NodeGroup
	for _, name := range names {
		ng := api.NewNodeGroup()
		ng.Name = name
		nodeGroups = append(nodeGroups, ng)
	}
	return nodeGroups
}

func makeManagedNodeGroups(names ...string) []*api.ManagedNodeGroup {
	var managedNodeGroups []*api.ManagedNodeGroup
	for _, name := range names {
		ng := api.NewManagedNodeGroup()
		ng.Name = name
		managedNodeGroups = append(managedNodeGroups, ng)
	}
	return managedNodeGroups
}
