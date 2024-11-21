package automode

import (
	"embed"
	"fmt"
	"io"
	"path"
)

//go:embed assets/node-rbac/*
var nodeRBACResources embed.FS

// RawClient represents a raw Kubernetes client.
type RawClient interface {
	CreateOrReplace(manifest []byte, plan bool) error
	Delete(manifest []byte) error
}

// RBACApplier applies RBAC resources.
type RBACApplier struct {
	// RawClient is a Kubernetes client for applying or deleting resources.
	RawClient RawClient
}

// ApplyRBACResources applies node RBAC resources to the cluster.
func (r *RBACApplier) ApplyRBACResources() error {
	return r.forEachResource(func(manifest []byte) error {
		if err := r.RawClient.CreateOrReplace(manifest, false); err != nil {
			return fmt.Errorf("applying node RBAC resource: %w", err)
		}
		return nil
	})
}

func (r *RBACApplier) DeleteRBACResources() error {
	return r.forEachResource(func(manifest []byte) error {
		if err := r.RawClient.Delete(manifest); err != nil {
			return fmt.Errorf("deleting node RBAC resource: %w", err)
		}
		return nil
	})
}

func (r *RBACApplier) forEachResource(f func(manifest []byte) error) error {
	baseDir := path.Join("assets", "node-rbac")
	entries, err := nodeRBACResources.ReadDir(baseDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		file, err := nodeRBACResources.Open(path.Join(baseDir, entry.Name()))
		if err != nil {
			return err
		}
		manifest, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		if err := f(manifest); err != nil {
			return err
		}
	}
	return nil
}
