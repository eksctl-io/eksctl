package integration

import (
	"fmt"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/testutils/aws"
)

// CleanupAws cleans up the created AWS infrastructure
func CleanupAws(clusterName string, region string) {
	session := aws.NewSession(region)

	if found, _ := aws.EksClusterExists(clusterName, session); found {
		if err := aws.EksClusterDelete(clusterName, session); err != nil {
			logger.Debug("EKS cluster couldn't be deleted: %v", err)
		}
	}

	stackName := fmt.Sprintf("eksctl-%s-cluster", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		if err := aws.DeleteStack(stackName, session); err != nil {
			logger.Debug("Cluster stack couldn't be deleted: %v", err)
		}
	}

	stackName = fmt.Sprintf("eksctl-%s-nodegroup-%d", clusterName, 0)
	if found, _ := aws.StackExists(stackName, session); found {
		if err := aws.DeleteStack(stackName, session); err != nil {
			logger.Debug("NodeGroup stack couldn't be deleted: %v", err)
		}
	}

}
