package get

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var (
	chunkSize int
	output    string
)

// Command will create the `get` commands
func Command(g *cmdutils.Grouping) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(getClusterCmd(g))
	cmd.AddCommand(getNodegroupCmd(g))
	cmd.AddCommand(getIAMIdentityMappingCmd(g))

	return cmd
}
