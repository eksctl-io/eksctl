package utils

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
)

// Command will create the `utils` commands
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utils",
		Short: "Various utils",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(waitNodesCmd())
	cmd.AddCommand(writeKubeconfigCmd())
	cmd.AddCommand(describeStacksCmd())

	return cmd
}
