package fargate

import (
	"fmt"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

func (m *Manager) Create() error {
	ctl := m.ctl
	cfg := m.cfg
	if ok, err := ctl.CanOperate(cfg); !ok {
		return errors.Wrap(err, "couldn't check cluster operable status")
	}

	supportsFargate, err := ctl.SupportsFargate(cfg)
	if err != nil {
		return errors.Wrap(err, "couldn't check fargate support")
	}
	if !supportsFargate {
		return fmt.Errorf("Fargate is not supported for this cluster version. Please update the cluster to be at least eks.%d", fargate.MinPlatformVersion)
	}

	hasClusterStack, err := m.stackManager.HasClusterStack()
	if err != nil {
		return errors.Wrap(err, "couldn't check stack")
	}

	if hasClusterStack {
		exists, err := m.fargateRoleExistsOnClusterStack()
		if err != nil {
			return err
		}

		if !exists {
			err := ensureFargateRoleStackExists(cfg, ctl.Provider, m.stackManager)
			if err != nil {
				return err
			}
		}

		if err := ctl.LoadClusterIntoSpecFromStack(cfg, m.stackManager); err != nil {
			return errors.Wrap(err, "couldn't load cluster into spec")
		}
	} else {
		if err := ensureFargateRoleStackExists(cfg, ctl.Provider, m.stackManager); err != nil {
			return errors.Wrap(err, "couldn't ensure unowned cluster is ready for fargate")
		}
	}

	if !api.IsSetAndNonEmptyString(cfg.IAM.FargatePodExecutionRoleARN) {
		if err := m.stackManager.RefreshFargatePodExecutionRoleARN(); err != nil {
			return errors.Wrap(err, "couldn't refresh role arn")
		}
	}

	manager := fargate.NewFromProvider(cfg.Metadata.Name, ctl.Provider)
	if err := eks.DoCreateFargateProfiles(cfg, &manager); err != nil {
		return err
	}
	clientSet, err := m.newStdClientSet()
	if err != nil {
		return err
	}
	return eks.ScheduleCoreDNSOnFargateIfRelevant(cfg, ctl, clientSet)
}

func (m *Manager) fargateRoleExistsOnClusterStack() (bool, error) {
	stack, err := m.stackManager.DescribeClusterStack()
	if err != nil {
		return false, err
	}

	for _, output := range stack.Outputs {
		if *output.OutputKey == outputs.FargatePodExecutionRoleARN {
			return true, nil
		}
	}

	return false, nil
}
