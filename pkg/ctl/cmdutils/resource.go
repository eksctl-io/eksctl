package cmdutils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// ResourceCmd holds attributes that most of the commands use
type ResourceCmd struct {
	Command      *cobra.Command
	FlagSetGroup *NamedFlagSetGroup

	Plan, Wait bool

	NameArg string

	ClusterConfigFile string

	ProviderConfig *api.ProviderConfig
	ClusterConfig  *api.ClusterConfig

	Include, Exclude []string
}

// AddResourceCmd create a registers a new command under the given verb command
func AddResourceCmd(flagGrouping *FlagGrouping, parentVerbCmd *cobra.Command, newResourceCmd func(*ResourceCmd)) {
	resource := &ResourceCmd{
		Command:        &cobra.Command{},
		ProviderConfig: &api.ProviderConfig{},

		Plan: true,  // always on by default
		Wait: false, // varies in some commands
	}
	resource.FlagSetGroup = flagGrouping.New(resource.Command)
	newResourceCmd(resource)
	resource.FlagSetGroup.AddTo(resource.Command)
	parentVerbCmd.AddCommand(resource.Command)
}

// SetDescription sets usage along with short and long descriptions as well as aliases
func (rc *ResourceCmd) SetDescription(use, short, long string, aliases ...string) {
	rc.Command.Use = use
	rc.Command.Short = short
	rc.Command.Long = long
	rc.Command.Aliases = aliases
}

// SetRunFunc registers a command function
func (rc *ResourceCmd) SetRunFunc(cmd func() error) {
	rc.Command.Run = func(_ *cobra.Command, _ []string) {
		run(cmd)
	}
}

// SetRunFuncWithNameArg registers a command function with an optional name argument
func (rc *ResourceCmd) SetRunFuncWithNameArg(cmd func() error) {
	rc.Command.Run = func(_ *cobra.Command, args []string) {
		rc.NameArg = GetNameArg(args)
		run(cmd)
	}
}

func run(cmd func() error) {
	if err := cmd(); err != nil {
		logger.Critical("%s\n", err.Error())
		os.Exit(1)
	}
}
