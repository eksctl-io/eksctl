package cmdutils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// Cmd holds attributes that are common between commands;
// not all commands use each attribute, but but they can if needed
type Cmd struct {
	CobraCommand *cobra.Command
	FlagSetGroup *NamedFlagSetGroup

	Plan, Wait, Validate bool

	NameArg string

	ClusterConfigFile string

	ProviderConfig *api.ProviderConfig
	ClusterConfig  *api.ClusterConfig

	Include, Exclude []string
}

// NewCtl performs common defaulting and validation and constructs a new
// instance of eks.ClusterProvider, it may return an error if configuration
// is invalid or region is not supported
func (c *Cmd) NewCtl() (*eks.ClusterProvider, error) {
	api.SetClusterConfigDefaults(c.ClusterConfig)

	if err := api.ValidateClusterConfig(c.ClusterConfig); err != nil {
		if c.Validate {
			return nil, err
		}
		logger.Warning("ignoring validation error: %s", err.Error())
	}

	for i, ng := range c.ClusterConfig.NodeGroups {
		if err := api.ValidateNodeGroup(i, ng); err != nil {
			if c.Validate {
				return nil, err
			}
			logger.Warning("ignoring validation error: %s", err.Error())
		}
		// defaulting of nodegroup currently depends on validation;
		// that may change, but at present that's how it's meant to work
		api.SetNodeGroupDefaults(i, ng)
	}

	ctl := eks.New(c.ProviderConfig, c.ClusterConfig)

	if !ctl.IsSupportedRegion() {
		return nil, ErrUnsupportedRegion(c.ProviderConfig)
	}

	return ctl, nil
}

// AddResourceCmd create a registers a new command under the given verb command
func AddResourceCmd(flagGrouping *FlagGrouping, parentVerbCmd *cobra.Command, newCmd func(*Cmd)) {
	c := &Cmd{
		CobraCommand:   &cobra.Command{},
		ProviderConfig: &api.ProviderConfig{},

		Plan:     true,  // always on by default
		Wait:     false, // varies in some commands
		Validate: true,  // also on by default
	}
	c.FlagSetGroup = flagGrouping.New(c.CobraCommand)
	newCmd(c)
	c.FlagSetGroup.AddTo(c.CobraCommand)
	parentVerbCmd.AddCommand(c.CobraCommand)
}

// SetDescription sets usage along with short and long descriptions as well as aliases
func (c *Cmd) SetDescription(use, short, long string, aliases ...string) {
	c.CobraCommand.Use = use
	c.CobraCommand.Short = short
	c.CobraCommand.Long = long
	c.CobraCommand.Aliases = aliases
}

// SetRunFunc registers a command function
func (c *Cmd) SetRunFunc(cmd func() error) {
	c.CobraCommand.Run = func(_ *cobra.Command, _ []string) {
		run(cmd)
	}
}

// SetRunFuncWithNameArg registers a command function with an optional name argument
func (c *Cmd) SetRunFuncWithNameArg(cmd func() error) {
	c.CobraCommand.Run = func(_ *cobra.Command, args []string) {
		c.NameArg = GetNameArg(args)
		run(cmd)
	}
}

func run(cmd func() error) {
	if err := cmd(); err != nil {
		logger.Critical("%s\n", err.Error())
		os.Exit(1)
	}
}
