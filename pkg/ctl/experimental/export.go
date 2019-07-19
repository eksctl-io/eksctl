package experimental

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

var (
	exportDir        string
	offline          bool
	withoutNodeGroup bool
	output           string
)

func exportTemplatesCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	rc.ClusterConfig = cfg

	rc.SetDescription("export-templates", "Export CloudFormation templates",
		"Generate and write CloudFormation templates to disk")

	rc.SetRunFuncWithNameArg(func() error {
		return doExport(rc)
	})

	exampleClusterName := cmdutils.ClusterName("", "")
	exampleNodeGroupName := cmdutils.NodeGroupName("", "")

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))
		fs.StringVarP(&output, "output", "o", "yaml", "specifies the output format (valid option: json, yaml)")
		fs.StringVar(&exportDir, "dir", defaultExportDir(), "the directory to save exported files into")
		fs.BoolVar(&offline, "offline", offline, "disable AWS EC2 AMI lookup and other AWS API interactions")
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
	})

	rc.FlagSetGroup.InFlagSet("Initial nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVar(&ng.Name, "nodegroup-name", "", fmt.Sprintf("name of the nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		fs.BoolVar(&withoutNodeGroup, "without-nodegroup", false, "if set, initial nodegroup will not be created")
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, true)
}

func doExport(rc *cmdutils.ResourceCmd) error {
	ngFilter := cmdutils.NewNodeGroupFilter()
	ngFilter.ExcludeAll = withoutNodeGroup

	// TODO: this should have a custom loader

	if err := cmdutils.NewCreateClusterLoader(rc, ngFilter).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig
	meta := rc.ClusterConfig.Metadata

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ngFilter.ValidateNodeGroupsAndSetDefaults(cfg.NodeGroups); err != nil {
		return err
	}

	if !offline {
		if err := ctl.CheckAuth(); err != nil {
			return err
		}
	} else {
		logger.Info("offline mode: skipping STS access check")
	}

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(rc.ProviderConfig)
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
		// TODO: should be possible to simplify this
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
			// TODO: handle AZs properly (offline and not)
			// if len(cfg.AvailabilityZones) == 0 {
			// 	cfg.AvailabilityZones = api.DefaultAvailabilityZones
			// }

			if err := vpc.SetSubnets(cfg); err != nil {
				return err
			}
			return nil
		}

		if !offline {
			if err := vpc.ImportAllSubnets(ctl.Provider, cfg); err != nil {
				return err
			}
		} else {
			logger.Info("offline mode: skip importing subnets")
		}

		if err := cfg.HasSufficientSubnets(); err != nil {
			logger.Critical("unable to use given %s", subnetInfo())
			return err
		}

		if err := ngFilter.ForEach(cfg.NodeGroups, canUseForPrivateNodeGroups); err != nil {
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
	err := ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
		// prepend bootstrap commands for cluster discovery
		{
			commonFlags := fmt.Sprintf("--region=%s --name=%s", meta.Region, meta.Name)

			ng.PreBootstrapCommands = append(ng.PreBootstrapCommands, "aws eks wait cluster-active "+commonFlags)

			describeCluster := "aws eks describe-cluster --output=text " + commonFlags

			ng.PreBootstrapCommands = append(ng.PreBootstrapCommands, describeCluster+" --query=cluster.certificateAuthority.data | base64 -d > /etc/eksctl/ca.crt")

			ng.PreBootstrapCommands = append(ng.PreBootstrapCommands, fmt.Sprintf(
				"kubectl --kubeconfig=/etc/eksctl/kubeconfig.yaml config set-cluster %s.%s.eksctl.io --server=\"$(%s --query=cluster.endpoint)\"",
				meta.Name, meta.Region, describeCluster))
		}

		// TODO: make sure SSH keys are handled correctly offline and not
		// resolve AMI
		if !offline {
			if err := ctl.EnsureAMI(meta.Version, ng); err != nil {
				return err
			}
		} else {
			logger.Info("offline mode: skipping AMI resolution")
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
	ngSubset, _ := ngFilter.MatchAll(cfg.NodeGroups)
	stackManager := ctl.NewStackManager(cfg)
	ngFilter.LogInfo(cfg.NodeGroups)

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

	// TODO: document steps one should use to deploy stacks, include auth configmap

	for name, template := range templates {
		if err := saveObj(printer, name, template); err != nil {
			return err
		}
	}
	logger.Success("all EKS cluster resource for %q had been exported", meta.Name)
	return nil
}

func defaultExportDir() string {
	return filepath.Join(".", "export") // TODO: name it by cluster name
}

func ensureDir(path string) error {
	// TODO: should consider creating a tarball instead, so we are clear what we exporting is all there is
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
