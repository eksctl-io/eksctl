package podidentityassociation

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// A StackLister lists and describes CloudFormation stacks.
type StackLister interface {
	ListPodIdentityStackNames(ctx context.Context) ([]string, error)
	DescribeStack(ctx context.Context, stack *manager.Stack) (*manager.Stack, error)
	GetStackTemplate(ctx context.Context, stackName string) (string, error)
	GetIAMServiceAccounts(ctx context.Context) ([]*api.ClusterIAMServiceAccount, error)
}

// A StackDeleter lists and deletes CloudFormation stacks.
type StackDeleter interface {
	StackLister
	DeleteStackBySpecSync(ctx context.Context, stack *cfntypes.Stack, errCh chan error) error
}

// APILister lists pod identity associations using the EKS API.
type APILister interface {
	ListPodIdentityAssociations(ctx context.Context, params *eks.ListPodIdentityAssociationsInput, optFns ...func(*eks.Options)) (*eks.ListPodIdentityAssociationsOutput, error)
}

// APIDeleter lists and deletes pod identity associations using the EKS API.
type APIDeleter interface {
	APILister
	DeletePodIdentityAssociation(ctx context.Context, params *eks.DeletePodIdentityAssociationInput, optFns ...func(*eks.Options)) (*eks.DeletePodIdentityAssociationOutput, error)
}

// A Deleter deletes pod identity associations.
type Deleter struct {
	// ClusterName is the cluster name.
	ClusterName string
	// StackDeleter is used to delete stacks.
	StackDeleter StackDeleter
	// APIDeleter deletes pod identity associations using the EKS API.
	APIDeleter APIDeleter
	// ClientSet is used to delete K8s service accounts.
	ClientSet kubeclient.Interface
}

// Identifier represents a pod identity association.
type Identifier struct {
	// Namespace is the namespace the service account belongs to.
	Namespace string
	// ServiceAccountName is the name of the Kubernetes ServiceAccount.
	ServiceAccountName string
}

func (i Identifier) IDString() string {
	return i.toString("/")
}

func (i Identifier) NameString() string {
	return i.toString("-")
}

func (i Identifier) toString(delimiter string) string {
	return i.Namespace + delimiter + i.ServiceAccountName
}

func NewDeleter(clusterName string, stackDeleter StackDeleter, apiDeleter APIDeleter, clientSet kubeclient.Interface) *Deleter {
	return &Deleter{
		ClusterName:  clusterName,
		StackDeleter: stackDeleter,
		APIDeleter:   apiDeleter,
		ClientSet:    clientSet,
	}
}

// Delete deletes the specified podIdentityAssociations.
func (d *Deleter) Delete(ctx context.Context, podIDs []Identifier) error {
	tasks, err := d.DeleteTasks(ctx, podIDs)
	if err != nil {
		return err
	}
	return runAllTasks(tasks)
}

func (d *Deleter) DeleteTasks(ctx context.Context, podIDs []Identifier) (*tasks.TaskTree, error) {
	roleStackNames, err := d.StackDeleter.ListPodIdentityStackNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing stack names for pod identity associations: %w", err)
	}
	taskTree := &tasks.TaskTree{Parallel: true}

	// this is true during cluster deletion, when no association identifier is given as user input,
	// instead we will delete all pod-identity-role stacks for the cluster
	if len(podIDs) == 0 {
		for _, stackName := range roleStackNames {
			name := strings.Clone(stackName)
			taskTree.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("deleting IAM resources stack %q", stackName),
				Doer: func() error {
					return d.deleteRoleStack(ctx, name)
				},
			})
		}
		return taskTree, nil
	}

	for _, podID := range podIDs {
		podID := podID
		piaDeletionTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		deleteTask, err := d.makeDeleteTask(ctx, podID, roleStackNames)
		if err != nil {
			return nil, err
		}
		piaDeletionTasks.Append(deleteTask)
		piaDeletionTasks.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("delete service account %q, if it exists and is managed by eksctl", podID.IDString()),
			Doer: func() error {
				if err := kubernetes.MaybeDeleteServiceAccount(d.ClientSet, v1.ObjectMeta{
					Name:      podID.ServiceAccountName,
					Namespace: podID.Namespace,
				}); err != nil {
					return fmt.Errorf("failed to delete service account %q: %w", podID.IDString(), err)
				}
				return nil
			},
		})
		taskTree.Append(piaDeletionTasks)
	}
	return taskTree, nil
}

