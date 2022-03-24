package manager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/weaveworks/eksctl/pkg/awsapi"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type createClusterTask struct {
	info                 string
	stackCollection      *StackCollection
	supportsManagedNodes bool
	ctx                  context.Context
}

func (t *createClusterTask) Describe() string { return t.info }

func (t *createClusterTask) Do(errorCh chan error) error {
	return t.stackCollection.createClusterTask(t.ctx, errorCh, t.supportsManagedNodes)
}

type nodeGroupTask struct {
	info              string
	ctx               context.Context
	nodeGroup         *api.NodeGroup
	forceAddCNIPolicy bool
	vpcImporter       vpc.Importer
	stackCollection   *StackCollection
}

func (t *nodeGroupTask) Describe() string { return t.info }
func (t *nodeGroupTask) Do(errs chan error) error {
	return t.stackCollection.createNodeGroupTask(t.ctx, errs, t.nodeGroup, t.forceAddCNIPolicy, t.vpcImporter)
}

type managedNodeGroupTask struct {
	info              string
	nodeGroup         *api.ManagedNodeGroup
	stackCollection   *StackCollection
	forceAddCNIPolicy bool
	vpcImporter       vpc.Importer
	ctx               context.Context
}

func (t *managedNodeGroupTask) Describe() string { return t.info }

func (t *managedNodeGroupTask) Do(errorCh chan error) error {
	return t.stackCollection.createManagedNodeGroupTask(t.ctx, errorCh, t.nodeGroup, t.forceAddCNIPolicy, t.vpcImporter)
}

type managedNodeGroupTagsToASGPropagationTask struct {
	info            string
	nodeGroup       *api.ManagedNodeGroup
	stackCollection *StackCollection
}

func (t *managedNodeGroupTagsToASGPropagationTask) Describe() string { return t.info }

func (t *managedNodeGroupTagsToASGPropagationTask) Do(errorCh chan error) error {
	return t.stackCollection.propagateManagedNodeGroupTagsToASGTask(errorCh, t.nodeGroup)
}

type clusterCompatTask struct {
	info            string
	stackCollection *StackCollection
	ctx             context.Context
}

func (t *clusterCompatTask) Describe() string { return t.info }

func (t *clusterCompatTask) Do(errorCh chan error) error {
	defer close(errorCh)
	return t.stackCollection.FixClusterCompatibility(t.ctx)
}

type taskWithClusterIAMServiceAccountSpec struct {
	info            string
	stackCollection *StackCollection
	serviceAccount  *api.ClusterIAMServiceAccount
	oidc            *iamoidc.OpenIDConnectManager
}

func (t *taskWithClusterIAMServiceAccountSpec) Describe() string { return t.info }
func (t *taskWithClusterIAMServiceAccountSpec) Do(errs chan error) error {
	return t.stackCollection.createIAMServiceAccountTask(context.TODO(), errs, t.serviceAccount, t.oidc)
}

type taskWithStackSpec struct {
	info  string
	stack *Stack
	call  func(context.Context, *Stack, chan error) error
}

func (t *taskWithStackSpec) Describe() string { return t.info }
func (t *taskWithStackSpec) Do(errs chan error) error {
	return t.call(context.TODO(), t.stack, errs)
}

type asyncTaskWithStackSpec struct {
	info  string
	stack *Stack
	call  func(context.Context, *Stack) (*Stack, error)
}

func (t *asyncTaskWithStackSpec) Describe() string { return t.info + " [async]" }
func (t *asyncTaskWithStackSpec) Do(errs chan error) error {
	_, err := t.call(context.TODO(), t.stack)
	close(errs)
	return err
}

type asyncTaskWithoutParams struct {
	info string
	call func() error
}

func (t *asyncTaskWithoutParams) Describe() string { return t.info }
func (t *asyncTaskWithoutParams) Do(errs chan error) error {
	err := t.call()
	close(errs)
	return err
}

type kubernetesTask struct {
	info       string
	kubernetes kubewrapper.ClientSetGetter
	objectMeta v1.ObjectMeta
	call       func(kubernetes.Interface, v1.ObjectMeta) error
}

func (t *kubernetesTask) Describe() string { return t.info }
func (t *kubernetesTask) Do(errs chan error) error {
	if t.kubernetes == nil {
		return fmt.Errorf("cannot start task %q as Kubernetes client configurtaion wasn't provided", t.Describe())
	}
	clientSet, err := t.kubernetes.ClientSet()
	if err != nil {
		return err
	}
	err = t.call(clientSet, t.objectMeta)
	close(errs)
	return err
}

type AssignIpv6AddressOnCreationTask struct {
	EC2API        awsapi.EC2
	Context       context.Context
	ClusterConfig *api.ClusterConfig
}

func (t *AssignIpv6AddressOnCreationTask) Describe() string {
	return "set AssignIpv6AddressOnCreation to true for public subnets"
}

func (t *AssignIpv6AddressOnCreationTask) Do(errs chan error) error {
	defer close(errs)
	if t.ClusterConfig.VPC.Subnets.Public != nil {
		for _, subnet := range t.ClusterConfig.VPC.Subnets.Public.WithIDs() {
			_, err := t.EC2API.ModifySubnetAttribute(t.Context, &ec2.ModifySubnetAttributeInput{
				AssignIpv6AddressOnCreation: &ec2types.AttributeBooleanValue{
					Value: aws.Bool(true),
				},
				SubnetId: aws.String(subnet),
			})
			if err != nil {
				return fmt.Errorf("failed to update public subnet %q: %v", subnet, err)
			}
		}
	}
	return nil
}
