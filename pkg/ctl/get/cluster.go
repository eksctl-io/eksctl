package get

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func getClusterCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	var listAllRegions bool

	params := &getCmdParams{}

	rc.SetDescription("cluster", "Get cluster(s)", "", "clusters")

	rc.SetRunFuncWithNameArg(func() error {
		return doGetCluster(rc, params, listAllRegions)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		fs.BoolVarP(&listAllRegions, "all-regions", "A", false, "List clusters across all supported regions")
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doGetCluster(rc *cmdutils.ResourceCmd, params *getCmdParams, listAllRegions bool) error {
	cfg := rc.ClusterConfig
	regionGiven := cfg.Metadata.Region != "" // eks.New resets this field, so we need to check if it was set in the fist place

	ctl := eks.New(rc.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(rc.ProviderConfig)
	}

	if regionGiven && listAllRegions {
		logger.Warning("--region=%s is ignored, as --all-regions is given", cfg.Metadata.Region)
	}

	if cfg.Metadata.Name != "" && rc.NameArg != "" {
		return cmdutils.ErrNameFlagAndArg(cfg.Metadata.Name, rc.NameArg)
	}

	if rc.NameArg != "" {
		cfg.Metadata.Name = rc.NameArg
	}

	if cfg.Metadata.Name != "" && listAllRegions {
		return fmt.Errorf("--all-regions is for listing all clusters, it must be used without cluster name flag/argument")
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	return ctl.ListClusters(cfg.Metadata.Name, params.chunkSize, params.output, listAllRegions)
}
