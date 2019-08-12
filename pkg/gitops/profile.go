package gitops

import (
	"context"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
	"os"
	"path/filepath"
	"strings"
)

const (
	cloneDirPrefix    = "quickstart-"
	templateExtension = ".tmpl"
)

// Profile represents a GitOps profile
type Profile struct {
	Processor fileprocessor.FileProcessor
	Path      string
	GitOpts   GitOptions
	GitCloner git.Cloner
	Fs        afero.Fs
	IO        afero.Afero
}

// GitOptions holds options for cloning a git repository
type GitOptions struct {
	URL    string
	Branch string
}

// Generate clones the specified Git repo in a base directory and generates overlays if the Git repo
// points to a profile repo
func (p *Profile) Generate(ctx context.Context) error {
	logger.Info("cloning repository %q:%s", p.GitOpts.URL, p.GitOpts.Branch)
	clonedDir, err := p.GitCloner.CloneRepo(cloneDirPrefix, p.GitOpts.Branch, p.GitOpts.URL)
	if err != nil {
		return errors.Wrapf(err, "error cloning repository %s", p.GitOpts.URL)
	}

	allManifests, err := p.loadFiles(clonedDir)
	if err != nil {
		return errors.Wrapf(err, "error loading files from repository %s", p.GitOpts.URL)
	}

	logger.Info("processing template files in repository")
	outputFiles, err := p.processFiles(allManifests, clonedDir)
	if err != nil {
		return errors.Wrapf(err, "error processing manifests from repository %s", p.GitOpts.URL)
	}

	if len(outputFiles) > 0 {
		logger.Info("writing new manifests to %q", p.Path)
	} else {
		logger.Info("no template files found, nothing to write")
	}

	err = p.writeFiles(outputFiles, p.Path)
	if err != nil {
		return errors.Wrapf(err, "error writing manifests to dir: %q", p.Path)
	}

	return nil
}

func (p *Profile) loadFiles(directory string) ([]fileprocessor.File, error) {
	files := make([]fileprocessor.File, 0)
	err := p.IO.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "cannot walk files in directory: %q", directory)
		}
		if info.IsDir() || strings.HasSuffix(".git", path) {
			return nil
		}

		logger.Debug("found file %q", path)
		fileContents, err := p.IO.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "cannot read file %q", path)
		}
		files = append(files, fileprocessor.File{
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

func (p *Profile) processFiles(files []fileprocessor.File, baseDir string) ([]fileprocessor.File, error) {
	outputFiles := make([]fileprocessor.File, 0, len(files))
	for _, file := range files {
		outputFile, err := p.Processor.ProcessFile(file, baseDir)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing file %q ", file.Name)
		}
		if outputFile == nil {
			continue
		}
		outputFiles = append(outputFiles, *outputFile)
	}
	return outputFiles, nil
}

func (p *Profile) writeFiles(manifests []fileprocessor.File, outputPath string) error {
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
