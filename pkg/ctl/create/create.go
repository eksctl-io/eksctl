package create

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command creates  the `create` commands.
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("create", "Create resource(s)", "")

	cmdFuncs := []func(*cmdutils.Cmd){
		createClusterCmd,
		createNodeGroupCmd,
		createIAMServiceAccountCmd,
		createIAMIdentityMappingCmd,
		createFargateProfile,
		createAddonCmd,
		createAccessEntryCmd,
		createPodIdentityAssociationCmd,
	}
	for _, cmdFunc := range cmdFuncs {
		cmdutils.AddResourceCmd(flagGrouping, verbCmd, cmdFunc)
	}
	return verbCmd
}
