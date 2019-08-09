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

	"github.com/weaveworks/eksctl/pkg/git"

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
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

func installFluxCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"flux",
		"Bootstrap Flux, installing it in the cluster and initializing its manifests in the specified Git repository",
		"",
	)
	var opts installFluxOpts
	cmd.SetRunFuncWithNameArg(func() error {
		installer, err := newFluxInstaller(context.Background(), cmd, &opts)
		if err != nil {
			return err
		}
		return installer.run(context.Background())
	})
	cmd.FlagSetGroup.InFlagSet("Flux installation", func(fs *pflag.FlagSet) {
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
		fs.StringVar(&opts.templateParams.Namespace, "namespace", "flux",
			"Cluster namespace where to install Flux")
		fs.BoolVar(&opts.amend, "amend", false,
			"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &opts.timeout, 20*time.Second)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	cmd.ProviderConfig.WaitTimeout = opts.timeout
}

type fluxInstaller struct {
	opts          *installFluxOpts
	cmd           *cmdutils.Cmd
	k8sConfig     *clientcmdapi.Config
	k8sRestConfig *rest.Config
	k8sClientSet  *kubeclient.Clientset
	gitClient     *git.Client
}

func newFluxInstaller(ctx context.Context, cmd *cmdutils.Cmd, opts *installFluxOpts) (*fluxInstaller, error) {
	if opts.templateParams.GitURL == "" {
		return nil, fmt.Errorf("please supply a valid --git-url argument")
	}
	if opts.templateParams.GitEmail == "" {
		return nil, fmt.Errorf("please supply a valid --git-email argument")
	}

	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return nil, err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return nil, err
	}

	if err := ctl.CheckAuth(); err != nil {
		return nil, err
	}

	if err := ctl.RefreshClusterConfig(cfg); err != nil {
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

	gitClient := git.NewGitClient(ctx, opts.templateParams.GitUser, opts.templateParams.GitEmail, opts.timeout)
	fi := &fluxInstaller{
		opts:          opts,
		cmd:           cmd,
		k8sConfig:     k8sConfig,
		k8sRestConfig: k8sRestConfig,
		k8sClientSet:  k8sClientSet,
		gitClient:     gitClient,
	}
	return fi, nil
}

func (fi *fluxInstaller) run(ctx context.Context) error {
	manifests, err := getFluxManifests(fi.opts.templateParams, fi.k8sClientSet)
	if err != nil {
		return fmt.Errorf("failed to create flux manifests: %s", err)
	}

	logger.Info("Cloning %s", fi.opts.templateParams.GitURL)
	cloneDir, err := fi.gitClient.CloneRepo("eksctl-install-flux-clone-", fi.opts.templateParams.GitBranch, fi.opts.templateParams.GitURL)
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %s", fi.opts.templateParams.GitURL, err)
	}
	cleanCloneDir := false
	defer func() {
		if cleanCloneDir {
			fi.gitClient.DeleteLocalRepo()
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
		// Re-read the manifests, as they may have changed:
		manifests, err = readFluxManifests(fluxManifestDir)
		if err != nil {
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
	if err := fi.addFilesToRepo(ctx, cloneDir); err != nil {
		return err
	}
	cleanCloneDir = true

	logger.Info("Flux will only operate properly once it has write-access to the Git repository")
	logger.Info("please configure %s so that the following Flux SSH public key has write access to it\n%s",
		fi.opts.templateParams.GitURL, fluxSSHKey.Key)
	return nil
}

func (fi *fluxInstaller) addFilesToRepo(ctx context.Context, cloneDir string) error {
	if err := fi.gitClient.Add(fi.opts.gitFluxPath); err != nil {
		return err
	}

	// Confirm there is something to commit, otherwise move on
	if err := fi.gitClient.Commit("Add Initial Flux configuration"); err != nil {
		return err
	}

	// git push
	if err := fi.gitClient.Push(); err != nil {
		return err
	}
	return nil
}

func (fi *fluxInstaller) applyManifests(manifestsMap map[string][]byte) error {
	// TODO: initialise the client elsewhere so that a mock client can easily be dependency-injected for testing purposes.
	client, err := kubernetes.NewRawClient(fi.k8sClientSet, fi.k8sRestConfig)
	if err != nil {
		return err
	}

	// If the namespace needs to be created, do it first before any other
	// resource which should potentially be created within this namespace.
	// Otherwise, creation of these resources will fail.
	if namespace, ok := manifestsMap[namespaceFileName]; ok {
		if err := client.CreateOrReplace(namespace, false); err != nil {
			return err
		}
		delete(manifestsMap, namespaceFileName)
	}

	var manifestValues [][]byte
	for _, manifest := range manifestsMap {
		manifestValues = append(manifestValues, manifest)
	}
	manifests := kubernetes.ConcatManifests(manifestValues...)
	return client.CreateOrReplace(manifests, false)
}

func getFluxManifests(params install.TemplateParameters, clientSet *kubeclient.Clientset) (map[string][]byte, error) {
	params.AdditionalFluxArgs = []string{"--sync-garbage-collection", "--manifest-generation"}
	manifests, err := install.FillInTemplates(params)
	if err != nil {
		return nil, err
	}
	created, err := kubernetes.CheckNamespaceExists(clientSet, params.Namespace)
	if err != nil {
		return nil, err
	}
	if !created {
		// Add the namespace to the manifests so that it later gets serialised,
		// added to the Git repository, and added to the cluster.
		manifests[namespaceFileName] = kubernetes.NewNamespaceYAML(params.Namespace)
	}
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

func readFluxManifests(baseDir string) (map[string][]byte, error) {
	manifestFiles, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list Flux manifest files in %s: %s", baseDir, err)
	}
	manifests := map[string][]byte{}
	for _, manifestFile := range manifestFiles {
		manifest, err := ioutil.ReadFile(manifestFile.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read Flux manifest file %s: %s", manifestFile.Name(), err)
		}
		manifests[manifestFile.Name()] = manifest
	}
	return manifests, nil
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
	// Return Flux's public SSH key as it later needs to be printed/logged for
	// the end-user to take action and add it to their Git repository.
	return fluxGitConfig.PublicSSHKey, nil
}
