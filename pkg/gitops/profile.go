package gitops

import (
	"bytes"
	"context"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/git"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const (
	cloneDirPrefix    = "quickstart-"
	templateExtension = ".templ"
)

// Profile represents a GitOps profile
type Profile struct {
	Params    TemplateParameters
	Path      string
	GitOpts   GitOptions
	gitCloner git.Cloner
	Fs        afero.Fs
	IO        afero.Afero
}

// GitOptions holds options for cloning a git repository
type GitOptions struct {
	URL    string
	Branch string
}

// TemplateParameters represents the API variables that can be used to template a profile. This set of variables will
// be applied to the go template files found. Templates filenames must end in .templ
type TemplateParameters struct {
	ClusterName string
}

// NewTemplateParams creates a set of variables for templating given a ClusterConfig object
func NewTemplateParams(clusterConfig *api.ClusterConfig) TemplateParameters {
	return TemplateParameters{
		ClusterName: clusterConfig.Metadata.Name,
	}
}

// ManifestFile represents a manifest file with data
type ManifestFile struct {
	Name string
	Data []byte
}

// Generate clones the specified Git repo in a base directory and generates overlays if the Git repo
// points to a profile repo
func (p *Profile) Generate(ctx context.Context) error {
	logger.Info("cloning repository %q:%s", p.GitOpts.URL, p.GitOpts.Branch)
	clonedDir, err := p.gitCloner.CloneRepo(cloneDirPrefix, p.GitOpts.Branch, p.GitOpts.URL)
	if err != nil {
		return errors.Wrapf(err, "error cloning repository %s", p.GitOpts.URL)
	}

	allManifests, err := p.loadFiles(clonedDir)
	if err != nil {
		return errors.Wrapf(err, "error loading files from repository %s", p.GitOpts.URL)
	}

	logger.Info("processing template files in repository")
	outputFiles, err := processFiles(allManifests, p.Params, clonedDir)
	if err != nil {
		return errors.Wrapf(err, "error processing manifests from repository %s", p.GitOpts.URL)
	}

	logger.Info("writing new manifests to %q", p.Path)
	err = p.writeFiles(outputFiles, p.Path)
	if err != nil {
		return errors.Wrapf(err, "error writing manifests to dir: %q", p.Path)
	}

	return nil
}

func (p *Profile) loadFiles(directory string) ([]ManifestFile, error) {
	files := make([]ManifestFile, 0)
	err := p.IO.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "cannot walk files in directory: %q", directory)
		}
		if info.IsDir() || !isGoTemplate(info.Name()) {
			return nil
		}

		logger.Debug("found template file %q", path)
		fileContents, err := p.IO.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "cannot read file %q", path)
		}
		files = append(files, ManifestFile{
			Name: path,
			Data: fileContents,
		})
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load files from directory %q", directory)
	}
	return files, nil
}

func processFiles(files []ManifestFile, params TemplateParameters, baseDir string) ([]ManifestFile, error) {
	outputFiles := make([]ManifestFile, 0, len(files))
	for _, file := range files {
		manifestTemplate, err := template.New(file.Name).Parse(string(file.Data))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse manifest template file %q", file.Name)
		}

		logger.Debug("executing template for file: %q", file.Name)
		out := bytes.NewBuffer(nil)
		if err = manifestTemplate.Execute(out, params); err != nil {
			return nil, errors.Wrapf(err, "cannot execute template for file %q", file.Name)
		}

		relPath, err := filepath.Rel(baseDir, file.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get relative path for file %q", file.Name)
		}
		newFileName := strings.TrimSuffix(relPath, templateExtension)
		outputFiles = append(outputFiles, ManifestFile{
			Data: out.Bytes(),
			Name: newFileName,
		})
	}
	return outputFiles, nil
}

func isGoTemplate(fileName string) bool {
	return strings.HasSuffix(fileName, templateExtension)
}

func (p *Profile) writeFiles(manifests []ManifestFile, outputPath string) error {
	for _, manifest := range manifests {
		filePath := filepath.Join(outputPath, manifest.Name)
		fileBase := filepath.Dir(filePath)

		if err := p.Fs.MkdirAll(fileBase, 0755); err != nil {
			return errors.Wrapf(err, "error creating output manifests dir: %q", outputPath)
		}

		err := p.IO.WriteFile(filePath, manifest.Data, 0644)
		if err != nil {
			return errors.Wrapf(err, "error writing manifest: %q", filePath)
		}
	}
	return nil
}
