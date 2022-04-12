package get

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/weaveworks/eksctl/pkg/actions/cluster"
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
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doGetCluster(cmd *cmdutils.Cmd, params *getCmdParams, listAllRegions bool) error {
	if err := cmdutils.NewGetClusterLoader(cmd).Load(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig
	regionGiven := cfg.Metadata.Region != "" // eks.New resets this field, so we need to check if it was set in the first place

	if params.output != printers.TableType {
		logger.Writer = os.Stderr
	}

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

	if params.output != printers.TableType {
		// log warnings and errors to stdout
		logger.Writer = os.Stderr
	}

	ctx := context.TODO()
	if cfg.Metadata.Name == "" {
		return getAndPrinterClusters(ctx, ctl, params, listAllRegions)
	}

	return getAndPrintCluster(ctx, cfg, ctl, params)
}

func getAndPrinterClusters(ctx context.Context, ctl *eks.ClusterProvider, params *getCmdParams, listAllRegions bool) error {
	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		addGetClustersSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	clusters, err := cluster.GetClusters(ctx, ctl.Provider, listAllRegions, params.chunkSize)
	if err != nil {
		return err
	}

	return printer.PrintObjWithKind("clusters", clusters, os.Stdout)
}

func addGetClustersSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c cluster.Description) string {
		return c.Name
	})
	printer.AddColumn("REGION", func(c cluster.Description) string {
		return c.Region
	})
	printer.AddColumn("EKSCTL CREATED", func(c cluster.Description) api.EKSCTLCreated {
		return c.Owned
	})
}

func getAndPrintCluster(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, params *getCmdParams) error {
	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		addGetClusterSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	cluster, err := ctl.GetCluster(ctx, cfg.Metadata.Name)

	if err != nil {
		return err
	}

	return printer.PrintObjWithKind("clusters", []*awseks.Cluster{cluster}, os.Stdout)
}

func addGetClusterSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c *awseks.Cluster) string {
		if c.Name == nil {
			return "-"
		}
		return *c.Name
	})
	printer.AddColumn("VERSION", func(c *awseks.Cluster) string {
		if c.Version == nil {
			return "-"
		}
		return *c.Version
	})
	printer.AddColumn("STATUS", func(c *awseks.Cluster) string {
		if c.Status == nil {
			return "-"
		}
		return *c.Status
	})
	printer.AddColumn("CREATED", func(c *awseks.Cluster) string {
		if c.CreatedAt == nil {
			return "-"
		}
		return c.CreatedAt.Format(time.RFC3339)
	})
	printer.AddColumn("VPC", func(c *awseks.Cluster) string {
		if c.ResourcesVpcConfig == nil {
			return "-"
		}
		return *c.ResourcesVpcConfig.VpcId
	})
	printer.AddColumn("SUBNETS", func(c *awseks.Cluster) string {
		if c.ResourcesVpcConfig == nil || c.ResourcesVpcConfig.SubnetIds == nil {
			return "-"
		}
		subnets := sets.NewString()
		for _, subnetid := range c.ResourcesVpcConfig.SubnetIds {
			if api.IsSetAndNonEmptyString(subnetid) {
				subnets.Insert(*subnetid)
			}
		}
		return strings.Join(subnets.List(), ",")
	})
	printer.AddColumn("SECURITYGROUPS", func(c *awseks.Cluster) string {
		if c.ResourcesVpcConfig == nil || c.ResourcesVpcConfig.SecurityGroupIds == nil {
			return "-"
		}
		groups := sets.NewString()
		for _, sg := range c.ResourcesVpcConfig.SecurityGroupIds {
			if api.IsSetAndNonEmptyString(sg) {
				groups.Insert(*sg)
			}
		}
		return strings.Join(groups.List(), ",")
	})

	printer.AddColumn("PROVIDER", func(c *awseks.Cluster) string {
		if c.ConnectorConfig != nil {
			return *c.ConnectorConfig.Provider
		}
		return "EKS"
	})
}
