package main

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/gitops"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/completion"
	"github.com/weaveworks/eksctl/pkg/ctl/create"
	"github.com/weaveworks/eksctl/pkg/ctl/delete"
	"github.com/weaveworks/eksctl/pkg/ctl/drain"
	"github.com/weaveworks/eksctl/pkg/ctl/get"
	"github.com/weaveworks/eksctl/pkg/ctl/install"
	"github.com/weaveworks/eksctl/pkg/ctl/scale"
	"github.com/weaveworks/eksctl/pkg/ctl/update"
	"github.com/weaveworks/eksctl/pkg/ctl/utils"
)

func addCommands(rootCmd *cobra.Command, flagGrouping *cmdutils.FlagGrouping) {
	rootCmd.AddCommand(create.Command(flagGrouping))
	rootCmd.AddCommand(get.Command(flagGrouping))
	rootCmd.AddCommand(update.Command(flagGrouping))
	rootCmd.AddCommand(delete.Command(flagGrouping))
	rootCmd.AddCommand(scale.Command(flagGrouping))
	rootCmd.AddCommand(drain.Command(flagGrouping))
	if os.Getenv("EKSCTL_EXPERIMENTAL") == "true" {
		rootCmd.AddCommand(install.Command(flagGrouping))
		rootCmd.AddCommand(gitops.Command(flagGrouping))
	}
	rootCmd.AddCommand(utils.Command(flagGrouping))
	rootCmd.AddCommand(completion.Command(rootCmd))
	rootCmd.AddCommand(versionCmd(flagGrouping))
}

func main() {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "eksctl [command]",
		Short: "The official CLI for Amazon EKS",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring cobra error %q", err.Error())
			}
		},
	}

	flagGrouping := cmdutils.NewGrouping()

	addCommands(rootCmd, flagGrouping)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "help for this command")
	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 3, "set log level, use 0 to silence, 4 for debugging and 5 for debugging with AWS debug logging")

	colorValue := rootCmd.PersistentFlags().StringP("color", "C", "true", "toggle colorized logs (valid options: true, false, fabulous)")

	cobra.OnInitialize(func() {
		// Control colored output
		logger.Color = *colorValue == "true"
		logger.Fabulous = *colorValue == "fabulous"
		// Add timestamps for debugging
		logger.Timestamps = logger.Level >= 4
	})

	rootCmd.SetUsageFunc(flagGrouping.Usage)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
