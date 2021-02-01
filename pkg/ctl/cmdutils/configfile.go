package cmdutils

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
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

	flagsIncompatibleWithConfigFile    sets.String
	flagsIncompatibleWithoutConfigFile sets.String
	validateWithConfigFile             func() error
	validateWithoutConfigFile          func() error
}

var (
	defaultFlagsIncompatibleWithConfigFile = [...]string{
		"name",
		"region",
		"version",
		"cluster",
		"namepace",
	}
	defaultFlagsIncompatibleWithoutConfigFile = [...]string{
		"only",
		"include",
		"exclude",
		"only-missing",
	}
)

func newCommonClusterConfigLoader(cmd *Cmd) *commonClusterConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &commonClusterConfigLoader{
		Cmd: cmd,

		validateWithConfigFile:             nilValidatorFunc,
		flagsIncompatibleWithConfigFile:    sets.NewString(defaultFlagsIncompatibleWithConfigFile[:]...),
		validateWithoutConfigFile:          nilValidatorFunc,
		flagsIncompatibleWithoutConfigFile: sets.NewString(defaultFlagsIncompatibleWithoutConfigFile[:]...),
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

	if l.flagsIncompatibleWithConfigFile.Has("name") && l.NameArg != "" {
		return ErrCannotUseWithConfigFile("name argument")
	}

	if meta.Name == "" {
		return ErrMustBeSet("metadata.name")
	}

	if meta.Region == "" {
		return ErrMustBeSet("metadata.region")
	}
	l.ProviderConfig.Region = meta.Region

	api.SetDefaultGitSettings(l.ClusterConfig)
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

// NewCreateClusterLoader will load config or use flags for 'eksctl create cluster'
func NewCreateClusterLoader(cmd *Cmd, ngFilter *filter.NodeGroupFilter, ng *api.NodeGroup, params *CreateClusterCmdParams) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	ngFilter.SetExcludeAll(params.WithoutNodeGroup)

	l.flagsIncompatibleWithConfigFile.Insert(
		"tags",
		"zones",
		"managed",
		"fargate",
		"spot",
		"instance-types",
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
		"enable-ssm",
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
		clusterConfig := l.ClusterConfig
		if clusterConfig.VPC == nil {
			clusterConfig.VPC = api.NewClusterVPC()
		}

		if clusterConfig.VPC.NAT == nil {
			clusterConfig.VPC.NAT = api.DefaultClusterNAT()
		}

		if !api.IsSetAndNonEmptyString(clusterConfig.VPC.NAT.Gateway) {
			*clusterConfig.VPC.NAT.Gateway = api.ClusterSingleNAT
		}

		if clusterConfig.PrivateCluster != nil && clusterConfig.PrivateCluster.Enabled {
			if clusterEndpoints := clusterConfig.VPC.ClusterEndpoints; clusterEndpoints != nil && (clusterEndpoints.PublicAccess != nil || clusterEndpoints.PrivateAccess != nil) {
				return errors.New("vpc.clusterEndpoints cannot be set for a fully-private cluster (privateCluster.enabled) as the endpoint access defaults to private-only")
			}
		}

		api.SetClusterEndpointAccessDefaults(clusterConfig.VPC)

		if !clusterConfig.HasClusterEndpointAccess() {
			return api.ErrClusterEndpointNoAccess
		}

		if clusterConfig.HasAnySubnets() && len(clusterConfig.AvailabilityZones) != 0 {
			return fmt.Errorf("vpc.subnets and availabilityZones cannot be set at the same time")
		}

		if clusterConfig.Git != nil {
			repo := clusterConfig.Git.Repo
			if repo == nil || repo.URL == "" {
				return ErrMustBeSet("git.repo.URL")
			}

			if repo.Email == "" {
				return ErrMustBeSet("git.repo.email")
			}

			profile := clusterConfig.Git.BootstrapProfile
			if profile != nil && profile.Source == "" {
				return ErrMustBeSet("git.bootstrapProfile.source")
			}
		}

		return nil
	}

	l.validateWithoutConfigFile = func() error {
		meta := l.ClusterConfig.Metadata

		// generate cluster name or use either flag or argument
		if names.ForCluster(meta.Name, l.NameArg) == "" {
			return ErrClusterFlagAndArg(l.Cmd, meta.Name, l.NameArg)
		}
		meta.Name = names.ForCluster(meta.Name, l.NameArg)

		if l.ClusterConfig.Status != nil {
			return fmt.Errorf("status fields are read-only")
		}

		if err := validateManagedNGFlags(l.CobraCommand, params.Managed); err != nil {
			return err
		}

		// prevent creation of invalid config object with irrelevant nodegroup
		// that may or may not be constructed correctly
		if !params.WithoutNodeGroup {
			if params.Managed {
				l.ClusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{makeManagedNodegroup(ng, params.CreateManagedNGOptions)}
			} else {
				l.ClusterConfig.NodeGroups = []*api.NodeGroup{ng}
			}
		}

		api.SetClusterEndpointAccessDefaults(l.ClusterConfig.VPC)

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

		return nil
	}

	return l
}

