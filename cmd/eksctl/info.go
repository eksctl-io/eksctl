package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/info"
)

func infoCmd(cmd *cmdutils.Cmd) {
	var output string

	cmd.SetDescription("info", "Output the version of eksctl, kubectl and OS info", "")
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&output, "output", "o", "", "specifies the output format (valid option: json)")
	})
	cmd.CobraCommand.Args = cobra.NoArgs
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		switch output {
		case "":
			version := info.GetInfo()
			fmt.Printf("eksctl version: %s\n", version.EksctlVersion)
			fmt.Printf("kubectl version: %s\n", version.KubectlVersion)
			fmt.Printf("OS: %s\n", version.OS)
		case "json":
			fmt.Printf("%s\n", info.String())
		default:
			return fmt.Errorf("unknown output: %s", output)
		}
		return nil
	}
}
