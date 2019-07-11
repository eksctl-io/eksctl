package eks

import (
	"fmt"
	"time"

	"github.com/weaveworks/eksctl/pkg/spotinst"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/iam"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func isNodeReady(node *corev1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func getNodes(clientSet kubernetes.Interface, ng *api.NodeGroup) (int, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(ng.ListOptions())
	if err != nil {
		return 0, err
	}
	logger.Info("nodegroup %q has %d node(s)", ng.Name, len(nodes.Items))
	counter := 0
	for _, node := range nodes.Items {
		// logger.Debug("node[%d]=%#v", n, node)
		ready := "not ready"
		if isNodeReady(&node) {
			ready = "ready"
			counter++
		}
		logger.Info("node %q is %s", node.ObjectMeta.Name, ready)
	}
	return counter, nil
}

// WaitForNodes waits till the nodes are ready
func (c *ClusterProvider) WaitForNodes(clientSet kubernetes.Interface, ng *api.NodeGroup) error {
	if ng.MinSize == nil || *ng.MinSize == 0 {
		return nil
	}
	timer := time.After(c.Provider.WaitTimeout())
	timeout := false
	readyNodes := sets.NewString()
	watcher, err := clientSet.CoreV1().Nodes().Watch(ng.ListOptions())
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}

	counter, err := getNodes(clientSet, ng)
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}

	logger.Info("waiting for at least %d node(s) to become ready in %q", *ng.MinSize, ng.Name)
	for !timeout && counter < *ng.MinSize {
		select {
		case event := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if event.Object != nil && event.Type != watch.Deleted {
				if node, ok := event.Object.(*corev1.Node); ok {
					if isNodeReady(node) {
						readyNodes.Insert(node.Name)
						counter = readyNodes.Len()
						logger.Debug("node %q is ready in %q", node.Name, ng.Name)
					} else {
						logger.Debug("node %q seen in %q, but not ready yet", node.Name, ng.Name)
						logger.Debug("node = %#v", *node)
					}
				}
			}
		case <-timer:
			timeout = true
		}
	}
	watcher.Stop()
	if timeout {
		return fmt.Errorf("timed out (after %s) waiting for at least %d nodes to join the cluster and become ready in %q", c.Provider.WaitTimeout(), *ng.MinSize, ng.Name)
	}

	if _, err = getNodes(clientSet, ng); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}

// GetNodeGroupIAM retrieves the IAM configuration of the given nodegroup
func (c *ClusterProvider) GetNodeGroupIAM(stackManager *manager.StackCollection, spec *api.ClusterConfig, ng *api.NodeGroup) error {
	stacks, err := stackManager.DescribeNodeGroupStacks()
	if err != nil {
		return err
	}

	for _, s := range stacks {
		if stackManager.GetNodeGroupName(s) == ng.Name {
			if !stackManager.StackStatusIsNotTransitional(s) {
				return fmt.Errorf("nodegroup %q is in transitional state (%q)", ng.Name, *s.StackStatus)
			}
			return iam.UseFromNodeGroup(c.Provider, s, ng)
		}
	}

	return fmt.Errorf("stack not found for nodegroup %q", ng.Name)
}

// GetNodeGroupSpotinstOceanId retrieves the Spotinst Ocean cluster identifier.
func (c *ClusterProvider) GetNodeGroupSpotinstOceanId(stackManager *manager.StackCollection, spec *api.ClusterConfig, ng *api.NodeGroup) error {
	logger.Debug("attempting to find spotinst ocean nodegroup stack")

	stacks, err := stackManager.DescribeNodeGroupStacks()
	if err != nil {
		return err
	}

	for _, s := range stacks {
		if stackManager.GetNodeGroupName(s) == "ocean" {
			if !stackManager.StackStatusIsNotTransitional(s) {
				return fmt.Errorf("nodegroup %q is in transitional state (%q)", ng.Name, *s.StackStatus)
			}

			logger.Debug("nodegroup %q will join to an existing spotinst ocean", ng.Name)
			return spotinst.UseNodeGroupSpotinstOceanId(c.Provider, s, ng)
		}
	}

	return fmt.Errorf("stack not found for spotinst ocean nodegroup")
}
