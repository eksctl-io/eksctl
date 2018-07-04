package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eksctl",
	Short: "a CLI for Amazon EKS",
	Run: func(c *cobra.Command, _ []string) {
		c.Help()
	},
}

func init() {

	addCommands()

	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 3, "set log level, use 0 to silence, 4 for debugging and 5 for debugging with AWS debug logging")
	rootCmd.PersistentFlags().BoolVarP(&logger.Color, "color", "C", true, "toggle colorized logs")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func addCommands() {
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(createCmd())
	rootCmd.AddCommand(deleteCmd())
	rootCmd.AddCommand(getCmd())
	rootCmd.AddCommand(utilsCmd())
}

func getNameArg(args []string) string {
	if len(args) > 1 {
		logger.Critical("only one argument is allowed to be used as a name")
		os.Exit(1)
	}
	if len(args) == 1 {
		return (strings.TrimSpace(args[0]))
	}
	return ""
}
