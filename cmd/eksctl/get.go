package main

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

const (
	DEFAULT_CHUNK_SIZE = 100
)

var (
	chunkSize int
	output    string
)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(getClusterCmd())

	return cmd
}

func getClusterCmd() *cobra.Command {
	cfg := &api.ClusterConfig{}

	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Get cluster(s)",
		Aliases: []string{"clusters"},
		Run: func(_ *cobra.Command, args []string) {
			if err := doGetCluster(cfg, getNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name")
	fs.IntVar(&chunkSize, "chunk-size", DEFAULT_CHUNK_SIZE, "Return large lists in chunks rather than all at once. Pass 0 to disable.")

	fs.StringVarP(&cfg.Region, "region", "r", api.DEFAULT_EKS_REGION, "AWS region")
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
