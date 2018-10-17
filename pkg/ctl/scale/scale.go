package scale

import (
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
)

// Command will create the `scale` commands
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scale",
		Short: "Scale resources(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(scaleNodeGroupCmd())

	return cmd
}
