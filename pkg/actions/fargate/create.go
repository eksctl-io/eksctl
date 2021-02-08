package fargate

import (
	"fmt"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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

	stackManager := ctl.NewStackManager(cfg)

	hasClusterStack, err := stackManager.HasClusterStack()
	if err != nil {
		return errors.Wrap(err, "couldn't check stack")
	}

	if !hasClusterStack {
		if err := EnsureUnownedClusterReadyForFargate(cfg, ctl.Provider, stackManager); err != nil {
			return errors.Wrap(err, "couldn't ensure unowned cluster is ready for fargate")
		}
	} else {
		if err := ctl.LoadClusterIntoSpec(cfg); err != nil {
			return errors.Wrap(err, "couldn't load cluster into spec")
		}
	}

	if !api.IsSetAndNonEmptyString(cfg.IAM.FargatePodExecutionRoleARN) {
		// Read back the default Fargate pod execution role ARN from CloudFormation:
		if err := ctl.NewStackManager(cfg).RefreshFargatePodExecutionRoleARN(); err != nil {
			return errors.Wrap(err, "couldn't refresh role arn")
		}
	}

	manager := fargate.NewFromProvider(cfg.Metadata.Name, ctl.Provider)
	if err := eks.DoCreateFargateProfiles(cfg, &manager); err != nil {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}
	return eks.ScheduleCoreDNSOnFargateIfRelevant(cfg, ctl, clientSet)
}
