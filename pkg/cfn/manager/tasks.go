package manager

import (
	"context"
	"fmt"

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
	ctx             context.Context
}

func (t *managedNodeGroupTagsToASGPropagationTask) Describe() string { return t.info }

func (t *managedNodeGroupTagsToASGPropagationTask) Do(errorCh chan error) error {
	return t.stackCollection.propagateManagedNodeGroupTagsToASGTask(t.ctx, errorCh, t.nodeGroup, t.stackCollection.PropagateManagedNodeGroupTagsToASG)
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

type deleteStackTask struct {
	info            string
	stack           *Stack
	stackCollection *StackCollection
	ctx             context.Context
	wait            bool
}

func (t *deleteStackTask) Describe() string {
	if t.wait {
		return t.info
	}
	return t.info + " [async]"
}
func (t *deleteStackTask) Do(errs chan error) error {
	options := DeleteStackOptions{
		Stack: t.stack,
	}
	if t.wait {
		options.ErrCh = errs
	} else {
		defer close(errs)
	}
	return t.stackCollection.DeleteStack(t.ctx, options)
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
