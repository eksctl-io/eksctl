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
	*Cmd

	flagsIncompatibleWithConfigFile, flagsIncompatibleWithoutConfigFile sets.String

	validateWithConfigFile, validateWithoutConfigFile func() error
}

func newCommonClusterConfigLoader(cmd *Cmd) *commonClusterConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &commonClusterConfigLoader{
		Cmd: cmd,

		validateWithConfigFile: nilValidatorFunc,
		flagsIncompatibleWithConfigFile: sets.NewString(
			"name",
			"region",
			"version",
			"cluster",
			"namepace",
		),
		validateWithoutConfigFile: nilValidatorFunc,
		flagsIncompatibleWithoutConfigFile: sets.NewString(
			"only",
			"include",
			"exclude",
			"only-missing",
		),
	}
}

// Load ClusterConfig or use flags
func (l *commonClusterConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.ClusterConfigFile == "" {
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.CobraCommand.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.validateWithoutConfigFile()
	}

	var err error

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
	if l.ClusterConfig, err = eks.LoadConfigFromFile(l.ClusterConfigFile); err != nil {
		return err
	}
	meta := l.ClusterConfig.Metadata

	if meta == nil {
		return ErrMustBeSet("metadata")
	}

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.CobraCommand.Flag(f); flag != nil && flag.Changed {
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

func (l *commonClusterConfigLoader) validateMetadataWithoutConfigFile() error {
	meta := l.ClusterConfig.Metadata

	if meta.Name != "" && l.NameArg != "" {
		return ErrClusterFlagAndArg(l.Cmd, meta.Name, l.NameArg)
	}

	if l.NameArg != "" {
		meta.Name = l.NameArg
	}

	if meta.Name == "" {
		return ErrMustBeSet(ClusterNameFlag(l.Cmd))
	}

	return nil
}

// NewMetadataLoader handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fields, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func NewMetadataLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithoutConfigFile = l.validateMetadataWithoutConfigFile

	return l
}

// NewEnableProfileLoader handles loading of clusterConfigFile vs using flags for gitops commands
func NewEnableProfileLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithoutConfigFile = func() error {
		meta := l.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet("--cluster")
		}
		return nil
	}

	return l
}

// NewInstallFluxLoader handles loading of clusterConfigFile vs using flags for install commands
func NewInstallFluxLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithoutConfigFile = func() error {
		meta := l.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet("--cluster")
		}
		return nil
	}

	return l
}

// NewCreateClusterLoader will load config or use flags for 'eksctl create cluster'
func NewCreateClusterLoader(cmd *Cmd, ngFilter *NodeGroupFilter, ng *api.NodeGroup, withoutNodeGroup bool) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	ngFilter.ExcludeAll = withoutNodeGroup

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
		"vpc-private-subnets",
		"vpc-public-subnets",
		"vpc-cidr",
		"vpc-nat-mode",
		"vpc-from-kops-cluster",
	)

	l.flagsIncompatibleWithoutConfigFile.Insert("install-vpc-controllers")

	l.validateWithConfigFile = func() error {
		if l.ClusterConfig.VPC == nil {
			l.ClusterConfig.VPC = api.NewClusterVPC()
		}

		if l.ClusterConfig.VPC.NAT == nil {
			l.ClusterConfig.VPC.NAT = api.DefaultClusterNAT()
		}

		if !api.IsSetAndNonEmptyString(l.ClusterConfig.VPC.NAT.Gateway) {
			*l.ClusterConfig.VPC.NAT.Gateway = api.ClusterSingleNAT
		}

		if l.ClusterConfig.VPC.ClusterEndpoints == nil {
			l.ClusterConfig.VPC.ClusterEndpoints = api.ClusterEndpointAccessDefaults()
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
			return ErrClusterFlagAndArg(l.Cmd, meta.Name, l.NameArg)
		}
		meta.Name = ClusterName(meta.Name, l.NameArg)

		if l.ClusterConfig.Status != nil {
			return fmt.Errorf("status fields are read-only")
		}

		// prevent creation of invalid config object with irrelevant nodegroup
		// that may or may not be constructed correctly
		if !withoutNodeGroup {
			l.ClusterConfig.NodeGroups = append(l.ClusterConfig.NodeGroups, ng)
		}

		return ngFilter.ForEach(l.ClusterConfig.NodeGroups, func(i int, ng *api.NodeGroup) error {
			// generate nodegroup name or use flag
			ng.Name = NodeGroupName(ng.Name, "")
			return normalizeNodeGroup(ng, l)
		})
	}

	return l
}

// NewCreateNodeGroupLoader will load config or use flags for 'eksctl create nodegroup'
func NewCreateNodeGroupLoader(cmd *Cmd, ngFilter *NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
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
		return ngFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.NodeGroups)
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		return ngFilter.ForEach(l.ClusterConfig.NodeGroups, func(i int, ng *api.NodeGroup) error {
			// generate nodegroup name or use either flag or argument
			ngName := NodeGroupName(ng.Name, l.NameArg)
			if ngName == "" {
				return ErrClusterFlagAndArg(l.Cmd, ng.Name, l.NameArg)
			}
			ng.Name = ngName
			return normalizeNodeGroup(ng, l)
		})
	}

	return l
}

