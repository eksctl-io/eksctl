package fargate

import (
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/logger"
)

type createFargateStackTask struct {
	cfg      *api.ClusterConfig
	provider api.ClusterProvider
}

func (t *createFargateStackTask) Describe() string { return "create fargate IAM stacK" }

func makeClusterStackName(clusterName string) string {
	return "eksctl-" + clusterName + "-fargate"
}

func (t *createFargateStackTask) Do(errs chan error) error {
	stackCollection := manager.NewStackCollection(t.provider, t.cfg)
	rs := builder.NewFargateResourceSet(t.provider, t.cfg)
	if err := rs.AddAllResources(); err != nil {
		return errors.Wrap(err, "couldn't add all resources to fargate resource set")
	}
	return stackCollection.CreateStack(makeClusterStackName(t.cfg.Metadata.Name), rs, nil, nil, errs)
}

// EnsureUnownedClusterReadyForFargate creates fargate IAM resources if they
// don't exist and are needed.
func EnsureUnownedClusterReadyForFargate(
	cfg *api.ClusterConfig, provider api.ClusterProvider, stackManager manager.StackManager,
) error {
	if api.IsSetAndNonEmptyString(cfg.IAM.FargatePodExecutionRoleARN) {
		return nil
	}

	fargateStack, err := stackManager.GetFargateStack()
	if err != nil {
		return err
	}
	if fargateStack == nil {
		var taskTree tasks.TaskTree
		taskTree.Append(&createFargateStackTask{
			cfg:      cfg,
			provider: provider,
		})
		errs := taskTree.DoAllSync()
		for _, e := range errs {
			logger.Critical("%s\n", e.Error())
		}
		if len(errs) > 0 {
			return errors.New("couldn't create fargate stack")
		}
	}
	return nil
}
