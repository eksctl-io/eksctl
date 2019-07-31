package gitops

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"sigs.k8s.io/kustomize/pkg/gvk"
	"sigs.k8s.io/kustomize/pkg/patch"
	"sigs.k8s.io/kustomize/pkg/types"
	"sigs.k8s.io/yaml"
)

const (
	overlaysDir = "overlays/cluster-components"
	baseDir     = "base"

	defaultProfileGitHost = "github.com"
	defaultProfileGitPath = "weaveworks/eks-gitops-example.git"
)

// GitOptions holds options for cloning a git repository
type GitOptions struct {
	URL    string
	Branch string
}

// CloneFunc is a function that clones a Git repository
type CloneFunc func(ctx context.Context, path string, options *git.CloneOptions) error

// Profile represents a GitOps profile
type Profile struct {
	ClusterName string
	Path        string
	CloneFunc   CloneFunc
	Fs          afero.Fs
	IO          afero.Afero
}

type jsonPatch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// ManifestFile represents a manifest file with data
type ManifestFile struct {
	Name string
	Data []byte
}

// Generate clones the specified Git repo in a base directory and generates overlays if the Git repo
// points to a profile repo
func (p *Profile) Generate(ctx context.Context, o GitOptions) error {
	err := p.CloneFunc(ctx, filepath.Join(p.Path, baseDir), &git.CloneOptions{
		URL:           o.URL,
		ReferenceName: plumbing.NewBranchReferenceName(o.Branch),
	})

	if err != nil {
		return errors.Wrap(err, "error cloning repository")
	}

	gitEndpoint, err := transport.NewEndpoint(o.URL)
	if err != nil {
		return errors.Wrap(err, "error parsing Git URL")
	}

	if !matchesDefaultProfile(gitEndpoint) {
		return nil
	}

	manifestFiles, err := p.MakeOverlays()
	if err != nil {
		return err
	}

	overlaysPath := filepath.Join(p.Path, overlaysDir)

	if err := p.Fs.MkdirAll(overlaysPath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "error creating overlays dir: %q", overlaysDir)
	}

	for _, manifestFile := range manifestFiles {
		err := p.IO.WriteFile(filepath.Join(overlaysPath, manifestFile.Name), manifestFile.Data, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "error writing overlay manifest: %q", manifestFile.Name)
		}
	}
	return nil
}

func matchesDefaultProfile(e *transport.Endpoint) bool {
	return e.Host == defaultProfileGitHost && ((e.Protocol == "ssh" && e.Path == defaultProfileGitPath) ||
		(e.Protocol == "https" && e.Path == "/"+e.Path))
}

type kustomizePatch struct {
	filename    string
	patches     []jsonPatch
	patchTarget *patch.Target
}

// MakeOverlays generates overlays for the default profile
func (p *Profile) MakeOverlays() ([]ManifestFile, error) {
	kustomizePatches := []kustomizePatch{
		{
			filename: "alb-patch.yaml",
			patches: []jsonPatch{
				{
					Op:    "add",
					Path:  "/spec/template/spec/containers/0/args/-",
					Value: "--cluster-name=" + p.ClusterName,
				},
			},
			patchTarget: &patch.Target{
				Name: "alb-ingress-controller",
				Gvk: gvk.Gvk{
					Group:   "apps",
					Kind:    "Deployment",
					Version: "v1",
				},
			},
		},
		{
			filename: "grafana-patch.yaml",
			patches: []jsonPatch{
				{
					Op:   "add",
					Path: "/spec/values/grafana",
					Value: map[string]interface{}{
						"enabled":                  true,
						"defaultDashboardsEnabled": true,
					},
				},
			},
			patchTarget: &patch.Target{
				Name: "prometheus-operator",
				Gvk: gvk.Gvk{
					Group:   "flux.weave.works",
					Kind:    "HelmRelease",
					Version: "v1beta1",
				},
			},
		},
	}

	relBaseDir, err := filepath.Rel(overlaysDir, baseDir)
	if err != nil {
		return nil, err
	}

	kustomization := types.Kustomization{
		TypeMeta: types.TypeMeta{
			Kind:       types.KustomizationKind,
			APIVersion: types.KustomizationVersion,
		},
		Bases: []string{relBaseDir},
	}

	var overlayManifests []ManifestFile

	for _, kp := range kustomizePatches {
		kustomization.PatchesJson6902 = append(kustomization.PatchesJson6902, patch.Json6902{
			Path:   kp.filename,
			Target: kp.patchTarget,
		})

		bytes, err := yaml.Marshal(kp.patches)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal YAML: %q", kp.filename)
		}
		overlayManifests = append(overlayManifests, ManifestFile{
			Name: kp.filename,
			Data: bytes,
		})
	}

	bytes, err := yaml.Marshal(kustomization)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal kustomization config")
	}

	overlayManifests = append(overlayManifests, ManifestFile{
		Name: "kustomization.yaml",
		Data: bytes,
	})

	return overlayManifests, nil

}
