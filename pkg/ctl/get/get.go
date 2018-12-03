package get

import (
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	defaultChunkSize = 100
)

var (
	chunkSize int
	output    string
)

// Command will create the `get` commands
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(getClusterCmd())
	cmd.AddCommand(getNodegroupCmd())

	return cmd
}
