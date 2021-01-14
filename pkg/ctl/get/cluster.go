package get

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func getClusterCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var listAllRegions bool

	params := &getCmdParams{}

	cmd.SetDescription("cluster", "Get cluster(s)", "", "clusters")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doGetCluster(cmd, params, listAllRegions)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		fs.BoolVarP(&listAllRegions, "all-regions", "A", false, "List clusters across all supported regions")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doGetCluster(cmd *cmdutils.Cmd, params *getCmdParams, listAllRegions bool) error {
	cfg := cmd.ClusterConfig
	regionGiven := cfg.Metadata.Region != "" // eks.New resets this field, so we need to check if it was set in the fist place

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if regionGiven && listAllRegions {
		logger.Warning("--region=%s is ignored, as --all-regions is given", cfg.Metadata.Region)
	}

	if cfg.Metadata.Name != "" && cmd.NameArg != "" {
		return cmdutils.ErrClusterFlagAndArg(cmd, cfg.Metadata.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		cfg.Metadata.Name = cmd.NameArg
	}

	if cfg.Metadata.Name != "" && listAllRegions {
		return fmt.Errorf("--all-regions is for listing all clusters, it must be used without cluster name flag/argument")
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if params.output == "table" && !listAllRegions {
		cmdutils.LogRegionAndVersionInfo(cmd.ClusterConfig.Metadata)
	}

	if cfg.Metadata.Name == "" {
		return getAndPrinterClusters(ctl, params, listAllRegions)
	}

	return getAndPrintCluster(cfg, ctl, params)
}

func getAndPrinterClusters(ctl *eks.ClusterProvider, params *getCmdParams, listAllRegions bool) error {
	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == "table" {
		addGetClustersSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	clusters, err := ctl.ListClusters(params.chunkSize, listAllRegions)
	if err != nil {
		return err
	}

	return printer.PrintObjWithKind("clusters", clusters, os.Stdout)
}

func addGetClustersSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c *api.ClusterConfig) string {
		return c.Metadata.Name
	})
	printer.AddColumn("REGION", func(c *api.ClusterConfig) string {
		return c.Metadata.Region
	})
	printer.AddColumn("EKSCTL CREATED", func(c *api.ClusterConfig) api.EKSCTLCreated {
		return c.Status.EKSCTLCreated
	})
}

func getAndPrintCluster(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, params *getCmdParams) error {
	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == "table" {
		addGetClusterSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	cluster, err := ctl.GetCluster(cfg.Metadata.Name)

	if err != nil {
		return err
	}

	if err := printer.PrintObjWithKind("clusters", []*awseks.Cluster{cluster}, os.Stdout); err != nil {
		return err
	}

	return nil
}

func addGetClusterSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c *awseks.Cluster) string {
		return *c.Name
	})
	printer.AddColumn("VERSION", func(c *awseks.Cluster) string {
		return *c.Version
	})
	printer.AddColumn("STATUS", func(c *awseks.Cluster) string {
		return *c.Status
	})
	printer.AddColumn("CREATED", func(c *awseks.Cluster) string {
		return c.CreatedAt.Format(time.RFC3339)
	})
	printer.AddColumn("VPC", func(c *awseks.Cluster) string {
		return *c.ResourcesVpcConfig.VpcId
	})
	printer.AddColumn("SUBNETS", func(c *awseks.Cluster) string {
		subnets := sets.NewString()
		for _, subnetid := range c.ResourcesVpcConfig.SubnetIds {
			if api.IsSetAndNonEmptyString(subnetid) {
				subnets.Insert(*subnetid)
			}
		}
		return strings.Join(subnets.List(), ",")
	})
	printer.AddColumn("SECURITYGROUPS", func(c *awseks.Cluster) string {
		groups := sets.NewString()
		for _, sg := range c.ResourcesVpcConfig.SecurityGroupIds {
			if api.IsSetAndNonEmptyString(sg) {
				groups.Insert(*sg)
			}
		}
		return strings.Join(groups.List(), ",")
	})
}
