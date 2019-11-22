package eks

import (
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/utils"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func isNodeReady(node *corev1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func getNodes(clientSet kubernetes.Interface, ng KubeNodeGroup) (int, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(ng.ListOptions())
	if err != nil {
		return 0, err
	}
	logger.Info("nodegroup %q has %d node(s)", ng.NameString(), len(nodes.Items))
	counter := 0
	for _, node := range nodes.Items {
		// logger.Debug("node[%d]=%#v", n, node)
		ready := "not ready"
		if isNodeReady(&node) {
			ready = "ready"
			counter++
		}
		logger.Info("node %q is %s", node.ObjectMeta.Name, ready)
	}
	return counter, nil
}

// ValidateFeatureCompatibility validates whether the cluster version supports the features specified in the
// ClusterConfig. Support for Managed Nodegroups or Windows requires the EKS cluster version to be 1.14 and above.
// If the version requirement isn't met, an error is returned
func ValidateFeatureCompatibility(clusterConfig *api.ClusterConfig, kubeNodeGroups []KubeNodeGroup) error {
	if err := ValidateManagedNodesSupport(clusterConfig); err != nil {
		return err
	}
	return ValidateWindowsCompatibility(kubeNodeGroups, clusterConfig.Metadata.Version)
}

// ValidateManagedNodesSupport validates support for Managed Nodegroups
func ValidateManagedNodesSupport(clusterConfig *api.ClusterConfig) error {
	if len(clusterConfig.ManagedNodeGroups) > 0 {
		minRequiredVersion := api.Version1_14
		supportsManagedNodes, err := VersionSupportsManagedNodes(clusterConfig.Metadata.Version)
		if err != nil {
			return err
		}
		if !supportsManagedNodes {
			return fmt.Errorf("Managed Nodegroups are only supported on EKS version %s and above", minRequiredVersion)
		}
	}
	return nil
}

// VersionSupportsManagedNodes reports whether the control plane version can support Managed Nodes
func VersionSupportsManagedNodes(controlPlaneVersion string) (bool, error) {
	minRequiredVersion := api.Version1_14
	supportsManagedNodes, err := utils.IsMinVersion(minRequiredVersion, controlPlaneVersion)
	if err != nil {
		return false, err
	}
	return supportsManagedNodes, nil
}

// ValidateWindowsCompatibility validates Windows compatibility
func ValidateWindowsCompatibility(kubeNodeGroups []KubeNodeGroup, controlPlaneVersion string) error {
	if !hasWindowsNode(kubeNodeGroups) {
		return nil
	}

	supportsWindows, err := utils.IsMinVersion(api.Version1_14, controlPlaneVersion)
	if err != nil {
		return err
	}
	if !supportsWindows {
		return errors.New("Windows nodes are only supported on Kubernetes 1.14 and above")
	}
	return nil
}

// SupportsWindowsWorkloads reports whether nodeGroups can support running Windows workloads
func SupportsWindowsWorkloads(nodeGroups []KubeNodeGroup) bool {
	return hasWindowsNode(nodeGroups) && hasAmazonLinux2Node(nodeGroups)
}

// hasWindowsNode reports whether there's at least one Windows node in nodeGroups
func hasWindowsNode(nodeGroups []KubeNodeGroup) bool {
	for _, ng := range nodeGroups {
		if api.IsWindowsImage(ng.GetAMIFamily()) {
			return true
		}
	}
	return false
}

// hasAmazonLinux2Node reports whether there's at least one Windows node in nodeGroups
func hasAmazonLinux2Node(nodeGroups []KubeNodeGroup) bool {
	for _, ng := range nodeGroups {
		if ng.GetAMIFamily() == api.NodeImageFamilyAmazonLinux2 {
			return true
		}
	}
	return false
}

// LogWindowsCompatibility logs Windows compatibility messages
func LogWindowsCompatibility(nodeGroups []KubeNodeGroup, clusterMeta *api.ClusterMeta) {
	if hasWindowsNode(nodeGroups) {
		if !hasAmazonLinux2Node(nodeGroups) {
			logger.Warning("a Linux node group is required to support Windows workloads")
			logger.Warning("add it using 'eksctl create nodegroup --cluster=%s --node-ami-family=%s'", clusterMeta.Name, api.NodeImageFamilyAmazonLinux2)
		}
		logger.Warning("Windows VPC resource controller is required to run Windows workloads")
		logger.Warning("install it using 'eksctl utils install-vpc-controllers --name=%s --region=%s --approve'", clusterMeta.Name, clusterMeta.Region)
	}
}

// KubeNodeGroup defines a set of Kubernetes Nodes
type KubeNodeGroup interface {
	// NameString returns the name
	NameString() string
	// Size returns the number of the nodes (desired capacity)
	Size() int
	// ListOptions returns the selector for listing nodes in this nodegroup
	ListOptions() metav1.ListOptions
	// GetAMIFamily returns the AMI family
	GetAMIFamily() string
}

// WaitForNodes waits till the nodes are ready
func (c *ClusterProvider) WaitForNodes(clientSet kubernetes.Interface, ng KubeNodeGroup) error {
	minSize := ng.Size()
	if minSize == 0 {
		return nil
	}
	timer := time.After(c.Provider.WaitTimeout())
	timeout := false
	readyNodes := sets.NewString()
	watcher, err := clientSet.CoreV1().Nodes().Watch(ng.ListOptions())
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}

	counter, err := getNodes(clientSet, ng)
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}

	logger.Info("waiting for at least %d node(s) to become ready in %q", minSize, ng.NameString())
	for !timeout && counter < minSize {
		select {
		case event := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if event.Object != nil && event.Type != watch.Deleted {
				if node, ok := event.Object.(*corev1.Node); ok {
					if isNodeReady(node) {
						readyNodes.Insert(node.Name)
						counter = readyNodes.Len()
						logger.Debug("node %q is ready in %q", node.Name, ng.NameString())
					} else {
						logger.Debug("node %q seen in %q, but not ready yet", node.Name, ng.NameString())
						logger.Debug("node = %#v", *node)
					}
				}
			}
		case <-timer:
			timeout = true
		}
	}
	watcher.Stop()
	if timeout {
		return fmt.Errorf("timed out (after %s) waiting for at least %d nodes to join the cluster and become ready in %q", c.Provider.WaitTimeout(), minSize, ng.NameString())
	}

	if _, err = getNodes(clientSet, ng); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}

// GetNodeGroupIAM retrieves the IAM configuration of the given nodegroup
func (c *ClusterProvider) GetNodeGroupIAM(stackManager *manager.StackCollection, spec *api.ClusterConfig, ng *api.NodeGroup) error {
	stacks, err := stackManager.DescribeNodeGroupStacks()
	if err != nil {
		return err
	}

	for _, s := range stacks {
		if stackManager.GetNodeGroupName(s) == ng.Name {
			if !stackManager.StackStatusIsNotTransitional(s) {
				return fmt.Errorf("nodegroup %q is in transitional state (%q)", ng.Name, *s.StackStatus)
			}
			return iam.UseFromNodeGroup(c.Provider, s, ng)
		}
	}

	return fmt.Errorf("stack not found for nodegroup %q", ng.Name)
}
