package irsa

import (
	"context"
	"fmt"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// StackManager manages CloudFormation stacks for IAM Service Accounts.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	CreateStack(ctx context.Context, name string, stack builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
	DeleteStackBySpec(ctx context.Context, s *cfntypes.Stack) (*cfntypes.Stack, error)
	DeleteStackBySpecSync(ctx context.Context, s *cfntypes.Stack, errs chan error) error
	DescribeIAMServiceAccountStacks(ctx context.Context) ([]*cfntypes.Stack, error)
	GetIAMServiceAccounts(ctx context.Context) ([]*api.ClusterIAMServiceAccount, error)
	GetStackTemplate(ctx context.Context, stackName string) (string, error)
	UpdateStack(ctx context.Context, options manager.UpdateStackOptions) error
}

type Manager struct {
	clusterName  string
	oidcManager  *iamoidc.OpenIDConnectManager
	stackManager StackManager
	clientSet    kubeclient.Interface
}

type action string

const (
	actionCreate action = "create"
	actionDelete action = "delete"
	actionUpdate action = "update"
)

func New(clusterName string, stackManager StackManager, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) *Manager {
	return &Manager{
		clusterName:  clusterName,
		oidcManager:  oidcManager,
		stackManager: stackManager,
		clientSet:    clientSet,
	}
}

func doTasks(taskTree *tasks.TaskTree, action action) error {
	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		logger.Info("%d error(s) occurred and IAM Role stacks haven't been %sd properly, you may wish to check CloudFormation console", len(errs), action)
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to %s iamserviceaccount(s)", action)
	}
	return nil
}

// logPlanModeWarning will log a message to inform user that they are in plan-mode
func logPlanModeWarning(plan bool) {
	if plan {
		logger.Warning("no changes were applied, run again with '--approve' to apply the changes")
	}
}
