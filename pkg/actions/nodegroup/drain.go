package nodegroup

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/weaveworks/eksctl/pkg/drain"
)

type DrainInput struct {
	NodeGroups          []eks.KubeNodeGroup
	Plan                bool
	MaxGracePeriod      time.Duration
	NodeDrainWaitPeriod time.Duration
	Undo                bool
	DisableEviction     bool
}

func (m *Manager) Drain(input *DrainInput) error {
	if !input.Plan {
		for _, n := range input.NodeGroups {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(m.clientSet, n, m.ctl.Provider.WaitTimeout(), input.MaxGracePeriod, input.NodeDrainWaitPeriod, input.Undo, input.DisableEviction)
			if err := nodeGroupDrainer.Drain(); err != nil {
				return err
			}
		}
	}
	return nil
}
