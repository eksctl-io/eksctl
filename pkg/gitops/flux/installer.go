package flux

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	fluxinstall "github.com/fluxcd/flux/pkg/install"
	helmopinstall "github.com/fluxcd/helm-operator/pkg/install"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	fluxNamespaceFileName       = "flux-namespace.yaml"
	fluxPrivateSSHKeyFileName   = "flux-secret.yaml"
	fluxPrivateSSHKeySecretName = "flux-git-deploy"
	fluxHelmVersions            = "v3"
)

// Installer installs Flux
type Installer struct {
	cluster       *api.ClusterMeta
	opts          *api.Git
	timeout       time.Duration
	k8sRestConfig *rest.Config
	k8sClientSet  kubeclient.Interface
	gitClient     *git.Client
}

// NewInstaller creates a new Flux installer
func NewInstaller(k8sRestConfig *rest.Config, k8sClientSet kubeclient.Interface, cfg *api.ClusterConfig, timeout time.Duration) (*Installer, error) {
	if cfg.Git == nil {
		return nil, errors.New("expected git configuration in cluster configuration but found nil")
	}
	if cfg.Git.Repo == nil {
		return nil, errors.New("expected git.repo in cluster configuration but found nil")
	}
	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: cfg.Git.Repo.PrivateSSHKeyPath,
	})
	fi := &Installer{
		opts:          cfg.Git,
		k8sRestConfig: k8sRestConfig,
		k8sClientSet:  k8sClientSet,
		gitClient:     gitClient,
		timeout:       timeout,
	}
	return fi, nil
}

// Run runs the Flux installer
func (fi *Installer) Run(ctx context.Context) (string, error) {

	logger.Info("Generating manifests")
	manifests, err := fi.getManifests()
	if err != nil {
		return "", err
	}

	logger.Info("Cloning %s", fi.opts.Repo.URL)
	options := git.CloneOptions{
		URL:       fi.opts.Repo.URL,
		Branch:    fi.opts.Repo.Branch,
		Bootstrap: true,
	}
	cloneDir, err := fi.gitClient.CloneRepoInTmpDir("eksctl-install-flux-clone-", options)
	if err != nil {
		return "", errors.Wrapf(err, "cannot clone repository %s", fi.opts.Repo.URL)
	}
	cleanCloneDir := false
	defer func() {
		if cleanCloneDir {
			_ = fi.gitClient.DeleteLocalRepo()
		} else {
			logger.Critical("You may find the local clone of %s used by eksctl at %s",
				fi.opts.Repo.URL,
				cloneDir)
		}
	}()
	logger.Info("Writing Flux manifests")
	fluxManifestDir := filepath.Join(cloneDir, fi.opts.Repo.FluxPath)
	if err := writeFluxManifests(fluxManifestDir, manifests); err != nil {
		return "", err
	}

	if err := fi.createFluxNamespaceIfMissing(manifests); err != nil {
		return "", err
	}

	logger.Info("Applying manifests")
	if err := fi.applyManifests(manifests); err != nil {
		return "", err
	}

	if api.IsEnabled(fi.opts.Operator.WithHelm) {
		logger.Info("Waiting for Helm Operator to start")
		if err := waitForHelmOpToStart(ctx, fi.opts.Operator.Namespace, fi.timeout, fi.k8sRestConfig, fi.k8sClientSet); err != nil {
			return "", err
		}
		logger.Info("Helm Operator started successfully")
		logger.Info("see https://docs.fluxcd.io/projects/helm-operator for details on how to use the Helm Operator")
	}

	logger.Info("Waiting for Flux to start")
	fluxSSHKey, err := waitForFluxToStart(ctx, fi.opts.Operator.Namespace, fi.timeout, fi.k8sRestConfig, fi.k8sClientSet)
	if err != nil {
		return "", err
	}
	logger.Info("Flux started successfully")
	logger.Info("see https://docs.fluxcd.io/projects/flux for details on how to use Flux")

	logger.Info("Committing and pushing manifests to %s", fi.opts.Repo.URL)
	if err := fi.addFilesToRepo(); err != nil {
		return "", err
	}
	cleanCloneDir = true

	logger.Info("Flux will only operate properly once it has write-access to the Git repository")
	instruction := fmt.Sprintf("please configure %s so that the following Flux SSH public key has write access to it\n%s",
		fi.opts.Repo.URL, fluxSSHKey.Key)
	return instruction, nil
}

// IsFluxInstalled returns an error if Flux is not installed in the cluster. To determine that it looks for the flux
// pod
func (fi *Installer) IsFluxInstalled() (bool, error) {
	_, err := fi.k8sClientSet.AppsV1().Deployments(fi.opts.Operator.Namespace).Get("flux", metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("flux deployment was not found")
			return false, nil
		}
		return false, errors.Wrapf(err, "error while looking for flux pod")
	}
	return true, nil
}

func (fi *Installer) addFilesToRepo() error {
	if err := fi.gitClient.Add(fi.opts.Repo.FluxPath); err != nil {
		return err
	}

	// Confirm there is something to commit, otherwise move on
	if err := fi.gitClient.Commit("Add Initial Flux configuration", fi.opts.Repo.User, fi.opts.Repo.Email); err != nil {
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
		if existence[fi.opts.Operator.Namespace][fluxPrivateSSHKeySecretName] {
			delete(manifestsMap, fluxPrivateSSHKeyFileName)
		}
	}

	// do not recreate Flux's cache if exists
	for m, r := range manifestsMap {
		if strings.Contains(m, "memcache") {
			existence, err := client.Exists(r)
			if err != nil {
				return err
			}
			for _, found := range existence[fi.opts.Operator.Namespace] {
				if found {
					delete(manifestsMap, m)
				}
			}
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

func (fi *Installer) getManifests() (map[string][]byte, error) {
	var manifests map[string][]byte

	// Flux
	var err error
	if manifests, err = getFluxManifests(fi.opts, fi.k8sClientSet); err != nil {
		return nil, err
	}

	// Helm Operator
	if !api.IsEnabled(fi.opts.Operator.WithHelm) {
		return manifests, nil
	}
	helmOpManifests, err := getHelmOpManifests(fi.opts.Operator.Namespace)
	if err != nil {
		return nil, err
	}
	manifests = mergeMaps(manifests, helmOpManifests)

	return manifests, nil
}

func getFluxManifests(opts *api.Git, cs kubeclient.Interface) (map[string][]byte, error) {
	manifests := map[string][]byte{}
	fluxNSExists, err := kubernetes.CheckNamespaceExists(cs, opts.Operator.Namespace)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot check if namespace %s exists", opts.Operator.Namespace)
	}
	if !fluxNSExists {
		manifests[fluxNamespaceFileName] = kubernetes.NewNamespaceYAML(opts.Operator.Namespace)
	}
	fluxParameters := fluxinstall.TemplateParameters{
		GitURL:             opts.Repo.URL,
		GitBranch:          opts.Repo.Branch,
		GitPaths:           opts.Repo.Paths,
		GitLabel:           opts.Operator.Label,
		GitUser:            opts.Repo.User,
		GitEmail:           opts.Repo.Email,
		GitReadOnly:        false,
		Namespace:          opts.Operator.Namespace,
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
