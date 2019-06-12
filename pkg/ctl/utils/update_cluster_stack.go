package utils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateClusterStackCmd(rc *cmdutils.ResourceCmd) {
	rc.SetDescription("update-cluster-stack", "DEPRECATED: Use 'eksctl update cluster' instead", "")

	rc.Command.Run = func(cmd *cobra.Command, _ []string) {
		logger.Critical(cmd.Short)
		os.Exit(1)
	}
}
