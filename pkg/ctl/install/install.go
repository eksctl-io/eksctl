package install

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	portforward "github.com/justinbarrick/go-k8s-portforward"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/riywo/loginshell"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	fluxapi "github.com/weaveworks/flux/api/v6"
	transport "github.com/weaveworks/flux/http"
	"github.com/weaveworks/flux/http/client"
	"github.com/weaveworks/flux/install"
	"github.com/weaveworks/flux/ssh"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

const namespaceFileName = "namespace.yaml"

// Command will create the `install` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("install", "Install components in a cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, installFluxCmd)

	return verbCmd
}

type installFluxOpts struct {
	templateParams install.TemplateParameters
	gitFluxPath    string
	timeout        time.Duration
	amend          bool
}

func installFluxCmd(rc *cmdutils.ResourceCmd) {
	rc.ClusterConfig = api.NewClusterConfig()
	rc.SetDescription(
		"flux",
		"Bootstrap Flux, installing it in the cluster and initializing its manifests in the Git repository",
		"",
	)
	var opts installFluxOpts
	rc.SetRunFuncWithNameArg(func() error {
		installer, err := newFluxInstaller(rc, &opts)
		if err != nil {
			return err
		}
		return installer.run(context.Background())
	})
	rc.FlagSetGroup.InFlagSet("Flux installation", func(fs *pflag.FlagSet) {
		fs.StringVar(&opts.templateParams.GitURL, "git-url", "",
			"URL of the Git repository to be used by Flux, e.g. git@github.com:<your username>/flux-get-started")
		fs.StringVar(&opts.templateParams.GitBranch, "git-branch", "master",
			"Git branch to be used by Flux")
		fs.StringSliceVar(&opts.templateParams.GitPaths, "git-paths", []string{},
			"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
		fs.StringVar(&opts.templateParams.GitLabel, "git-label", "flux",
			"Git label to keep track of Flux's sync progress; overrides both --git-sync-tag and --git-notes-ref")
		fs.StringVar(&opts.templateParams.GitUser, "git-user", "Flux",
			"Username to use as Git committer")
		fs.StringVar(&opts.templateParams.GitEmail, "git-email", "",
			"Email to use as Git committer")
		fs.StringVar(&opts.gitFluxPath, "git-flux-subdir", "flux/",
			"Directory within the Git repository where to commit the Flux manifests")
		fs.DurationVar(&opts.timeout, "timeout", 20*time.Second,
			"Timeout for I/O operations")
		fs.StringVar(&opts.templateParams.Namespace, "namespace", "flux",
			"Cluster namespace where to install Flux")
		fs.BoolVar(&opts.amend, "amend", false,
			"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
	})
	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, rc.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
	})
}

type fluxInstaller struct {
	opts          *installFluxOpts
	resourceCmd   *cmdutils.ResourceCmd
	k8sConfig     *clientcmdapi.Config
	k8sRestConfig *rest.Config
	k8sClientSet  *kubeclient.Clientset
}

func newFluxInstaller(rc *cmdutils.ResourceCmd, opts *installFluxOpts) (*fluxInstaller, error) {
	if opts.templateParams.GitURL == "" {
		return nil, fmt.Errorf("please supply a valid --git-url argument")
	}
	if opts.templateParams.GitEmail == "" {
		return nil, fmt.Errorf("please supply a valid --git-email argument")
	}

	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return nil, err
	}
	cfg := rc.ClusterConfig
	ctl := eks.New(rc.ProviderConfig, cfg)
	if !ctl.IsSupportedRegion() {
		return nil, cmdutils.ErrUnsupportedRegion(rc.ProviderConfig)
	}
	if err := ctl.GetCredentials(cfg); err != nil {
		return nil, err
	}
	kubernetesClientConfigs, err := ctl.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	k8sConfig := kubernetesClientConfigs.Config

	k8sRestConfig, err := clientcmd.NewDefaultClientConfig(*k8sConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Errorf("cannot create Kubernetes client configuration: %s", err)
	}
	k8sClientSet, err := kubeclient.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, errors.Errorf("cannot create Kubernetes client set: %s", err)
	}
	fi := &fluxInstaller{
		opts:          opts,
		resourceCmd:   rc,
		k8sConfig:     k8sConfig,
		k8sRestConfig: k8sRestConfig,
		k8sClientSet:  k8sClientSet,
	}
	return fi, nil
}

