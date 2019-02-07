package eks

import (
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
)

// CreateOrUpdateNodeGroupAuthConfigMap creates or updates the auth config map for the given nodegroup
func (c *ClusterProvider) CreateOrUpdateNodeGroupAuthConfigMap(clientSet *clientset.Clientset, ng *api.NodeGroup) error {
	cm := &corev1.ConfigMap{}
	client := clientSet.CoreV1().ConfigMaps(utils.AuthConfigMapNamespace)
	create := false

	if existing, err := client.Get(utils.AuthConfigMapName, metav1.GetOptions{}); err != nil {
		if kerr.IsNotFound(err) {
			create = true
		} else {
			return errors.Wrapf(err, "getting auth ConfigMap")
		}
	} else {
		*cm = *existing
	}

	if create {
		cm, err := utils.NewAuthConfigMap(ng.IAM.InstanceRoleARN)
		if err != nil {
			return errors.Wrap(err, "constructing auth ConfigMap")
		}
		if _, err := client.Create(cm); err != nil {
			return errors.Wrap(err, "creating auth ConfigMap")
		}
		logger.Debug("created auth ConfigMap for %s", ng.Name)
		return nil
	}

	if err := utils.UpdateAuthConfigMap(cm, ng.IAM.InstanceRoleARN); err != nil {
		return errors.Wrap(err, "creating an update for auth ConfigMap")
	}
	if _, err := client.Update(cm); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap")
	}
	logger.Debug("updated auth ConfigMap for %s", ng.Name)
	return nil
}

func isNodeReady(node *corev1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func getNodes(clientSet *clientset.Clientset, ng *api.NodeGroup) (int, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(ng.ListOptions())
	if err != nil {
		return 0, err
	}
	logger.Info("nodegroup %q has %d node(s)", ng.Name, len(nodes.Items))
	for _, node := range nodes.Items {
		// logger.Debug("node[%d]=%#v", n, node)
		ready := "not ready"
		if isNodeReady(&node) {
			ready = "ready"
		}
		logger.Info("node %q is %s", node.ObjectMeta.Name, ready)
	}
	return len(nodes.Items), nil
}

// WaitForNodes waits till the nodes are ready
func (c *ClusterProvider) WaitForNodes(clientSet *clientset.Clientset, ng *api.NodeGroup) error {
	if ng.MinSize == nil || *ng.MinSize == 0 {
		return nil
	}
	timer := time.After(c.Provider.WaitTimeout())
	timeout := false
	watcher, err := clientSet.CoreV1().Nodes().Watch(ng.ListOptions())
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}

	counter, err := getNodes(clientSet, ng)
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}

	logger.Info("waiting for at least %d node(s) to become ready in %q", *ng.MinSize, ng.Name)
	for !timeout && counter <= *ng.MinSize {
		select {
		case event := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if event.Object != nil && event.Type != watch.Deleted {
				if node, ok := event.Object.(*corev1.Node); ok {
					if isNodeReady(node) {
						counter++
						logger.Debug("node %q is ready in %q", node.ObjectMeta.Name, ng.Name)
					} else {
						logger.Debug("node %q seen in %q, but not ready yet", node.ObjectMeta.Name, ng.Name)
						logger.Debug("node = %#v", *node)
					}
				}
			}
		case <-timer:
			timeout = true
		}
	}
	if timeout {
		return fmt.Errorf("timed out (after %s) waitiing for at least %d nodes to join the cluster and become ready in %q", c.Provider.WaitTimeout(), *ng.MinSize, ng.Name)
	}

	if _, err = getNodes(clientSet, ng); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}
