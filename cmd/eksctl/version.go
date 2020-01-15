package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/version"
)

var output string

func versionCmd(_ *cmdutils.FlagGrouping) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Output the version of eksctl",
		RunE: func(_ *cobra.Command, _ []string) error {
			switch output {
			case "":
				fmt.Printf("%s\n", version.GetVersion())
			case "json":
				fmt.Printf("%#v\n", version.GetVersionInfo())
			default:
				return fmt.Errorf("unknown output: %s", output)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "specifies the output format (valid option: json)")
	return cmd
}
