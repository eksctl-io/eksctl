package create

import (
	"github.com/kris-nova/logger"
	"github.com/lithammer/dedent"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/iam"
)

func createIAMIdentityMappingCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	id := &iam.Identity{}

	rc.SetDescription("iamidentitymapping", "Create an IAM identity mapping",
		dedent.Dedent(`Creates a mapping from IAM role or user to Kubernetes user and groups.

			Note aws-iam-authenticator only considers the last entry for any given
			role. If you create a duplicate entry it will shadow all the previous
			username and groups mapping.
		`),
	)

	rc.SetRunFunc(func() error {
		return doCreateIAMIdentityMapping(rc, id)
	})

	var arn iam.ARN
	var username string
	var groups []string

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.Var(&arn, "arn", "ARN of the IAM role or user to create")
		fs.StringVar(&username, "username", "", "User name within Kubernetes to map to IAM role")
		fs.StringArrayVar(&groups, "group", []string{}, "Group within Kubernetes to which IAM role is mapped")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
	})

	id, err := iam.NewIdentity(arn, username, groups)
	if err != nil {
		// What does one do here?
	}

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doCreateIAMIdentityMapping(rc *cmdutils.ResourceCmd, id *iam.Identity) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}
	if err := id.Valid(); err != nil {
		return err
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}
	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}

	// Check whether role already exists.
	identities, err := acm.Identities()
	if err != nil {
		return err
	}

	arn, _ := id.ARN() // The call to Valid above makes sure this cannot error
	duplicates := 0
	for _, identity := range identities {
		_arn, err := identity.ARN()
		if err != nil {
			return err
		}

		if arn.String() == _arn.String() {
			duplicates++
		}
	}

	if duplicates > 0 {
		logger.Warning("found %d mappings with same arn %q (which will be shadowed by your new mapping)", duplicates, arn)
	}

	if err := acm.AddIdentity(*id); err != nil {
		return err
	}
	return acm.Save()
}
