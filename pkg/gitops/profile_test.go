package gitops_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
	"text/template"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/eksctl/pkg/gitops"
	"gopkg.in/src-d/go-git.v4"
)

func TestOverlays(t *testing.T) {

	overlayTests := []struct {
		clusterName string
	}{
		{clusterName: "eks-cluster"},
		{clusterName: "gitops-cluster"},
		{clusterName: "dev-eks"},
		{clusterName: "prod-eks-1"},
	}

	expectedKustFile := `apiVersion: kustomize.config.k8s.io/v1beta1
bases:
- ../../base
kind: Kustomization
patchesJson6902:
- path: alb-patch.yaml
  target:
    group: apps
    kind: Deployment
    name: alb-ingress-controller
    version: v1
- path: grafana-patch.yaml
  target:
    group: flux.weave.works
    kind: HelmRelease
    name: prometheus-operator
    version: v1beta1
`

	expectedGrafanaPatch := `- op: add
  path: /spec/values/grafana
  value:
    defaultDashboardsEnabled: true
    enabled: true
`

	for _, tt := range overlayTests {
		t.Run(tt.clusterName, func(t *testing.T) {
			profile := &gitops.Profile{
				ClusterName: tt.clusterName,
			}

			overlayFiles, err := profile.MakeOverlays()

			if err != nil {
				t.Fatal(err)
			}

			expectedPatches := map[string]string{
				"alb-patch.yaml":     generateALBPatch(t, tt.clusterName),
				"grafana-patch.yaml": expectedGrafanaPatch,
				"kustomization.yaml": expectedKustFile,
			}
			assert.Equal(t, len(expectedPatches), len(overlayFiles))

			for i, overlayFile := range overlayFiles {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					assert := assert.New(t)
					expectedPatch, ok := expectedPatches[overlayFile.Name]
					if !ok {
						t.Fatalf("unexpected file: %q", overlayFile.Name)
					}

					assert.Equal(expectedPatch, string(overlayFile.Data))
				})
			}
		})

	}
}

func TestGenerateProfile(t *testing.T) {
	profileTests := []struct {
		clusterName      string
		gitURL           string
		expectedOverlays bool
		expectedError    error
	}{
		{
			clusterName:      "eks-default-profile",
			gitURL:           "git@github.com:weaveworks/eks-gitops-example.git",
			expectedOverlays: true,
		},
		{
			clusterName:      "eks-default-profile-2",
			gitURL:           "https://github.com/weaveworks/eks-gitops-example.git",
			expectedOverlays: true,
		},
		{
			clusterName:      "eks-phoney-profile",
			gitURL:           "https://github.com/_weaveworks/eks-gitops-example.git",
			expectedOverlays: false,
		},
		{
			clusterName:      "eks-phoney-profile-2",
			gitURL:           "git@github.com:weaveworks/eks-gitops-example_.git",
			expectedOverlays: false,
		},
		{
			clusterName:      "eks-custom-profile",
			gitURL:           "git@github.com:aws/gitops-profile.git",
			expectedOverlays: false,
		},
	}

	for i, tt := range profileTests {
		memFs := afero.NewMemMapFs()
		af := afero.Afero{Fs: memFs}

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				cloneCalled = false
				dir         = fmt.Sprint("/tmp/gitops-profile-", i)
			)

			profile := &gitops.Profile{
				ClusterName: tt.clusterName,
				CloneFunc: func(ctx context.Context, path string, options *git.CloneOptions) error {
					cloneCalled = true
					return nil
				},
				Fs:   memFs,
				IO:   af,
				Path: dir,
			}

			err := profile.Generate(context.Background(), gitops.GitOptions{
				URL:    tt.gitURL,
				Branch: "master",
			})

			assert := assert.New(t)
			assert.Nil(err)
			assert.True(cloneCalled, "expected CloneFunc to be called")

			dirExists, err := af.DirExists(filepath.Join(dir, "overlays"))
			assert.Nil(err)
			assert.True(dirExists == tt.expectedOverlays)

			if !tt.expectedOverlays {
				return
			}

			generatedFiles, err := afero.Glob(memFs, filepath.Join(dir, "*"))
			assert.Nil(err)
			expectedFiles := []string{".flux.yaml", "overlays"}
			for i, generatedFile := range generatedFiles {
				assert.Equal(expectedFiles[i], filepath.Base(generatedFile))
			}

			expectedFiles = []string{"alb-patch.yaml", "grafana-patch.yaml", "kustomization.yaml"}

			matches, err := afero.Glob(memFs, filepath.Join(dir, "overlays/cluster-components/*"))
			assert.Nil(err)

			assert.Equal(len(expectedFiles), len(matches))

			for i, file := range matches {
				assert.Equal(expectedFiles[i], filepath.Base(file))
			}

			albPatch, err := af.ReadFile(filepath.Join(dir, "overlays/cluster-components/alb-patch.yaml"))
			assert.Nil(err)
			expectedPatch := generateALBPatch(t, tt.clusterName)
			assert.Equal(expectedPatch, string(albPatch))
			assert.Nil(err)
		})
	}

}

func generateALBPatch(t *testing.T, clusterName string) string {
	tpl, err := template.New("test").Parse(`- op: add
  path: /spec/template/spec/containers/0/args/-
  value: --cluster-name={{.}}
`)

	if err != nil {
		t.Fatalf("error parsing template: %+v", err)
	}

	var data bytes.Buffer

	if err := tpl.Execute(&data, clusterName); err != nil {
		t.Fatalf("error executing template: %+v", err)
	}

	return data.String()
}
