package eks

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	addons "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	nodes, err := clientSet.CoreV1().Nodes().List(context.TODO(), ng.ListOptions())
	if err != nil {
		return 0, err
	}
	logger.Info("nodegroup %q has %d node(s)", ng.NameString(), len(nodes.Items))
	counter := 0
	for _, node := range nodes.Items {
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
// Bottlerocket nodegroups are only supported on EKS version 1.15 and above
// If the version requirement isn't met, an error is returned
func ValidateFeatureCompatibility(clusterConfig *api.ClusterConfig, kubeNodeGroups []KubeNodeGroup) error {
	if err := ValidateKMSSupport(clusterConfig, clusterConfig.Metadata.Version); err != nil {
		return err
	}
	if err := ValidateManagedNodesSupport(clusterConfig); err != nil {
		return err
	}
	if err := ValidateBottlerocketSupport(clusterConfig.Metadata.Version, kubeNodeGroups); err != nil {
		return err
	}

	return ValidateWindowsCompatibility(kubeNodeGroups, clusterConfig.Metadata.Version)
}

// ValidateBottlerocketSupport validates support for Bottlerocket nodegroups
func ValidateBottlerocketSupport(controlPlaneVersion string, kubeNodeGroups []KubeNodeGroup) error {
	const minSupportedVersion = api.Version1_15

	supportsBottlerocket, err := utils.IsMinVersion(minSupportedVersion, controlPlaneVersion)
	if err != nil {
		return err
	}
	if supportsBottlerocket {
		return nil
	}

	for _, ng := range kubeNodeGroups {
		if ng.GetAMIFamily() == api.NodeImageFamilyBottlerocket {
			return errors.Errorf("Bottlerocket is only supported on EKS version %s and above", minSupportedVersion)
		}
	}
	return nil
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

// ValidateKMSSupport validates support for KMS encryption
func ValidateKMSSupport(clusterConfig *api.ClusterConfig, eksVersion string) error {
	if clusterConfig.SecretsEncryption == nil {
		return nil
	}

	const minReqVersion = api.Version1_13
	supportsKMS, err := utils.IsMinVersion(minReqVersion, eksVersion)
	if err != nil {
		return errors.Wrap(err, "error validating KMS support")
	}
	if !supportsKMS {
		return fmt.Errorf("secrets encryption with KMS is only supported for EKS version %s and above", minReqVersion)
	}

	if _, err := arn.Parse(clusterConfig.SecretsEncryption.KeyARN); err != nil {
		return errors.Wrapf(err, "invalid ARN in secretsEncryption.keyARN: %q", clusterConfig.SecretsEncryption.KeyARN)
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

//go:generate "${GOBIN}/mockery" --name=KubeNodeGroup --output=mocks/
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

// GetNodeGroupIAM retrieves the IAM configuration of the given nodegroup
func (c *ClusterProvider) GetNodeGroupIAM(stackManager manager.StackManager, ng *api.NodeGroup) error {
	stacks, err := stackManager.DescribeNodeGroupStacks()
	if err != nil {
		return err
	}

	for _, s := range stacks {
		if stackManager.GetNodeGroupName(s) == ng.Name {
			err := iam.UseFromNodeGroup(s, ng)
			// An empty InstanceRoleARN likely also points to an error
			if err == nil && ng.IAM.InstanceRoleARN == "" {
				err = errors.New("InstanceRoleARN empty")
			}
			if err != nil {
				return errors.Wrapf(
					err, "couldn't get iam configuration for nodegroup %q (perhaps state %q is transitional)",
					ng.Name, *s.StackStatus,
				)
			}
			return nil
		}
	}

	return fmt.Errorf("stack not found for nodegroup %q", ng.Name)
}

func getAWSNodeSAARNAnnotation(clientSet kubernetes.Interface) (string, error) {
	clusterDaemonSet, err := clientSet.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Get(context.TODO(), addons.AWSNode, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", addons.AWSNode)
			return "", nil
		}
		return "", errors.Wrapf(err, "getting %q", addons.AWSNode)
	}

	return clusterDaemonSet.Annotations[api.AnnotationEKSRoleARN], nil
}

// DoesAWSNodeUseIRSA evaluates whether an aws-node uses IRSA
func (n *NodeGroupService) DoesAWSNodeUseIRSA(provider api.ClusterProvider, clientSet kubernetes.Interface) (bool, error) {
	roleArn, err := getAWSNodeSAARNAnnotation(clientSet)
	if err != nil {
		return false, errors.Wrap(err, "error retrieving aws-node arn")
	}
	if roleArn == "" {
		return false, nil
	}
	arnParts := strings.Split(roleArn, "/")
	if len(arnParts) <= 1 {
		return false, errors.Errorf("invalid ARN %s", roleArn)
	}
	input := awsiam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(arnParts[len(arnParts)-1]),
	}
	policies, err := provider.IAM().ListAttachedRolePolicies(&input)
	if err != nil {
		return false, errors.Wrap(err, "error listing attached policies")
	}
	logger.Debug("found following policies attached to role annotated on aws-node service account: %s", policies.AttachedPolicies)
	for _, p := range policies.AttachedPolicies {
		if *p.PolicyName == api.IAMPolicyAmazonEKSCNIPolicy {
			return true, nil
		}
	}
	return false, nil
}

type suspendProcesses struct {
	asg             autoscalingiface.AutoScalingAPI
	nodegroup       *api.NodeGroupBase
	stackCollection manager.StackManager
}

func (t *suspendProcesses) Describe() string {
	return fmt.Sprintf("suspend ASG processes for nodegroup %s", t.nodegroup.Name)
}

func (t *suspendProcesses) Do() error {
	ngStack, err := t.stackCollection.DescribeNodeGroupStack(t.nodegroup.Name)
	if err != nil {
		return errors.Wrapf(err, "couldn't describe nodegroup stack for nodegroup %s", t.nodegroup.Name)
	}
	asgName, err := t.stackCollection.GetAutoScalingGroupName(ngStack)
	if err != nil {
		return errors.Wrapf(err, "couldn't get autoscalinggroup name nodegroup %s", t.nodegroup.Name)
	}
	_, err = t.asg.SuspendProcesses(&autoscaling.ScalingProcessQuery{
		AutoScalingGroupName: aws.String(asgName),
		ScalingProcesses:     aws.StringSlice(t.nodegroup.ASGSuspendProcesses),
	})
	logger.Info("suspended ASG processes %v for %s", t.nodegroup.ASGSuspendProcesses, t.nodegroup.Name)
	return err
}

// newSuspendProcesses returns a task that suspends the given processes for this
// AutoScalingGroup
func newSuspendProcesses(c *ClusterProvider, spec *api.ClusterConfig, nodegroup *api.NodeGroupBase) tasks.Task {
	return tasks.SynchronousTask{
		SynchronousTaskIface: &suspendProcesses{
			asg:             c.Provider.ASG(),
			stackCollection: c.NewStackManager(spec),
			nodegroup:       nodegroup,
		},
	}
}
