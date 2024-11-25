package update

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	automodeactions "github.com/weaveworks/eksctl/pkg/actions/automode"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/automode"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type autoModeOptions struct {
	drainNodeGroups      bool
	drainParallel        int
	ignoreMissingSubnets bool
}

func updateAutoModeConfigCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"auto-mode-config",
		"Update the Auto Mode config",
		"Enable or disable Auto Mode on an existing cluster",
	)

	var options autoModeOptions
	cmd.FlagSetGroup.InFlagSet("Auto Mode", func(fs *pflag.FlagSet) {
		fs.BoolVar(&options.drainNodeGroups, "drain-all-nodegroups", false, "Drains nodegroups after enabling "+
			"Auto Mode in the cluster so that workloads on existing nodegroups are moved to the Auto Mode")
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
		return updateAutoMode(cmd, options)
	}
}

func updateAutoMode(cmd *cmdutils.Cmd, options autoModeOptions) error {
	if err := cmdutils.NewAutoModeLoader(cmd).Load(); err != nil {
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
	clusterRoleName, err := getClusterRoleName(ctl.Status.ClusterInfo.Cluster)
	if err != nil {
		return err
	}
	autoModeUpdater := &automodeactions.Updater{
		EKSUpdater: ctl.AWSProvider.EKS(),
		RoleManager: &roleManager{
			RoleCreator: &automode.RoleCreator{
				StackCreator: stackManager,
			},
			RoleDeleter: &automode.RoleDeleter{
				Cluster:      ctl.Status.ClusterInfo.Cluster,
				StackDeleter: stackManager,
			},
		},
		ClusterRoleManager: &automode.ClusterRoleManager{
			StackManager:    stackManager,
			IAMRoleManager:  ctl.AWSProvider.IAM(),
			ClusterRoleName: clusterRoleName,
			Region:          cmd.ClusterConfig.Metadata.Region,
		},
		PodsGetter: clientSet.CoreV1(),
	}
	if options.drainNodeGroups {
		autoModeUpdater.Drainer = &nodeGroupDrainer{
			nodeGroupLister: &eks.NodeGroupLister{
				NodeGroupStackLister: stackManager,
			},
			clientSet:     clientSet,
			drainParallel: options.drainParallel,
		}
	}
	return autoModeUpdater.Update(ctx, cmd.ClusterConfig, ctl.Status.ClusterInfo.Cluster)
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

func getClusterRoleName(cluster *ekstypes.Cluster) (string, error) {
	roleARN, err := arn.Parse(*cluster.RoleArn)
	if err != nil {
		return "", fmt.Errorf("parsing cluster role ARN %q: %w", *cluster.RoleArn, err)
	}
	parts := strings.Split(roleARN.Resource, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("expected role to have pattern %q; got %q: %w", "role/role-name", roleARN.Resource, err)
	}
	return parts[1], nil
}

type roleManager struct {
	*automode.RoleCreator
	*automode.RoleDeleter
}
