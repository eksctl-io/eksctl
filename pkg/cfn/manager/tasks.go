package manager

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
)

type createClusterTask struct {
	info                 string
	stackCollection      *StackCollection
	supportsManagedNodes bool
}

func (t *createClusterTask) Describe() string { return t.info }

func (t *createClusterTask) Do(errorCh chan error) error {
	return t.stackCollection.createClusterTask(errorCh, t.supportsManagedNodes)
}

type nodeGroupTask struct {
	info                 string
	nodeGroup            *api.NodeGroup
	supportsManagedNodes bool
	forceAddCNIPolicy    bool
	stackCollection      *StackCollection
}

func (t *nodeGroupTask) Describe() string { return t.info }
func (t *nodeGroupTask) Do(errs chan error) error {
	return t.stackCollection.createNodeGroupTask(errs, t.nodeGroup, t.supportsManagedNodes, t.forceAddCNIPolicy)
}

type managedNodeGroupTask struct {
	info              string
	nodeGroup         *api.ManagedNodeGroup
	stackCollection   *StackCollection
	forceAddCNIPolicy bool
}

func (t *managedNodeGroupTask) Describe() string { return t.info }

func (t *managedNodeGroupTask) Do(errorCh chan error) error {
	return t.stackCollection.createManagedNodeGroupTask(errorCh, t.nodeGroup, t.forceAddCNIPolicy)
}

type clusterCompatTask struct {
	info            string
	stackCollection *StackCollection
}

func (t *clusterCompatTask) Describe() string { return t.info }

func (t *clusterCompatTask) Do(errorCh chan error) error {
	defer close(errorCh)
	return t.stackCollection.FixClusterCompatibility()
}

type taskWithClusterIAMServiceAccountSpec struct {
	info            string
	stackCollection *StackCollection
	serviceAccount  *api.ClusterIAMServiceAccount
	oidc            *iamoidc.OpenIDConnectManager
}

func (t *taskWithClusterIAMServiceAccountSpec) Describe() string { return t.info }
func (t *taskWithClusterIAMServiceAccountSpec) Do(errs chan error) error {
	return t.stackCollection.createIAMServiceAccountTask(errs, t.serviceAccount, t.oidc)
}

type taskWithStackSpec struct {
	info  string
	stack *Stack
	call  func(*Stack, chan error) error
}

func (t *taskWithStackSpec) Describe() string { return t.info }
func (t *taskWithStackSpec) Do(errs chan error) error {
	return t.call(t.stack, errs)
}

type asyncTaskWithStackSpec struct {
	info  string
	stack *Stack
	call  func(*Stack) (*Stack, error)
}

func (t *asyncTaskWithStackSpec) Describe() string { return t.info + " [async]" }
func (t *asyncTaskWithStackSpec) Do(errs chan error) error {
	_, err := t.call(t.stack)
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
	call       func(kubernetes.Interface) error
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
	err = t.call(clientSet)
	close(errs)
	return err
}
