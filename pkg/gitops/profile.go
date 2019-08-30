package gitops

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
)

const (
	cloneDirPrefix = "quickstart-"
)

// Profile represents a GitOps profile
type Profile struct {
	Processor fileprocessor.FileProcessor
	Path      string
	GitOpts   git.Options
	GitCloner git.TmpCloner
	FS        afero.Fs
	IO        afero.Afero
	clonedDir string
}

// Generate clones the specified Git repo in a base directory and generates overlays if the Git repo
// points to a profile repo
func (p *Profile) Generate(ctx context.Context) error {
	logger.Info("cloning repository %q:%s", p.GitOpts.URL, p.GitOpts.Branch)
	options := git.CloneOptions{
		URL:    p.GitOpts.URL,
		Branch: p.GitOpts.Branch,
	}
	clonedDir, err := p.GitCloner.CloneRepoInTmpDir(cloneDirPrefix, options)
	if err != nil {
		return errors.Wrapf(err, "error cloning repository %s", p.GitOpts.URL)
	}
	p.clonedDir = clonedDir

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
		return nil
	}

	err = p.writeFiles(outputFiles, p.Path)
	if err != nil {
		return errors.Wrapf(err, "error writing manifests to dir: %q", p.Path)
	}

	return nil
}

// DeleteClonedDirectory deletes the directory where the repository was cloned
func (p *Profile) DeleteClonedDirectory() {
	if p.clonedDir == "" {
		logger.Debug("no cloned directory to delete")
		return
	}
	logger.Debug("deleting cloned directory %q", p.clonedDir)
	if err := p.IO.RemoveAll(p.clonedDir); err != nil {
		logger.Warning("unable to delete cloned directory %q", p.clonedDir)
	}
}

func (p *Profile) loadFiles(directory string) ([]fileprocessor.File, error) {
	files := make([]fileprocessor.File, 0)
	err := p.IO.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "cannot walk files in directory: %q", directory)
		}
		if info.IsDir() || isGitFile(directory, path) {
			return nil
		}

		logger.Debug("found file %q", path)
		fileContents, err := p.IO.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "cannot read file %q", path)
		}
		files = append(files, fileprocessor.File{
			Path: path,
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
		outputFile, err := p.Processor.ProcessFile(file)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing file %q ", file.Path)
		}

		// Rewrite the path to a relative path from the root of the repo
		relPath, err := filepath.Rel(baseDir, outputFile.Path)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get relative path for file %q", file.Path)
		}
		outputFile.Path = relPath
		outputFiles = append(outputFiles, outputFile)
	}
	return outputFiles, nil
}

func (p *Profile) writeFiles(manifests []fileprocessor.File, outputPath string) error {
	for _, manifest := range manifests {
		filePath := filepath.Join(outputPath, manifest.Path)
		fileBase := filepath.Dir(filePath)

		if err := p.FS.MkdirAll(fileBase, 0755); err != nil {
			return errors.Wrapf(err, "error creating output manifests dir: %q", outputPath)
		}

		logger.Debug("writing file %q", filePath)
		err := p.IO.WriteFile(filePath, manifest.Data, 0644)
		if err != nil {
			return errors.Wrapf(err, "error writing manifest: %q", filePath)
		}
	}
	return nil
}

func isGitFile(baseDir string, path string) bool {
	return strings.HasPrefix(path, filepath.Join(baseDir, ".git"))
}
