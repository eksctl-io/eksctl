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
	corev1 "k8s.io/api/core/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	tillerinstall "k8s.io/helm/cmd/helm/installer"
	"sigs.k8s.io/yaml"
)

const (
	fluxNamespaceFileName       = "flux-namespace.yaml"
	fluxPrivateSSHKeyFileName   = "flux-secret.yaml"
	fluxPrivateSSHKeySecretName = "flux-git-deploy"
	helmTLSValidFor             = 5 * 365 * 24 * time.Hour // 5 years
	tillerManifestPrefix        = "tiller-"
	tillerServiceName           = "tiller-deploy" // do not change at will, hardcoded in Tiller's manifest generation API
	tillerServiceAccountName    = "tiller"
	tillerImageSpec             = "gcr.io/kubernetes-helm/tiller:v2.14.3"
	tillerRBACTemplate          = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: %[1]s
  namespace: %[2]s
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: %[1]s
    namespace: %[2]s

---
# Helm client serviceaccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: helm
  namespace: %[2]s
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: tiller-user
  namespace: %[2]s
rules:
- apiGroups:
  - ""
  resources:
  - pods/portforward
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - list
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: tiller-user-binding
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: tiller-user
subjects:
- kind: ServiceAccount
  name: helm
  namespace: %[2]s
`
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
	pki, pkiPaths, err := fi.setupPKI()
	if err != nil {
		return "", err
	}

	logger.Info("Generating manifests")
	manifests, secrets, err := fi.getManifestsAndSecrets(pki, pkiPaths)
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

	if len(secrets) > 0 {
		logger.Info("Applying Helm TLS Secret(s)")
		if err := fi.applySecrets(secrets); err != nil {
			return "", err
		}
		logger.Warning("Note: certificate secrets aren't added to the Git repository for security reasons")
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

func (fi *Installer) setupPKI() (*publicKeyInfrastructure, *publicKeyInfrastructurePaths, error) {
	if !fi.opts.WithHelm {
		return nil, nil, nil
	}

	logger.Info("Generating public key infrastructure for the Helm Operator and Tiller")
	logger.Info("  this may take up to a minute, please be patient")
	tillerHost := tillerServiceName + "." + fi.opts.Namespace
	pki, err := newPKI(tillerHost, helmTLSValidFor, 4096)
	if err != nil {
		return nil, nil, err
	}
	baseDir, err := ioutil.TempDir(os.TempDir(), "eksctl-helm-pki")
	if err != nil {
		return nil, nil, errors.Errorf("cannot create temporary directory %q to output PKI files", baseDir)
	}
	pkiPaths := &publicKeyInfrastructurePaths{
		caKey:             filepath.Join(baseDir, "ca-key.pem"),
		caCertificate:     filepath.Join(baseDir, "ca.pem"),
		serverKey:         filepath.Join(baseDir, "key.pem"),
		serverCertificate: filepath.Join(baseDir, "cert.pem"),
		clientKey:         filepath.Join(baseDir, "client-key.pem"),
		clientCertificate: filepath.Join(baseDir, "client-cert.pem"),
	}
	if err = pki.saveTo(pkiPaths); err != nil {
		return nil, nil, err
	}
	logger.Warning("Public key infrastructure files were written into directory %q", baseDir)
	logger.Warning("please move the files into a safe place or delete them")
	return pki, pkiPaths, nil
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

func (fi *Installer) applySecrets(secrets []*corev1.Secret) error {
	secretMap := map[string][]byte{}
	for _, secret := range secrets {
		id := fmt.Sprintf("secret/%s/%s", secret.Namespace, secret.Name)
		secretBytes, err := yaml.Marshal(secret)
		if err != nil {
			return errors.Wrapf(err, "cannot serialize secret %s", id)
		}
		secretMap[id] = secretBytes
	}
	return fi.applyManifests(secretMap)
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

func (fi *Installer) getManifestsAndSecrets(pki *publicKeyInfrastructure,
	pkiPaths *publicKeyInfrastructurePaths) (map[string][]byte, []*corev1.Secret, error) {
	manifests := map[string][]byte{}
	secrets := []*corev1.Secret{}

	// Flux
	var err error
	if manifests, err = getFluxManifests(fi.opts, fi.k8sClientSet); err != nil {
		return nil, nil, err
	}

	// Helm Operator
	if !fi.opts.WithHelm {
		return manifests, secrets, nil
	}
	helmOpManifests, helmOpSecrets, err := getHelmOpManifestsAndSecrets(fi.opts.Namespace, pki)
	if err != nil {
		return nil, nil, err
	}
	manifests = mergeMaps(manifests, helmOpManifests)
	secrets = append(secrets, helmOpSecrets...)

	// Tiller
	tillerManifests, tillerSecrets, err := getTillerManifestsAndSecrets(fi.opts.Namespace, fi.k8sClientSet, pkiPaths)
	if err != nil {
		return nil, nil, err
	}
	manifests = mergeMaps(manifests, tillerManifests)
	secrets = append(secrets, tillerSecrets...)

	return manifests, secrets, nil
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

func getHelmOpManifestsAndSecrets(namespace string, pki *publicKeyInfrastructure) (map[string][]byte, []*corev1.Secret, error) {
	var secrets []*corev1.Secret
	helmOpParameters := helmopinstall.TemplateParameters{
		Namespace:       namespace,
		TillerNamespace: namespace,
		SSHSecretName:   "flux-git-deploy", // determined by the generated Flux manifests
	}
	if pki != nil {
		helmOpParameters.EnableTillerTLS = true
		helmOpParameters.TillerTLSCACertContent = string(pki.caCertificate)
		helmOpTLSSecretName := "flux-helm-tls-cert"
		helmOpParameters.TillerTLSCertSecretName = helmOpTLSSecretName
		tlsSecret := &corev1.Secret{
			Type: corev1.SecretTypeTLS,
			Data: map[string][]byte{
				corev1.TLSCertKey:       pki.clientCertificate,
				corev1.TLSPrivateKeyKey: pki.clientKey,
			},
		}
		tlsSecret.Kind = "Secret"
		tlsSecret.APIVersion = "v1"
		tlsSecret.Name = helmOpTLSSecretName
		tlsSecret.Namespace = namespace
		secrets = append(secrets, tlsSecret)
	}
	manifests, err := helmopinstall.FillInTemplates(helmOpParameters)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Helm Operator Manifests")
	}
	return manifests, secrets, nil
}

func getTillerManifestsAndSecrets(namespace string, cs kubeclient.Interface,
	pkiPaths *publicKeyInfrastructurePaths) (map[string][]byte, []*corev1.Secret, error) {
	manifests := map[string][]byte{}
	tillerOptions := &tillerinstall.Options{
		Namespace:                    namespace,
		ServiceAccount:               tillerServiceAccountName,
		AutoMountServiceAccountToken: true,
		MaxHistory:                   10,
		EnableTLS:                    true,
		VerifyTLS:                    true,
		TLSKeyFile:                   pkiPaths.serverKey,
		TLSCertFile:                  pkiPaths.serverCertificate,
		TLSCaCertFile:                pkiPaths.caCertificate,
		UseCanary:                    false,
		ImageSpec:                    tillerImageSpec,
	}
	tillerDeployment, err := tillerinstall.Deployment(tillerOptions)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate Tiller's Deployment")
	}
	tillerDeploymentBytes, err := yaml.Marshal(tillerDeployment)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to serialize Tiller's Deployment")
	}
	manifests[tillerManifestPrefix+"dep.yaml"] = tillerDeploymentBytes
	tillerService := tillerinstall.Service(namespace)
	tillerServiceBytes, err := yaml.Marshal(tillerService)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to serialize Tiller's Deployment")
	}
	manifests[tillerManifestPrefix+"svc.yaml"] = tillerServiceBytes
	tillerRBACManifests := fmt.Sprintf(tillerRBACTemplate, tillerServiceAccountName, namespace)
	manifests[tillerManifestPrefix+"rbac.yaml"] = []byte(tillerRBACManifests)
	tillerSecret, err := tillerinstall.Secret(tillerOptions)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate Tiller's Secret")
	}
	return manifests, []*corev1.Secret{tillerSecret}, nil
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
