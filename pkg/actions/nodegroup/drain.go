package nodegroup

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/logger"
)

func (m *Manager) Drain(nodeGroups []eks.KubeNodeGroup, plan bool, maxGracePeriod time.Duration) error {
	cmdutils.LogIntendedAction(plan, "drain %d nodegroup(s) in cluster %q", len(nodeGroups), m.cfg.Metadata.Name)

	if !plan {
		for _, n := range nodeGroups {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(m.clientSet, n, m.ctl.Provider.WaitTimeout(), maxGracePeriod, false)
			if err := nodeGroupDrainer.Drain(); err != nil {
				logger.Warning("error occurred during drain, to skip drain use '--drain=false' flag")
				return err
			}
		}
	}
	return nil
}
