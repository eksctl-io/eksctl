package main

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "eksctl",
	Run: func(c *cobra.Command, _ []string) {
		c.Help()
	},
}

func init() {

	addCommands()

	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 3, "set log level")
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
