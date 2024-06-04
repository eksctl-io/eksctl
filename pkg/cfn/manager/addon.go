package manager

import "fmt"

// MakeAddonStackName creates a stack name for clusterName and addonName.
func MakeAddonStackName(clusterName, addonName string) string {
	return fmt.Sprintf("eksctl-%s-addon-%s", clusterName, addonName)
}
