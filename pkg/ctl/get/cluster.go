package get

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)


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
	fs.IntVar(&chunkSize, "chunk-size", defaultChunkSize, "Return large lists in chunks rather than all at once. Pass 0 to disable.")

	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.StringVarP(&output, "output", "o", "table", "Specifies the output format. Choose from table,json,yaml. Defaults to table.")
	return cmd
}

func doGetCluster(cfg *api.ClusterConfig, name string) error {
	ctl := eks.New(cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.ClusterName != "" && name != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.ClusterName, name)
	}

	if name != "" {
		cfg.ClusterName = name
	}

	if err := ctl.ListClusters(chunkSize, output); err != nil {
		return err
	}

	return nil
}

func getSupportedRegionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supported-regions",
		Short: "List the regions EKS is supported",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("EKS is supported in the following regions:")
			fmt.Println("------------------------------------------")
			fmt.Println("us-east-1 - US East (N. Virginia)")
			fmt.Println("us-west-2 - US West (Oregon)")
			fmt.Println("eu-west-1 - EU (Ireland)")
			fmt.Println("------------------------------------------")
			fmt.Println("For the most up to date information, please check:")
			fmt.Println("https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services/")
		},
	}
	return cmd
}
