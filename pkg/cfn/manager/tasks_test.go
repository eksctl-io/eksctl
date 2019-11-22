package manager

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

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

		Context("With various sets of nested dummy tasks", func() {

			It("should have nice description", func() {
				{
					tasks := &TaskTree{Parallel: false}
					tasks.Append(&TaskTree{Parallel: false})
					Expect(tasks.Describe()).To(Equal("1 task: { no tasks }"))
					tasks.IsSubTask = true
					tasks.PlanMode = true
					tasks.Append(&TaskTree{Parallel: false, IsSubTask: true})
					Expect(tasks.Describe()).To(Equal("(plan) 2 sequential sub-tasks: { no tasks, no tasks }"))
				}

				{
					tasks := &TaskTree{Parallel: false}
					subTask1 := &TaskTree{Parallel: false, IsSubTask: true}
					subTask1.Append(&taskWithoutParams{
						info: "t1.1",
					})
					tasks.Append(subTask1)

					Expect(tasks.Describe()).To(Equal("1 task: { t1.1 }"))

					subTask2 := &TaskTree{Parallel: false, IsSubTask: true}
					subTask2.Append(&taskWithoutParams{
						info: "t2.1",
					})
					subTask3 := &TaskTree{Parallel: true, IsSubTask: true}
					subTask3.Append(&taskWithoutParams{
						info: "t3.1",
					})
					subTask3.Append(&taskWithoutParams{
						info: "t3.2",
					})
					tasks.Append(subTask2)
					subTask1.Append(subTask3)

					Expect(tasks.Describe()).To(Equal("2 sequential tasks: { 2 sequential sub-tasks: { t1.1, 2 parallel sub-tasks: { t3.1, t3.2 } }, t2.1 }"))
				}
			})

			It("should execute orderly", func() {
				{
					var status struct {
						messages  []string
						mutex     sync.Mutex
						startTime time.Time
					}

					status.messages = []string{}

					updateStatus := func(msg string) {
						status.mutex.Lock()
						ts := time.Since(status.startTime).Round(50 * time.Millisecond).String()
						status.messages = append(status.messages,
							fmt.Sprintf("%s: %s", ts, msg),
						)
						status.mutex.Unlock()
					}

					tasks := &TaskTree{Parallel: false}
					subTask1 := &TaskTree{Parallel: false, IsSubTask: true}
					subTask1.Append(&taskWithoutParams{
						info: "t1.1",
						call: func(errs chan error) error {
							updateStatus("started t1.1")
							go func() {
								time.Sleep(100 * time.Millisecond)
								updateStatus("finished t1.1")
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})
					tasks.Append(subTask1)

					subTask2 := &TaskTree{Parallel: false, IsSubTask: true}
					subTask2.Append(&taskWithoutParams{
						info: "t2.1",
						call: func(errs chan error) error {
							go func() {
								errs <- fmt.Errorf("never happens")
								close(errs)
							}()
							return nil
						},
					})
					tasks.Append(subTask2)

					subTask3 := &TaskTree{Parallel: true, IsSubTask: true}
					subTask3.Append(&taskWithoutParams{
						info: "t3.1",
						call: func(errs chan error) error {
							updateStatus("started t3.1")
							go func() {
								time.Sleep(200 * time.Millisecond)
								updateStatus("finished t3.1")
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})
					subTask3.Append(&taskWithoutParams{
						info: "t3.2",
						call: func(errs chan error) error {
							updateStatus("started t3.2")
							go func() {
								time.Sleep(350 * time.Millisecond)
								updateStatus("finished t3.2")
								errs <- fmt.Errorf("t3.2 always fails")
								close(errs)
							}()
							return nil
						},
					})
					subTask1.Append(subTask3)

					Expect(tasks.Describe()).To(Equal("2 sequential tasks: { 2 sequential sub-tasks: { t1.1, 2 parallel sub-tasks: { t3.1, t3.2 } }, t2.1 }"))

					status.startTime = time.Now()
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(1))
					Expect(errs[0].Error()).To(Equal("t3.2 always fails"))

					Expect(status.messages).To(HaveLen(6))

					Expect(status.messages[0]).To(
						Equal("0s: started t1.1"),
					)
					Expect(status.messages[1]).To(
						Equal("100ms: finished t1.1"),
					)
					// t3.1 and t3.2 run in parallel, so may start approximately at the same time
					Expect(status.messages[2]).To(
						HavePrefix("100ms: started t3."),
					)
					Expect(status.messages[3]).To(
						HavePrefix("100ms: started t3."),
					)
					Expect(status.messages[4]).To(Equal(
						"300ms: finished t3.1",
					))
					Expect(status.messages[5]).To(Equal(
						"450ms: finished t3.2",
					))
				}

				{
					var status struct {
						messages  []string
						mutex     sync.Mutex
						startTime time.Time
					}

					status.messages = []string{}

					updateStatus := func(msg string) {
						status.mutex.Lock()
						ts := time.Since(status.startTime).Round(50 * time.Millisecond).String()
						status.messages = append(status.messages,
							fmt.Sprintf("%s: %s", ts, msg),
						)
						status.mutex.Unlock()
					}

					tasks := &TaskTree{Parallel: false}
					subTask1 := &TaskTree{Parallel: false, IsSubTask: true}
					subTask1.Append(&taskWithoutParams{
						info: "t1.1",
						call: func(errs chan error) error {
							updateStatus("started t1.1")
							go func() {
								time.Sleep(100 * time.Millisecond)
								updateStatus("finished t1.1")
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})
					tasks.Append(subTask1)

					subTask2 := &TaskTree{Parallel: false, IsSubTask: true}
					subTask2.Append(&taskWithoutParams{
						info: "t2.1",
						call: func(errs chan error) error {
							updateStatus("started t2.1")
							go func() {
								time.Sleep(150 * time.Millisecond)
								updateStatus("finished t2.1")
								errs <- fmt.Errorf("t2.1 always fails")
								close(errs)
							}()
							return nil
						},
					})
					tasks.Append(subTask2)

					subTask3 := &TaskTree{Parallel: true, IsSubTask: true}
					subTask3.Append(&taskWithoutParams{
						info: "t3.1",
						call: func(errs chan error) error {
							updateStatus("started t3.1")
							go func() {
								time.Sleep(200 * time.Millisecond)
								updateStatus("finished t3.1")
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})
					subTask3.Append(&taskWithoutParams{
						info: "t3.2",
						call: func(errs chan error) error {
							updateStatus("started t3.2")
							go func() {
								time.Sleep(350 * time.Millisecond)
								updateStatus("finished t3.2")
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})
					subTask1.Append(subTask3)

					Expect(tasks.Describe()).To(Equal("2 sequential tasks: { 2 sequential sub-tasks: { t1.1, 2 parallel sub-tasks: { t3.1, t3.2 } }, t2.1 }"))

					status.startTime = time.Now()
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(1))
					Expect(errs[0].Error()).To(Equal("t2.1 always fails"))

					Expect(status.messages).To(HaveLen(8))

					Expect(status.messages[0]).To(
						Equal("0s: started t1.1"),
					)
					Expect(status.messages[1]).To(
						Equal("100ms: finished t1.1"),
					)
					// t3.1 and t3.2 run in parallel, so may start approximately at the same time
					Expect(status.messages[2]).To(
						HavePrefix("100ms: started t3."),
					)
					Expect(status.messages[3]).To(
						HavePrefix("100ms: started t3."),
					)
					Expect(status.messages[4]).To(Equal(
						"300ms: finished t3.1",
					))
					Expect(status.messages[5]).To(Equal(
						"450ms: finished t3.2",
					))
					Expect(status.messages[6]).To(Equal(
						"450ms: started t2.1",
					))
					Expect(status.messages[7]).To(Equal(
						"600ms: finished t2.1",
					))
				}

				{
					tasks := &TaskTree{Parallel: false}
					Expect(tasks.DoAllSync()).To(HaveLen(0))
				}

				{
					tasks := &TaskTree{Parallel: false}
					tasks.Append(&TaskTree{Parallel: false})
					tasks.Append(&TaskTree{Parallel: true})
					Expect(tasks.DoAllSync()).To(HaveLen(0))
				}

				{
					tasks := &TaskTree{Parallel: true}

					counter := int32(0)

					tasks.Append(&taskWithoutParams{
						info: "t1.0",
						call: func(errs chan error) error {
							close(errs)
							atomic.AddInt32(&counter, 1)
							return fmt.Errorf("t1.0 does not even bother and always returns an immediate error")
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.1",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(10 * time.Millisecond)
								errs <- nil
								close(errs)
								atomic.AddInt32(&counter, 1)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.2",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(100 * time.Millisecond)
								errs <- fmt.Errorf("t1.2 always fails")
								close(errs)
								atomic.AddInt32(&counter, 1)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.3",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(50 * time.Microsecond)
								errs <- fmt.Errorf("t1.3 always fails")
								close(errs)
								atomic.AddInt32(&counter, 1)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.4",
						call: func(errs chan error) error {
							time.Sleep(150 * time.Millisecond)
							close(errs)
							atomic.AddInt32(&counter, 1)
							return fmt.Errorf("t1.4 does busy work and always returns an immediate error")
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.5",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(15 * time.Millisecond)
								errs <- nil
								close(errs)
								atomic.AddInt32(&counter, 1)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.6",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(15 * time.Millisecond)
								errs <- nil
								close(errs)
								atomic.AddInt32(&counter, 1)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.7",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(215 * time.Millisecond)
								errs <- nil
								close(errs)
								atomic.AddInt32(&counter, 1)
							}()
							return nil
						},
					})

					tasks.PlanMode = true

					Expect(tasks.DoAllSync()).To(HaveLen(0))

					tasks.PlanMode = false
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(4))
					Expect(errs[0].Error()).To(Equal("t1.0 does not even bother and always returns an immediate error"))
					Expect(errs[1].Error()).To(Equal("t1.3 always fails"))
					Expect(errs[2].Error()).To(Equal("t1.2 always fails"))
					Expect(errs[3].Error()).To(Equal("t1.4 does busy work and always returns an immediate error"))

					Expect(atomic.LoadInt32(&counter)).To(Equal(int32(8)))
				}

				{
					tasks := &TaskTree{Parallel: false}

					counter := int32(0)

					tasks.Append(&taskWithoutParams{
						info: "t1",
						call: func(errs chan error) error {
							close(errs)
							atomic.AddInt32(&counter, 1)
							return fmt.Errorf("t1.0 does not even bother and always returns an immediate error")
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t2",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(10 * time.Millisecond)
								errs <- nil
								close(errs)
								atomic.AddInt32(&counter, 1)
							}()
							return nil
						},
					})

					tasks.PlanMode = false
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(1))
					Expect(errs[0].Error()).To(Equal("t1.0 does not even bother and always returns an immediate error"))

					Expect(atomic.LoadInt32(&counter)).To(Equal(int32(1)))
				}

				{
					tasks := &TaskTree{Parallel: true}

					tasks.Append(&taskWithoutParams{
						info: "t1.1",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(100 * time.Millisecond)
								errs <- fmt.Errorf("t1.1 always fails")
								close(errs)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.3",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(150 * time.Millisecond)
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.3",
						call: func(errs chan error) error {
							go func() {
								errs <- fmt.Errorf("t1.3 always fails")
								close(errs)
							}()
							return nil
						},
					})

					tasks.PlanMode = true

					Expect(tasks.DoAllSync()).To(HaveLen(0))

					tasks.PlanMode = false
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(2))
					Expect(errs[0].Error()).To(Equal("t1.3 always fails"))
					Expect(errs[1].Error()).To(Equal("t1.1 always fails"))
				}

				{
					tasks := &TaskTree{Parallel: false}

					tasks.Append(&taskWithoutParams{
						info: "t1.1",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(100 * time.Millisecond)
								errs <- fmt.Errorf("t1.1 always fails")
								close(errs)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.3",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(150 * time.Millisecond)
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.3",
						call: func(errs chan error) error {
							go func() {
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})

					tasks.PlanMode = true

					Expect(tasks.DoAllSync()).To(HaveLen(0))

					tasks.PlanMode = false
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(1))
					Expect(errs[0].Error()).To(Equal("t1.1 always fails"))
				}

				{
					tasks := &TaskTree{Parallel: false}

					tasks.Append(&taskWithoutParams{
						info: "t1.1",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(100 * time.Millisecond)
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.3",
						call: func(errs chan error) error {
							go func() {
								time.Sleep(150 * time.Millisecond)
								errs <- nil
								close(errs)
							}()
							return nil
						},
					})

					tasks.Append(&taskWithoutParams{
						info: "t1.3",
						call: func(errs chan error) error {
							go func() {
								errs <- fmt.Errorf("t1.3 always fails")
								close(errs)
							}()
							return nil
						},
					})

					tasks.PlanMode = true

					Expect(tasks.DoAllSync()).To(HaveLen(0))

					tasks.PlanMode = false
					errs := tasks.DoAllSync()
					Expect(errs).To(HaveLen(1))
					Expect(errs[0].Error()).To(Equal("t1.3 always fails"))
				}
			})
		})

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
				// TODO use DescribeTable

				// The supportsManagedNodes argument has no effect on the Describe call, so the values are alternated
				// in these tests
				{
					tasks := stackManager.NewTasksToCreateNodeGroups(makeNodeGroups("bar", "foo"), true)
					Expect(tasks.Describe()).To(Equal(`2 parallel tasks: { create nodegroup "bar", create nodegroup "foo" }`))
				}
				{
					tasks := stackManager.NewTasksToCreateNodeGroups(makeNodeGroups("bar"), false)
					Expect(tasks.Describe()).To(Equal(`1 task: { create nodegroup "bar" }`))
				}
				{
					tasks := stackManager.NewTasksToCreateNodeGroups(makeNodeGroups("foo"), true)
					Expect(tasks.Describe()).To(Equal(`1 task: { create nodegroup "foo" }`))
				}
				{
					tasks := stackManager.NewTasksToCreateNodeGroups(nil, false)
					Expect(tasks.Describe()).To(Equal(`no tasks`))
				}
				{
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar", "foo"), nil, true)
					Expect(tasks.Describe()).To(Equal(`2 sequential tasks: { create cluster control plane "test-cluster", 2 parallel sub-tasks: { create nodegroup "bar", create nodegroup "foo" } }`))
				}
				{
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar"), nil, false)
					Expect(tasks.Describe()).To(Equal(`2 sequential tasks: { create cluster control plane "test-cluster", create nodegroup "bar" }`))
				}
				{
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(nil, nil, true)
					Expect(tasks.Describe()).To(Equal(`1 task: { create cluster control plane "test-cluster" }`))
				}
				{
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("bar", "foo"), makeManagedNodeGroups("m1", "m2"), false)
					Expect(tasks.Describe()).To(Equal(`2 sequential tasks: { create cluster control plane "test-cluster", 4 parallel sub-tasks: { create nodegroup "bar", create nodegroup "foo", create managed nodegroup "m1", create managed nodegroup "m2" } }`))
				}
				{
					tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(makeNodeGroups("foo"), makeManagedNodeGroups("m1"), true)
					Expect(tasks.Describe()).To(Equal(`2 sequential tasks: { create cluster control plane "test-cluster", 2 parallel sub-tasks: { create nodegroup "foo", create managed nodegroup "m1" } }`))
				}
			})
		})

	})
})
