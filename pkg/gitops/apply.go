package gitops

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
)

// GitOps can set up a repo as a gitops repo with flux
type GitOps struct {
	ClusterConfig    *api.ClusterConfig
	UsersRepoOpts    git.Options
	QuickstartName   string
	FluxInstaller    *flux.Installer
	ProfileGenerator *Profile
	GitClient        *git.Client
}

// Apply setups gitops in a repository and a cluster and installs flux, helm, tiller and a quickstart into the cluster
func (g *GitOps) Apply(ctx context.Context) error {

	// Install Flux, Helm and Tiller. Clones the user's repo
	err := g.FluxInstaller.Run(context.Background())
	if err != nil {
		return err
	}

	// Add quickstart components to user's repo. Clones the quickstart base repo
	err = g.ProfileGenerator.Generate(context.Background())
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	// Git add, commit and push component files
	if err = g.GitClient.Add("."); err != nil {
		return err
	}

	commitMsg := fmt.Sprintf("Add %s quickstart components", g.QuickstartName)
	if err = g.GitClient.Commit(commitMsg, g.UsersRepoOpts.User, g.UsersRepoOpts.Email); err != nil {
		return err
	}

	if err = g.GitClient.Push(); err != nil {
		return err
	}

	// Delete temporary clone of the quickstart profile repo
	g.FluxInstaller.DeleteCloneDir()
	return nil
}
