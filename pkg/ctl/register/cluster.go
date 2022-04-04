package register

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/afero"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/connector"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func registerClusterCmd(cmd *cmdutils.Cmd) {
	cmd.SetDescription("cluster", "Register a non-EKS Kubernetes cluster", "")

	var cluster connector.ExternalCluster

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		return registerCluster(cmd, cluster)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cluster.Name, "name", "", "EKS cluster name")
		fs.StringVar(&cluster.Provider, "provider", "", fmt.Sprintf("Kubernetes provider name (one of %s)", strings.Join(connector.ValidProviders(), ", ")))
		fs.StringVar(&cluster.ConnectorRoleARN, "role-arn", "", "EKS Connector role ARN")

		requiredFlags := []string{"name", "provider"}
		for _, f := range requiredFlags {
			_ = cobra.MarkFlagRequired(fs, f)
		}
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

}

func registerCluster(cmd *cmdutils.Cmd, cluster connector.ExternalCluster) error {
	clusterProvider, err := eks.New(context.TODO(), &cmd.ProviderConfig, nil)
	if err != nil {
		return err
	}

	manifestTemplate, err := connector.GetManifestTemplate()
	if err != nil {
		return errors.Wrap(err, "error getting manifests for EKS Connector")
	}

	c := connector.EKSConnector{
		Provider:         clusterProvider.Provider,
		ManifestTemplate: manifestTemplate,
	}
	resourceList, err := c.RegisterCluster(context.TODO(), cluster)
	if err != nil {
		return errors.Wrap(err, "error registering cluster")
	}

	logger.Info("registered cluster %q successfully", cluster.Name)

	// TODO consider providing a manifests-dir argument to allow writing EKS Connector resources to a specific directory.
	return connector.WriteResources(afero.NewOsFs(), resourceList)
}
