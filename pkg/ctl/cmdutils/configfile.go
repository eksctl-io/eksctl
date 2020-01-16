package cmdutils

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/names"
)

// AddConfigFileFlag adds common --config-file flag
func AddConfigFileFlag(fs *pflag.FlagSet, path *string) {
	fs.StringVarP(path, "config-file", "f", "", "load configuration from a file (or stdin if set to '-')")
}

// ClusterConfigLoader is an interface that loaders should implement
type ClusterConfigLoader interface {
	Load() error
}

type commonClusterConfigLoader struct {
	*Cmd

	flagsIncompatibleWithConfigFile, flagsIncompatibleWithoutConfigFile sets.String

	validateWithConfigFile, validateWithoutConfigFile func() error
}

var (
	defaultFlagsIncompatibleWithConfigFile = sets.NewString(
		"name",
		"region",
		"version",
		"cluster",
		"namespace",
	)
	defaultFlagsIncompatibleWithoutConfigFile = sets.NewString(
		"only",
		"include",
		"exclude",
		"only-missing",
	)
)

var multiErr *multierror.Error

func newCommonClusterConfigLoader(cmd *Cmd) *commonClusterConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &commonClusterConfigLoader{
		Cmd: cmd,

		validateWithConfigFile:             nilValidatorFunc,
		flagsIncompatibleWithConfigFile:    defaultFlagsIncompatibleWithConfigFile,
		validateWithoutConfigFile:          nilValidatorFunc,
		flagsIncompatibleWithoutConfigFile: defaultFlagsIncompatibleWithoutConfigFile,
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
				multiErr = multierror.Append(fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f))
				return multiErr.ErrorOrNil()
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
		multiErr = multierror.Append(ErrMustBeSet("metadata"))
	}

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.CobraCommand.Flag(f); flag != nil && flag.Changed {
			multiErr = multierror.Append(ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f)))
		}
	}

	if l.NameArg != "" {
		multiErr = multierror.Append(ErrCannotUseWithConfigFile(fmt.Sprintf("name argument %q", l.NameArg)))
	}

	if meta.Name == "" {
		multiErr = multierror.Append(ErrMustBeSet("metadata.name"))
	}

	if meta.Region == "" {
		multiErr = multierror.Append(ErrMustBeSet("metadata.region"))
	}

	if multiErr != nil {
		return multiErr
	}

	l.ProviderConfig.Region = meta.Region

	return l.validateWithConfigFile()
}

func (l *commonClusterConfigLoader) validateMetadataWithoutConfigFile() error {
	meta := l.ClusterConfig.Metadata

	if meta.Name != "" && l.NameArg != "" {
		multiErr = multierror.Append(ErrClusterFlagAndArg(l.Cmd, meta.Name, l.NameArg))
	}

	if l.NameArg != "" {
		meta.Name = l.NameArg
	}

	if meta.Name == "" {
		multiErr = multierror.Append(ErrMustBeSet(ClusterNameFlag(l.Cmd)))
	}

	return multiErr.ErrorOrNil()
}

// NewMetadataLoader handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fields, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func NewMetadataLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithoutConfigFile = l.validateMetadataWithoutConfigFile

	return l
}

