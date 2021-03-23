package apply

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type Reconciler struct {
	cfg          *api.ClusterConfig
	plan         bool
	ctl          *eks.ClusterProvider
	irsaManager  IRSAManager
	stackManager manager.StackManager
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubernetes.Interface, plan bool) *Reconciler {
	logger.Info("plan: %v", plan)
	return &Reconciler{
		cfg:          cfg,
		plan:         plan,
		ctl:          ctl,
		irsaManager:  irsa.New(cfg.Metadata.Name, stackManager, oidcManager, clientSet),
		stackManager: stackManager,
	}
}

func (r *Reconciler) Reconcile() error {
	createTasks, updateTasks, deleteTasks, err := r.ReconcileIAMServiceAccounts()
	if err != nil {
		return err
	}

	logger.Info("Creating: %s", createTasks.Describe())
	logger.Info("Updating: %s", updateTasks.Describe())
	logger.Info("Deleting: %s", deleteTasks.Describe())

	if !r.plan {
		createErrs := createTasks.DoAllSync()
		updateErrs := updateTasks.DoAllSync()
		deleteErrs := deleteTasks.DoAllSync()
		logErrors(createErrs)
		logErrors(updateErrs)
		logErrors(deleteErrs)
		if len(createErrs) != 0 || len(updateErrs) != 0 || len(deleteErrs) != 0 {
			return fmt.Errorf("failed to reconcile cluster")
		}
	}
	return nil
}

func logErrors(errs []error) {
	for _, err := range errs {
		logger.Info("%v", err)
	}
}
