package main

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/completion"
	"github.com/weaveworks/eksctl/pkg/ctl/create"
	"github.com/weaveworks/eksctl/pkg/ctl/delete"
	"github.com/weaveworks/eksctl/pkg/ctl/get"
	"github.com/weaveworks/eksctl/pkg/ctl/scale"
	"github.com/weaveworks/eksctl/pkg/ctl/utils"
)

var (
	rootCmd = &cobra.Command{
		Use:   "eksctl",
		Short: "a CLI for Amazon EKS",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}
	colorValue string
)

func init() {
	addCommands()
	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 3, "set log level, use 0 to silence, 4 for debugging and 5 for debugging with AWS debug logging")
	rootCmd.PersistentFlags().StringVarP(&colorValue, "color", "C", "true", "toggle colorized logs (true,false,fabulous)")

	cobra.OnInitialize(func() {
		// Control colored output
		color := true
		fabulous := false
		switch colorValue {
		case "false":
			color = false
		case "fabulous":
			color = false
			fabulous = true
		}
		logger.Color = color
		logger.Fabulous = fabulous

		// Add timestamps for debugging
		logger.Timestamps = false
		if logger.Level >= 4 {
			logger.Timestamps = true
		}
	})
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
	rootCmd.AddCommand(completion.Command(rootCmd))
}
