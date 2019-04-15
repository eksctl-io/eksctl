package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

var (
	exportDir        string
	offline          bool
	withoutNodeGroup bool
)

func exportCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export CloudFormation stacks for a given cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doExport(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	exampleClusterName := cmdutils.ClusterName("", "")
	exampleNodeGroupName := cmdutils.NodeGroupName("", "")

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))
		fs.StringVarP(&output, "output", "o", "yaml", "specifies the output format (valid option: json, yaml)")
		fs.StringVar(&exportDir, "dir", defaultExportDir(), "the directory to save exported files into")
		fs.BoolVar(&offline, "offline", offline, "disable AWS EC2 AMI lookup and other AWS API interactions")
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
	})

	group.InFlagSet("Initial nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVar(&ng.Name, "nodegroup-name", "", fmt.Sprintf("name of the nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		fs.BoolVar(&withoutNodeGroup, "without-nodegroup", false, "if set, initial nodegroup will not be created")
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)

	return cmd
}

func doExport(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
	ngFilter := cmdutils.NewNodeGroupFilter()
	ngFilter.SkipAll = withoutNodeGroup

	if err := cmdutils.NewCreateClusterLoader(p, cfg, clusterConfigFile, nameArg, cmd, ngFilter).Load(); err != nil {
		return err
	}

	if err := ngFilter.ValidateNodeGroupsAndSetDefaults(cfg.NodeGroups); err != nil {
		return err
	}

	meta := cfg.Metadata
	ctl := eks.New(p, cfg)

	if !offline {
		if err := ctl.CheckAuth(); err != nil {
			return err
		}
	}

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if cfg.Metadata.Version == "" {
		cfg.Metadata.Version = api.LatestVersion
	}
	if cfg.Metadata.Version != api.LatestVersion {
		validVersion := false
		for _, v := range api.SupportedVersions() {
			if cfg.Metadata.Version == v {
				validVersion = true
			}
		}
		if !validVersion {
			return fmt.Errorf("invalid version, supported values: %s", strings.Join(api.SupportedVersions(), ", "))
		}
	}
	if cfg.Metadata.Region == "" {
		cfg.Metadata.Region = api.DefaultRegion
	}

	// VPC and subnets
	createOrImportVPC := func() error {

		subnetInfo := func() string {
			return fmt.Sprintf("VPC (%s) and subnets (private:%v public:%v)",
				cfg.VPC.ID, cfg.PrivateSubnetIDs(), cfg.PublicSubnetIDs())
		}

		customNetworkingNotice := "custom VPC/subnets will be used; if resulting cluster doesn't function as expected, make sure to review the configuration of VPC/subnets"

		canUseForPrivateNodeGroups := func(_ int, ng *api.NodeGroup) error {
			if ng.PrivateNetworking && !cfg.HasSufficientPrivateSubnets() {
				return fmt.Errorf("none or too few private subnets to use with --node-private-networking")
			}
			return nil
		}

		if !cfg.HasAnySubnets() {
			// default: create dedicated VPC
			if len(cfg.AvailabilityZones) == 0 {
				cfg.AvailabilityZones = api.DefaultAvailabilityZones
			}

			if err := vpc.SetSubnets(cfg); err != nil {
				return err
			}
			return nil
		}

		if err := vpc.ImportAllSubnets(ctl.Provider, cfg); err != nil {
			return err
		}

		if err := cfg.HasSufficientSubnets(); err != nil {
			logger.Critical("unable to use given %s", subnetInfo())
			return err
		}

		if err := ngFilter.CheckEachNodeGroup(cfg.NodeGroups, canUseForPrivateNodeGroups); err != nil {
			return err
		}

		logger.Success("using existing %s", subnetInfo())
		logger.Warning(customNetworkingNotice)
		return nil
	}

	if err := createOrImportVPC(); err != nil {
		return err
	}

	// Finalize node groups configuration
	err := ngFilter.CheckEachNodeGroup(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
		// default override user data
		if ng.OverrideUserData == nil {
			logger.Info("setting default user data override")
			ng.OverrideUserData = nodebootstrap.DefaultOverrideUserData(ng)
		}

		// resolve AMI
		if !offline {
			if err := ctl.EnsureAMI(meta.Version, ng); err != nil {
				return err
			}
		}
		logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, cfg.Metadata.Version)

		if err := ctl.SetNodeLabels(ng, meta); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("exporting %s", meta.LogString())

	// core action
	ngSubset := ngFilter.MatchAll(cfg)
	stackManager := ctl.NewStackManager(cfg)
	ngFilter.LogInfo(cfg)

	templates, errs := stackManager.ExportClusterWithNodeGroups(ngSubset)
	if len(errs) > 0 {
		logger.Info("%d error(s) occurred and cluster hasn't been exported properly", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to export cluster %q", meta.Name)
	}
	if len(templates) == 0 {
		logger.Info("no stacks to export")
	}

	err = ensureDir(exportDir)
	if err != nil {
		return err
	}

	printer, err := printers.NewPrinter(output)
	if err != nil {
		return err
	}

	configName := makeClusterConfigName(cfg.Metadata.Name)
	if err := saveObj(printer, configName, cfg); err != nil {
		return err
	}

	for name, template := range templates {
		if err := saveObj(printer, name, template); err != nil {
			return err
		}
	}
	logger.Success("all EKS cluster resource for %q had been exported", meta.Name)
	return nil
}

func defaultExportDir() string {
	return filepath.Join(".", "export")
}

func ensureDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return err
}

func saveObj(printer printers.OutputPrinter, name string, obj interface{}) error {
	path := filepath.Join(exportDir, name+"."+output)
	logger.Info("saving export file: %q", path)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	w := bufio.NewWriter(f)
	if err := printer.PrintObj(obj, w); err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}

func makeClusterConfigName(name string) string {
	return "eksctl-" + name + "-config"
}
