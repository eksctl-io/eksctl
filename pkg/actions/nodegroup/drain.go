package nodegroup

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/kris-nova/logger"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks"
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

// A Drainer drains nodegroups.
type Drainer struct {
	ClientSet kubernetes.Interface
}

// Drain drains nodegroups.
func (d *Drainer) Drain(ctx context.Context, input *DrainInput) error {
	parallelLimit := int64(input.Parallel)
	sem := semaphore.NewWeighted(parallelLimit)
	logger.Info("starting parallel draining, max in-flight of %d", parallelLimit)

	if input.Plan {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, nodegroup := range input.NodeGroups {
		nodegroup := nodegroup
		g.Go(func() error {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(d.ClientSet, nodegroup, input.MaxGracePeriod, input.NodeDrainWaitPeriod, input.PodEvictionWaitPeriod, input.Undo, input.DisableEviction, input.Parallel)
			return nodeGroupDrainer.Drain(ctx, sem)
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("draining nodegroups: %w", err)
	}
	return nil
}