// NewCreateNodeGroupLoader will load config or use flags for 'eksctl create nodegroup'
func NewCreateNodeGroupLoader(cmd *Cmd, ng *api.NodeGroup, ngFilter *filter.NodeGroupFilter, mngOptions CreateManagedNGOptions) ClusterConfigLoader {
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
		"spot",
		"instance-types",
		"ssh-access",
		"ssh-public-key",
		"enable-ssm",
		"node-private-networking",
		"node-security-groups",
		"node-labels",
		"node-zones",
		"asg-access",
		"external-dns-access",
		"full-ecr-access",
	)

	l.validateWithConfigFile = func() error {
		return ngFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.GetAllNodeGroupNames())
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if err := validateManagedNGFlags(l.CobraCommand, mngOptions.Managed); err != nil {
			return err
		}
		if mngOptions.Managed {
			l.ClusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{makeManagedNodegroup(ng, mngOptions)}

		} else {
			l.ClusterConfig.NodeGroups = []*api.NodeGroup{ng}
		}

		// Validate both filtered and unfiltered nodegroups
		if mngOptions.Managed {
			for _, ng := range l.ClusterConfig.ManagedNodeGroups {
				ngName := names.ForNodeGroup(ng.Name, l.NameArg)
				if ngName == "" {
					return ErrFlagAndArg("--name", ng.Name, l.NameArg)
				}
				ng.Name = ngName
			}
		} else {
			for _, ng := range l.ClusterConfig.NodeGroups {
				// generate nodegroup name or use either flag or argument
				ngName := names.ForNodeGroup(ng.Name, l.NameArg)
				if ngName == "" {
					return ErrFlagAndArg("--name", ng.Name, l.NameArg)
				}
				ng.Name = ngName
				if err := normalizeNodeGroup(ng, l); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return l
}

func makeManagedNodegroup(nodeGroup *api.NodeGroup, options CreateManagedNGOptions) *api.ManagedNodeGroup {
	ngBase := *nodeGroup.NodeGroupBase
	if ngBase.SecurityGroups != nil {
		ngBase.SecurityGroups = &api.NodeGroupSGs{
			AttachIDs: ngBase.SecurityGroups.AttachIDs,
		}
	}
	return &api.ManagedNodeGroup{
		NodeGroupBase: &ngBase,
		Spot:          options.Spot,
		InstanceTypes: options.InstanceTypes,
	}
}

func validateManagedNGFlags(cmd *cobra.Command, managed bool) error {
	if managed {
		for _, f := range incompatibleManagedNodesFlags() {
			if flag := cmd.Flag(f); flag != nil && flag.Changed {
				return ErrUnsupportedManagedFlag(fmt.Sprintf("--%s", f))
			}
		}
		return nil
	}
	flagsValidOnlyWithMNG := []string{"spot", "instance-types"}
	for _, f := range flagsValidOnlyWithMNG {
		if flag := cmd.Flag(f); flag != nil && flag.Changed {
			return errors.Errorf("--%s is only valid with managed nodegroups (--managed)", f)
		}
	}
	return nil
}

func normalizeNodeGroup(ng *api.NodeGroup, l *commonClusterConfigLoader) error {
	if flag := l.CobraCommand.Flag("ssh-public-key"); flag != nil && flag.Changed {
		if *ng.SSH.PublicKeyPath == "" {
			return fmt.Errorf("--ssh-public-key must be non-empty string")
		}
		if flag := l.CobraCommand.Flag("ssh-access"); flag == nil || !flag.Changed {
			ng.SSH.Allow = api.Enabled()
		}
	} else {
		ng.SSH.PublicKeyPath = nil
	}

	if *ng.VolumeType == api.NodeVolumeTypeIO1 {
		return fmt.Errorf("%s volume type is not supported via flag --node-volume-type, please use a config file", api.NodeVolumeTypeIO1)
	}

	return nil
}

// NewDeleteNodeGroupLoader will load config or use flags for 'eksctl delete nodegroup'
func NewDeleteNodeGroupLoader(cmd *Cmd, ng *api.NodeGroup, ngFilter *filter.NodeGroupFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		return ngFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.GetAllNodeGroupNames())
	}

	l.flagsIncompatibleWithoutConfigFile.Insert(
		"approve",
	)

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		if ng.Name != "" && l.NameArg != "" {
			return ErrFlagAndArg("--name", ng.Name, l.NameArg)
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

// NewUtilsEnableEndpointAccessLoader will load config or use flags for 'eksctl utils update-cluster-endpoints'.
func NewUtilsEnableEndpointAccessLoader(cmd *Cmd, privateAccess, publicAccess bool) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"private-access",
		"public-access",
	)
	l.validateWithoutConfigFile = func() error {
		if err := l.validateMetadataWithoutConfigFile(); err != nil {
			return err
		}

		if flag := l.CobraCommand.Flag("private-access"); flag != nil && flag.Changed {
			cmd.ClusterConfig.VPC.ClusterEndpoints.PrivateAccess = &privateAccess
		} else {
			cmd.ClusterConfig.VPC.ClusterEndpoints.PrivateAccess = nil
		}

		if flag := l.CobraCommand.Flag("public-access"); flag != nil && flag.Changed {
			cmd.ClusterConfig.VPC.ClusterEndpoints.PublicAccess = &publicAccess
		} else {
			cmd.ClusterConfig.VPC.ClusterEndpoints.PublicAccess = nil
		}

		return nil
	}
	l.validateWithConfigFile = func() error {
		if l.ClusterConfig.VPC == nil {
			l.ClusterConfig.VPC = api.NewClusterVPC()
		}
		api.SetClusterEndpointAccessDefaults(l.ClusterConfig.VPC)
		return nil
	}

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
		if l.ClusterConfig.IAM == nil || api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
			return fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile)
		}
		return nil
	}

	return l
}

