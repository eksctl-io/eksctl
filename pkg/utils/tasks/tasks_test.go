package tasks

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TaskTree", func() {
	Context("With various sets of nested dummy tasks", func() {

		It("should have nice description", func() {
			{
				tasks := &TaskTree{Parallel: false}
				tasks.Append(&TaskTree{Parallel: false})
				Expect(tasks.Describe()).To(Equal("1 task: { no tasks }"))
				tasks.IsSubTask = true
				tasks.PlanMode = true
				tasks.Append(&TaskTree{Parallel: false, IsSubTask: true})
				fmt.Println(tasks.Describe())
				expected := []byte(`(plan) 
    2 sequential sub-tasks: { 
        no tasks,
        no tasks,
    }
`)
				Expect([]byte(tasks.Describe())).To(Equal(expected))
			}
			{
				tasks := &TaskTree{Parallel: false}
				subTask1 := &TaskTree{Parallel: false, IsSubTask: true}
				subTask1.Append(&TaskWithoutParams{
					Info: "t1.1",
				})
				tasks.Append(subTask1)

				Expect(tasks.Describe()).To(Equal("1 task: { t1.1 }"))

				subTask2 := &TaskTree{Parallel: false, IsSubTask: true}
				subTask2.Append(&TaskWithoutParams{
					Info: "t2.1",
				})
				subTask3 := &TaskTree{Parallel: true, IsSubTask: true}
				subTask3.Append(&TaskWithoutParams{
					Info: "t3.1",
				})
				subTask3.Append(&TaskWithoutParams{
					Info: "t3.2",
				})
				tasks.Append(subTask2)
				subTask1.Append(subTask3)

				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { 
    2 sequential sub-tasks: { 
        t1.1,
        2 parallel sub-tasks: { 
            t3.1,
            t3.2,
        },
    }, t2.1 
}
`))
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
				subTask1.Append(&TaskWithoutParams{
					Info: "t1.1",
					Call: func(errs chan error) error {
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
				subTask2.Append(&TaskWithoutParams{
					Info: "t2.1",
					Call: func(errs chan error) error {
						go func() {
							errs <- fmt.Errorf("never happens")
							close(errs)
						}()
						return nil
					},
				})
				tasks.Append(subTask2)

				subTask3 := &TaskTree{Parallel: true, IsSubTask: true}
				subTask3.Append(&TaskWithoutParams{
					Info: "t3.1",
					Call: func(errs chan error) error {
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
				subTask3.Append(&TaskWithoutParams{
					Info: "t3.2",
					Call: func(errs chan error) error {
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

				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { 
    2 sequential sub-tasks: { 
        t1.1,
        2 parallel sub-tasks: { 
            t3.1,
            t3.2,
        },
    }, t2.1 
}
`))

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
				subTask1.Append(&TaskWithoutParams{
					Info: "t1.1",
					Call: func(errs chan error) error {
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
				subTask2.Append(&TaskWithoutParams{
					Info: "t2.1",
					Call: func(errs chan error) error {
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
				subTask3.Append(&TaskWithoutParams{
					Info: "t3.1",
					Call: func(errs chan error) error {
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
				subTask3.Append(&TaskWithoutParams{
					Info: "t3.2",
					Call: func(errs chan error) error {
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

				Expect(tasks.Describe()).To(Equal(`
2 sequential tasks: { 
    2 sequential sub-tasks: { 
        t1.1,
        2 parallel sub-tasks: { 
            t3.1,
            t3.2,
        },
    }, t2.1 
}
`))

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

				tasks.Append(&TaskWithoutParams{
					Info: "t1.0",
					Call: func(errs chan error) error {
						close(errs)
						atomic.AddInt32(&counter, 1)
						return fmt.Errorf("t1.0 does not even bother and always returns an immediate error")
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.1",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(10 * time.Millisecond)
							errs <- nil
							close(errs)
							atomic.AddInt32(&counter, 1)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.2",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(100 * time.Millisecond)
							errs <- fmt.Errorf("t1.2 always fails")
							close(errs)
							atomic.AddInt32(&counter, 1)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.3",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(50 * time.Microsecond)
							errs <- fmt.Errorf("t1.3 always fails")
							close(errs)
							atomic.AddInt32(&counter, 1)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.4",
					Call: func(errs chan error) error {
						time.Sleep(150 * time.Millisecond)
						close(errs)
						atomic.AddInt32(&counter, 1)
						return fmt.Errorf("t1.4 does busy work and always returns an immediate error")
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.5",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(15 * time.Millisecond)
							errs <- nil
							close(errs)
							atomic.AddInt32(&counter, 1)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.6",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(15 * time.Millisecond)
							errs <- nil
							close(errs)
							atomic.AddInt32(&counter, 1)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.7",
					Call: func(errs chan error) error {
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

				tasks.Append(&TaskWithoutParams{
					Info: "t1",
					Call: func(errs chan error) error {
						close(errs)
						atomic.AddInt32(&counter, 1)
						return fmt.Errorf("t1.0 does not even bother and always returns an immediate error")
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t2",
					Call: func(errs chan error) error {
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

				tasks.Append(&TaskWithoutParams{
					Info: "t1.1",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(100 * time.Millisecond)
							errs <- fmt.Errorf("t1.1 always fails")
							close(errs)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.3",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(150 * time.Millisecond)
							errs <- nil
							close(errs)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.3",
					Call: func(errs chan error) error {
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

				tasks.Append(&TaskWithoutParams{
					Info: "t1.1",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(100 * time.Millisecond)
							errs <- fmt.Errorf("t1.1 always fails")
							close(errs)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.3",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(150 * time.Millisecond)
							errs <- nil
							close(errs)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.3",
					Call: func(errs chan error) error {
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

				tasks.Append(&TaskWithoutParams{
					Info: "t1.1",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(100 * time.Millisecond)
							errs <- nil
							close(errs)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.3",
					Call: func(errs chan error) error {
						go func() {
							time.Sleep(150 * time.Millisecond)
							errs <- nil
							close(errs)
						}()
						return nil
					},
				})

				tasks.Append(&TaskWithoutParams{
					Info: "t1.3",
					Call: func(errs chan error) error {
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
})
