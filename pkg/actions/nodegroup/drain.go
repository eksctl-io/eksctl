package nodegroup

import (
	"context"
	"time"

	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/weaveworks/eksctl/pkg/drain"
)

type DrainInput struct {
	NodeGroups            []eks.KubeNodeGroup
	Plan                  bool
	MaxGracePeriod        time.Duration
	NodeDrainWaitPeriod   time.Duration
	PodEvictionWaitPeriod time.Duration
	Undo                  bool
	DisableEviction       bool
	Parallel              int
}

func (m *Manager) Drain(ctx context.Context, input *DrainInput) error {
	if input.Plan {
		return nil
	}
	for _, n := range input.NodeGroups {
		nodeGroupDrainer := drain.NewNodeGroupDrainer(m.clientSet, n, input.MaxGracePeriod, input.NodeDrainWaitPeriod, input.PodEvictionWaitPeriod, input.Undo, input.DisableEviction, input.Parallel)
		if err := nodeGroupDrainer.Drain(ctx); err != nil {
			return err
		}
	}
	return nil
}
