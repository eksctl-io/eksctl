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

	stackName := fmt.Sprintf("EKS-%s-DefaultNodeGroup", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		if err := aws.DeleteStack(stackName, session); err != nil {
			logger.Debug("DefaultNodeGroup stack couldn't be deleted: %v", err)
		}
	}

	stackName = fmt.Sprintf("EKS-%s-VPC", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		if err := aws.DeleteStack(stackName, session); err != nil {
			logger.Debug("VPC stack couldn't be deleted: %v", err)
		}
	}

	stackName = fmt.Sprintf("EKS-%s-ControlPlane", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		if err := aws.DeleteStack(stackName, session); err != nil {
			logger.Debug("ControlPlane stack couldn't be deleted: %v", err)
		}
	}

	stackName = fmt.Sprintf("EKS-%s-ServiceRole", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		if err := aws.DeleteStack(stackName, session); err != nil {
			logger.Debug("ServiceRole stack couldn't be deleted: %v", err)
		}
	}
}
