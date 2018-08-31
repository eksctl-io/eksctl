package integration

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/testutils/aws"
)

func CleanupAws(clusterName string, region string) {
	session := aws.NewSession(region)

	if found, _ := aws.EksClusterExists(clusterName, session); found {
		aws.EksClusterDelete(clusterName, session)
	}

	stackName := fmt.Sprintf("EKS-%s-DefaultNodeGroup", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		aws.DeleteStack(stackName, session)
	}

	stackName = fmt.Sprintf("EKS-%s-VPC", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		aws.DeleteStack(stackName, session)
	}

	stackName = fmt.Sprintf("EKS-%s-ControlPlane", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		aws.DeleteStack(stackName, session)
	}

	stackName = fmt.Sprintf("EKS-%s-ServiceRole", clusterName)
	if found, _ := aws.StackExists(stackName, session); found {
		aws.DeleteStack(stackName, session)
	}
}
