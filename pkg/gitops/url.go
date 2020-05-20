package gitops

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/git"
)

// RepositoryURL returns the full Git repository URL corresponding to the
// provided "quickstart profile" mnemonic a.k.a. short name. If a valid Git URL
// is provided, this function returns it as-is.
func RepositoryURL(quickstartArgument string) (string, error) {
	if git.IsGitURL(quickstartArgument) {
		return quickstartArgument, nil
	}
	if quickstartArgument == "app-dev" {
		return "https://github.com/weaveworks/eks-quickstart-app-dev", nil
	}
	if quickstartArgument == "appmesh" {
		return "https://github.com/weaveworks/eks-appmesh-profile", nil
	}
	return "", fmt.Errorf("invalid URL or unknown Quick Start profile: %s", quickstartArgument)
}
