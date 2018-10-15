package integration

// // CleanupStacks cleans up the created AWS infrastructure
// func CleanupStacks(clusterName string, region string) {
// 	session := aws.NewSession(region)

// 	stackName := fmt.Sprintf("eksctl-%s-cluster", clusterName)
// 	if found, _ := aws.StackExists(stackName, session); found {
// 		if err := aws.DeleteStackWait(stackName, session); err != nil {
// 			logger.Debug("Cluster stack couldn't be deleted: %v", err)
// 		}
// 	}

// 	stackName = fmt.Sprintf("eksctl-%s-nodegroup-%d", clusterName, 0)
// 	if found, _ := aws.StackExists(stackName, session); found {
// 		if err := aws.DeleteStack(stackName, session); err != nil {
// 			logger.Debug("NodeGroup stack couldn't be deleted: %v", err)
// 		}
// 	}
// }
