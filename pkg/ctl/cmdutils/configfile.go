package cmdutils

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// AddConfigFileFlag adds common --config-file flag
func AddConfigFileFlag(fs *pflag.FlagSet, path *string) {
	fs.StringVarP(path, "config-file", "f", "", "load configuration from a file (or stdin if set to '-')")
}

// ClusterConfigLoader is an inteface that loaders should implement
type ClusterConfigLoader interface {
	Load() error
}

type commonClusterConfigLoader struct {
	*ResourceCmd

	flagsIncompatibleWithConfigFile, flagsIncompatibleWithoutConfigFile sets.String

	validateWithConfigFile, validateWithoutConfigFile func() error
}

func newCommonClusterConfigLoader(rc *ResourceCmd) *commonClusterConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &commonClusterConfigLoader{
		ResourceCmd: rc,

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

	if l.ClusterConfigFile == "" {
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.Command.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.validateWithoutConfigFile()
	}

	if err := eks.LoadConfigFromFile(l.ClusterConfigFile, l.ClusterConfig); err != nil {
		return err
	}
	meta := l.ClusterConfig.Metadata

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.Command.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	if l.NameArg != "" {
		return ErrCannotUseWithConfigFile(fmt.Sprintf("name argument %q", l.NameArg))
	}

	if meta.Name == "" {
		return ErrMustBeSet("metadata.name")
	}

	if meta.Region == "" {
		return ErrMustBeSet("metadata.region")
	}
	l.ProviderConfig.Region = meta.Region

	return l.validateWithConfigFile()
}

// NewMetadataLoader handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fileds, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func NewMetadataLoader(rc *ResourceCmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(rc)

	l.validateWithoutConfigFile = func() error {
		meta := l.ClusterConfig.Metadata

		if meta.Name != "" && l.NameArg != "" {
			return ErrNameFlagAndArg(meta.Name, l.NameArg)
		}

		if l.NameArg != "" {
			meta.Name = l.NameArg
		}

		if meta.Name == "" {
			return ErrMustBeSet("--name")
		}

		return nil
	}

	return l
}

// NewCreateClusterLoader will laod config or use flags for 'eksctl create cluster'
func NewCreateClusterLoader(rc *ResourceCmd, ngFilter *NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(rc)

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
		"vpc-nat-mode",
		"vpc-from-kops-cluster",
	)

	l.validateWithConfigFile = func() error {
		if l.ClusterConfig.VPC == nil {
			l.ClusterConfig.VPC = api.NewClusterVPC()
		}

		if l.ClusterConfig.HasAnySubnets() && len(l.ClusterConfig.AvailabilityZones) != 0 {
			return fmt.Errorf("vpc.subnets and availabilityZones cannot be set at the same time")
		}

		return nil
	}

	l.validateWithoutConfigFile = func() error {
		meta := l.ClusterConfig.Metadata

		// generate cluster name or use either flag or argument
		if ClusterName(meta.Name, l.NameArg) == "" {
			return ErrNameFlagAndArg(meta.Name, l.NameArg)
		}
		meta.Name = ClusterName(meta.Name, l.NameArg)

		if l.ClusterConfig.Status != nil {
			return fmt.Errorf("status fields are read-only")
		}

		return ngFilter.ForEach(l.ClusterConfig.NodeGroups, func(i int, ng *api.NodeGroup) error {
			if l.Command.Flag("ssh-public-key").Changed {
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
func NewCreateNodeGroupLoader(rc *ResourceCmd, ngFilter *NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(rc)

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
		if err := ngFilter.AppendGlobs(l.IncludeNodeGroups, l.ExcludeNodeGroups, l.ClusterConfig.NodeGroups); err != nil {
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
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		return ngFilter.ForEach(l.ClusterConfig.NodeGroups, func(i int, ng *api.NodeGroup) error {

			if l.Command.Flag("ssh-public-key").Changed {
				if *ng.SSH.PublicKeyPath == "" {
					return fmt.Errorf("--ssh-public-key must be non-empty string")
				}
				ng.SSH.Allow = api.Enabled()
			} else {
				ng.SSH.PublicKeyPath = nil
			}

			// generate nodegroup name or use either flag or argument
			if NodeGroupName(ng.Name, l.NameArg) == "" {
				return ErrNameFlagAndArg(ng.Name, l.NameArg)
			}
			ng.Name = NodeGroupName(ng.Name, l.NameArg)

			return nil
		})
	}

	return l
}

// NewDeleteNodeGroupLoader will laod config or use flags for 'eksctl delete nodegroup'
func NewDeleteNodeGroupLoader(rc *ResourceCmd, ng *api.NodeGroup, ngFilter *NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(rc)

	l.flagsIncompatibleWithConfigFile.Insert(
		"cluster",
	)

	l.validateWithConfigFile = func() error {
		return ngFilter.AppendGlobs(l.IncludeNodeGroups, l.ExcludeNodeGroups, l.ClusterConfig.NodeGroups)
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"only",
		"include",
		"exclude",
		"only-missing",
		"approve",
	)

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		if ng.Name != "" && l.NameArg != "" {
			return ErrNameFlagAndArg(ng.Name, l.NameArg)
		}

		if l.NameArg != "" {
			ng.Name = l.NameArg
		}

		if ng.Name == "" {
			return ErrMustBeSet("--name")
		}

		ngFilter.AppendIncludeNames(ng.Name)

		l.Plan = false

		return nil
	}

	return l
}
