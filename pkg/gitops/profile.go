package gitops

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
)

const (
	cloneDirPrefix = "quickstart-"

	eksctlIgnoreFilename = ".eksctlignore"
)

// Profile represents a gitops profile
type Profile struct {
	Processor         fileprocessor.FileProcessor
	ProfileCloner     git.TmpCloner
	UserRepoGitClient *git.Client
	FS                afero.Fs
	IO                afero.Afero
}

// Install installs the profile specified in the cluster config into a user's repository
func (p *Profile) Install(clusterConfig *api.ClusterConfig) error {
	if !clusterConfig.HasBootstrapProfile() {
		return nil
	}

	userRepo := clusterConfig.Git.Repo

	// Clone user's repo to apply Quick Start profile
	gitCfg := clusterConfig.Git
	bootstrapProfile := clusterConfig.Git.BootstrapProfile

	// Clone user's repo to apply Quick Start profile
	usersRepoName, err := git.RepoName(gitCfg.Repo.URL)
	if err != nil {
		return err
	}
	usersRepoDir, err := ioutil.TempDir("", usersRepoName+"-")
	if err != nil {
		return errors.Wrapf(err, "unable to create temporary directory for %q", usersRepoName)
	}
	bootstrapProfile.OutputPath = filepath.Join(usersRepoDir, "base")
	logger.Debug("directory %s will be used to clone the configuration repository and install the profile", usersRepoDir)

	err = p.UserRepoGitClient.CloneRepoInPath(
		usersRepoDir,
		git.CloneOptions{
			URL:       gitCfg.Repo.URL,
			Branch:    gitCfg.Repo.Branch,
			Bootstrap: true,
		},
	)
	if err != nil {
		return err
	}

	// Add quickstart components to user's repo. Clones the quickstart base repo
	err = p.Generate(*bootstrapProfile)
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	// Git add, commit and push component files
	if err = p.UserRepoGitClient.Add("."); err != nil {
		return err
	}

	commitMsg := fmt.Sprintf("Add %s quickstart components", bootstrapProfile.Source)
	if err = p.UserRepoGitClient.Commit(commitMsg, userRepo.User, userRepo.Email); err != nil {
		return err
	}

	if err = p.UserRepoGitClient.Push(); err != nil {
		return err
	}

	logger.Debug("deleting cloned directory %q", usersRepoDir)
	if err := p.IO.RemoveAll(usersRepoDir); err != nil {
		logger.Warning("unable to delete cloned directory %q", usersRepoDir)
	}
	return nil
}

// Generate clones the specified Git repo in a base directory and generates overlays if the Git repo
// points to a profile repo
func (p *Profile) Generate(profile api.Profile) error {
	// Translate the profile name to a URL
	var err error
	profile.Source, err = RepositoryURL(profile.Source)
	if err != nil {
		return errors.Wrap(err, "please supply a valid Quick Start name or URL")
	}

	logger.Info("cloning repository %q:%s", profile.Source, profile.Revision)
	options := git.CloneOptions{
		URL:    profile.Source,
		Branch: profile.Revision,
	}
	clonedDir, err := p.ProfileCloner.CloneRepoInTmpDir(cloneDirPrefix, options)
	if err != nil {
		return errors.Wrapf(err, "error cloning repository %s", profile.Source)
	}

	if err := p.ignoreFiles(clonedDir); err != nil {
		return errors.Wrapf(err, "error ignoring files of repository %s", profile.Source)
	}

	allManifests, err := p.loadFiles(clonedDir)
	if err != nil {
		return errors.Wrapf(err, "error loading files from repository %s", profile.Source)
	}

	logger.Info("processing template files in repository")
	outputFiles, err := p.processFiles(allManifests, clonedDir)
	if err != nil {
		return errors.Wrapf(err, "error processing manifests from repository %s", profile.Source)
	}

	if len(outputFiles) > 0 {
		logger.Info("writing new manifests to %q", profile.OutputPath)
	} else {
		logger.Info("no template files found, nothing to write")
		return nil
	}

	err = p.writeFiles(outputFiles, profile.OutputPath)
	if err != nil {
		return errors.Wrapf(err, "error writing manifests to dir: %q", profile.OutputPath)
	}

	// Delete temporary clone of the quickstart repo
	if clonedDir == "" {
		logger.Debug("no cloned directory to delete")
		return nil
	}
	logger.Debug("deleting cloned directory %q", clonedDir)
	if err := p.IO.RemoveAll(clonedDir); err != nil {
		logger.Warning("unable to delete cloned directory %q", clonedDir)
	}

	return nil
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

func (p *Profile) ignoreFiles(baseDir string) error {
	ignoreFilePath := path.Join(baseDir, eksctlIgnoreFilename)
	if exists, _ := p.IO.Exists(ignoreFilePath); exists {
		logger.Info("ignoring files declared in %s", eksctlIgnoreFilename)
		file, err := p.IO.Open(ignoreFilePath)
		if err != nil {
			return err
		}
		pathsToIgnores, err := parseDotIgnorefile(file)
		// Need to close the ignore file here as it is also deleted
		file.Close()
		if err != nil {
			return err
		}

		for _, pathToIgnore := range pathsToIgnores {
			err := p.IO.RemoveAll(path.Join(baseDir, pathToIgnore))
			if err != nil {
				return err
			}
			logger.Info("ignored %q", pathToIgnore)
		}

		// Remove the ignore file after finish
		if err := p.IO.Remove(ignoreFilePath); err != nil {
			return err
		}
	}
	return nil
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

func parseDotIgnorefile(reader io.Reader) ([]string, error) {
	result := []string{}
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile(`(?ms)^\s*(?P<pathToIgnore>[^\s#]+).*$`)
	for scanner.Scan() {
		groups := re.FindStringSubmatch(scanner.Text())
		if len(groups) != 2 {
			continue
		}
		pathToIgnore := groups[1]
		result = append(result, pathToIgnore)
	}
	return result, nil
}