func (d *Deleter) makeDeleteTask(ctx context.Context, p Identifier, roleStackNames []string) (tasks.Task, error) {
	podIdentityAssociationID := p.IDString()
	output, err := d.APIDeleter.ListPodIdentityAssociations(ctx, &eks.ListPodIdentityAssociationsInput{
		ClusterName:    aws.String(d.ClusterName),
		Namespace:      aws.String(p.Namespace),
		ServiceAccount: aws.String(p.ServiceAccountName),
	})
	if err != nil {
		return nil, fmt.Errorf("listing pod identity associations: %w", err)
	}

	switch len(output.Associations) {
	default:
		return nil, fmt.Errorf("expected to find only 1 pod identity association for %q; got %d", podIdentityAssociationID, len(output.Associations))
	case 0:
		logger.Warning("pod identity association %q not found", podIdentityAssociationID)
	case 1:
		association := output.Associations[0]
		if association.OwnerArn != nil {
			return nil, fmt.Errorf("cannot delete podidentityassociation %s as it is in use by addon %s; "+
				"please use `eksctl update addon` or `eksctl delete addon` instead", p.IDString(), *association.OwnerArn)
		}
	}
	return &tasks.GenericTask{
		Description: fmt.Sprintf("delete pod identity association %q", podIdentityAssociationID),
		Doer: func() error {
			if len(output.Associations) == 1 {
				if _, err := d.APIDeleter.DeletePodIdentityAssociation(ctx, &eks.DeletePodIdentityAssociationInput{
					ClusterName:   aws.String(d.ClusterName),
					AssociationId: output.Associations[0].AssociationId,
				}); err != nil {
					return fmt.Errorf("deleting pod identity association: %w", err)
				}
			}
			if err := d.deleteIAMResources(ctx, p, roleStackNames); err != nil {
				return fmt.Errorf("error deleting pod identity association %q: %w", podIdentityAssociationID, err)
			}
			return nil
		},
	}, nil
}

func (d *Deleter) deleteIAMResources(ctx context.Context, p Identifier, roleStackNames []string) error {
	stackName, hasStack := getIAMResourcesStack(roleStackNames, p)
	if !hasStack {
		return nil
	}
	logger.Info("deleting IAM resources stack %q for pod identity association %q", stackName, p.IDString())
	return d.deleteRoleStack(ctx, stackName)
}

func (d *Deleter) deleteRoleStack(ctx context.Context, stackName string) error {
	stack, err := d.StackDeleter.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return fmt.Errorf("describing stack %q: %w", stackName, err)
	}

	deleteStackCh := make(chan error)
	if err := d.StackDeleter.DeleteStackBySpecSync(ctx, stack, deleteStackCh); err != nil {
		return fmt.Errorf("deleting stack %q for IAM role: %w", stackName, err)
	}
	select {
	case err := <-deleteStackCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timed out waiting for deletion of pod identity association: %w", ctx.Err())
	}
}

// ToIdentifiers maps a list of PodIdentityAssociations to a list of Identifiers.
func ToIdentifiers(podIdentityAssociations []api.PodIdentityAssociation) []Identifier {
	identifiers := make([]Identifier, len(podIdentityAssociations))
	for i, p := range podIdentityAssociations {
		identifiers[i] = Identifier{
			Namespace:          p.Namespace,
			ServiceAccountName: p.ServiceAccountName,
		}
	}
	return identifiers
}

func getIAMResourcesStack(stackNames []string, p Identifier) (string, bool) {
	for _, name := range stackNames {
		if strings.Contains(name, p.NameString()) {
			return name, true
		}
	}
	return "", false
}
