package utils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/thapakazi/easyssh-go/library"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func listNodesCmd(g *cmdutils.Grouping) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-nodes",
		Short: "List the nodes in the cluster in ssh config format",
		Run: func(_ *cobra.Command, args []string) {
			logger.Info("listing the worker nodes...")
			generateSSHConfig()
		},
	}
	return cmd
}

// generateSSHConfig populates the ssh config based on cluster tags
func generateSSHConfig() {

	// TODO: discover these base on cluster nodes
	username := "ec2-user"
	port := "22"
	var tags = make(map[string][]string)
	tags["tag:eksctl.cluster.k8s.io/v1alpha1/cluster-name"] = []string{"myeks"}
	tags["instance-state-name"] = []string{"running", "pending"}

	response, err := library.FetchIps(tags)
	if err != nil {
		logger.Critical("Error fetching ips: %q", err.Error())
		os.Exit(-1)
	}
	err = library.GenerateConfig(username, port, response, os.Stdout)
	if err != nil {
		logger.Critical("Error generating ssh config: %q", err.Error())
		os.Exit(-1)
	}
}
