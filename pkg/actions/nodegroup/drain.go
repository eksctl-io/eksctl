package nodegroup

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/weaveworks/eksctl/pkg/drain"
)

func (m *Manager) Drain(nodeGroups []eks.KubeNodeGroup, plan bool, maxGracePeriod time.Duration, undo bool, disableEviction bool) error {
	if !plan {
		for _, n := range nodeGroups {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(m.clientSet, n, m.ctl.Provider.WaitTimeout(), maxGracePeriod, undo, disableEviction)
			if err := nodeGroupDrainer.Drain(); err != nil {
				return err
			}
		}
	}
	return nil
}
