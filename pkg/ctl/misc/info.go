// Package misc provides miscellaneous commands for eksctl
package misc

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/info"
)

// infoCmd implements the 'info' command which displays system information
// including eksctl version, kubectl version, and OS details
// It supports both plain text and JSON output formats
func infoCmd(cmd *cmdutils.Cmd) {
	var output string

	// Set up the command description and help text
	cmd.SetDescription("info", "Output the version of eksctl, kubectl and OS info", "")

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
			version := info.GetInfo()
			fmt.Printf("eksctl version: %s\n", version.EksctlVersion)
			fmt.Printf("kubectl version: %s\n", version.KubectlVersion)
			fmt.Printf("OS: %s\n", version.OS)
		case "json":
			// JSON output format
			fmt.Printf("%s\n", info.String())
		default:
			// Handle invalid output format
			return fmt.Errorf("unknown output: %s", output)
		}
		return nil
	}
}
