// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	fluxinstall "github.com/fluxcd/flux/pkg/install"
	"github.com/fluxcd/go-git-providers/gitprovider"
	helmopinstall "github.com/fluxcd/helm-operator/pkg/install"
	portforward "github.com/justinbarrick/go-k8s-portforward"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/deploykey"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	fluxNamespaceFileName       = "flux-namespace.yaml"
	fluxPrivateSSHKeyFileName   = "flux-secret.yaml"
	fluxPrivateSSHKeySecretName = "flux-git-deploy"
	fluxHelmVersions            = "v3"
	portForwardingTimeout       = 120 * time.Second
	portForwardingRetryPeriod   = 2 * time.Second
)

type Installer struct {
	Cfg           *api.ClusterConfig
	Opts          *api.Git
	Timeout       time.Duration
	K8sRestConfig *rest.Config
	K8sClientSet  kubeclient.Interface
	GitClient     *git.Client
}

func New(
	k8sRestConfig *rest.Config, k8sClientSet kubeclient.Interface,
	cfg *api.ClusterConfig, timeout time.Duration,
) (*Installer, error) {
	if cfg.Git == nil {
		return nil, errors.New("expected git configuration in cluster configuration but found nil")
	}
	if cfg.Git.Repo == nil {
		return nil, errors.New("expected git.repo in cluster configuration but found nil")
	}
	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: cfg.Git.Repo.PrivateSSHKeyPath,
	})
	installer := &Installer{
		Cfg:           cfg,
		Opts:          cfg.Git,
		K8sRestConfig: k8sRestConfig,
		K8sClientSet:  k8sClientSet,
		GitClient:     gitClient,
		Timeout:       timeout,
	}
	return installer, nil
}

func (fi *Installer) Run() error {
	fluxIsInstalled, err := fi.isFluxInstalled()
	if err != nil {
		// Continue with installation
		logger.Warning(err.Error())
	}

	if fluxIsInstalled {
		logger.Warning("found existing flux deployment in namespace %q. Skipping installation", fi.Opts.Operator.Namespace)
		return nil
	}

	userInstructions, err := fi.installFlux(context.Background())
	if err != nil {
		logger.Critical("unable to set up gitops repo: %s", err.Error())
		return err
	}
	logger.Info(userInstructions)

	return nil
}

func (fi *Installer) isFluxInstalled() (bool, error) {
	_, err := fi.K8sClientSet.AppsV1().Deployments(fi.Opts.Operator.Namespace).Get(context.TODO(), "flux", metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("flux deployment was not found")
			return false, nil
		}
		return false, errors.Wrapf(err, "error while looking for flux pod")
	}
	return true, nil
}

