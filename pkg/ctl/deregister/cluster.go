package deregister

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/connector"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func deregisterClusterCmd(cmd *cmdutils.Cmd) {
	cmd.SetDescription("cluster", "Deregister a non-EKS Kubernetes cluster", "")

	var clusterName string

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		return deregisterCluster(cmd, clusterName)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&clusterName, "name", "", "EKS cluster name")
		_ = cobra.MarkFlagRequired(fs, "name")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

}

func deregisterCluster(cmd *cmdutils.Cmd, clusterName string) error {
	clusterProvider, err := eks.New(&cmd.ProviderConfig, nil)
	if err != nil {
		return err
	}

	c := connector.EKSConnector{
		Provider: clusterProvider.Provider,
	}

	if err := c.DeregisterCluster(clusterName); err != nil {
		return errors.Wrap(err, "error deregistering cluster")
	}

	logger.Info("unregistered cluster %q successfully", clusterName)
	manifestFilenames, err := connector.GetManifestFilenames()
	if err != nil {
		return err
	}
	logger.Info("run `kubectl delete -f %s` on your cluster to remove EKS Connector resources", strings.Join(manifestFilenames, ","))
	return nil
}
