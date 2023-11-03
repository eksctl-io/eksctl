package kubeconfig_test

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	eksctlapi "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/utils/kubectl"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	kubectlfakes "github.com/weaveworks/eksctl/pkg/utils/kubectl/fakes"
)

var errGeneric = fmt.Errorf("generic error")

var _ = Describe("Kubeconfig", func() {
	var configFile *os.File

	var contextName = "test-context"

	var testConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"test-user": {Token: "test-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"test-cluster": {Server: "https://127.0.0.1:8443"}},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {AuthInfo: "test-user", Cluster: "test-cluster", Namespace: "test-ns"}},
		CurrentContext: contextName,
	}
	var exampleSSHKeyPath = "~/.ssh/id_rsa.pub"

	BeforeEach(func() {
		configFile, _ = os.CreateTemp("", "")
	})

	AfterEach(func() {
		os.Remove(configFile.Name())
	})

	var writeConfig = func(filename string) error {
		minikubeSample, err := os.ReadFile("testdata/minikube_sample.golden")
		if err != nil {
			GinkgoT().Fatalf("failed reading .golden: %s", err)
		}

		return os.WriteFile(filename, minikubeSample, os.FileMode(0755))
	}

	It("creating new Kubeconfig", func() {
		filename, err := kubeconfig.Write(configFile.Name(), testConfig, false)
		Expect(err).To(BeNil())

		readConfig, err := clientcmd.LoadFromFile(filename)
		Expect(err).To(BeNil())

		Expect(len(readConfig.Clusters)).To(Equal(1))
		Expect(readConfig.Clusters["test-cluster"].Server).To(Equal(testConfig.Clusters["test-cluster"].Server))

		Expect(len(readConfig.AuthInfos)).To(Equal(1))
		Expect(readConfig.AuthInfos["test-user"].Token).To(Equal(testConfig.AuthInfos["test-user"].Token))

		Expect(len(readConfig.Contexts)).To(Equal(1))
		Expect(readConfig.Contexts["test-context"].Namespace).To(Equal(testConfig.Contexts["test-context"].Namespace))
	})

	It("creating new Kubeconfig with directory", func() {
		filename, err := kubeconfig.Write("/", testConfig, false)
		Expect(err).NotTo(BeNil())
		Expect(filename).To(BeEmpty())
	})

	It("creating new Kubeconfig in non-existent directory", func() {
		tempDir, _ := os.MkdirTemp("", "")
		filename, err := kubeconfig.Write(path.Join(tempDir, "nonexistentdir", "kubeconfig"), testConfig, false)
		Expect(err).To(BeNil())
		Expect(filename).NotTo(BeEmpty())
	})

	It("sets new Kubeconfig context", func() {
		testConfigContext := testConfig
		testConfigContext.CurrentContext = "test-context"

		filename, err := kubeconfig.Write(configFile.Name(), testConfigContext, true)
		Expect(err).To(BeNil())

		readConfig, err := clientcmd.LoadFromFile(filename)
		Expect(err).To(BeNil())

		Expect(readConfig.CurrentContext).To(Equal("test-context"))
	})

	It("merge new Kubeconfig", func() {
		err := writeConfig(configFile.Name())
		Expect(err).To(BeNil())

		filename, err := kubeconfig.Write(configFile.Name(), testConfig, false)
		Expect(err).To(BeNil())

		readConfig, err := clientcmd.LoadFromFile(filename)
		Expect(err).To(BeNil())

		Expect(len(readConfig.Clusters)).To(Equal(2))
		Expect(readConfig.Clusters["test-cluster"].Server).To(Equal(testConfig.Clusters["test-cluster"].Server))
		Expect(readConfig.Clusters["minikube"].Server).To(Equal("https://192.168.64.19:8443"))

		Expect(len(readConfig.AuthInfos)).To(Equal(2))
		Expect(readConfig.AuthInfos["test-user"].Token).To(Equal(testConfig.AuthInfos["test-user"].Token))
		Expect(readConfig.AuthInfos["minikube"].ClientCertificate).To(Equal("/Users/bob/.minikube/client.crt"))

		Expect(len(readConfig.Contexts)).To(Equal(2))
		Expect(readConfig.Contexts["test-context"].Namespace).To(Equal(testConfig.Contexts["test-context"].Namespace))
		Expect(readConfig.Contexts["minikube"].Cluster).To(Equal("minikube"))
	})

	It("merge sets context", func() {
		err := writeConfig(configFile.Name())
		Expect(err).To(BeNil())

		testConfigContext := testConfig
		testConfigContext.CurrentContext = "test-context"

		filename, err := kubeconfig.Write(configFile.Name(), testConfigContext, true)
		Expect(err).To(BeNil())

		readConfig, err := clientcmd.LoadFromFile(filename)
		Expect(err).To(BeNil())

		Expect(readConfig.CurrentContext).To(Equal("test-context"))
	})

	It("merge does not sets context", func() {
		err := writeConfig(configFile.Name())
		Expect(err).To(BeNil())

		testConfigContext := testConfig
		testConfigContext.CurrentContext = "test-context"

		filename, err := kubeconfig.Write(configFile.Name(), testConfigContext, false)
		Expect(err).To(BeNil())

		readConfig, err := clientcmd.LoadFromFile(filename)
		Expect(err).To(BeNil())

		Expect(readConfig.CurrentContext).To(Equal("minikube"))
	})

	var (
		kubeconfigPathToRestore string
		hasKubeconfigPath       bool
	)

	ChangeKubeconfig := func() {
		if _, err := os.Stat(configFile.Name()); os.IsNotExist(err) {
			GinkgoT().Fatal(err)
		}

		kubeconfigPathToRestore, hasKubeconfigPath = os.LookupEnv("KUBECONFIG")
		os.Setenv("KUBECONFIG", configFile.Name())
	}

	Context("delete config", func() {
		// Default cluster name is 'foo' and region is 'us-west-2'
		var apiClusterConfigSample = eksctlapi.ClusterConfig{
			Metadata: &eksctlapi.ClusterMeta{
				Region: "us-west-2",
				Name:   "foo",
				Tags:   map[string]string{},
			},
			NodeGroups: []*eksctlapi.NodeGroup{
				{
					NodeGroupBase: &eksctlapi.NodeGroupBase{
						InstanceType:      "m5.large",
						AvailabilityZones: []string{"us-west-2b", "us-west-2a", "us-west-2c"},
						PrivateNetworking: false,
						SSH: &eksctlapi.NodeGroupSSH{
							Allow:         eksctlapi.Disabled(),
							PublicKeyPath: &exampleSSHKeyPath,
							PublicKey:     nil,
							PublicKeyName: nil,
						},
						ScalingConfig: &eksctlapi.ScalingConfig{},
						IAM: &eksctlapi.NodeGroupIAM{
							AttachPolicyARNs: []string(nil),
							InstanceRoleARN:  "",
							InstanceRoleName: "",
						},
						AMI:            "",
						MaxPodsPerNode: 0,
					},
				},
			},
			VPC: &eksctlapi.ClusterVPC{
				Network: eksctlapi.Network{
					ID:   "",
					CIDR: nil,
				},
				SecurityGroup: "",
			},
			AvailabilityZones: []string{"us-west-2b", "us-west-2a", "us-west-2c"},
		}

		var (
			emptyClusterAsBytes               []byte
			oneClusterAsBytes                 []byte
			twoClustersAsBytes                []byte
			oneClusterWithoutContextAsBytes   []byte
			oneClusterWithStaleContextAsBytes []byte
		)

		// Returns an ClusterConfig with a cluster name equal to the provided clusterName.
		GetClusterConfig := func(clusterName string) *eksctlapi.ClusterConfig {
			apiClusterConfig := apiClusterConfigSample
			apiClusterConfig.Metadata.Name = clusterName
			return &apiClusterConfig
		}

		RestoreKubeconfig := func() {
			if hasKubeconfigPath {
				os.Setenv("KUBECONFIG", kubeconfigPathToRestore)
			} else {
				os.Unsetenv("KUBECONFIG")
			}
		}

		BeforeEach(func() {
			ChangeKubeconfig()

			var err error

			if emptyClusterAsBytes, err = os.ReadFile("testdata/empty_cluster.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if oneClusterAsBytes, err = os.ReadFile("testdata/one_cluster.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if twoClustersAsBytes, err = os.ReadFile("testdata/two_clusters.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if oneClusterWithoutContextAsBytes, err = os.ReadFile("testdata/one_cluster_without_context.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if oneClusterWithStaleContextAsBytes, err = os.ReadFile("testdata/one_cluster_with_stale_context.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			_, err = configFile.Write(twoClustersAsBytes)
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			RestoreKubeconfig()
		})

		It("removes the only current cluster from the kubeconfig if the kubeconfig file includes the cluster", func() {
			_, err := configFile.Write(oneClusterAsBytes)
			Expect(err).To(BeNil())

			existingClusterConfig := GetClusterConfig("cluster-one")
			kubeconfig.MaybeDeleteConfig(existingClusterConfig.Metadata)

			configFileAsBytes, err := os.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(emptyClusterAsBytes), "Failed to delete cluster from config")
		})

		It("removes current cluster from the kubeconfig if the kubeconfig file includes the cluster", func() {
			existingClusterConfig := GetClusterConfig("cluster-one")
			kubeconfig.MaybeDeleteConfig(existingClusterConfig.Metadata)

			configFileAsBytes, err := os.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(oneClusterWithoutContextAsBytes), "Failed to delete cluster from config")
		})

		It("removes current cluster from the kubeconfig and clears stale context", func() {
			_, err := configFile.Write(oneClusterWithStaleContextAsBytes)
			Expect(err).To(BeNil())

			existingClusterConfig := GetClusterConfig("cluster-one")
			kubeconfig.MaybeDeleteConfig(existingClusterConfig.Metadata)

			configFileAsBytes, err := os.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(oneClusterWithoutContextAsBytes), "Updated config")

		})

		It("removes a secondary cluster from the kubeconfig if the kubeconfig file includes the cluster", func() {
			existingClusterConfig := GetClusterConfig("cluster-two")
			kubeconfig.MaybeDeleteConfig(existingClusterConfig.Metadata)

			configFileAsBytes, err := os.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(oneClusterAsBytes), "Failed to delete cluster from config")
		})

		It("not change the kubeconfig if the kubeconfig does not include the cluster", func() {
			nonExistentClusterConfig := GetClusterConfig("not-a-cluster")
			kubeconfig.MaybeDeleteConfig(nonExistentClusterConfig.Metadata)

			configFileAsBytes, err := os.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(twoClustersAsBytes), "Should not change")
		})
	})

	It("safely handles concurrent read-modify-write operations", func() {
		var (
			oneCluster  *clientcmdapi.Config
			twoClusters *clientcmdapi.Config
		)
		ChangeKubeconfig()

		var err error
		tmp, err := os.CreateTemp("", "")
		Expect(err).To(BeNil())

		{
			if oneCluster, err = clientcmd.LoadFromFile("testdata/one_cluster.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if twoClusters, err = clientcmd.LoadFromFile("testdata/two_clusters.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}
		}

		var wg sync.WaitGroup
		multiplier := 3
		iters := 10
		for i := 0; i < multiplier; i++ {
			for k := 0; k < iters; k++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := kubeconfig.Write(tmp.Name(), *oneCluster, false)
					Expect(err).To(BeNil())
				}()
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := kubeconfig.Write(tmp.Name(), *twoClusters, false)
					Expect(err).To(BeNil())
				}()
			}
		}

		wg.Wait()

	})
	Context("AppendAuthenticator", func() {
		var (
			config      *clientcmdapi.Config
			clusterInfo *eksctlapi.ClusterConfig
		)
		BeforeEach(func() {
			config = &clientcmdapi.Config{
				AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
				CurrentContext: "test",
			}
			clusterInfo = &eksctlapi.ClusterConfig{
				Metadata: &eksctlapi.ClusterMeta{
					Region: "us-west-2",
					Name:   "name",
				},
			}
			kubeconfig.SetNewVersionManager(func() kubectl.KubernetesVersionManager {
				fakeClient := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeClient.ClientVersionReturns("", errGeneric)
				return fakeClient
			})
		})
		AfterEach(func() {
			kubeconfig.SetNewVersionManager(kubectl.NewVersionManager)
		})

		It("writes the right api version if aws-iam-authenticator version is below 0.5.3", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `{"Version":"0.5.1","Commit":"85e50980d9d916ae95882176c18f14ae145f916f"}`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSIAMAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1alpha1"))
		})
		It("writes the right api version if aws-iam-authenticator version is above 0.5.3", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `{"Version":"0.5.5","Commit":"85e50980d9d916ae95882176c18f14ae145f916f"}`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSIAMAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1beta1"))
		})
		It("writes the right api version if aws-iam-authenticator version equals 0.5.3", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `{"Version":"0.5.3","Commit":"85e50980d9d916ae95882176c18f14ae145f916f"}`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSIAMAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1beta1"))
		})
		It("defaults to alpha1 if we fail to detect aws-iam-authenticator version", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), "fail")
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSIAMAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1alpha1"))
		})
		It("defaults to alpha1 if we fail to parse the output", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), "not-json-output")
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSIAMAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1alpha1"))
		})
		It("defaults to alpha1 if we can't parse the version because it's a dev version", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `{"Version":"git-85e50980","Commit":"85e50980d9d916ae95882176c18f14ae145f916f"}`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSIAMAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1alpha1"))
		})
		It("defaults to beta1 if we detect kubectl 1.24.0 or above", func() {
			kubeconfig.SetNewVersionManager(func() kubectl.KubernetesVersionManager {
				fakeClient := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeClient.ClientVersionReturns("1.24.0", nil)
				return fakeClient
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSEKSAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1beta1"))
		})
		It("doesn't default to beta1 if we detect kubectl version lower than 1.24.0", func() {
			kubeconfig.SetNewVersionManager(func() kubectl.KubernetesVersionManager {
				fakeClient := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeClient.ClientVersionReturns("1.23.6", nil)
				return fakeClient
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSIAMAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1alpha1"))
		})
		It("defaults to beta1 if we detect aws-cli v1 is at or above 1.23.9", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `aws-cli/1.23.9 Python/3.8.8 Linux/5.4.181-109.354.amzn2int.x86_64 exe/x86_64.amzn.2 prompt/off`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSEKSAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1beta1"))
		})
		It("doesn't default to beta1 if we detect aws-cli v1 is below 1.23.9", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `aws-cli/1.21.9 Python/3.8.8 Linux/5.4.181-109.354.amzn2int.x86_64 exe/x86_64.amzn.2 prompt/off`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSEKSAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1alpha1"))
		})
		It("defaults to beta1 if we detect aws-cli v2 is at or above 2.6.3", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `aws-cli/2.6.3 Python/3.8.8 Linux/5.4.181-109.354.amzn2int.x86_64 exe/x86_64.amzn.2 prompt/off`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSEKSAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1beta1"))
		})
		It("doesn't default to beta1 if we detect aws-cli v2 below 2.6.3", func() {
			kubeconfig.SetExecCommand(func(name string, arg ...string) *exec.Cmd {
				return exec.Command(filepath.Join("testdata", "fake-version"), `aws-cli/2.4.3 Python/3.8.8 Linux/5.4.181-109.354.amzn2int.x86_64 exe/x86_64.amzn.2 prompt/off`)
			})
			kubeconfig.AppendAuthenticator(config, clusterInfo, kubeconfig.AWSEKSAuthenticator, "", "")
			Expect(config.AuthInfos["test"].Exec.APIVersion).To(Equal("client.authentication.k8s.io/v1alpha1"))
		})
	})

	type checkAllCommandsEntry struct {
		kubeconfigPath                  string
		mockKubernetesVersionManager    func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func())
		overrideExecLookPath            func(file string) (string, error)
		overrideAuthenticatorLookupMock func() (string, bool)
		expectedErr                     string
	}

	DescribeTable("CheckAllCommands", func(e checkAllCommandsEntry) {
		kubeconfig.SetExecLookPath(func(file string) (string, error) {
			return "path", nil
		})
		if e.overrideExecLookPath != nil {
			kubeconfig.SetExecLookPath(e.overrideExecLookPath)
		}

		kubeconfig.SetLookupAuthenticator(func() (string, bool) {
			return kubeconfig.AWSIAMAuthenticator, true
		})
		if e.overrideAuthenticatorLookupMock != nil {
			kubeconfig.SetLookupAuthenticator(e.overrideAuthenticatorLookupMock)
		}

		mockManagerFunc, assertFakeManagerCalls := e.mockKubernetesVersionManager()
		kubeconfig.SetNewVersionManager(mockManagerFunc)

		err := kubeconfig.CheckAllCommands(
			e.kubeconfigPath,
			/* isContextSet */ false,
			/* contextName */ "ctx-name",
			/* env */ []string{"env"},
		)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
		assertFakeManagerCalls()
	},
		Entry("fails to lookup authenticator", checkAllCommandsEntry{
			overrideAuthenticatorLookupMock: func() (string, bool) {
				return "", false
			},
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				return mockManagerFunc, assertFakeManagerCalls
			},
			expectedErr: "could not find any of the authenticator commands",
		}),

		Entry("fails to lookup kubectl", checkAllCommandsEntry{
			overrideExecLookPath: func(file string) (string, error) {
				return "", errGeneric
			},
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				return mockManagerFunc, assertFakeManagerCalls
			},
			expectedErr: "kubectl not found, v1.10.0 or newer is required",
		}),

		Entry("fails to fetch client version", checkAllCommandsEntry{
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				fakeManager := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeManager.ClientVersionReturns("", errGeneric)
				mockManagerFunc = func() kubectl.KubernetesVersionManager {
					return fakeManager
				}
				return mockManagerFunc, assertFakeManagerCalls
			},
			expectedErr: "getting kubectl version",
		}),

		Entry("fails to validate client version", checkAllCommandsEntry{
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				fakeManager := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeManager.ClientVersionReturns("v1.27.4", nil)
				fakeManager.ValidateVersionReturns(errGeneric)
				mockManagerFunc = func() kubectl.KubernetesVersionManager {
					return fakeManager
				}
				return mockManagerFunc, assertFakeManagerCalls
			},
			expectedErr: "validating kubectl version",
		}),

		Entry("finishes successfully when kubeconfigPath is not set", checkAllCommandsEntry{
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				fakeManager := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeManager.ClientVersionReturns("v1.27.4", nil)
				fakeManager.ValidateVersionReturns(nil)
				mockManagerFunc = func() kubectl.KubernetesVersionManager {
					return fakeManager
				}
				assertFakeManagerCalls = func() {
					Expect(fakeManager.ServerVersionCallCount()).To(Equal(0))
				}
				return mockManagerFunc, assertFakeManagerCalls
			},
		}),

		Entry("fails to fetch server version", checkAllCommandsEntry{
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				fakeManager := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeManager.ClientVersionReturns("v1.27.4", nil)
				fakeManager.ValidateVersionReturns(nil)
				fakeManager.ServerVersionReturns("", errGeneric)
				mockManagerFunc = func() kubectl.KubernetesVersionManager {
					return fakeManager
				}
				return mockManagerFunc, assertFakeManagerCalls
			},
			kubeconfigPath: "path",
			expectedErr:    "getting Kubernetes version on EKS cluster",
		}),

		Entry("fails to validate server version", checkAllCommandsEntry{
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				fakeManager := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeManager.ClientVersionReturns("v1.27.4", nil)
				fakeManager.ValidateVersionReturnsOnCall(0, nil)
				fakeManager.ServerVersionReturns("1.24.7-eks-2f45561", nil)
				fakeManager.ValidateVersionReturnsOnCall(1, errGeneric)
				mockManagerFunc = func() kubectl.KubernetesVersionManager {
					return fakeManager
				}
				return mockManagerFunc, assertFakeManagerCalls
			},
			kubeconfigPath: "path",
			expectedErr:    "validating Kubernetes version returned by EKS API",
		}),

		Entry("finishes successfully when kubeconfigPath is set", checkAllCommandsEntry{
			mockKubernetesVersionManager: func() (mockManagerFunc func() kubectl.KubernetesVersionManager, assertFakeManagerCalls func()) {
				fakeManager := &kubectlfakes.FakeKubernetesVersionManager{}
				fakeManager.ClientVersionReturns("v1.27.4", nil)
				fakeManager.ValidateVersionReturnsOnCall(0, nil)
				fakeManager.ServerVersionReturns("v1.26.3", nil)
				fakeManager.ValidateVersionReturnsOnCall(1, nil)
				mockManagerFunc = func() kubectl.KubernetesVersionManager {
					return fakeManager
				}
				assertFakeManagerCalls = func() {
					env, args := fakeManager.ServerVersionArgsForCall(0)
					Expect(env).To(Equal([]string{"env"}))
					Expect(args).To(Equal([]string{"--kubeconfig=path", "--context=ctx-name"}))
				}
				return mockManagerFunc, assertFakeManagerCalls
			},
			kubeconfigPath: "path",
		}),
	)
})
