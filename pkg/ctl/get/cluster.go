package get

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

var listAllRegions bool

func getClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:     "cluster",
		Short:   "Get cluster(s)",
		Aliases: []string{"clusters"},
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doGetCluster(cp); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		fs.BoolVarP(&listAllRegions, "all-regions", "A", false, "List clusters across all supported regions")
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &chunkSize, &output)
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, false)

	group.AddTo(cp.Command)
	return cp.Command
}

func doGetCluster(cp *cmdutils.CommonParams) error {
	cfg := cp.ClusterConfig
	regionGiven := cfg.Metadata.Region != "" // eks.New resets this field, so we need to check if it was set in the fist place

	ctl := eks.New(cp.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(cp.ProviderConfig)
	}

	if regionGiven && listAllRegions {
		logger.Warning("--region=%s is ignored, as --all-regions is given", cfg.Metadata.Region)
	}

	if cfg.Metadata.Name != "" && cp.NameArg != "" {
		return cmdutils.ErrNameFlagAndArg(cfg.Metadata.Name, cp.NameArg)
	}

	if cp.NameArg != "" {
		cfg.Metadata.Name = cp.NameArg
	}

	if cfg.Metadata.Name != "" && listAllRegions {
		return fmt.Errorf("--all-regions is for listing all clusters, it must be used without cluster name flag/argument")
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	return ctl.ListClusters(cfg.Metadata.Name, chunkSize, output, listAllRegions)
}
