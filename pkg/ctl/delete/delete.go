package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
)

var (
	waitDelete bool
)

// Command will create the `delete` commands
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(deleteClusterCmd())

	return cmd
}
