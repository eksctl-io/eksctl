// Package misc provides miscellaneous commands for eksctl
package misc

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/version"
)

// versionCmd implements the 'version' command which displays the eksctl version
// It supports both plain text and JSON output formats
func versionCmd(cmd *cmdutils.Cmd) {
	var output string

	// Set up the command description and help text
	cmd.SetDescription("version", "Output the version of eksctl", "")

	// Define command flags
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		// Add -o/--output flag to specify output format
		fs.StringVarP(&output, "output", "o", "", "specifies the output format (valid option: json)")
	})

	// This command doesn't accept any arguments
	cmd.CobraCommand.Args = cobra.NoArgs

	// Define the command execution logic
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		switch output {
		case "":
			// Default output format (plain text)
			fmt.Printf("%s\n", version.GetVersion())
		case "json":
			// JSON output format
			fmt.Printf("%s\n", version.String())
		default:
			// Handle invalid output format
			return fmt.Errorf("unknown output: %s", output)
		}
		return nil
	}
}
