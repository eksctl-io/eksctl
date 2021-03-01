package create

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/utils/names"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func createNodeGroupCmd(cmd *cmdutils.Cmd) {
	createNodeGroupCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options nodegroup.CreateOpts, mngOptions cmdutils.CreateManagedNGOptions) error {
		ngFilter := filter.NewNodeGroupFilter()
		if err := cmdutils.NewCreateNodeGroupLoader(cmd, ng, ngFilter, mngOptions).Load(); err != nil {
			return errors.Wrap(err, "couldn't create node group filter from command line options")
		}
		ctl, err := cmd.NewCtl()
		if err != nil {
			return errors.Wrap(err, "couldn't create cluster provider from command line options")
		}
		cmdutils.LogRegionAndVersionInfo(cmd.ClusterConfig.Metadata)

		if ok, err := ctl.CanOperate(cmd.ClusterConfig); !ok {
			return err
		}

		clientSet, err := ctl.NewStdClientSet(cmd.ClusterConfig)
		if err != nil {
			return err
		}

		manager := nodegroup.New(cmd.ClusterConfig, ctl, clientSet)
		return manager.Create(options, *ngFilter)
	})
}

type runFn func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options nodegroup.CreateOpts, mngOptions cmdutils.CreateManagedNGOptions) error

func createNodeGroupCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc runFn) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var (
		options    nodegroup.CreateOpts
		mngOptions cmdutils.CreateManagedNGOptions
	)

	cfg.Metadata.Version = "auto"

	cmd.SetDescription("nodegroup", "Create a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng, options, mngOptions)
	}

	exampleNodeGroupName := names.ForNodeGroup("", "")

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to add the nodegroup to")
		cmdutils.AddStringToStringVarPFlag(fs, &cfg.Metadata.Tags, "tags", "", map[string]string{}, "Used to tag the AWS resources")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, `for nodegroups "auto" and "latest" can be used to automatically inherit version from the control plane or force latest`)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		cmdutils.AddUpdateAuthConfigMap(fs, &options.UpdateAuthConfigMap, "Add nodegroup IAM role to aws-auth configmap")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("New nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", fmt.Sprintf("name of the new nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		cmdutils.AddCommonCreateNodeGroupFlags(fs, cmd, ng, &mngOptions)
	})

	cmd.FlagSetGroup.InFlagSet("Addons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupIAMAddonsFlags(fs, ng)
		fs.BoolVarP(&options.InstallNeuronDevicePlugin, "install-neuron-plugin", "", true, "install Neuron plugin for Inferentia nodes")
		fs.BoolVarP(&options.InstallNvidiaDevicePlugin, "install-nvidia-plugin", "", true, "install Nvidia plugin for GPU nodes")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, true)
}
