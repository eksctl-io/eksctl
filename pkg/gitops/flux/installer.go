package flux

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	fluxinstall "github.com/fluxcd/flux/pkg/install"
	helmopinstall "github.com/fluxcd/helm-operator/pkg/install"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/riywo/loginshell"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	fluxNamespaceFileName       = "flux-namespace.yaml"
	fluxPrivateSSHKeyFileName   = "flux-secret.yaml"
	fluxPrivateSSHKeySecretName = "flux-git-deploy"
	fluxHelmVersions            = "v3"
)

// InstallOpts are the installation options for Flux
type InstallOpts struct {
	GitOptions  git.Options
	GitPaths    []string
	GitLabel    string
	GitFluxPath string
	Namespace   string
	Timeout     time.Duration
	Amend       bool // TODO: remove, as we eventually no longer want to support this mode?
	WithHelm    bool
}

// Installer installs Flux
type Installer struct {
	opts          *InstallOpts
	k8sRestConfig *rest.Config
	k8sClientSet  kubeclient.Interface
	gitClient     *git.Client
}

// NewInstaller creates a new Flux installer
func NewInstaller(k8sRestConfig *rest.Config, k8sClientSet kubeclient.Interface, opts *InstallOpts) *Installer {
	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: opts.GitOptions.PrivateSSHKeyPath,
	})
	fi := &Installer{
		opts:          opts,
		k8sRestConfig: k8sRestConfig,
		k8sClientSet:  k8sClientSet,
		gitClient:     gitClient,
	}
	return fi
}

// Run runs the Flux installer
func (fi *Installer) Run(ctx context.Context) (string, error) {

	logger.Info("Generating manifests")
	manifests, err := fi.getManifests()
	if err != nil {
		return "", err
	}

	logger.Info("Cloning %s", fi.opts.GitOptions.URL)
	options := git.CloneOptions{
		URL:       fi.opts.GitOptions.URL,
		Branch:    fi.opts.GitOptions.Branch,
		Bootstrap: true,
	}
	cloneDir, err := fi.gitClient.CloneRepoInTmpDir("eksctl-install-flux-clone-", options)
	if err != nil {
		return "", errors.Wrapf(err, "cannot clone repository %s", fi.opts.GitOptions.URL)
	}
	cleanCloneDir := false
	defer func() {
		if cleanCloneDir {
			_ = fi.gitClient.DeleteLocalRepo()
		} else {
			logger.Critical("You may find the local clone of %s used by eksctl at %s",
				fi.opts.GitOptions.URL,
				cloneDir)
		}
	}()
	logger.Info("Writing Flux manifests")
	fluxManifestDir := filepath.Join(cloneDir, fi.opts.GitFluxPath)
	if err := writeFluxManifests(fluxManifestDir, manifests); err != nil {
		return "", err
	}

	if fi.opts.Amend {
		logger.Info("Stopping to amend the the Flux manifests, please exit the shell when done.")
		if err := runShell(fluxManifestDir); err != nil {
			return "", err
		}
		// Re-read the manifests, as they may have changed:
		manifests, err = readFluxManifests(fluxManifestDir)
		if err != nil {
			return "", err
		}
	}

	if err := fi.createFluxNamespaceIfMissing(manifests); err != nil {
		return "", err
	}

	logger.Info("Applying manifests")
	if err := fi.applyManifests(manifests); err != nil {
		return "", err
	}

	if fi.opts.WithHelm {
		logger.Info("Waiting for Helm Operator to start")
		if err := waitForHelmOpToStart(ctx, fi.opts.Namespace, fi.opts.Timeout, fi.k8sRestConfig, fi.k8sClientSet); err != nil {
			return "", err
		}
		logger.Info("Helm Operator started successfully")
		logger.Info("see https://docs.fluxcd.io/projects/helm-operator for details on how to use the Helm Operator")
	}

	logger.Info("Waiting for Flux to start")
	fluxSSHKey, err := waitForFluxToStart(ctx, fi.opts.Namespace, fi.opts.Timeout, fi.k8sRestConfig, fi.k8sClientSet)
	if err != nil {
		return "", err
	}
	logger.Info("Flux started successfully")
	logger.Info("see https://docs.fluxcd.io/projects/flux for details on how to use Flux")

	logger.Info("Committing and pushing manifests to %s", fi.opts.GitOptions.URL)
	if err := fi.addFilesToRepo(ctx, cloneDir); err != nil {
		return "", err
	}
	cleanCloneDir = true

	logger.Info("Flux will only operate properly once it has write-access to the Git repository")
	instruction := fmt.Sprintf("please configure %s so that the following Flux SSH public key has write access to it\n%s",
		fi.opts.GitOptions.URL, fluxSSHKey.Key)
	return instruction, nil
}

func (fi *Installer) addFilesToRepo(ctx context.Context, cloneDir string) error {
	if err := fi.gitClient.Add(fi.opts.GitFluxPath); err != nil {
		return err
	}

	// Confirm there is something to commit, otherwise move on
	if err := fi.gitClient.Commit("Add Initial Flux configuration", fi.opts.GitOptions.User, fi.opts.GitOptions.Email); err != nil {
		return err
	}

	// git push
	if err := fi.gitClient.Push(); err != nil {
		return err
	}
	return nil
}

