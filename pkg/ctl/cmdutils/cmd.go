package cmdutils

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/outposts"
)

// Cmd holds attributes that are common between commands;
// not all commands use each attribute, but they can if needed
type Cmd struct {
	CobraCommand *cobra.Command
	FlagSetGroup *NamedFlagSetGroup

	Plan, Wait, Validate bool

	NameArg string

	ClusterConfigFile string

	ProviderConfig api.ProviderConfig
	ClusterConfig  *api.ClusterConfig

	Include, Exclude []string
}

// NewCtl performs common defaulting and validation and constructs a new
// instance of eks.ClusterProvider, it may return an error if configuration
// is invalid or region is not supported
func (c *Cmd) NewCtl() (*eks.ClusterProvider, error) {
	if err := c.InitializeClusterConfig(); err != nil {
		return nil, err
	}
	ctl, err := eks.New(context.TODO(), &c.ProviderConfig, c.ClusterConfig)
	if err != nil {
		return nil, err
	}

	if !ctl.IsSupportedRegion() {
		return nil, ErrUnsupportedRegion(&c.ProviderConfig)
	}

	return ctl, nil
}

// InitializeClusterConfig validates and initializes the ClusterConfig.
func (c *Cmd) InitializeClusterConfig() error {
	api.SetClusterConfigDefaults(c.ClusterConfig)

	if err := api.ValidateClusterConfig(c.ClusterConfig); err != nil {
		if c.Validate {
			return err
		}
		logger.Warning("ignoring validation error: %s", err.Error())
	}

	for i, ng := range c.ClusterConfig.NodeGroups {
		if err := api.ValidateNodeGroup(i, ng, c.ClusterConfig); err != nil {
			if c.Validate {
				return err
			}
			logger.Warning("ignoring validation error: %s", err.Error())
		}
		// defaulting of nodegroup currently depends on validation;
		// that may change, but at present that's how it's meant to work
		api.SetNodeGroupDefaults(ng, c.ClusterConfig.Metadata, c.ClusterConfig.IsControlPlaneOnOutposts())
	}

	for i, ng := range c.ClusterConfig.ManagedNodeGroups {
		api.SetManagedNodeGroupDefaults(ng, c.ClusterConfig.Metadata, c.ClusterConfig.IsControlPlaneOnOutposts())
		if err := api.ValidateManagedNodeGroup(i, ng); err != nil {
			return err
		}
	}
	return nil
}

// NewProviderForExistingCluster is a wrapper for NewCtl that also validates that the cluster exists and is not a
// registered/connected cluster.
func (c *Cmd) NewProviderForExistingCluster(ctx context.Context) (*eks.ClusterProvider, error) {
	return c.NewProviderForExistingClusterHelper(ctx, func(ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
		return nil
	})
}

// NewProviderForExistingClusterHelper allows formating cluster K8s version to a standard value before doing nodegroup validations and initializations
func (c *Cmd) NewProviderForExistingClusterHelper(ctx context.Context, standardizeClusterVersionFormat func(ctl *eks.ClusterProvider, meta *api.ClusterMeta) error) (*eks.ClusterProvider, error) {
	clusterProvider, err := eks.New(ctx, &c.ProviderConfig, c.ClusterConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create cluster provider from options: %w", err)
	}
	if !clusterProvider.IsSupportedRegion() {
		return nil, ErrUnsupportedRegion(&c.ProviderConfig)
	}
	if err := clusterProvider.RefreshClusterStatus(ctx, c.ClusterConfig); err != nil {
		return nil, err
	}

	if err := standardizeClusterVersionFormat(clusterProvider, c.ClusterConfig.Metadata); err != nil {
		return nil, err
	}

	if err := c.InitializeClusterConfig(); err != nil {
		return nil, err
	}
	if c.ClusterConfig.IsControlPlaneOnOutposts() {
		clusterProvider.AWSProvider = outposts.WrapClusterProvider(clusterProvider.AWSProvider)
	}
	return clusterProvider, nil
}

// AddResourceCmd create a registers a new command under the given verb command
func AddResourceCmd(flagGrouping *FlagGrouping, parentVerbCmd *cobra.Command, newCmd func(*Cmd)) {
	c := &Cmd{
		CobraCommand: &cobra.Command{},
		ProviderConfig: api.ProviderConfig{
			WaitTimeout: api.DefaultWaitTimeout,
		},

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