func (fi *fluxInstaller) run(ctx context.Context) error {
	manifests, err := getFluxManifests(fi.opts.templateParams, fi.k8sClientSet)
	if err != nil {
		return fmt.Errorf("failed to create flux manifests: %s", err)
	}

	logger.Info("Cloning %s", fi.opts.templateParams.GitURL)
	cloneDir, err := fi.cloneRepo(ctx)
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %s", fi.opts.templateParams.GitURL, err)
	}
	cleanCloneDir := false
	defer func() {
		if cleanCloneDir {
			os.RemoveAll(cloneDir)
		} else {
			logger.Critical("You may find the local clone of %s used by eksctl at %s",
				fi.opts.templateParams.GitURL,
				cloneDir)
		}
	}()

	logger.Info("Writing Flux manifests")
	fluxManifestDir := filepath.Join(cloneDir, fi.opts.gitFluxPath)
	if err := writeFluxManifests(fluxManifestDir, manifests); err != nil {
		return err
	}

	if fi.opts.amend {
		logger.Info("Stopping to amend the the Flux manifests, please exit the shell when done.")
		if err := runShell(fluxManifestDir); err != nil {
			return err
		}
	}

	// If we need to create a Namespace, we do it synchronously before invoking kubectl apply.
	// Otherwise there is a race between the creation of the namespace and the other resources living in it
	if _, ok := manifests[namespaceFileName]; ok {
		logger.Info("Creating namespace %s", fi.opts.templateParams.Namespace)
		if err := createNamespaceSynchronously(fi.k8sClientSet, fi.opts.templateParams.Namespace); err != nil {
			return err
		}
	}

	logger.Info("Installing Flux into the cluster")
	if err := fi.applyManifests(manifests); err != nil {
		return err
	}

	logger.Info("Waiting for Flux to start")
	fluxSSHKey, err := waitForFluxToStart(ctx, fi.opts, fi.k8sRestConfig, fi.k8sClientSet)
	if err != nil {
		return err
	}
	logger.Info("Flux started successfully")

	logger.Info("Committing and pushing Flux manifests to %s", fi.opts.templateParams.GitURL)
	if err := fi.addCommitAndPushFluxManifests(ctx, cloneDir); err != nil {
		return err
	}
	cleanCloneDir = true

	logger.Info("Flux will operate properly only once it has SSH access to: %s", fi.opts.templateParams.GitURL)
	logger.Info("please add the following Flux public SSH key to your repository:\n%s", fluxSSHKey.Key)
	return nil
}

func (fi *fluxInstaller) cloneRepo(ctx context.Context) (string, error) {
	cloneDir, err := ioutil.TempDir(os.TempDir(), "eksctl-install-flux-clone")
	if err != nil {
		return "", fmt.Errorf("cannot create temporary directory: %s", err)
	}
	cloneCtx, cloneCtxCancel := context.WithTimeout(ctx, fi.opts.timeout)
	defer cloneCtxCancel()
	args := []string{"clone", "-b", fi.opts.templateParams.GitBranch, fi.opts.templateParams.GitURL, cloneDir}
	err = runGitCmd(cloneCtx, cloneDir, args...)
	return cloneDir, err
}

func (fi *fluxInstaller) addCommitAndPushFluxManifests(ctx context.Context, cloneDir string) error {
	// Add
	addCtx, addCtxCancel := context.WithTimeout(ctx, fi.opts.timeout)
	defer addCtxCancel()
	if err := runGitCmd(addCtx, cloneDir, "add", "--", fi.opts.gitFluxPath); err != nil {
		return err
	}

	// Confirm there is something to commit, otherwise move on
	diffCtx, diffCtxCancel := context.WithTimeout(ctx, fi.opts.timeout)
	defer diffCtxCancel()
	if err := runGitCmd(diffCtx, cloneDir, "diff", "--cached", "--quiet", "--", fi.opts.gitFluxPath); err == nil {
		logger.Info("Nothing to commit (the repository contained identical manifests), moving on")
		return nil
	} else if _, ok := err.(*exec.ExitError); !ok {
		return err
	}

	// Commit
	commitCtx, commitCtxCancel := context.WithTimeout(ctx, fi.opts.timeout)
	defer commitCtxCancel()
	args := []string{"commit",
		"-m", "Add Initial Flux configuration",
		fmt.Sprintf("--author=%s <%s>", fi.opts.templateParams.GitUser, fi.opts.templateParams.GitEmail),
	}
	if err := runGitCmd(commitCtx, cloneDir, args...); err != nil {
		return err
	}

	// Push
	pushCtx, pushCtxCancel := context.WithTimeout(ctx, fi.opts.timeout)
	defer pushCtxCancel()
	return runGitCmd(pushCtx, cloneDir, "push")
}

func (fi *fluxInstaller) applyManifests(manifestsMap map[string][]byte) error {
	// TODO: initialise the client elsewhere so that a mock client can easily be dependency-injected for testing purposes.
	client, err := kubernetes.NewRawClient(fi.k8sClientSet, fi.k8sRestConfig)
	if err != nil {
		return err
	}

	manifests := kubernetes.JoinManifestValues(manifestsMap)
	objects, err := kubernetes.NewRawExtensions(manifests)
	if err != nil {
		return err
	}
	for _, object := range objects {
		resource, err := client.NewRawResource(object)
		if err != nil {
			return err
		}
		status, err := resource.CreateOrReplace(false)
		if err != nil {
			return err
		}
		logger.Info(status)
	}
	return nil
}

