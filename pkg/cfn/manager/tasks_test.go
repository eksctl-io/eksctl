package manager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
				tasks := stackManager.NewUnmanagedNodeGroupTask(context.Background(), makeNodeGroups("bar", "foo"), false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`
2 parallel tasks: { create nodegroup "bar", create nodegroup "foo" 
}
`))
			}
			{
				tasks := stackManager.NewUnmanagedNodeGroupTask(context.Background(), makeNodeGroups("bar"), false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`1 task: { create nodegroup "bar" }`))
			}
			{
				tasks := stackManager.NewUnmanagedNodeGroupTask(context.Background(), makeNodeGroups("foo"), false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`1 task: { create nodegroup "foo" }`))
			}
			{
				tasks := stackManager.NewUnmanagedNodeGroupTask(context.Background(), nil, false, fakeVPCImporter)
				Expect(tasks.Describe()).To(Equal(`no tasks`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(context.Background(), makeNodeGroups("bar", "foo"), nil)
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
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(context.Background(), makeNodeGroups("bar"), nil)
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", create nodegroup "bar" 
}
`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(context.Background(), nil, nil)
				Expect(tasks.Describe()).To(Equal(`1 task: { create cluster control plane "test-cluster" }`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(context.Background(), makeNodeGroups("bar", "foo"), makeManagedNodeGroups("m1", "m2"))
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", 
    2 parallel sub-tasks: { 
        2 parallel sub-tasks: { 
            create nodegroup "bar",
            create nodegroup "foo",
        },
        2 parallel sub-tasks: { 
            create managed nodegroup "m1",
            create managed nodegroup "m2",
        },
    } 
}
`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(context.Background(), makeNodeGroups("bar", "foo"), makeManagedNodeGroupsWithPropagatedTags("m1", "m2"))
				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", 
    2 parallel sub-tasks: { 
        2 parallel sub-tasks: { 
            create nodegroup "bar",
            create nodegroup "foo",
        },
        2 parallel sub-tasks: { 
            2 sequential sub-tasks: { 
                create managed nodegroup "m1",
                propagate tags to ASG for managed nodegroup "m1",
            },
            2 sequential sub-tasks: { 
                create managed nodegroup "m2",
                propagate tags to ASG for managed nodegroup "m2",
            },
        },
    } 
}
`))
			}
			{
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(context.Background(), makeNodeGroups("foo"), makeManagedNodeGroups("m1"))
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
				tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(context.Background(), makeNodeGroups("bar"), nil, &task{id: 1})
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
				p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{}, nil)
				ng := api.NewManagedNodeGroup()
				fakeVPCImporter := new(vpcfakes.FakeImporter)
				tasks := stackManager.NewManagedNodeGroupTask(context.Background(), []*api.ManagedNodeGroup{ng}, false, fakeVPCImporter)
				errs := tasks.DoAllSync()
				Expect(errs).To(HaveLen(1))
				Expect(errs[0]).To(MatchError(ContainSubstring("managed nodegroups cannot be created on IPv6 unowned clusters")))
			})
			When("finding the stack fails", func() {
				It("returns the stack error", func() {
					p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
					ng := api.NewManagedNodeGroup()
					fakeVPCImporter := new(vpcfakes.FakeImporter)
					tasks := stackManager.NewManagedNodeGroupTask(context.Background(), []*api.ManagedNodeGroup{ng}, false, fakeVPCImporter)
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(1))
					Expect(errs[0]).To(MatchError(ContainSubstring("not found")))
				})
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

func makeManagedNodeGroupsWithPropagatedTags(names ...string) []*api.ManagedNodeGroup {
	propagate := true
	var managedNodeGroups []*api.ManagedNodeGroup
	for _, name := range names {
		ng := api.NewManagedNodeGroup()
		ng.Name = name
		ng.PropagateASGTags = &propagate
		managedNodeGroups = append(managedNodeGroups, ng)
	}
	return managedNodeGroups
}
