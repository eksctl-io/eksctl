// Package misc provides miscellaneous commands for eksctl that don't fit into other categories
package misc

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command registers all the miscellaneous commands to the root command
// These include utility commands like version and info
func Command(flagGrouping *cmdutils.FlagGrouping, rootCmd *cobra.Command) *cobra.Command {
	// Register the info command to display system information
	cmdutils.AddResourceCmd(flagGrouping, rootCmd, infoCmd)
	// Register the version command to display eksctl version
	cmdutils.AddResourceCmd(flagGrouping, rootCmd, versionCmd)
	return rootCmd
}
