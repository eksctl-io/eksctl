package main

import (
	"github.com/kubicorn/kubicorn/pkg/logger"

	cloudformation "github.com/weaveworks/eksctl/pkg/cfn"
)

const (
	DEFAULT_EKS_REGION   = "us-west-2"
	DEFAULT_CLUSTER_NAME = "cluster-1"
	DEFAULT_NODE_COUNT   = 2
	DEFAULT_NODE_TYPE    = "m4.large" // TODO check kops
)

func main() {
	logger.Level = 4

	config := &cloudformation.Config{
		ClusterName: DEFAULT_CLUSTER_NAME,
		Region:      DEFAULT_EKS_REGION,

		KeyName: "ilya", // TODO

		MinNodes: DEFAULT_NODE_COUNT,
		MaxNodes: DEFAULT_NODE_COUNT,
		NodeType: DEFAULT_NODE_TYPE,
	}

	cfn := cloudformation.New(config)

	if err := cfn.CheckAuth(); err != nil {
		logger.Critical("%s", err)
		return
	}

	// TODO waitgroups
	{
		done := make(chan error)
		if err := cfn.CreateStackServiceRole(done); err != nil {
			logger.Critical("%s", err)
			return
		}
		if err := <-done; err != nil {
			logger.Critical("%s", err)
			return
		}
	}
	{
		done := make(chan error)
		if err := cfn.CreateStackVPC(done); err != nil {
			logger.Critical("%s", err)
			return
		}
		if err := <-done; err != nil {
			logger.Critical("%s", err)
			return
		}
	}
}
