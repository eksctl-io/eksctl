package eks

import (
	"fmt"
	"time"

	"github.com/ghodss/yaml"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
)

func (c *ClusterProvider) newNodeAuthConfigMap(ng *api.NodeGroup) (*corev1.ConfigMap, error) {
	mapRoles := make([]map[string]interface{}, 1)
	mapRoles[0] = make(map[string]interface{})

	mapRoles[0]["rolearn"] = ng.InstanceRoleARN
	mapRoles[0]["username"] = "system:node:{{EC2PrivateDNSName}}"
	mapRoles[0]["groups"] = []string{
		"system:bootstrappers",
		"system:nodes",
	}

	mapRolesBytes, err := yaml.Marshal(mapRoles)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-auth",
			Namespace: "kube-system",
		},
		Data: map[string]string{
			"mapRoles": string(mapRolesBytes),
		},
	}

	return cm, nil
}

// CreateNodeGroupAuthConfigMap creates the auth config map for the default node group
func (c *ClusterProvider) CreateNodeGroupAuthConfigMap(clientSet *clientset.Clientset, ng *api.NodeGroup) error {
	cm, err := c.newNodeAuthConfigMap(ng)
	if err != nil {
		return errors.Wrap(err, "constructing auth ConfigMap for DefaultNodeGroup")
	}
	if _, err := clientSet.CoreV1().ConfigMaps("kube-system").Create(cm); err != nil {
		return errors.Wrap(err, "creating auth ConfigMap for DefaultNodeGroup")
	}
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

func getNodes(clientSet *clientset.Clientset) (int, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	logger.Info("the cluster has %d nodes", len(nodes.Items))
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
	if ng.MinSize == 0 {
		return nil
	}
	timer := time.After(c.Provider.WaitTimeout())
	timeout := false
	watcher, err := clientSet.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}

	counter, err := getNodes(clientSet)
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}

	logger.Info("waiting for at least %d nodes to become ready", ng.MinSize)
	for !timeout && counter <= ng.MinSize {
		select {
		case event := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if event.Object != nil && event.Type != watch.Deleted {
				if node, ok := event.Object.(*corev1.Node); ok {
					if isNodeReady(node) {
						counter++
						logger.Debug("node %q is ready", node.ObjectMeta.Name)
					} else {
						logger.Debug("node %q seen, but not ready yet", node.ObjectMeta.Name)
						logger.Debug("node = %#v", *node)
					}
				}
			}
		case <-timer:
			timeout = true
		}
	}
	if timeout {
		return fmt.Errorf("timed out (after %s) waitiing for at least %d nodes to join the cluster and become ready", c.Provider.WaitTimeout(), ng.MinSize)
	}

	if _, err = getNodes(clientSet); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}
