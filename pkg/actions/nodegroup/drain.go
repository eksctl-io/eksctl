package nodegroup

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/drain"
)

func (ng *Manager) Drain(nodeGroups []eks.KubeNodeGroup, plan bool, maxGracePeriod time.Duration) error {
	cmdutils.LogIntendedAction(plan, "drain %d nodegroup(s) in cluster %q", len(nodeGroups), ng.cfg.Metadata.Name)

	if !plan {
		for _, n := range nodeGroups {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(ng.clientSet, n, ng.ctl.Provider.WaitTimeout(), maxGracePeriod, false)
			if err := nodeGroupDrainer.Drain(); err != nil {
				logger.Warning("error occurred during drain, to skip drain use '--drain=false' flag")
				return err
			}
		}
	}
	return nil
}
