package harness

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (test *Test) listNodes(options metav1.ListOptions) (*v1.NodeList, error) {
	return test.harness.kubeClient.Core().Nodes().List(options)
}

// ListNodes returns all nodes that are part of the cluster.
func (test *Test) ListNodes(options metav1.ListOptions) *v1.NodeList {
	pl, err := test.listNodes(options)
	test.err(err)
	return pl
}

func (test *Test) nodeReady(node *v1.Node) (bool, error) {
	switch node.Status.Phase {
	case v1.NodeTerminated, v1.NodePending:
		return false, nil
	case v1.NodeRunning:
	default:
		for _, cond := range node.Status.Conditions {
			if cond.Type != v1.NodeReady {
				continue
			}
			return cond.Status == v1.ConditionTrue, nil
		}
		return false, fmt.Errorf("node ready condition not found")
	}
	return false, nil
}

// NodeReady returns whether a node is ready.
func (test *Test) NodeReady(node *v1.Node) bool {
	ready, err := test.nodeReady(node)
	test.err(err)
	return ready
}

// waitForDeploymentReady waits until all replica pods are running and ready.
func (test *Test) waitForNodesReady(expectedNodes int, timeout time.Duration) error {
	numReady := 0

	test.Infof("waiting for %d nodes to be ready", expectedNodes)

	return wait.Poll(time.Second, timeout, func() (bool, error) {
		current, err := test.listNodes(metav1.ListOptions{})
		if err != nil {
			test.Debugf("api server not ready: %v", err)
			return false, nil
		}

		currentNumReady := 0
		for i := range current.Items {
			node := &current.Items[i]
			ready, err := test.nodeReady(node)
			if err != nil {
				return false, err
			}
			if ready {
				currentNumReady++
			}
		}

		if numReady != currentNumReady {
			numReady = currentNumReady
			test.Debugf("nodes ready: %d/%d", numReady, expectedNodes)
		}
		if currentNumReady == expectedNodes {
			return true, nil
		}

		return false, nil
	})
}

// WaitForNodesReady waits until the specified number of nodes are running and
// ready. The function waits until the exact number is of expected node is matched.
func (test *Test) WaitForNodesReady(expectedNodes int, timeout time.Duration) {
	err := test.waitForNodesReady(expectedNodes, timeout)
	test.err(err)
}
