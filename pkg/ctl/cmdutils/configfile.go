package cmdutils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// AddConfigFileFlag adds common --config-file flag
func AddConfigFileFlag(path *string, fs *pflag.FlagSet) {
	fs.StringVarP(path, "config-file", "f", "", "load configuration from a file")
}

// ClusterConfigLoader is an inteface that loaders should implement
type ClusterConfigLoader interface {
	Load() error
}

type commonClusterConfigLoader struct {
	provider *api.ProviderConfig
	path     string
	cmd      *cobra.Command
	spec     *api.ClusterConfig
	nameArg  string

	flagsIncompatibleWithConfigFile, flagsIncompatibleWithoutConfigFile sets.String

	validateWithConfigFile, validateWithoutConfigFile func() error
}

func newCommonClusterConfigLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile string, cmd *cobra.Command) *commonClusterConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &commonClusterConfigLoader{
		provider: provider,
		path:     clusterConfigFile,
		cmd:      cmd,
		spec:     spec,

		validateWithConfigFile:             nilValidatorFunc,
		flagsIncompatibleWithConfigFile:    sets.NewString("name", "region", "version"),
		validateWithoutConfigFile:          nilValidatorFunc,
		flagsIncompatibleWithoutConfigFile: sets.NewString(),
	}
}

// Load ClusterConfig or use flags
func (l *commonClusterConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.path == "" {
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.cmd.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.validateWithoutConfigFile()
	}

	if err := eks.LoadConfigFromFile(l.path, l.spec); err != nil {
		return err
	}
	meta := l.spec.Metadata

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.cmd.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	if l.nameArg != "" {
		return ErrCannotUseWithConfigFile(fmt.Sprintf("name argument %q", l.nameArg))
	}

	if meta.Name == "" {
		return ErrMustBeSet("metadata.name")
	}

	if meta.Region == "" {
		return ErrMustBeSet("metadata.region")
	}
	l.provider.Region = meta.Region

	return l.validateWithConfigFile()
}

// NewMetadataLoader handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fileds, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func NewMetadataLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile, nameArg string, cmd *cobra.Command) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(provider, spec, clusterConfigFile, cmd)

	l.nameArg = nameArg

	l.validateWithoutConfigFile = func() error {
		meta := l.spec.Metadata

		if meta.Name != "" && l.nameArg != "" {
			return ErrNameFlagAndArg(meta.Name, l.nameArg)
		}

		if l.nameArg != "" {
			meta.Name = l.nameArg
		}

		if meta.Name == "" {
			return ErrMustBeSet("--name")
		}

		return nil
	}

	return l
}

// NewCreateClusterLoader will laod config or use flags for 'eksctl create cluster'
func NewCreateClusterLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile, nameArg string, cmd *cobra.Command, ngFilter *NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(provider, spec, clusterConfigFile, cmd)

	l.nameArg = nameArg

	l.flagsIncompatibleWithConfigFile.Insert(
		"tags",
		"zones",
		"nodes",
		"nodes-min",
		"nodes-max",
		"node-type",
		"node-volume-size",
		"node-volume-type",
		"max-pods-per-node",
		"node-ami",
		"node-ami-family",
		"ssh-access",
		"ssh-public-key",
		"node-private-networking",
		"node-security-groups",
		"node-labels",
		"node-zones",
		"asg-access",
		"external-dns-access",
		"full-ecr-access",
		"storage-class",
		"vpc-private-subnets",
		"vpc-public-subnets",
		"vpc-cidr",
		"vpc-from-kops-cluster",
	)

	l.validateWithConfigFile = func() error {
		if l.spec.VPC == nil {
			l.spec.VPC = api.NewClusterVPC()
		}

		if l.spec.HasAnySubnets() && len(l.spec.AvailabilityZones) != 0 {
			return fmt.Errorf("vpc.subnets and availabilityZones cannot be set at the same time")
		}

		return nil
	}

	l.validateWithoutConfigFile = func() error {
		meta := l.spec.Metadata

		// generate cluster name or use either flag or argument
		if ClusterName(meta.Name, l.nameArg) == "" {
			return ErrNameFlagAndArg(meta.Name, l.nameArg)
		}
		meta.Name = ClusterName(meta.Name, l.nameArg)

		if l.spec.Status != nil {
			return fmt.Errorf("status fields are read-only")
		}

		return ngFilter.ForEach(l.spec.NodeGroups, func(i int, ng *api.NodeGroup) error {
			sshPublicKey := cmd.Flag("ssh-public-key")
			if sshPublicKey != nil && sshPublicKey.Changed {
				if *ng.SSH.PublicKeyPath == "" {
					return fmt.Errorf("--ssh-public-key must be non-empty string")
				}
				ng.SSH.Allow = api.Enabled()
			} else {
				ng.SSH.PublicKeyPath = nil
			}

			// generate nodegroup name or use flag
			ng.Name = NodeGroupName(ng.Name, "")

			return nil
		})
	}

	return l
}

