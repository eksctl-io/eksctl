package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var (
	wait = false
	plan = true

	clusterConfigFile = ""
)

// Command will create the `delete` commands
func Command(g *cmdutils.Grouping) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(deleteClusterCmd(g))
	cmd.AddCommand(deleteNodeGroupCmd(g))
	cmd.AddCommand(deleteIAMIdentityMappingCmd(g))

	return cmd
}
