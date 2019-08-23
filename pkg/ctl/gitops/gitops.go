package gitops

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `gitops` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("gitops", "Helps setting up GitOps in a cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, applyGitops)

	return verbCmd
}
