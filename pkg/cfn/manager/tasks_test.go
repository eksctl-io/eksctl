package manager

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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

		stackManager *StackCollection
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
		Context("With real tasks", func() {

			BeforeEach(func() {

				p = mockprovider.NewMockProvider()

				cfg = newClusterConfig("test-cluster")

				stackManager = NewStackCollection(p, cfg)
			})

			It("should have nice description", func() {
				makeNodeGroups := func(names ...string) []*api.NodeGroup {
					var nodeGroups []*api.NodeGroup
					for _, name := range names {
						ng := api.NewNodeGroup()
						ng.Name = name
						nodeGroups = append(nodeGroups, ng)
					}
					return nodeGroups
				}

				makeManagedNodeGroups := func(names ...string) []*api.ManagedNodeGroup {
					var managedNodeGroups []*api.ManagedNodeGroup
					for _, name := range names {
						ng := api.NewManagedNodeGroup()
						ng.Name = name
						managedNodeGroups = append(managedNodeGroups, ng)
					}
					return managedNodeGroups
				}

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
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar", "foo"), nil, true)
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
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar"), nil, false)
					Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { create cluster control plane "test-cluster", create nodegroup "bar" 
}
`))
				}
				{
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(nil, nil, true)
					Expect(tasks.Describe()).To(Equal(`1 task: { create cluster control plane "test-cluster" }`))
				}
				{
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar", "foo"), makeManagedNodeGroups("m1", "m2"), false)
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
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("foo"), makeManagedNodeGroups("m1"), true)
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
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar"), nil, false, &task{id: 1})
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
	})
})
