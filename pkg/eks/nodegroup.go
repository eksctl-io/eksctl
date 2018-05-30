package eks

import (
	"fmt"
	"time"

	"github.com/ghodss/yaml"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
)

func (c *Config) newNodeAuthConfigMap() (*corev1.ConfigMap, error) {
	mapRoles := make([]map[string]interface{}, 1)
	mapRoles[0] = make(map[string]interface{})

	mapRoles[0]["rolearn"] = c.nodeInstanceRoleARN
	mapRoles[0]["username"] = "system:node:{{EC2PrivateDNSName}}"
	mapRoles[0]["groups"] = []string{
		"system:bootstrappers",
		"system:nodes",
		"system:nodes",
	}

	mapRolesBytes, err := yaml.Marshal(mapRoles)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-auth",
			Namespace: "default",
		},
		BinaryData: map[string][]byte{
			"mapRoles": mapRolesBytes,
		},
	}

	return cm, nil
}

func (c *Config) CreateDefaultNodeGroupAuthConfigMap(clientSet *clientset.Clientset) error {
	cm, err := c.newNodeAuthConfigMap()
	if err != nil {
		return errors.Wrap(err, "contructing auth ConfigMap for DefaultNodeGroup")
	}
	if _, err := clientSet.CoreV1().ConfigMaps("default").Create(cm); err != nil {
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
	nodes, err := clientSet.Core().Nodes().List(metav1.ListOptions{})
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

func (c *Config) WaitForNodes(clientSet *clientset.Clientset) error {
	timeoutAfter := 20 * time.Minute
	timer := time.After(timeoutAfter)
	timeout := false
	watcher, err := clientSet.Core().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}

	counter, err := getNodes(clientSet)
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}

	logger.Info("waiting for at least %d nodes to become ready", c.MinNodes)
	for !timeout && counter <= c.MinNodes {
		select {
		case event, _ := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if event.Type == watch.Added {
				node := event.Object.(*corev1.Node)
				if isNodeReady(node) {
					counter++
					logger.Debug("node %q is ready", node.ObjectMeta.Name)
				} else {
					logger.Debug("node %q seen, but not ready yet", node.ObjectMeta.Name)
					logger.Debug("node = %#v", *node)
				}
			}
		case <-timer:
			timeout = true
		}
	}
	if timeout {
		return fmt.Errorf("timed out (after %s) waitiing for at least %d nodes to join the cluster and become ready", timeoutAfter, c.MinNodes)
	}

	if _, err = getNodes(clientSet); err != nil {
		errors.Wrap(err, "re-listing nodes")
	}

	return nil
}
