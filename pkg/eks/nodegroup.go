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

func (c *Config) WaitForNodes(clientSet *clientset.Clientset) error {
	timer := time.After(5 * time.Minute)
	timeout := false
	watcher, err := clientSet.Core().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}

	nodes, err := clientSet.Core().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}
	counter := len(nodes.Items)

	logger.Info("waiting for at least %d nodes to become ready", c.MinNodes)
	for !timeout && counter <= c.MinNodes {
		select {
		case event, _ := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if event.Type == watch.Added {
				// TODO(p1): check readiness
				counter++
			}
		case <-timer:
			timeout = true
		}
	}
	if timeout {
		// TODO(p2): doesn't have to be fatal
		return fmt.Errorf("timed out waitiing for nodes")
	}
	logger.Info("the cluster has %d nodes", counter, c.ClusterName)

	nodes, err = clientSet.Core().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}
	for n, node := range nodes.Items {
		logger.Debug("node[%n]=%#v", n, node)
	}

	return nil
}
