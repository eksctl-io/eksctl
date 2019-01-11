package eks

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
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

	utils.UpdateAuthConfigMap(cm, ng.IAM.InstanceRoleARN)
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
	if ng.MinSize == 0 {
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

	logger.Info("waiting for at least %d node(s) to become ready in %q", ng.MinSize, ng.Name)
	for !timeout && counter <= ng.MinSize {
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
		return fmt.Errorf("timed out (after %s) waitiing for at least %d nodes to join the cluster and become ready in %q", c.Provider.WaitTimeout(), ng.MinSize, ng.Name)
	}

	if _, err = getNodes(clientSet, ng); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}

// ValidateConfigForExistingNodeGroups looks at each of the existing nodegroups and
// validates configuration, if it find issues it logs messages
func (c *ClusterProvider) ValidateConfigForExistingNodeGroups(cfg *api.ClusterConfig) error {
	stackManager := c.NewStackManager(cfg)
	resourcesByNodeGroup, err := stackManager.DescribeResourcesOfNodeGroupStacks()
	if err != nil {
		return errors.Wrap(err, "getting resources for of all nodegroup stacks")
	}

	{
		securityGroupIDs := []string{}
		securityGroupIDsToNodeGroup := map[string]string{}
		for ng, resources := range resourcesByNodeGroup {
			for n := range resources.StackResources {
				r := resources.StackResources[n]
				if *r.ResourceType == "AWS::EC2::SecurityGroup" {
					securityGroupIDs = append(securityGroupIDs, *r.PhysicalResourceId)
					securityGroupIDsToNodeGroup[*r.PhysicalResourceId] = ng
				}
			}
		}

		input := &ec2.DescribeSecurityGroupsInput{
			GroupIds: aws.StringSlice(securityGroupIDs),
		}
		output, err := c.Provider.EC2().DescribeSecurityGroups(input)
		if err != nil {
			return errors.Wrap(err, "describing security groups")
		}

		for _, sg := range output.SecurityGroups {
			id := *sg.GroupId
			ng := securityGroupIDsToNodeGroup[id]
			logger.Debug("%s/%s = %#v", ng, id, sg)
			hasDNS := 0
			for _, p := range sg.IpPermissions {
				if p.FromPort != nil && *p.FromPort == 53 && p.ToPort != nil && *p.ToPort == 53 {
					if *p.IpProtocol == "tcp" || *p.IpProtocol == "udp" {
						// we cannot check p.IpRanges as we don't have VPC CIDR info when
						// we create the nodegroup, it may become important, but for now
						// it does't appear critical to check it
						hasDNS++
					}
				}
			}
			if hasDNS != 2 {
				logger.Critical("nodegroup %q may not have DNS port open to other nodegroups, so cluster DNS maybe be broken", ng)
				logger.Critical("it's recommended to delete the nodegroup %q and use new one instead")
				logger.Critical("check/update %q ingress rules - port 53 (TCP & UDP) has to be open for all sources inside the VPC", sg)
			}
		}
	}

	return nil
}
