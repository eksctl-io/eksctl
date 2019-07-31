package gitops

import (
	"context"

	"gopkg.in/src-d/go-git.v4"
)

// CloneContext clones a Git repo
func CloneContext(ctx context.Context, path string, o *git.CloneOptions) error {
	_, err := git.PlainCloneContext(ctx, path, false, o)
	return err
}