func (fi *Installer) installFlux(ctx context.Context) (string, error) {
	logger.Info("Generating manifests")
	manifests, err := fi.GetManifests()
	if err != nil {
		return "", err
	}

	logger.Info("Cloning %s", fi.Opts.Repo.URL)
	options := git.CloneOptions{
		URL:       fi.Opts.Repo.URL,
		Branch:    fi.Opts.Repo.Branch,
		Bootstrap: true,
	}
	cloneDir, err := fi.GitClient.CloneRepoInTmpDir("eksctl-install-flux-clone-", options)
	if err != nil {
		return "", errors.Wrapf(err, "cannot clone repository %s", fi.Opts.Repo.URL)
	}
	cleanCloneDir := false
	defer func() {
		if cleanCloneDir {
			_ = fi.GitClient.DeleteLocalRepo()
		} else {
			logger.Critical("You may find the local clone of %s used by eksctl at %s",
				fi.Opts.Repo.URL,
				cloneDir)
		}
	}()
	logger.Info("Writing Flux manifests")
	fluxManifestDir := filepath.Join(cloneDir, fi.Opts.Repo.FluxPath)
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

	if api.IsEnabled(fi.Opts.Operator.WithHelm) {
		logger.Info("Waiting for Helm Operator to start")
		if err := waitForHelmOpToStart(fi.Opts.Operator.Namespace, fi.Timeout, fi.K8sClientSet); err != nil {
			return "", err
		}
		logger.Info("Helm Operator started successfully")
		logger.Info("see https://docs.fluxcd.io/projects/helm-operator for details on how to use the Helm Operator")
	}

	logger.Info("Waiting for Flux to start")
	err = waitForFluxToStart(fi.Opts.Operator.Namespace, fi.Timeout, fi.K8sClientSet)
	if err != nil {
		return "", err
	}
	logger.Info("fetching public SSH key from Flux")
	fluxSSHKey, err := getPublicKeyFromFlux(ctx, fi.Opts.Operator.Namespace, fi.Timeout, fi.K8sRestConfig, fi.K8sClientSet)
	if err != nil {
		return "", err
	}

	logger.Info("Flux started successfully")
	logger.Info("see https://docs.fluxcd.io/projects/flux for details on how to use Flux")

	if api.IsEnabled(fi.Opts.Operator.CommitOperatorManifests) {
		logger.Info("Committing and pushing manifests to %s", fi.Opts.Repo.URL)
		if err = fi.addFilesToRepo(); err != nil {
			return "", err
		}
	}
	cleanCloneDir = true

	logger.Info("Flux will only operate properly once it has write-access to the Git repository")
	instruction := fmt.Sprintf("please configure %s so that the following Flux SSH public key has write access to it\n%s",
		fi.Opts.Repo.URL, fluxSSHKey.Key)

	client, err := deploykey.GetDeployKeyClient(ctx, fi.Opts.Repo.URL)
	if err != nil {
		logger.Warning(
			"could not find git provider implementation for url %q: %q. Skipping authorization of SSH key",
			fi.Opts.Repo.URL,
			err.Error(),
		)
		return instruction, nil
	}

	keyTitle := KeyTitle(*fi.Cfg.Metadata)
	_, _, err = client.Reconcile(ctx, gitprovider.DeployKeyInfo{
		Name:     keyTitle,
		Key:      []byte(fluxSSHKey.Key),
		ReadOnly: &fi.Opts.Operator.ReadOnly,
	})
	if err != nil {
		return instruction, errors.Wrapf(err, "could not authorize SSH key")
	}
	instruction = fmt.Sprintf("Flux SSH key with name %q authorized to access the repo", keyTitle)

	return instruction, nil
}

func (fi *Installer) addFilesToRepo() error {
	if err := fi.GitClient.Add(fi.Opts.Repo.FluxPath); err != nil {
		return err
	}

	// Confirm there is something to commit, otherwise move on
	if err := fi.GitClient.Commit("Add Initial Flux configuration", fi.Opts.Repo.User, fi.Opts.Repo.Email); err != nil {
		return err
	}

	// git push
	if err := fi.GitClient.Push(); err != nil {
		return err
	}
	return nil
}