// NewCreateClusterLoader will load config or use flags for 'eksctl create cluster'
func NewCreateClusterLoader(cmd *Cmd, ngFilter *NodeGroupFilter, ng *api.NodeGroup, params *CreateClusterCmdParams) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	ngFilter.ExcludeAll = params.WithoutNodeGroup

	l.flagsIncompatibleWithConfigFile.Insert(
		"tags",
		"zones",
		"managed",
		"fargate",
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
			multiErr = multierror.Append(fmt.Errorf("vpc.subnets and availabilityZones cannot be set at the same time"))
		}

		return multiErr.ErrorOrNil()
	}

	l.validateWithoutConfigFile = func() error {
		meta := l.ClusterConfig.Metadata

		// generate cluster name or use either flag or argument
		if names.ForCluster(meta.Name, l.NameArg) == "" {
			multiErr = multierror.Append(ErrClusterFlagAndArg(l.Cmd, meta.Name, l.NameArg))
		}
		meta.Name = names.ForCluster(meta.Name, l.NameArg)

		if l.ClusterConfig.Status != nil {
			multiErr = multierror.Append(fmt.Errorf("status fields are read-only"))
		}

		if params.Managed {
			for _, f := range incompatibleManagedNodesFlags() {
				if flag := l.CobraCommand.Flag(f); flag != nil && flag.Changed {
					multiErr = multierror.Append(ErrUnsupportedManagedFlag(fmt.Sprintf("--%s", f)))
				}
			}
		}

		// prevent creation of invalid config object with irrelevant nodegroup
		// that may or may not be constructed correctly
		if !params.WithoutNodeGroup {
			if params.Managed {
				l.ClusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{makeManagedNodegroup(ng)}
			} else {
				l.ClusterConfig.NodeGroups = []*api.NodeGroup{ng}
			}
		}

		if params.Fargate {
			l.ClusterConfig.SetDefaultFargateProfile()
			// A Fargate-only cluster should NOT have any un-managed node group:
			l.ClusterConfig.NodeGroups = []*api.NodeGroup{}
		}

		for _, ng := range l.ClusterConfig.NodeGroups {
			// generate nodegroup name or use flag
			ng.Name = names.ForNodeGroup(ng.Name, "")
			if err := normalizeNodeGroup(ng, l); err != nil {
				return err
			}
		}

		for _, ng := range l.ClusterConfig.ManagedNodeGroups {
			ng.Name = names.ForNodeGroup(ng.Name, "")
		}

		return multiErr.ErrorOrNil()
	}

	return l
}

// NewCreateNodeGroupLoader will load config or use flags for 'eksctl create nodegroup'
func NewCreateNodeGroupLoader(cmd *Cmd, ng *api.NodeGroup, ngFilter *NodeGroupFilter, managedNodeGroup bool) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"managed",
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
		return ngFilter.AppendGlobs(l.Include, l.Exclude, getAllNodeGroupNames(l.ClusterConfig))
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			multiErr = multierror.Append(ErrMustBeSet(ClusterNameFlag(cmd)))
		}
		if managedNodeGroup {
			for _, f := range incompatibleManagedNodesFlags() {
				if flag := l.CobraCommand.Flag(f); flag != nil && flag.Changed {
					multiErr = multierror.Append(ErrUnsupportedManagedFlag(fmt.Sprintf("--%s", f)))
				}
			}
			l.ClusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{makeManagedNodegroup(ng)}
		} else {
			l.ClusterConfig.NodeGroups = []*api.NodeGroup{ng}
		}

		// Validate both filtered and unfiltered nodegroups
		if managedNodeGroup {
			for _, ng := range l.ClusterConfig.ManagedNodeGroups {
				ngName := names.ForNodeGroup(ng.Name, l.NameArg)
				if ngName == "" {
					multiErr = multierror.Append(ErrClusterFlagAndArg(l.Cmd, ng.Name, l.NameArg))
				}
				ng.Name = ngName
			}
		} else {
			for _, ng := range l.ClusterConfig.NodeGroups {
				// generate nodegroup name or use either flag or argument
				ngName := names.ForNodeGroup(ng.Name, l.NameArg)
				if ngName == "" {
					multiErr = multierror.Append(ErrClusterFlagAndArg(l.Cmd, ng.Name, l.NameArg))
				}
				ng.Name = ngName
				if err := normalizeNodeGroup(ng, l); err != nil {
					return err
				}
			}
		}
		return multiErr.ErrorOrNil()
	}

	return l
}

func makeManagedNodegroup(nodeGroup *api.NodeGroup) *api.ManagedNodeGroup {
	return &api.ManagedNodeGroup{
		AvailabilityZones: nodeGroup.AvailabilityZones,
		Name:              nodeGroup.Name,
		IAM:               nodeGroup.IAM,
		SSH:               nodeGroup.SSH,
		InstanceType:      nodeGroup.InstanceType,
		Labels:            nodeGroup.Labels,
		Tags:              nodeGroup.Tags,
		AMIFamily:         nodeGroup.AMIFamily,
		VolumeSize:        nodeGroup.VolumeSize,
		ScalingConfig: &api.ScalingConfig{
			MinSize:         nodeGroup.MinSize,
			MaxSize:         nodeGroup.MaxSize,
			DesiredCapacity: nodeGroup.DesiredCapacity,
		},
	}
}