func normalizeNodeGroup(ng *api.NodeGroup, l *commonClusterConfigLoader) error {
	if flag := l.CobraCommand.Flag("ssh-public-key"); flag != nil && flag.Changed {
		if *ng.SSH.PublicKeyPath == "" {
			return fmt.Errorf("--ssh-public-key must be non-empty string")
		}
		ng.SSH.Allow = api.Enabled()
	} else {
		ng.SSH.PublicKeyPath = nil
	}

	if *ng.VolumeType == api.NodeVolumeTypeIO1 {
		return fmt.Errorf("%s volume type is not supported via flag --node-volume-type, please use a config file", api.NodeVolumeTypeIO1)
	}

	return nil
}

// NewDeleteNodeGroupLoader will load config or use flags for 'eksctl delete nodegroup'
func NewDeleteNodeGroupLoader(cmd *Cmd, ng *api.NodeGroup, ngFilter *NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		return ngFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.NodeGroups)
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"approve",
	)

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		if ng.Name != "" && l.NameArg != "" {
			return ErrClusterFlagAndArg(l.Cmd, ng.Name, l.NameArg)
		}

		if l.NameArg != "" {
			ng.Name = l.NameArg
		}

		if ng.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		ngFilter.AppendIncludeNames(ng.Name)

		l.Plan = false

		return nil
	}

	return l
}

// NewUtilsEnableLoggingLoader will load config or use flags for 'eksctl utils update-cluster-logging'
func NewUtilsEnableLoggingLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"enable-types",
		"disable-types",
	)

	l.validateWithoutConfigFile = l.validateMetadataWithoutConfigFile

	return l
}

// NewUtilsEnableEndpointAccessLoader will load config or use flags for 'eksctl utils vpc-cluster-api-access
func NewUtilsEnableEndpointAccessLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"private-access",
		"public-access",
	)

	l.validateWithoutConfigFile = l.validateMetadataWithoutConfigFile

	return l
}

// NewUtilsAssociateIAMOIDCProviderLoader will load config or use flags for 'eksctl utils associal-iam-oidc-provider'
func NewUtilsAssociateIAMOIDCProviderLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithoutConfigFile = func() error {
		l.ClusterConfig.IAM.WithOIDC = api.Enabled()
		return l.validateMetadataWithoutConfigFile()
	}

	l.validateWithConfigFile = func() error {
		if api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
			return fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile)
		}
		return nil
	}

	return l
}

// NewCreateIAMServiceAccountLoader will laod config or use flags for 'eksctl create iamserviceaccount'
func NewCreateIAMServiceAccountLoader(cmd *Cmd, saFilter *IAMServiceAccountFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"policy-arn",
	)

	l.validateWithConfigFile = func() error {
		return saFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.IAM.ServiceAccounts)
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		if len(l.ClusterConfig.IAM.ServiceAccounts) != 1 {
			return fmt.Errorf("unexpected number of service accounts")
		}

		serviceAccount := l.ClusterConfig.IAM.ServiceAccounts[0]

		if serviceAccount.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		if len(serviceAccount.AttachPolicyARNs) == 0 {
			return ErrMustBeSet("--attach-policy-arn")
		}

		return nil
	}

	return l
}

// NewGetIAMServiceAccountLoader will load config or use flags for 'eksctl get iamserviceaccount'
func NewGetIAMServiceAccountLoader(cmd *Cmd, sa *api.ClusterIAMServiceAccount) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		if api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
			return fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile)
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		sa.AttachPolicyARNs = []string{""} // force to pass general validation

		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		if l.NameArg != "" {
			sa.Name = l.NameArg
		}

		if sa.Name == "" {
			l.ClusterConfig.IAM.ServiceAccounts = nil
		}

		l.Plan = false

		return nil
	}

	return l
}

// NewDeleteIAMServiceAccountLoader will load config or use flags for 'eksctl delete iamserviceaccount'
func NewDeleteIAMServiceAccountLoader(cmd *Cmd, sa *api.ClusterIAMServiceAccount, saFilter *IAMServiceAccountFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		if api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
			return fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile)
		}
		return saFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.IAM.ServiceAccounts)
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"approve",
	)

	l.validateWithoutConfigFile = func() error {
		sa.AttachPolicyARNs = []string{""} // force to pass general validation

		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet("--cluster")
		}

		if sa.Name != "" && l.NameArg != "" {
			return ErrClusterFlagAndArg(l.Cmd, sa.Name, l.NameArg)
		}

		if l.NameArg != "" {
			sa.Name = l.NameArg
		}

		if sa.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		l.Plan = false

		return nil
	}

	return l
}
