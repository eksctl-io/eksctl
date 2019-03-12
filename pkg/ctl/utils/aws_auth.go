package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

var (
	awsAuthAddAdminRoles  []string
	awsAuthRemoveRoles    []string
	awsAuthAddAccounts    []string
	awsAuthRemoveAccounts []string
)

func addAWSAuthCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	cmd := &cobra.Command{
		Use:   "aws-auth",
		Short: "Manipulates the aws-auth ConfigMap which maps IAM entities to Kubernetes groups",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doAWSAuth(p, cfg, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name (required)")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
		fs.StringArrayVarP(&awsAuthAddAdminRoles, "add-admin-role", "", []string{}, "IAM role to grant system:masters permissions to")
		fs.StringArrayVarP(&awsAuthRemoveRoles, "remove-role", "", []string{}, "IAM role to remove from auth ConfigMap")
		fs.StringArrayVarP(&awsAuthAddAccounts, "add-account", "", []string{}, "IAM role to grant system:masters permissions to")
		fs.StringArrayVarP(&awsAuthRemoveAccounts, "remove-account", "", []string{}, "IAM role to remove from auth ConfigMap")
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)

	return cmd
}

func doAWSAuth(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && nameArg != "" {
		return cmdutils.ErrNameFlagAndArg(cfg.Metadata.Name, nameArg)
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	client := clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace)

	obj, err := client.Get(authconfigmap.ObjectName, metav1.GetOptions{})
	// It is fine for the configmap not to exist. Any other error is fatal.
	if err != nil && !kerr.IsNotFound(err) {
		return errors.Wrapf(err, "getting auth ConfigMap")
	}

	logger.Debug("aws-auth = %s", awsutil.Prettify(obj))

	if len(awsAuthAddAdminRoles) == 0 && len(awsAuthRemoveRoles) == 0 &&
		len(awsAuthAddAccounts) == 0 && len(awsAuthRemoveAccounts) == 0 {
		logger.Info("no actions given, use --{add,remove}-{role-account}")
		return nil
	}

	acm := authconfigmap.New(obj)
	// Roles.
	for _, r := range awsAuthAddAdminRoles {
		if err := acm.AddRole(r, []string{authconfigmap.GroupMasters}); err != nil {
			return errors.Wrap(err, "adding role to auth ConfigMap")
		}
	}
	for _, r := range awsAuthRemoveRoles {
		if err := acm.RemoveRole(r); err != nil {
			return errors.Wrap(err, "removing role from auth ConfigMap")
		}
	}

	// Accounts.
	for _, a := range awsAuthAddAccounts {
		if err := acm.AddAccount(a); err != nil {
			return errors.Wrap(err, "adding account to auth ConfigMap")
		}
	}
	for _, a := range awsAuthRemoveAccounts {
		if err := acm.RemoveAccount(a); err != nil {
			return errors.Wrap(err, "removing account from auth ConfigMap")
		}
	}

	if err := acm.Save(client); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap")
	}
	logger.Success("saved auth ConfigMap")

	return nil
}
