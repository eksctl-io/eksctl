package create

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/names"
)

type nodegroupOptions struct {
	cmdutils.CreateNGOptions
	cmdutils.CreateManagedNGOptions
	UpdateAuthConfigMap     bool
	SkipOutdatedAddonsCheck bool
	SubnetIDs               []string
}

func createNodeGroupCmd(cmd *cmdutils.Cmd) {
	createNodeGroupCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options nodegroupOptions) error {
		if ng.Name != "" && api.IsInvalidNameArg(ng.Name) {
			return api.ErrInvalidName(ng.Name)
		}

		if options.SubnetIDs != nil {
			ng.Subnets = append(ng.Subnets, options.SubnetIDs...)
		}

		ngFilter := filter.NewNodeGroupFilter()
		if err := cmdutils.NewCreateNodeGroupLoader(cmd, ng, ngFilter, options.CreateNGOptions, options.CreateManagedNGOptions).Load(); err != nil {
			return errors.Wrap(err, "couldn't create node group filter from command line options")
		}

		if options.DryRun {
			originalWriter := logger.Writer
			logger.Writer = io.Discard
			defer func() {
				logger.Writer = originalWriter
			}()
		}

		ctx := context.Background()
		ctl, err := cmd.NewProviderForExistingClusterHelper(ctx, checkNodeGroupVersion)
		if err != nil {
			return fmt.Errorf("could not create cluster provider from options: %w", err)
		}

		if ok, err := ctl.CanOperate(cmd.ClusterConfig); !ok {
			return err
		}

		clientSet, err := ctl.NewStdClientSet(cmd.ClusterConfig)
		if err != nil {
			return err
		}

		instanceSelector, err := selector.New(ctx, ctl.AWSProvider.AWSConfig())
		if err != nil {
			return err
		}

		manager := nodegroup.New(cmd.ClusterConfig, ctl, clientSet, instanceSelector)
		return manager.Create(ctx, nodegroup.CreateOpts{
			InstallNeuronDevicePlugin: options.InstallNeuronDevicePlugin,
			InstallNvidiaDevicePlugin: options.InstallNvidiaDevicePlugin,
			UpdateAuthConfigMap:       options.UpdateAuthConfigMap,
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    options.DryRun,
				OutStream: cmd.CobraCommand.OutOrStdout(),
			},
			SkipOutdatedAddonsCheck: options.SkipOutdatedAddonsCheck,
			ConfigFileProvided:      cmd.ClusterConfigFile != "",
		}, ngFilter)
	})
}

type runFn func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options nodegroupOptions) error

func createNodeGroupCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc runFn) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var options nodegroupOptions

	cfg.Metadata.Version = "auto"

	cmd.SetDescription("nodegroup", "Create a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng, options)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddStringToStringVarPFlag(fs, &cfg.Metadata.Tags, "tags", "", map[string]string{}, "Used to tag the AWS resources")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, `for nodegroups "auto" and "latest" can be used to automatically inherit version from the control plane or force latest`)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		cmdutils.AddUpdateAuthConfigMap(fs, &options.UpdateAuthConfigMap, "Add nodegroup IAM role to aws-auth configmap")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddSubnetIDs(fs, &options.SubnetIDs, "Define an optional list of subnet IDs to create the nodegroup in")
		fs.BoolVarP(&options.DryRun, "dry-run", "", false, "Dry-run mode that skips nodegroup creation and outputs a ClusterConfig")
		fs.BoolVarP(&options.SkipOutdatedAddonsCheck, "skip-outdated-addons-check", "", false, "whether the creation of ARM nodegroups should proceed when the cluster addons are outdated")
	})

	cmd.FlagSetGroup.InFlagSet("New nodegroup", func(fs *pflag.FlagSet) {
		exampleNodeGroupName := names.ForNodeGroup("", "")
		fs.StringVarP(&ng.Name, "name", "n", "", fmt.Sprintf("name of the new nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		cmdutils.AddCommonCreateNodeGroupFlags(fs, cmd, ng, &options.CreateManagedNGOptions)
	})

	cmd.FlagSetGroup.InFlagSet("Addons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupAddonsFlags(fs, ng, &options.CreateNGOptions)
	})

	cmdutils.AddInstanceSelectorOptions(cmd.FlagSetGroup, ng)

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, true)
}

func checkNodeGroupVersion(ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
	switch meta.Version {
	case "auto":
		break
	case "":
		meta.Version = "auto"
	case "default":
		meta.Version = api.DefaultVersion
		logger.Info("will use default version (%s) for new nodegroup(s)", meta.Version)
	case "latest":
		meta.Version = api.LatestVersion
		logger.Info("will use latest version (%s) for new nodegroup(s)", meta.Version)
	default:
		if !api.IsSupportedVersion(meta.Version) {
			if api.IsDeprecatedVersion(meta.Version) {
				return fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", meta.Version, strings.Join(api.SupportedVersions(), ", "))
			}
			return fmt.Errorf("invalid version %s, supported values: auto, default, latest, %s", meta.Version, strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if v := ctl.ControlPlaneVersion(); v == "" {
		return fmt.Errorf("unable to get control plane version")
	} else if meta.Version == "auto" {
		meta.Version = v
		logger.Info("will use version %s for new nodegroup(s) based on control plane version", meta.Version)
	} else if meta.Version != v {
		hint := "--version=auto"
		logger.Warning("will use version %s for new nodegroup(s), while control plane version is %s; to automatically inherit the version use %q", meta.Version, v, hint)
	}

	return nil
}
