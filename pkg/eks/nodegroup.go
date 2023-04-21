package eks

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	addons "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/iam"
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

func GetNodes(clientSet kubernetes.Interface, ng KubeNodeGroup) (int, error) {
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
//
//go:generate "${GOBIN}/mockery" --name=KubeNodeGroup --output=mocks/
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
func (c *ClusterProvider) GetNodeGroupIAM(ctx context.Context, stackManager manager.StackManager, ng *api.NodeGroup) error {
	stacks, err := stackManager.ListNodeGroupStacks(ctx)
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
					ng.Name, s.StackStatus,
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
func DoesAWSNodeUseIRSA(ctx context.Context, provider api.ClusterProvider, clientSet kubernetes.Interface) (bool, error) {
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
	input := &awsiam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(arnParts[len(arnParts)-1]),
	}
	policies, err := provider.IAM().ListAttachedRolePolicies(ctx, input)
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
	asg             awsapi.ASG
	ctx             context.Context
	nodegroup       *api.NodeGroupBase
	stackCollection manager.StackManager
}

func (t *suspendProcesses) Describe() string {
	return fmt.Sprintf("suspend ASG processes for nodegroup %s", t.nodegroup.Name)
}

func (t *suspendProcesses) Do() error {
	ngStack, err := t.stackCollection.DescribeNodeGroupStack(context.TODO(), t.nodegroup.Name)
	if err != nil {
		return errors.Wrapf(err, "couldn't describe nodegroup stack for nodegroup %s", t.nodegroup.Name)
	}
	asgName, err := t.stackCollection.GetAutoScalingGroupName(context.TODO(), ngStack)
	if err != nil {
		return errors.Wrapf(err, "couldn't get autoscalinggroup name nodegroup %s", t.nodegroup.Name)
	}
	_, err = t.asg.SuspendProcesses(t.ctx, &autoscaling.SuspendProcessesInput{
		AutoScalingGroupName: aws.String(asgName),
		ScalingProcesses:     t.nodegroup.ASGSuspendProcesses,
	})
	logger.Info("suspended ASG processes %v for %s", t.nodegroup.ASGSuspendProcesses, t.nodegroup.Name)
	return err
}

// newSuspendProcesses returns a task that suspends the given processes for this
// AutoScalingGroup
func newSuspendProcesses(c *ClusterProvider, spec *api.ClusterConfig, nodegroup *api.NodeGroupBase) tasks.Task {
	return tasks.SynchronousTask{
		SynchronousTaskIface: &suspendProcesses{
			ctx:             context.Background(),
			asg:             c.AWSProvider.ASG(),
			stackCollection: c.NewStackManager(spec),
			nodegroup:       nodegroup,
		},
	}
}
