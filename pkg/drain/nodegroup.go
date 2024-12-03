package drain

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/logger"
	cmap "github.com/orcaman/concurrent-map"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/drain/evictor"
	"github.com/weaveworks/eksctl/pkg/eks"
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
	clientSet             kubernetes.Interface
	evictor               Evictor
	ng                    eks.KubeNodeGroup
	nodeDrainWaitPeriod   time.Duration
	podEvictionWaitPeriod time.Duration
	undo                  bool
	parallel              int
}

func NewNodeGroupDrainer(clientSet kubernetes.Interface, ng eks.KubeNodeGroup, maxGracePeriod, nodeDrainWaitPeriod time.Duration, podEvictionWaitPeriod time.Duration, undo, disableEviction bool, parallel int) NodeGroupDrainer {
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
		evictor:               evictor.New(clientSet, maxGracePeriod, ignoreDaemonSets, disableEviction),
		clientSet:             clientSet,
		ng:                    ng,
		nodeDrainWaitPeriod:   nodeDrainWaitPeriod,
		podEvictionWaitPeriod: podEvictionWaitPeriod,
		undo:                  undo,
		parallel:              parallel,
	}
}

// Drain drains a nodegroup
func (n *NodeGroupDrainer) Drain(ctx context.Context, sem *semaphore.Weighted) error {
	if err := n.evictor.CanUseEvictions(); err != nil {
		return fmt.Errorf("checking if cluster implements policy API: %w", err)
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

	drainedNodes := cmap.New()

	var evictErr error
	// loop until all nodes are drained to handle accidental scale-up
	// or any other changes in the ASG
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for nodegroup %q to be drained", n.ng.NameString())
		default:
			if evictErr != nil {
				return evictErr
			}
			nodes, err := n.clientSet.CoreV1().Nodes().List(context.TODO(), listOptions)
			if err != nil {
				return err
			}
			n.toggleCordon(true, nodes)

			newPendingNodes := sets.New[string]()

			for _, node := range nodes.Items {
				if !drainedNodes.Has(node.Name) {
					newPendingNodes.Insert(node.Name)
				}
			}

			if newPendingNodes.Len() == 0 {
				logger.Success("drained all nodes: %v", mapToList(drainedNodes.Items()))
				return nil // no new nodes were seen
			}

			logger.Debug("already drained: %v", mapToList(drainedNodes.Items()))
			logger.Debug("will drain: %v", sets.List(newPendingNodes))

			g, ctx := errgroup.WithContext(ctx)
			for _, node := range sets.List(newPendingNodes) {
				node := node
				g.Go(func() error {
					if err := sem.Acquire(ctx, 1); err != nil {
						return fmt.Errorf("failed to acquire semaphore: %w", err)
					}
					defer sem.Release(1)

					drainedNodes.Set(node, nil)
					logger.Debug("starting drain of node %s", node)
					if err := n.evictPods(ctx, node); err != nil {
						logger.Warning("pod eviction error (%q) on node %s", err, node)
						time.Sleep(retryDelay)
						return err
					}

					drainedNodes.Set(node, nil)

					if n.nodeDrainWaitPeriod > 0 {
						logger.Debug("waiting for %.0f seconds before draining next node", n.nodeDrainWaitPeriod.Seconds())
						time.Sleep(n.nodeDrainWaitPeriod)
					}
					return nil
				})
			}
			// We need to loop even if this is an error to check whether the error was a
			// context timeout or something else.  This lets us log timout errors consistently
			evictErr = g.Wait()
		}
	}
}

func mapToList(m map[string]interface{}) []string {
	var list []string
	for key := range m {
		list = append(list, key)
	}
	return list
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

func (n *NodeGroupDrainer) evictPods(ctx context.Context, node string) error {
	// Loop until context times out.  We want to continually try to remove pods
	// from the node as their eviction status changes.
	previousReportTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for node %q to be drained", node)
		default:
			list, errs := n.evictor.GetPodsForEviction(node)
			if len(errs) > 0 {
				return fmt.Errorf("errs: %v", errs) // TODO: improve formatting
			}
			if list == nil || len(list.Pods()) == 0 {
				return nil
			}
			pods := list.Pods()
			if w := list.Warnings(); w != "" {
				logger.Warning(w)
			}
			// This log message is important but can be noisy.  It's useful to only
			// update on a node every minute.
			if time.Since(previousReportTime) > time.Minute*1 && len(pods) > 0 {
				logger.Warning("%d pods are unevictable from node %s", len(pods), node)
				previousReportTime = time.Now()
			}
			logger.Debug("%d pods to be evicted from %s", pods, node)
			failedEvictions := false
			for _, pod := range pods {
				if err := n.evictor.EvictOrDeletePod(pod); err != nil {
					if !isEvictionErrorRecoverable(err) {
						return fmt.Errorf("unrecoverable error evicting pod: %s/%s: %w", pod.Namespace, pod.Name, err)
					}
					logger.Debug("recoverable pod eviction failure: %q", err)
					failedEvictions = true
				}
			}
			if failedEvictions {
				time.Sleep(n.podEvictionWaitPeriod)
			}
		}
	}
}

func cordonStatus(desired bool) string {
	if desired {
		return "cordon"
	}
	return "uncordon"
}

func isEvictionErrorRecoverable(err error) bool {
	var recoverableCheckerFuncs []func(error) bool
	recoverableCheckerFuncs = append(
		recoverableCheckerFuncs,
		apierrors.IsGone,
		apierrors.IsNotFound,
		apierrors.IsResourceExpired,
		apierrors.IsServerTimeout,
		apierrors.IsServiceUnavailable,
		apierrors.IsTimeout,
		// IsTooManyRequests also captures PDB errors
		apierrors.IsTooManyRequests,
	)

	for _, f := range recoverableCheckerFuncs {
		if f(err) {
			return true
		}
	}
	return false
}
