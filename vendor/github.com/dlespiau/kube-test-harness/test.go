package harness

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"github.com/dlespiau/kube-test-harness/logger"
	"github.com/dlespiau/kube-test-harness/testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"golang.org/x/sync/errgroup"
)

type finalizer func() error

// Test is a single test running in a kubernetes cluster.
type Test struct {
	// ID is a unique identifier for the test, defined from the test function name.
	ID string
	// Namespace is name of the namespace automatically crafted by Setup for the
	// test to run in.
	Namespace string

	nextObjectID uint64
	harness      *Harness
	t            testing.T
	logger       logger.Logger
	inError      bool
	namespaces   []string // List of namespaces created by the test
	cleanUpFns   []finalizer
}

func testLogger(l logger.Logger, t testing.T) logger.Logger {
	return l.ForTest(t)
}

// NewTest creates a new test. Call Close() to free kubernetes resources
// allocated during the test.
func (h *Harness) NewTest(t testing.T) *Test {
	// TestCtx is used among others for namespace names where '/' is forbidden
	prefix := strings.TrimPrefix(
		strings.Replace(
			t.Name(),
			"/",
			"-",
			-1,
		),
		"Test",
	)

	id := toSnake(prefix) + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	test := &Test{
		ID:      id,
		harness: h,
		t:       t,
		logger:  testLogger(h.options.Logger, t),
	}
	test.Namespace = test.getObjID("ns")

	test.Infof("using API server %s", h.apiServer)

	return test
}

// getObjID returns an unique ID that can be used to name kubernetes objects. We
// also encode the object type in the name.
func (t *Test) getObjID(objectType string) string {
	id := atomic.AddUint64(&t.nextObjectID, 1)
	return t.ID + "-" + objectType + "-" + fmt.Sprintf("%d", id)
}

// Setup setups the test to be run in the Test.Namespace temporary namespace.
func (t *Test) Setup() *Test {
	t.CreateNamespace(t.Namespace)
	return t
}

func podContainersReady(pod *v1.Pod) (numReady int, numContainers int) {
	for _, cond := range pod.Status.Conditions {
		if cond.Type != v1.PodReady {
			continue
		}
		numContainers++
		if cond.Status == v1.ConditionTrue {
			numReady++
		}
	}
	return numReady, numContainers
}

// XXX: Maybe a bit too simplistic to synthesize a status.
func podStatus(pod *v1.Pod) (status, containerName string) {
	for _, cs := range pod.Status.ContainerStatuses {
		if state := cs.State.Waiting; state != nil {
			return state.Reason, cs.Name
		}
	}
	return "Ready", ""
}

type dumpLogs struct {
	pod           v1.Pod
	containerName string
}

// DumpNamespace writes to w information about the pods in a namespace.
func (t *Test) DumpNamespace(w io.Writer, ns string) {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	fmt.Fprintf(w, "\n=== pods, namespace=%s\n\n", ns)

	fmt.Fprintln(tw, "NAME\t  READY\t  STATUS")

	var logs []dumpLogs
	for _, pod := range t.ListPods(ns, metav1.ListOptions{}).Items {
		numReady, numContainers := podContainersReady(&pod)
		status, containerName := podStatus(&pod)
		if status != "Ready" {
			logs = append(logs, dumpLogs{pod, containerName})
		}

		fmt.Fprintf(tw, "%s\t  %d/%d\t  %s\n",
			pod.Name,
			numReady, numContainers,
			status,
		)
	}

	tw.Flush()

	for _, l := range logs {
		fmt.Fprintf(w, "\n=== logs, pod=%s, container=%s\n\n", l.pod.Name, l.containerName)
		if err := t.PodLogs(w, &l.pod, l.containerName); err != nil {
			fmt.Println(err)

		}
	}
}

// DumpTestState writes to w information about the objects created by the test.
func (t *Test) DumpTestState(w io.Writer) {
	// kube-system is interesting because it has pods that could make tests fail
	// (eg. kube-dns)
	namespaces := make([]string, len(t.namespaces)+1)
	namespaces[0] = "kube-system"
	copy(namespaces[1:], t.namespaces)

	for _, ns := range namespaces {
		t.DumpNamespace(w, ns)
	}
}

func (t *Test) dumpTestState() {
	t.DumpTestState(os.Stderr)
	fmt.Fprintln(os.Stderr)
}

func (t *Test) fatal(args ...interface{}) {
	t.t.Fatal(args...)
}

// Close frees all kubernetes resources allocated during the test.
func (t *Test) Close() {
	// We're being called while panicking, don't cleanup!
	if r := recover(); r != nil {
		t.dumpTestState()
		panic(r)
	}
	if (t.t.Failed()) || t.inError {
		t.dumpTestState()
		return
	}

	if t.harness.options.NoCleanup {
		return
	}

	var eg errgroup.Group

	for i := len(t.cleanUpFns) - 1; i >= 0; i-- {
		eg.Go(t.cleanUpFns[i])
	}

	if err := eg.Wait(); err != nil {
		t.fatal(err)
	}
}

func (t *Test) err(err error) {
	if err != nil {
		t.inError = true
		t.fatal(err)
	}
}

func (t *Test) addNamespace(ns string) {
	t.namespaces = append(t.namespaces, ns)
}

func (t *Test) removeNamespace(ns string) {
	new := make([]string, 0, len(t.namespaces)-1)
	for _, s := range t.namespaces {
		if s == ns {
			continue
		}
		new = append(new, s)
	}
	t.namespaces = new
}

func (t *Test) addFinalizer(fn finalizer) {
	t.cleanUpFns = append(t.cleanUpFns, fn)
}

// Debug prints a debug message.
func (t *Test) Debug(msg string) {
	t.t.Helper()
	t.logger.Logf(logger.Debug, msg)
}

// Debugf prints a debug message with a format string.
func (t *Test) Debugf(f string, args ...interface{}) {
	t.t.Helper()
	t.logger.Logf(logger.Debug, f, args...)
}

// Info prints an informational message.
func (t *Test) Info(msg string) {
	t.t.Helper()
	t.logger.Log(logger.Info, msg)
}

// Infof prints a informational message with a format string.
func (t *Test) Infof(f string, args ...interface{}) {
	t.t.Helper()
	t.logger.Logf(logger.Info, f, args...)
}
