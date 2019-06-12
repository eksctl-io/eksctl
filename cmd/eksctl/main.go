package main

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/completion"
	"github.com/weaveworks/eksctl/pkg/ctl/create"
	"github.com/weaveworks/eksctl/pkg/ctl/delete"
	"github.com/weaveworks/eksctl/pkg/ctl/drain"
	"github.com/weaveworks/eksctl/pkg/ctl/get"
	"github.com/weaveworks/eksctl/pkg/ctl/scale"
	"github.com/weaveworks/eksctl/pkg/ctl/update"
	"github.com/weaveworks/eksctl/pkg/ctl/utils"
)

var rootCmd = &cobra.Command{
	Use:   "eksctl [command]",
	Short: "a CLI for Amazon EKS",
	Run: func(c *cobra.Command, _ []string) {
		if err := c.Help(); err != nil {
			logger.Debug("ignoring error %q", err.Error())
		}
	},
}

func init() {

	var colorValue string

	flagGrouping := cmdutils.NewGrouping()

	addCommands(flagGrouping)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "help for this command")
	rootCmd.PersistentFlags().StringVarP(&colorValue, "color", "C", "true", "toggle colorized logs (true,false,fabulous)")
	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 3, "set log level, use 0 to silence, 4 for debugging and 5 for debugging with AWS debug logging")

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

	rootCmd.SetUsageFunc(flagGrouping.Usage)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) // outputs cobra errors
		os.Exit(-1)
	}
}

func addCommands(flagGrouping *cmdutils.FlagGrouping) {
	rootCmd.AddCommand(versionCmd(flagGrouping))
	rootCmd.AddCommand(create.Command(flagGrouping))
	rootCmd.AddCommand(delete.Command(flagGrouping))
	rootCmd.AddCommand(get.Command(flagGrouping))
	rootCmd.AddCommand(update.Command(flagGrouping))
	rootCmd.AddCommand(scale.Command(flagGrouping))
	rootCmd.AddCommand(drain.Command(flagGrouping))
	rootCmd.AddCommand(utils.Command(flagGrouping))
	rootCmd.AddCommand(completion.Command(rootCmd))
}
