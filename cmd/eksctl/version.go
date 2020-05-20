package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/version"
)

func versionCmd(cmd *cmdutils.Cmd) {
	var output string

	cmd.SetDescription("version", "Output the version of eksctl", "")
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&output, "output", "o", "", "specifies the output format (valid option: json)")
	})
	cmd.CobraCommand.Args = cobra.NoArgs
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		switch output {
		case "":
			fmt.Printf("%s\n", version.GetVersion())
		case "json":
			fmt.Printf("%s\n", version.String())
		default:
			return fmt.Errorf("unknown output: %s", output)
		}
		return nil
	}
}
