package create

import (
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
)

// Command will create the `create` commands
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(createClusterCmd())

	return cmd
}
