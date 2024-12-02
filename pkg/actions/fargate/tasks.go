package fargate

import (
	"context"
	"errors"
	"fmt"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type createFargateStackTask struct {
	cfg          *api.ClusterConfig
	provider     api.ClusterProvider
	stackManager manager.StackManager
}

func (t *createFargateStackTask) Describe() string { return "create fargate IAM stacK" }

func makeClusterStackName(clusterName string) string {
	return "eksctl-" + clusterName + "-fargate"
}

func (t *createFargateStackTask) Do(errs chan error) error {
	rs := builder.NewFargateResourceSet(t.cfg)
	if err := rs.AddAllResources(); err != nil {
		return fmt.Errorf("couldn't add all resources to fargate resource set: %w", err)
	}
	return t.stackManager.CreateStack(context.TODO(), makeClusterStackName(t.cfg.Metadata.Name), rs, nil, nil, errs)
}

// ensureFargateRoleStackExists creates fargate IAM resources if they
func ensureFargateRoleStackExists(ctx context.Context, cfg *api.ClusterConfig, provider api.ClusterProvider, stackManager manager.StackManager) error {
	if api.IsSetAndNonEmptyString(cfg.IAM.FargatePodExecutionRoleARN) {
		return nil
	}

	fargateStack, err := stackManager.GetFargateStack(ctx)
	if err != nil {
		return err
	}

	if fargateStack == nil {
		var taskTree tasks.TaskTree
		taskTree.Append(&createFargateStackTask{
			cfg:          cfg,
			provider:     provider,
			stackManager: stackManager,
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
