package utils

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `utils` commands
func Command(g *cmdutils.Grouping) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utils",
		Short: "Various utils",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(waitNodesCmd(g))
	cmd.AddCommand(writeKubeconfigCmd(g))
	cmd.AddCommand(describeStacksCmd(g))
	cmd.AddCommand(updateClusterStackCmd(g))
	cmd.AddCommand(updateKubeProxyCmd(g))
	cmd.AddCommand(updateAWSNodeCmd(g))
	cmd.AddCommand(updateCoreDNSCmd(g))
	cmd.AddCommand(installCoreDNSCmd(g))

	return cmd
}