func runGitCmd(ctx context.Context, dir string, args ...string) error {
	gitCmd := exec.CommandContext(ctx, "git", args...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	gitCmd.Dir = dir
	return gitCmd.Run()
}

func getFluxManifests(params install.TemplateParameters, cs *kubeclient.Clientset) (map[string][]byte, error) {
	params.AdditionalFluxArgs = []string{"--sync-garbage-collection", "--manifest-generation"}
	manifests, err := install.FillInTemplates(params)
	if err != nil {
		return nil, err
	}
	// Verify if the namespace to install Flux already exists and otherwise add a manifest for it
	_, err = cs.CoreV1().Namespaces().Get(params.Namespace, metav1.GetOptions{})
	if err == nil {
		return manifests, nil
	}
	if !k8serrors.IsNotFound(err) {
		return nil, fmt.Errorf("cannot check whether namespace %s exists: %s", params.Namespace, err)
	}
	nsTemplate := `---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    name: %s
  name: %s
`
	ns := fmt.Sprintf(nsTemplate, params.Namespace, params.Namespace)
	manifests[namespaceFileName] = []byte(ns)
	return manifests, nil
}

func writeFluxManifests(baseDir string, manifests map[string][]byte) error {
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return fmt.Errorf("cannot create Flux manifests directory (%s): %s", baseDir, err)
	}
	for fileName, contents := range manifests {
		fullPath := filepath.Join(baseDir, fileName)
		if err := ioutil.WriteFile(fullPath, contents, 0600); err != nil {
			return fmt.Errorf("failed to write Flux manifest file %s: %s", fullPath, err)
		}
	}
	return nil
}

func runShell(workDir string) error {
	shell, err := loginshell.Shell()
	if err != nil {
		return fmt.Errorf("failed to obtain login shell %s: %s", shell, err)
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

func waitForFluxToStart(ctx context.Context, opts *installFluxOpts, restConfig *rest.Config,
	cs *kubeclient.Clientset) (ssh.PublicKey, error) {
	fluxSelector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "name",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"flux"},
			},
		},
	}
	portforwarder := portforward.PortForward{
		Labels:          fluxSelector,
		Config:          restConfig,
		Clientset:       cs,
		DestinationPort: 3030,
		Namespace:       opts.templateParams.Namespace,
	}
	podDeadline := time.Now().Add(30 * time.Second)
	for ; time.Now().Before(podDeadline); time.Sleep(1 * time.Second) {
		err := portforwarder.Start()
		if err == nil {
			defer portforwarder.Stop()
			break
		}
		if !strings.Contains(err.Error(), "Could not find running pod for selector") {
			logger.Warning("Flux is not ready yet (%s), retrying ...", err)
		}
	}
	if time.Now().After(podDeadline) {
		return ssh.PublicKey{}, fmt.Errorf("timed out waiting for flux's pod to be created")
	}
	fluxURL := fmt.Sprintf("http://127.0.0.1:%d/api/flux", portforwarder.ListenPort)
	fluxClient := client.New(http.DefaultClient, transport.NewAPIRouter(), fluxURL, client.Token(""))

	// Obtain SSH key
	var fluxGitConfig fluxapi.GitConfig
	gitConfigDeadline := time.Now().Add(30 * time.Second)
	for ; time.Now().Before(gitConfigDeadline); time.Sleep(100 * time.Millisecond) {
		repoCtx, repoCtxCancel := context.WithTimeout(ctx, opts.timeout)
		var err error
		fluxGitConfig, err = fluxClient.GitRepoConfig(repoCtx, false)
		repoCtxCancel()
		if err == nil {
			defer portforwarder.Stop()
			break
		}
		logger.Warning("Flux is not ready yet (%s), retrying ...", err)
	}
	if time.Now().After(podDeadline) {
		return ssh.PublicKey{}, fmt.Errorf("timed out waiting for Flux to be operative")
	}
	return fluxGitConfig.PublicSSHKey, nil
}

func createNamespaceSynchronously(cs *kubeclient.Clientset, namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: map[string]string{"name": namespace},
		},
	}
	if _, err := cs.CoreV1().Namespaces().Create(ns); err != nil {
		return fmt.Errorf("cannot create namespace %s: %s", namespace, err)
	}
	nsDeadline := time.Now().Add(30 * time.Second)
	for ; time.Now().Before(nsDeadline); time.Sleep(100 * time.Millisecond) {
		_, err := cs.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
		if err == nil {
			break
		}
		if !k8serrors.IsNotFound(err) {
			return fmt.Errorf("cannot check whether namespace %s exists: %s", namespace, err)
		}
	}
	if time.Now().After(nsDeadline) {
		return fmt.Errorf("timed out waiting for namespace %s to be created", namespace)
	}
	return nil
}
