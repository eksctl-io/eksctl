package main

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
)

// TODO (alpha release)
// - add flags and flag defaults
// - basic support for addons
// - other key items from the readme

var rootCmd = &cobra.Command{
	Use: "eksctl",
	Run: func(c *cobra.Command, _ []string) {
		c.Help()
	},
}

func init() {

	addCommands()

	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 4, "set log level")
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
}
