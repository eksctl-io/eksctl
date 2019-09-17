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
		// TODO: handle API rate limitter error
		if err := drainer.EvictOrDeletePod(pod); err != nil {
			return pending, err
		}
	}
	return pending, nil
}

// NodeGroup drains a nodegroup
func NodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup, waitTimeout time.Duration, undo bool) error {
	evictError := 0
	drainer := &Helper{
		Client: clientSet,

		// TODO: Force, DeleteLocalData & IgnoreAllDaemonSets shouldn't
		// be enabled by default, we need flags to control thes, but that
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
		return errors.Wrapf(err, "checking if cluster implements policy API")
	}

	drainedNodes := sets.NewString()
	// loop until all nodes are drained to handle accidental scale-up
	// or any other changes in the ASG
	timer := time.After(waitTimeout)
	timeout := false
	for !timeout || evictError > 0 {
		select {
		case <-timer:
			timeout = true
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
				desired := CordonNode
				if undo {
					desired = UncordonNode
				}
				c := NewCordonHelper(&node, desired)
				if c.IsUpdateRequired() {
					err, patchErr := c.PatchOrReplace(clientSet)
					if patchErr != nil {
						logger.Warning(patchErr.Error())
					}
					if err != nil {
						logger.Critical(err.Error())
					}
					logger.Info("%s node %q", desired, node.Name)
				} else {
					logger.Debug("no need to %s node %q", desired, node.Name)
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

			for _, node := range nodes.Items {
				if newPendingNodes.Has(node.Name) {
					retry := 1
					pending, err := evictPods(drainer, &node)
					for pending > 0 {
						logger.Debug("%d pods to be evicted from %s", pending, node.Name)
						time.Sleep(5 * time.Second)
						pending, err = evictPods(drainer, &node)
						if err != nil && retry < podEvictionMaxRetries {
							logger.Warning("pod eviction error: \"%s\", on node: %s (retry in 5 sec, %d/%d)", err, node.Name, retry, podEvictionMaxRetries)
						} else if retry >= podEvictionMaxRetries {
							logger.Warning("pod eviction unable to complete after %d retries on node: %s", podEvictionMaxRetries, node.Name)
							evictError++
							break
						}
						retry++
					}
					drainedNodes.Insert(node.Name)
				}
			}
		}
	}

	if evictError > 0 {
		return fmt.Errorf("pod eviction error on nodegroup %q", ng.Name)
	}

	if timeout {
		return fmt.Errorf("timed out (after %s) waiting for nodedroup %q to be drain", waitTimeout, ng.Name)
	}

	return nil
}