func normalizeNodeGroup(ng *api.NodeGroup, l *commonClusterConfigLoader) error {
	if flag := l.CobraCommand.Flag("ssh-public-key"); flag != nil && flag.Changed {
		if *ng.SSH.PublicKeyPath == "" {
			multiErr = multierror.Append(fmt.Errorf("--ssh-public-key must be non-empty string"))
		}
		ng.SSH.Allow = api.Enabled()
	} else {
		ng.SSH.PublicKeyPath = nil
	}

	if *ng.VolumeType == api.NodeVolumeTypeIO1 {
		multiErr = multierror.Append(fmt.Errorf("%s volume type is not supported via flag --node-volume-type, please use a config file", api.NodeVolumeTypeIO1))
	}

	return multiErr.ErrorOrNil()
}

// NewDeleteNodeGroupLoader will load config or use flags for 'eksctl delete nodegroup'
func NewDeleteNodeGroupLoader(cmd *Cmd, ng *api.NodeGroup, ngFilter *NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		return ngFilter.AppendGlobs(l.Include, l.Exclude, getAllNodeGroupNames(l.ClusterConfig))
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"approve",
	)

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			multiErr = multierror.Append(ErrMustBeSet(ClusterNameFlag(cmd)))
		}

		if ng.Name != "" && l.NameArg != "" {
			multiErr = multierror.Append(ErrClusterFlagAndArg(l.Cmd, ng.Name, l.NameArg))
		}

		if l.NameArg != "" {
			ng.Name = l.NameArg
		}

		if ng.Name == "" {
			multiErr = multierror.Append(ErrMustBeSet("--name"))
		}

		ngFilter.AppendIncludeNames(ng.Name)

		l.Plan = false

		return multiErr.ErrorOrNil()
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
			multiErr = multierror.Append(fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile))
		}
		return multiErr.ErrorOrNil()
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
			multiErr = multierror.Append(ErrMustBeSet(ClusterNameFlag(cmd)))
		}

		if len(l.ClusterConfig.IAM.ServiceAccounts) != 1 {
			multiErr = multierror.Append(fmt.Errorf("unexpected number of service accounts"))
		}

		serviceAccount := l.ClusterConfig.IAM.ServiceAccounts[0]

		if serviceAccount.Name == "" {
			multiErr = multierror.Append(ErrMustBeSet(ClusterNameFlag(cmd)))
		}

		if len(serviceAccount.AttachPolicyARNs) == 0 {
			multiErr = multierror.Append(ErrMustBeSet("--attach-policy-arn"))
		}

		return multiErr.ErrorOrNil()
	}

	return l
}

// NewGetIAMServiceAccountLoader will load config or use flags for 'eksctl get iamserviceaccount'
func NewGetIAMServiceAccountLoader(cmd *Cmd, sa *api.ClusterIAMServiceAccount) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		if api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
			multiErr = multierror.Append(fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile))
			return multiErr.ErrorOrNil()
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		sa.AttachPolicyARNs = []string{""} // force to pass general validation

		if l.ClusterConfig.Metadata.Name == "" {
			multiErr = multierror.Append(ErrMustBeSet(ClusterNameFlag(cmd)))
		}

		if l.NameArg != "" {
			sa.Name = l.NameArg
		}

		if sa.Name == "" {
			l.ClusterConfig.IAM.ServiceAccounts = nil
		}

		l.Plan = false

		return multiErr.ErrorOrNil()
	}

	return l
}

// NewDeleteIAMServiceAccountLoader will load config or use flags for 'eksctl delete iamserviceaccount'
func NewDeleteIAMServiceAccountLoader(cmd *Cmd, sa *api.ClusterIAMServiceAccount, saFilter *IAMServiceAccountFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		if api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
			multiErr = multierror.Append(fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile))
			return multiErr.ErrorOrNil()
		}
		return saFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.IAM.ServiceAccounts)
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"approve",
	)

	l.validateWithoutConfigFile = func() error {
		sa.AttachPolicyARNs = []string{""} // force to pass general validation

		if l.ClusterConfig.Metadata.Name == "" {
			multiErr = multierror.Append(ErrMustBeSet(ClusterNameFlag(cmd)))
		}

		if sa.Name != "" && l.NameArg != "" {
			multiErr = multierror.Append(ErrClusterFlagAndArg(l.Cmd, sa.Name, l.NameArg))
		}

		if l.NameArg != "" {
			sa.Name = l.NameArg
		}

		if sa.Name == "" {
			multiErr = multierror.Append(ErrMustBeSet("--name"))
		}

		l.Plan = false

		return multiErr.ErrorOrNil()
	}

	return l
}
