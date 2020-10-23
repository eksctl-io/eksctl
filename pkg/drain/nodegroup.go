package drain

import (
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

//go:generate counterfeiter -o fakes/fake_drainer.go . Drainer
type Drainer interface {
	CanUseEvictions() error
	EvictOrDeletePod(pod corev1.Pod) error
	GetPodsForDeletion(nodeName string) (*evictor.PodDeleteList, []error)
}

type NodeGroupDrainer struct {
	clientSet   kubernetes.Interface
	drainer     Drainer
	ng          eks.KubeNodeGroup
	waitTimeout time.Duration
	undo        bool
}

func NewNodeGroupDrainer(clientSet kubernetes.Interface, ng eks.KubeNodeGroup, waitTimeout time.Duration, maxGracePeriod time.Duration, undo bool) NodeGroupDrainer {
	drainer := &evictor.Evictor{
		Client: clientSet,

		// TODO: Force, DeleteLocalData & IgnoreAllDaemonSets shouldn't
		// be enabled by default, we need flags to control these, but that
		// requires more improvements in the underlying drain package,
		// as it currently produces errors and warnings with references
		// to kubectl flags
		Force:               true,
		DeleteLocalData:     true,
		IgnoreAllDaemonSets: true,

		MaxGracePeriodSeconds: int(maxGracePeriod.Seconds()),

		// TODO: ideally only the list of well-known DaemonSets should
		// be set by default
		IgnoreDaemonSets: []metav1.ObjectMeta{
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
		},
	}

	return NodeGroupDrainer{
		drainer:     drainer,
		clientSet:   clientSet,
		ng:          ng,
		waitTimeout: waitTimeout,
		undo:        undo,
	}
}

// NodeGroup drains a nodegroup
func (n *NodeGroupDrainer) Drain() error {

	if err := n.drainer.CanUseEvictions(); err != nil {
		return errors.Wrap(err, "checking if cluster implements policy API")
	}

	drainedNodes := sets.NewString()
	// loop until all nodes are drained to handle accidental scale-up
	// or any other changes in the ASG
	timer := time.NewTimer(n.waitTimeout)
	defer timer.Stop()

	timeoutErr := fmt.Errorf("timed out (after %s) waiting for nodegroup %q to be drained", n.waitTimeout, n.ng.NameString())
	for {
		select {
		case <-timer.C:
			return timeoutErr
		default:
			nodes, err := n.clientSet.CoreV1().Nodes().List(n.ng.ListOptions())
			if err != nil {
				return err
			}

			if len(nodes.Items) == 0 {
				logger.Warning("no nodes found in nodegroup %q (label selector: %q)", n.ng.NameString(), n.ng.ListOptions().LabelSelector)
				return nil
			}

			newPendingNodes := sets.NewString()

			for _, node := range nodes.Items {
				if drainedNodes.Has(node.Name) {
					continue // already drained, get next one
				}
				newPendingNodes.Insert(node.Name)
				desired := !n.undo
				c := NewCordonHelper(&node, desired)
				if c.IsUpdateRequired() {
					err, patchErr := c.PatchOrReplace(n.clientSet)
					if patchErr != nil {
						logger.Warning(patchErr.Error())
					}
					if err != nil {
						logger.Critical(err.Error())
					}
					logger.Info("%s node %q", cordonStatus(desired), node.Name)
				} else {
					logger.Debug("no need to %s node %q", cordonStatus(desired), node.Name)
				}
			}

			if n.undo {
				return nil // no need to kill any pods
			}

			if drainedNodes.HasAll(newPendingNodes.List()...) {
				logger.Success("drained nodes: %v", drainedNodes.List())
				return nil // no new nodes were seen
			}

			logger.Debug("already drained: %v", drainedNodes.List())
			logger.Debug("will drain: %v", newPendingNodes.List())

			for _, node := range nodes.Items {
				if newPendingNodes.Has(node.Name) {
					select {
					case <-timer.C:
						return timeoutErr
					default:
						pending, err := n.evictPods(&node)
						if err != nil {
							logger.Warning("pod eviction error (%q) on node %s – will retry after delay of %s", err, node.Name, retryDelay)
							time.Sleep(retryDelay)
							continue
						}
						logger.Debug("%d pods to be evicted from %s", pending, node.Name)
						if pending == 0 {
							drainedNodes.Insert(node.Name)
						}
					}
				}
			}
		}
	}
}

func (n *NodeGroupDrainer) evictPods(node *corev1.Node) (int, error) {
	list, errs := n.drainer.GetPodsForDeletion(node.Name)
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
		if err := n.drainer.EvictOrDeletePod(pod); err != nil {
			return pending, err
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
