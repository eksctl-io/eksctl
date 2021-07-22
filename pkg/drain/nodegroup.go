package drain

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/pkg/drain/evictor"

	corev1 "k8s.io/api/core/v1"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/eks"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

// this is our custom addition, it's not part of the package
// we copied from Kubernetes

// retryDelay is how long is slept before retry after an error occurs during drainage
const retryDelay = 5 * time.Second

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_evictor.go . Evictor
type Evictor interface {
	CanUseEvictions() error
	EvictOrDeletePod(pod corev1.Pod) error
	GetPodsForEviction(nodeName string) (*evictor.PodDeleteList, []error)
}

type NodeGroupDrainer struct {
	clientSet   kubernetes.Interface
	evictor     Evictor
	ng          eks.KubeNodeGroup
	waitTimeout time.Duration
	undo        bool
}

func NewNodeGroupDrainer(clientSet kubernetes.Interface, ng eks.KubeNodeGroup, waitTimeout time.Duration, maxGracePeriod time.Duration, undo bool, disableEviction bool) NodeGroupDrainer {
	ignoreDaemonSets := []metav1.ObjectMeta{
		{
			Namespace: "kube-system",
			Name:      "aws-node",
		},
		{
			Namespace: "kube-system",
			Name:      "kube-proxy",
		},
		{
			Name: "node-exporter",
		},
		{
			Name: "prom-node-exporter",
		},
		{
			Name: "weave-scope",
		},
		{
			Name: "weave-scope-agent",
		},
		{
			Name: "weave-net",
		},
	}

	return NodeGroupDrainer{
		evictor:     evictor.New(clientSet, maxGracePeriod, ignoreDaemonSets, disableEviction),
		clientSet:   clientSet,
		ng:          ng,
		waitTimeout: waitTimeout,
		undo:        undo,
	}
}

// Drain drains a nodegroup
func (n *NodeGroupDrainer) Drain() error {
	if err := n.evictor.CanUseEvictions(); err != nil {
		return errors.Wrap(err, "checking if cluster implements policy API")
	}

	listOptions := n.ng.ListOptions()
	nodes, err := n.clientSet.CoreV1().Nodes().List(context.TODO(), listOptions)
	if err != nil {
		return err
	}

	if len(nodes.Items) == 0 {
		logger.Warning("no nodes found in nodegroup %q (label selector: %q)", n.ng.NameString(), n.ng.ListOptions().LabelSelector)
		return nil
	}

	if n.undo {
		n.toggleCordon(false, nodes)
		return nil // no need to kill any pods
	}

	drainedNodes := sets.NewString()
	// loop until all nodes are drained to handle accidental scale-up
	// or any other changes in the ASG
	timer := time.NewTimer(n.waitTimeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timed out (after %s) waiting for nodegroup %q to be drained", n.waitTimeout, n.ng.NameString())
		default:
			nodes, err := n.clientSet.CoreV1().Nodes().List(context.TODO(), listOptions)
			if err != nil {
				return err
			}
			n.toggleCordon(true, nodes)

			newPendingNodes := sets.NewString()

			for _, node := range nodes.Items {
				if !drainedNodes.Has(node.Name) {
					newPendingNodes.Insert(node.Name)
				}
			}

			if newPendingNodes.Len() == 0 {
				logger.Success("drained all nodes: %v", drainedNodes.List())
				return nil // no new nodes were seen
			}

			logger.Debug("already drained: %v", drainedNodes.List())
			logger.Debug("will drain: %v", newPendingNodes.List())

			for _, node := range newPendingNodes.List() {
				pending, err := n.evictPods(node)
				if err != nil {
					logger.Warning("pod eviction error (%q) on node %s", err, node)
					time.Sleep(retryDelay)
					continue
				}
				logger.Debug("%d pods to be evicted from %s", pending, node)
				if pending == 0 {
					drainedNodes.Insert(node)
				}
			}
		}
	}
}

func (n *NodeGroupDrainer) toggleCordon(cordon bool, nodes *corev1.NodeList) {
	for _, node := range nodes.Items {
		c := NewCordonHelper(&node, cordon)
		if c.IsUpdateRequired() {
			err, patchErr := c.PatchOrReplace(n.clientSet)
			if patchErr != nil {
				logger.Warning(patchErr.Error())
			}
			if err != nil {
				logger.Critical(err.Error())
			}
			logger.Info("%s node %q", cordonStatus(cordon), node.Name)
		} else {
			logger.Debug("no need to %s node %q", cordonStatus(cordon), node.Name)
		}
	}

}

func (n *NodeGroupDrainer) evictPods(node string) (int, error) {
	list, errs := n.evictor.GetPodsForEviction(node)
	if len(errs) > 0 {
		return 0, fmt.Errorf("errs: %v", errs) // TODO: improve formatting
	}
	if w := list.Warnings(); w != "" {
		logger.Warning(w)
	}
	pods := list.Pods()
	pending := len(pods)
	for _, pod := range pods {
		// TODO: handle API rate limiter error
		if err := n.evictor.EvictOrDeletePod(pod); err != nil {
			return pending, errors.Wrapf(err, "error evicting pod: %s/%s", pod.Namespace, pod.Name)
		}
	}
	return pending, nil
}

func cordonStatus(desired bool) string {
	if desired {
		return "cordon"
	}
	return "uncordon"
}
