package drain

import (
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// this is our custom addition, it's not part of the package
// we copied from Kubernetes

func evictPods(drainer *Helper, node *corev1.Node) (int, error) {
	list, errs := drainer.GetPodsForDeletion(node.Name)
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
		if err := drainer.EvictOrDeletePod(pod); err != nil {
			return pending, err
		}
	}
	return pending, nil
}

// NodeGroup drains a nodegroup
func NodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup, waitTimeout time.Duration, undo bool) error {
	drainer := &Helper{
		Client: clientSet,

		// TODO: Force, DeleteLocalData & IgnoreAllDaemonSets shouldn't
		// be enabled by default, we need flags to control these, but that
		// requires more improvements in the underlying drain package,
		// as it currently produces errors and warnings with references
		// to kubectl flags
		Force:               true,
		DeleteLocalData:     true,
		IgnoreAllDaemonSets: true,

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

	if err := drainer.CanUseEvictions(); err != nil {
		return errors.Wrap(err, "checking if cluster implements policy API")
	}

	drainedNodes := sets.NewString()
	// loop until all nodes are drained to handle accidental scale-up
	// or any other changes in the ASG
	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timed out (after %s) waiting for nodegroup %q to be drained", waitTimeout, ng.Name)
		default:
			nodes, err := clientSet.CoreV1().Nodes().List(ng.ListOptions())
			if err != nil {
				return err
			}

			if len(nodes.Items) == 0 {
				logger.Warning("no nodes found in nodegroup %q (label selector: %q)", ng.Name, ng.ListOptions().LabelSelector)
				return nil
			}

			newPendingNodes := sets.NewString()

			for _, node := range nodes.Items {
				if drainedNodes.Has(node.Name) {
					continue // already drained, get next one
				}
				newPendingNodes.Insert(node.Name)
				desired := !undo
				c := NewCordonHelper(&node, desired)
				if c.IsUpdateRequired() {
					err, patchErr := c.PatchOrReplace(clientSet)
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

			if undo {
				return nil // no need to kill any pods
			}

			if drainedNodes.HasAll(newPendingNodes.List()...) {
				logger.Success("drained nodes: %v", drainedNodes.List())
				return nil // no new nodes were seen
			}

			logger.Debug("already drained: %v", drainedNodes.List())
			logger.Debug("will drain: %v", newPendingNodes.List())

		evict_loop:
			for _, node := range nodes.Items {
				if newPendingNodes.Has(node.Name) {
					pending, err := evictPods(drainer, &node)
					if err != nil {
						logger.Warning("pod eviction error (%q) on node %s – will retry after delay of %s", err, node.Name, retryDelay)
						retryTimer := time.NewTimer(retryDelay)
						select {
						case <-retryTimer.C:
						case <-timer.C:
							retryTimer.Stop()
							break evict_loop
						}
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

	return nil
}

func cordonStatus(desired bool) string {
	if desired {
		return "cordon"
	}
	return "uncordon"
}
