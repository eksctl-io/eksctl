package get

import (
	"fmt"
	"os"
	"strings"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

var listAllRegions bool

func getClusterCmd() *cobra.Command {
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Get cluster(s)",
		Aliases: []string{"clusters"},
		Run: func(_ *cobra.Command, args []string) {
			if err := doGetCluster(cfg, ctl.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name")
	fs.BoolVarP(&listAllRegions, "all-regions", "A", false, "List clusters across all supported regions")
	fs.IntVar(&chunkSize, "chunk-size", defaultChunkSize, "Return large lists in chunks rather than all at once. Pass 0 to disable.")

	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.StringVarP(&output, "output", "o", "table", "Specifies the output format. Choose from table,json,yaml. Defaults to table.")
	return cmd
}

func doGetCluster(cfg *api.ClusterConfig, name string) error {
	regionGiven := cfg.Region != "" // eks.New resets this field, so we need to check if it was set in the fist place
	ctl := eks.New(cfg)

	if !cfg.IsSupportedRegion() {
		return fmt.Errorf("--region=%s is not supported - use one of: %s", cfg.Region, strings.Join(api.SupportedRegions(), ", "))
	}

	if regionGiven && listAllRegions {
		logger.Warning("--region=%s is ignored, as --all-regions is given", cfg.Region)
	}

	if cfg.ClusterName != "" && name != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.ClusterName, name)
	}

	if name != "" {
		cfg.ClusterName = name
	}

	if cfg.ClusterName != "" && listAllRegions {
		return fmt.Errorf("--all-regions is for listing all clusters, it must be used without cluster name flag/argument")
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	return ctl.ListClusters(chunkSize, output, listAllRegions)
}