// NewCreateNodeGroupLoader will laod config or use flags for 'eksctl create nodegroup'
func NewCreateNodeGroupLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile, nameArg string, cmd *cobra.Command, ngFilter *NodeGroupFilter, include, exclude []string) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(provider, spec, clusterConfigFile, cmd)

	l.nameArg = nameArg

	l.flagsIncompatibleWithConfigFile.Insert(
		"cluster",
		"nodes",
		"nodes-min",
		"nodes-max",
		"node-type",
		"node-volume-size",
		"node-volume-type",
		"max-pods-per-node",
		"node-ami",
		"node-ami-family",
		"ssh-access",
		"ssh-public-key",
		"node-private-networking",
		"node-security-groups",
		"node-labels",
		"node-zones",
		"asg-access",
		"external-dns-access",
		"full-ecr-access",
	)

	l.validateWithConfigFile = func() error {
		if err := ngFilter.AppendGlobs(include, exclude, spec.NodeGroups); err != nil {
			return err
		}
		return nil
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"only",
		"include",
		"exclude",
	)

	l.validateWithoutConfigFile = func() error {
		if spec.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		return ngFilter.ForEach(spec.NodeGroups, func(i int, ng *api.NodeGroup) error {
			sshPublicKey := cmd.Flag("ssh-public-key")
			if sshPublicKey != nil && sshPublicKey.Changed {
				if *ng.SSH.PublicKeyPath == "" {
					return fmt.Errorf("--ssh-public-key must be non-empty string")
				}
				ng.SSH.Allow = api.Enabled()
			} else {
				ng.SSH.PublicKeyPath = nil
			}

			// generate nodegroup name or use either flag or argument
			if NodeGroupName(ng.Name, l.nameArg) == "" {
				return ErrNameFlagAndArg(ng.Name, l.nameArg)
			}
			ng.Name = NodeGroupName(ng.Name, l.nameArg)

			return nil
		})
	}

	return l
}

// NewDeleteNodeGroupLoader will laod config or use flags for 'eksctl delete nodegroup'
func NewDeleteNodeGroupLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, ng *api.NodeGroup, clusterConfigFile, nameArg string, cmd *cobra.Command, ngFilter *NodeGroupFilter, include, exclude []string, plan *bool) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(provider, spec, clusterConfigFile, cmd)

	l.nameArg = nameArg

	l.flagsIncompatibleWithConfigFile.Insert(
		"cluster",
	)

	l.validateWithConfigFile = func() error {
		return ngFilter.AppendGlobs(include, exclude, spec.NodeGroups)
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"only",
		"include",
		"exclude",
		"only-missing",
		"approve",
	)

	l.validateWithoutConfigFile = func() error {
		if l.spec.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		if ng.Name != "" && nameArg != "" {
			return ErrNameFlagAndArg(ng.Name, nameArg)
		}

		if nameArg != "" {
			ng.Name = nameArg
		}

		if ng.Name == "" {
			return ErrMustBeSet("--name")
		}

		ngFilter.AppendIncludeNames(ng.Name)

		*plan = false

		return nil
	}

	return l
}