// NewUtilsPublicAccessCIDRsLoader loads config or uses flags for `eksctl utils set-public-access-cidrs <cidrs>`
func NewUtilsPublicAccessCIDRsLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		if cmd.NameArg != "" {
			return fmt.Errorf("config file and CIDR list argument %s", IncompatibleFlags)
		}
		if l.ClusterConfig.VPC == nil || l.ClusterConfig.VPC.PublicAccessCIDRs == nil {
			return errors.New("field vpc.publicAccessCIDRs is required")
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if cmd.NameArg == "" {
			return errors.New("a comma-separated CIDR list is required")
		}

		cidrs, err := parseList(cmd.NameArg)
		if err != nil {
			return err
		}
		l.ClusterConfig.VPC.PublicAccessCIDRs = cidrs
		return nil
	}
	return l
}

func parseList(arg string) ([]string, error) {
	reader := strings.NewReader(arg)
	csvReader := csv.NewReader(reader)
	return csvReader.Read()
}

// NewCreateIAMServiceAccountLoader will load config or use flags for 'eksctl create iamserviceaccount'
func NewCreateIAMServiceAccountLoader(cmd *Cmd, saFilter *filter.IAMServiceAccountFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"policy-arn",
	)

	l.validateWithConfigFile = func() error {
		if l.ClusterConfig.IAM == nil || l.ClusterConfig.IAM.ServiceAccounts == nil {
			return fmt.Errorf("'iam.serviceAccounts' is not defined in %q", l.ClusterConfigFile)
		}
		return saFilter.AppendGlobs(l.Include, l.Exclude, l.ClusterConfig.IAM.ServiceAccounts)
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		if len(l.ClusterConfig.IAM.ServiceAccounts) != 1 {
			return fmt.Errorf("unexpected number of service accounts")
		}

		serviceAccount := l.ClusterConfig.IAM.ServiceAccounts[0]

		if serviceAccount.Name != "" && l.NameArg != "" {
			return ErrFlagAndArg("--name", serviceAccount.Name, l.NameArg)
		}

		if l.NameArg != "" {
			serviceAccount.Name = l.NameArg
		}

		if serviceAccount.Name == "" {
			return ErrMustBeSet("--name")
		}

		if len(serviceAccount.AttachPolicyARNs) == 0 {
			return ErrMustBeSet("--attach-policy-arn")
		}

		return nil
	}

	return l
}

// NewGetIAMServiceAccountLoader will load config or use flags for 'eksctl get iamserviceaccount'
func NewGetIAMServiceAccountLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		if api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
			return fmt.Errorf("'iam.withOIDC' is not enabled in %q", l.ClusterConfigFile)
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		l.Plan = false

		return nil
	}

	return l
}

// NewDeleteIAMServiceAccountLoader will load config or use flags for 'eksctl delete iamserviceaccount'
func NewDeleteIAMServiceAccountLoader(cmd *Cmd, sa *api.ClusterIAMServiceAccount, saFilter *filter.IAMServiceAccountFilter) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithConfigFile = func() error {
		if l.ClusterConfig.IAM == nil || api.IsDisabled(l.ClusterConfig.IAM.WithOIDC) {
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
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		if sa.Name != "" && l.NameArg != "" {
			return ErrFlagAndArg("--name", sa.Name, l.NameArg)
		}

		if l.NameArg != "" {
			sa.Name = l.NameArg
		}

		if sa.Name == "" {
			return ErrMustBeSet("--name")
		}

		l.Plan = false

		return nil
	}

	return l
}
