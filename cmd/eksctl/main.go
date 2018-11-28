package main

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/create"
	"github.com/weaveworks/eksctl/pkg/ctl/delete"
	"github.com/weaveworks/eksctl/pkg/ctl/get"
	"github.com/weaveworks/eksctl/pkg/ctl/scale"
	"github.com/weaveworks/eksctl/pkg/ctl/utils"
)

var rootCmd = &cobra.Command{
	Use:   "eksctl",
	Short: "a CLI for Amazon EKS",
	Run: func(c *cobra.Command, _ []string) {
		if err := c.Help(); err != nil {
			logger.Debug("ignoring error %q", err.Error())
		}
	},
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(eksctl completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(eksctl completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {

	addCommands()

	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 3, "set log level, use 0 to silence, 4 for debugging and 5 for debugging with AWS debug logging")
	rootCmd.PersistentFlags().BoolVarP(&logger.Color, "color", "C", true, "toggle colorized logs")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) // outputs cobra errors
		os.Exit(-1)
	}
}

func addCommands() {
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(create.Command())
	rootCmd.AddCommand(delete.Command())
	rootCmd.AddCommand(get.Command())
	rootCmd.AddCommand(scale.Command())
	rootCmd.AddCommand(utils.Command())
	rootCmd.AddCommand(completionCmd)
}
