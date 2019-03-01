package utils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateClusterStackCmd(_ *cmdutils.Grouping) *cobra.Command {
	return &cobra.Command{
		Use:   "update-cluster-stack",
		Short: "DEPRECATED: Use 'eksctl update cluster' instead",
		Run: func(cmd *cobra.Command, _ []string) {
			logger.Critical(cmd.Short)
			os.Exit(1)
		},
	}
}