func (fi *Installer) createFluxNamespaceIfMissing(manifestsMap map[string][]byte) error {
	client, err := kubernetes.NewRawClient(fi.K8sClientSet, fi.K8sRestConfig)
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
	client, err := kubernetes.NewRawClient(fi.K8sClientSet, fi.K8sRestConfig)
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
		if existence[fi.Opts.Operator.Namespace][fluxPrivateSSHKeySecretName] {
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
			for _, found := range existence[fi.Opts.Operator.Namespace] {
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
		if err := os.WriteFile(fullPath, contents, 0600); err != nil {
			return errors.Wrapf(err, "failed to write Flux manifest file %s", fullPath)
		}
	}
	return nil
}

func (fi *Installer) GetManifests() (map[string][]byte, error) {
	var manifests map[string][]byte

	// Flux
	var err error
	if manifests, err = getFluxManifests(fi.Opts, fi.K8sClientSet); err != nil {
		return nil, err
	}

	// Helm Operator
	if !api.IsEnabled(fi.Opts.Operator.WithHelm) {
		return manifests, nil
	}
	helmOpManifests, err := getHelmOpManifests(fi.Opts.Operator.Namespace, fi.Opts.Operator.AdditionalHelmOperatorArgs)
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

	additionalFluxArgs := []string{"--sync-garbage-collection"}
	additionalFluxArgs = append(additionalFluxArgs, opts.Operator.AdditionalFluxArgs...)
	if opts.Operator.ReadOnly {
		additionalFluxArgs = append(additionalFluxArgs, "--registry-disable-scanning")
	}

	fluxParameters := fluxinstall.TemplateParameters{
		GitURL:             opts.Repo.URL,
		GitBranch:          opts.Repo.Branch,
		GitPaths:           opts.Repo.Paths,
		GitLabel:           opts.Operator.Label,
		GitUser:            opts.Repo.User,
		GitEmail:           opts.Repo.Email,
		GitReadOnly:        opts.Operator.ReadOnly,
		Namespace:          opts.Operator.Namespace,
		ManifestGeneration: true,
		AdditionalFluxArgs: additionalFluxArgs,
	}
	fluxManifests, err := fluxinstall.FillInTemplates(fluxParameters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Flux manifests")
	}
	return mergeMaps(manifests, fluxManifests), nil
}

func getHelmOpManifests(namespace string, additionalHelmOperatorArgs []string) (map[string][]byte, error) {
	helmOpParameters := helmopinstall.TemplateParameters{
		Namespace:      namespace,
		AdditionalArgs: additionalHelmOperatorArgs,
		HelmVersions:   fluxHelmVersions,
		SSHSecretName:  fluxPrivateSSHKeySecretName,
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

// PublicKey represents a public SSH key as it is returned by flux
type PublicKey struct {
	Key string `json:"key"`
}

func getPublicKeyFromFlux(
	ctx context.Context, namespace string, timeout time.Duration, restConfig *rest.Config,
	cs kubeclient.Interface,
) (PublicKey, error) {
	var deployKey PublicKey
	try := func(rootURL string) error {
		fluxURL := rootURL + "api/flux/v6/identity.pub"
		req, reqErr := http.NewRequest("GET", fluxURL, nil)
		if reqErr != nil {
			return fmt.Errorf("failed to create request: %s", reqErr)
		}
		repoCtx, repoCtxCancel := context.WithTimeout(ctx, timeout)
		defer repoCtxCancel()
		req = req.WithContext(repoCtx)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to query Flux API: %s", err)
		}
		if resp.Body == nil {
			return fmt.Errorf("failed to fetch Flux deploy key from: %s", fluxURL)
		}
		defer resp.Body.Close()

		jsonErr := json.NewDecoder(resp.Body).Decode(&deployKey)
		if jsonErr != nil {
			return fmt.Errorf("failed to decode Flux API response: %s", jsonErr)
		}

		if deployKey.Key == "" {
			return fmt.Errorf("failed to fetch Flux deploy key from: %s", fluxURL)
		}
		return nil
	}
	err := portForward(namespace, "flux", 3030, "Flux", restConfig, cs, try)
	return deployKey, err
}

type tryFunc func(rootURL string) error

func portForward(
	namespace string, nameLabelValue string, port int, name string,
	restConfig *rest.Config, cs kubeclient.Interface, try tryFunc,
) error {
	fluxSelector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "name",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{nameLabelValue},
			},
		},
	}
	portforwarder := portforward.PortForward{
		Labels:          fluxSelector,
		Config:          restConfig,
		Clientset:       cs,
		DestinationPort: port,
		Namespace:       namespace,
	}
	podDeadline := time.Now().Add(portForwardingTimeout)
	for ; time.Now().Before(podDeadline); time.Sleep(portForwardingRetryPeriod) {
		err := portforwarder.Start(context.TODO())
		if err == nil {
			defer portforwarder.Stop()
			break
		}
		if !strings.Contains(err.Error(), "Could not find running pod for selector") {
			logger.Warning("%s is not ready yet (%s), retrying ...", name, err)
		}
	}
	if time.Now().After(podDeadline) {
		return fmt.Errorf("timed out waiting for %s's pod to be created", name)
	}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d/", portforwarder.ListenPort)
	// Make sure it's alive
	retryDeadline := time.Now().Add(30 * time.Second)
	for ; time.Now().Before(retryDeadline); time.Sleep(2 * time.Second) {
		err := try(baseURL)
		if err == nil {
			break
		}
		logger.Warning("%s is not ready yet (%s), retrying ...", name, err)
	}
	if time.Now().After(retryDeadline) {
		return fmt.Errorf("timed out waiting for %s to be operative", name)
	}
	return nil
}

// KeyTitle returns the title for the SSH key used to access the gitops repo
func KeyTitle(clusterMeta api.ClusterMeta) string {
	return fmt.Sprintf("eksctl-flux-%s-%s", clusterMeta.Region, clusterMeta.Name)
}
