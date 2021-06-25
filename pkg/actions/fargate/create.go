package fargate

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

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

	clusterStack, err := m.stackManager.DescribeClusterStack()
	if err != nil {
		return errors.Wrap(err, "couldn't check cluster stack")
	}

	fargateRoleNeeded := false

	for _, profile := range cfg.FargateProfiles {
		if profile.PodExecutionRoleARN == "" {
			fargateRoleNeeded = true
			break
		}
	}

	if fargateRoleNeeded {
		if clusterStack != nil {
			if !m.fargateRoleExistsOnClusterStack(clusterStack) {
				err := ensureFargateRoleStackExists(cfg, ctl.Provider, m.stackManager)
				if err != nil {
					return errors.Wrap(err, "couldn't ensure fargate role exists")
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
			// Read back the default Fargate pod execution role ARN from CloudFormation:
			if err := m.stackManager.RefreshFargatePodExecutionRoleARN(); err != nil {
				return errors.Wrap(err, "couldn't refresh role arn")
			}
		}
	}

	fargateClient := fargate.NewFromProvider(cfg.Metadata.Name, ctl.Provider, m.stackManager)
	if err := eks.DoCreateFargateProfiles(cfg, &fargateClient); err != nil {
		return errors.Wrap(err, "could not create fargate profiles")
	}
	clientSet, err := m.newStdClientSet()
	if err != nil {
		return errors.Wrap(err, "couldn't create kubernetes client")
	}
	return eks.ScheduleCoreDNSOnFargateIfRelevant(cfg, ctl, clientSet)
}

func (m *Manager) fargateRoleExistsOnClusterStack(clusterStack *manager.Stack) bool {
	for _, output := range clusterStack.Outputs {
		if *output.OutputKey == outputs.FargatePodExecutionRoleARN {
			return true
		}
	}
	return false
}
