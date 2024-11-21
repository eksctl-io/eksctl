package update

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"

	autonomousmodeactions "github.com/weaveworks/eksctl/pkg/actions/autonomousmode"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/autonomousmode"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type autonomousModeOptions struct {
	drainNodeGroups      bool
	drainParallel        int
	ignoreMissingSubnets bool
}

func updateAutonomousModeConfigCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"autonomous-mode-config",
		"Update the Autonomous Mode config",
		"Enable or disable Autonomous Mode on an existing cluster",
	)

	var options autonomousModeOptions
	cmd.FlagSetGroup.InFlagSet("Autonomous Mode", func(fs *pflag.FlagSet) {
		fs.BoolVar(&options.drainNodeGroups, "drain-all-nodegroups", false, "Drains nodegroups after enabling "+
			"Autonomous Mode in the cluster so that workloads on existing nodegroups are moved to the Autonomous Mode")
		fs.IntVar(&options.drainParallel, "drain-parallel", 1, "Specifies the number of nodes to drain in parallel")
		fs.BoolVar(&options.ignoreMissingSubnets, "ignore-missing-subnets", false, "If the cluster's CloudFormation stack "+
			"contains any subnets that no longer exist, eksctl fails with an error. Specifying this flag suppresses the error.")
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, 40*time.Minute)
	})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		return updateAutonomousMode(cmd, options)
	}
}

func updateAutonomousMode(cmd *cmdutils.Cmd, options autonomousModeOptions) error {
	if err := cmdutils.NewAutonomousModeLoader(cmd).Load(); err != nil {
		return err
	}
	if options.drainParallel < 0 {
		return fmt.Errorf("invalid value %v for --drain-parallel", options.drainParallel)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cmd.ProviderConfig.WaitTimeout)
	defer cancel()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cmd.ClusterConfig)
	if err != nil {
		return err
	}
	stackManager := ctl.NewStackManager(cmd.ClusterConfig)
	rawClient, err := ctl.NewRawClient(cmd.ClusterConfig)
	if err != nil {
		return err
	}
	autonomousModeUpdater := &autonomousmodeactions.Updater{
		EKSUpdater: ctl.AWSProvider.EKS(),
		RoleManager: &roleManager{
			RoleCreator: &autonomousmode.RoleCreator{
				StackCreator: stackManager,
			},
			RoleDeleter: &autonomousmode.RoleDeleter{
				Cluster:      ctl.Status.ClusterInfo.Cluster,
				StackDeleter: stackManager,
			},
		},
		CoreV1Interface: clientSet.CoreV1(),
		RBACApplier: &autonomousmode.RBACApplier{
			RawClient: rawClient,
		},
	}
	if options.drainNodeGroups {
		autonomousModeUpdater.Drainer = &nodeGroupDrainer{
			nodeGroupLister: &eks.NodeGroupLister{
				NodeGroupStackLister: stackManager,
			},
			clientSet:     clientSet,
			drainParallel: options.drainParallel,
		}
	}
	return autonomousModeUpdater.Update(ctx, cmd.ClusterConfig, ctl.Status.ClusterInfo.Cluster)
}

type nodeGroupDrainer struct {
	nodeGroupLister *eks.NodeGroupLister
	drainParallel   int
	clientSet       kubernetes.Interface
}

func (d *nodeGroupDrainer) Drain(ctx context.Context) error {
	kubeNodeGroups, err := d.nodeGroupLister.List(ctx)
	if err != nil {
		return err
	}
	drainer := &nodegroup.Drainer{
		ClientSet: d.clientSet,
	}
	// TODO: configurable timeout.
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	return drainer.Drain(ctx, &nodegroup.DrainInput{
		NodeGroups:            kubeNodeGroups,
		MaxGracePeriod:        10 * time.Minute,
		PodEvictionWaitPeriod: 10 * time.Second,
		Parallel:              d.drainParallel,
	})
}

type roleManager struct {
	*autonomousmode.RoleCreator
	*autonomousmode.RoleDeleter
}