func (fi *Installer) createFluxNamespaceIfMissing(manifestsMap map[string][]byte) error {
	client, err := kubernetes.NewRawClient(fi.k8sClientSet, fi.k8sRestConfig)
	if err != nil {
		return err
	}

	// If the flux namespace needs to be created, do it first, before any other
	// resource which should potentially be created within the namespace.
	// Otherwise, creation of these resources will fail.
	if namespace, ok := manifestsMap[fluxNamespaceFileName]; ok {
		if err := client.CreateOrReplace(namespace, false); err != nil {
			return err
		}
		delete(manifestsMap, fluxNamespaceFileName)
	}
	return nil
}

func (fi *Installer) applyManifests(manifestsMap map[string][]byte) error {
	client, err := kubernetes.NewRawClient(fi.k8sClientSet, fi.k8sRestConfig)
	if err != nil {
		return err
	}

	if fluxSecret, ok := manifestsMap[fluxPrivateSSHKeyFileName]; ok {
		existence, err := client.Exists(fluxSecret)
		if err != nil {
			return err
		}
		// We do NOT want to recreate the flux-git-deploy Secret object inside
		// the flux-secret.yaml file, as it contains Flux's private SSH key,
		// and deleting it would force the user to set Flux's permissions up
		// again in their Git repository, which is not very "friendly".
		if existence[fi.opts.Namespace][fluxPrivateSSHKeySecretName] {
			delete(manifestsMap, fluxPrivateSSHKeyFileName)
		}
	}

	var manifestValues [][]byte
	for _, manifest := range manifestsMap {
		manifestValues = append(manifestValues, manifest)
	}
	manifests := kubernetes.ConcatManifests(manifestValues...)
	return client.CreateOrReplace(manifests, false)
}

func writeFluxManifests(baseDir string, manifests map[string][]byte) error {
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return errors.Wrapf(err, "cannot create Flux manifests directory (%s)", baseDir)
	}
	for fileName, contents := range manifests {
		fullPath := filepath.Join(baseDir, fileName)
		if err := ioutil.WriteFile(fullPath, contents, 0600); err != nil {
			return errors.Wrapf(err, "failed to write Flux manifest file %s", fullPath)
		}
	}
	return nil
}

func readFluxManifests(baseDir string) (map[string][]byte, error) {
	manifestFiles, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list Flux manifest files in %s", baseDir)
	}
	manifests := map[string][]byte{}
	for _, manifestFile := range manifestFiles {
		manifest, err := ioutil.ReadFile(manifestFile.Name())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read Flux manifest file %s", manifestFile.Name())
		}
		manifests[manifestFile.Name()] = manifest
	}
	return manifests, nil
}

func runShell(workDir string) error {
	shell, err := loginshell.Shell()
	if err != nil {
		return errors.Wrapf(err, "failed to obtain login shell %s", shell)
	}
	shellCmd := exec.Cmd{
		Path:   shell,
		Dir:    workDir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return shellCmd.Run()
}

func (fi *Installer) getManifests() (map[string][]byte, error) {
	manifests := map[string][]byte{}

	// Flux
	var err error
	if manifests, err = getFluxManifests(fi.opts, fi.k8sClientSet); err != nil {
		return nil, err
	}

	// Helm Operator
	if !fi.opts.WithHelm {
		return manifests, nil
	}
	helmOpManifests, err := getHelmOpManifests(fi.opts.Namespace)
	if err != nil {
		return nil, err
	}
	manifests = mergeMaps(manifests, helmOpManifests)

	return manifests, nil
}

func getFluxManifests(opts *InstallOpts, cs kubeclient.Interface) (map[string][]byte, error) {
	manifests := map[string][]byte{}
	fluxNSExists, err := kubernetes.CheckNamespaceExists(cs, opts.Namespace)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot check if namespace %s exists", opts.Namespace)
	}
	if !fluxNSExists {
		manifests[fluxNamespaceFileName] = kubernetes.NewNamespaceYAML(opts.Namespace)
	}
	fluxParameters := fluxinstall.TemplateParameters{
		GitURL:             opts.GitOptions.URL,
		GitBranch:          opts.GitOptions.Branch,
		GitPaths:           opts.GitPaths,
		GitLabel:           opts.GitLabel,
		GitUser:            opts.GitOptions.User,
		GitEmail:           opts.GitOptions.Email,
		GitReadOnly:        false,
		RegistryScanning:   true,
		Namespace:          opts.Namespace,
		ManifestGeneration: true,
		AdditionalFluxArgs: []string{"--sync-garbage-collection"},
	}
	fluxManifests, err := fluxinstall.FillInTemplates(fluxParameters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Flux manifests")
	}
	return mergeMaps(manifests, fluxManifests), nil
}

func getHelmOpManifests(namespace string) (map[string][]byte, error) {
	helmOpParameters := helmopinstall.TemplateParameters{
		Namespace:     namespace,
		HelmVersions:  fluxHelmVersions,
		SSHSecretName: fluxPrivateSSHKeySecretName,
	}
	manifests, err := helmopinstall.FillInTemplates(helmOpParameters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Helm Operator Manifests")
	}
	return manifests, nil
}

func mergeMaps(m1 map[string][]byte, m2 map[string][]byte) map[string][]byte {
	result := make(map[string][]byte, len(m1)+len(m2))
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result
}
