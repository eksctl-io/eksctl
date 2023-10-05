package manager

import (
	"context"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
	skipEgressRules   bool
}

func (t *nodeGroupTask) Describe() string { return t.info }
func (t *nodeGroupTask) Do(errs chan error) error {
	return t.stackCollection.createNodeGroupTask(t.ctx, errs, t.nodeGroup, t.forceAddCNIPolicy, t.skipEgressRules, t.vpcImporter)
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
