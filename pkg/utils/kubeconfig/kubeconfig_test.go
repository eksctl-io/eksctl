package kubeconfig_test

import (
	"io/ioutil"
	"os"
	"path"
	"sync"

	eksctlapi "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubeconfig", func() {
	var configFile *os.File

	var contextName = "test-context"

	var testConfig = api.Config{
		AuthInfos: map[string]*api.AuthInfo{
			"test-user": {Token: "test-token"}},
		Clusters: map[string]*api.Cluster{
			"test-cluster": {Server: "https://127.0.0.1:8443"}},
		Contexts: map[string]*api.Context{
			contextName: {AuthInfo: "test-user", Cluster: "test-cluster", Namespace: "test-ns"}},
		CurrentContext: contextName,
	}
	var exampleSSHKeyPath = "~/.ssh/id_rsa.pub"

	BeforeEach(func() {
		configFile, _ = ioutil.TempFile("", "")
	})

	AfterEach(func() {
		os.Remove(configFile.Name())
	})

	var writeConfig = func(filename string) error {
		minikubeSample, err := ioutil.ReadFile("testdata/minikube_sample.golden")
		if err != nil {
			GinkgoT().Fatalf("failed reading .golden: %s", err)
		}

		return ioutil.WriteFile(filename, minikubeSample, os.FileMode(0755))
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
		tempDir, _ := ioutil.TempDir("", "")
		filename, err := kubeconfig.Write(path.Join(tempDir, "nonexistentdir", "kubeconfig"), testConfig, false)
		Expect(err).To(BeNil())
		Expect(filename).ToNot(BeEmpty())
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
					AMI:               "",
					InstanceType:      "m5.large",
					AvailabilityZones: []string{"us-west-2b", "us-west-2a", "us-west-2c"},
					PrivateNetworking: false,
					SSH: &eksctlapi.NodeGroupSSH{
						Allow:         eksctlapi.Disabled(),
						PublicKeyPath: &exampleSSHKeyPath,
						PublicKey:     nil,
						PublicKeyName: nil,
					},
					DesiredCapacity: nil,
					MinSize:         nil,
					MaxSize:         nil,
					MaxPodsPerNode:  0,
					IAM: &eksctlapi.NodeGroupIAM{
						AttachPolicyARNs: []string(nil),
						InstanceRoleARN:  "",
						InstanceRoleName: "",
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

			if emptyClusterAsBytes, err = ioutil.ReadFile("testdata/empty_cluster.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if oneClusterAsBytes, err = ioutil.ReadFile("testdata/one_cluster.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if twoClustersAsBytes, err = ioutil.ReadFile("testdata/two_clusters.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if oneClusterWithoutContextAsBytes, err = ioutil.ReadFile("testdata/one_cluster_without_context.golden"); err != nil {
				GinkgoT().Fatalf("failed reading .golden: %v", err)
			}

			if oneClusterWithStaleContextAsBytes, err = ioutil.ReadFile("testdata/one_cluster_with_stale_context.golden"); err != nil {
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

			configFileAsBytes, err := ioutil.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(emptyClusterAsBytes), "Failed to delete cluster from config")
		})

		It("removes current cluster from the kubeconfig if the kubeconfig file includes the cluster", func() {
			existingClusterConfig := GetClusterConfig("cluster-one")
			kubeconfig.MaybeDeleteConfig(existingClusterConfig.Metadata)

			configFileAsBytes, err := ioutil.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(oneClusterWithoutContextAsBytes), "Failed to delete cluster from config")
		})

		It("removes current cluster from the kubeconfig and clears stale context", func() {
			_, err := configFile.Write(oneClusterWithStaleContextAsBytes)
			Expect(err).To(BeNil())

			existingClusterConfig := GetClusterConfig("cluster-one")
			kubeconfig.MaybeDeleteConfig(existingClusterConfig.Metadata)

			configFileAsBytes, err := ioutil.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(oneClusterWithoutContextAsBytes), "Updated config")

		})

		It("removes a secondary cluster from the kubeconfig if the kubeconfig file includes the cluster", func() {
			existingClusterConfig := GetClusterConfig("cluster-two")
			kubeconfig.MaybeDeleteConfig(existingClusterConfig.Metadata)

			configFileAsBytes, err := ioutil.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(oneClusterAsBytes), "Failed to delete cluster from config")
		})

		It("not change the kubeconfig if the kubeconfig does not include the cluster", func() {
			nonExistentClusterConfig := GetClusterConfig("not-a-cluster")
			kubeconfig.MaybeDeleteConfig(nonExistentClusterConfig.Metadata)

			configFileAsBytes, err := ioutil.ReadFile(configFile.Name())
			Expect(err).To(BeNil())
			Expect(configFileAsBytes).To(MatchYAML(twoClustersAsBytes), "Should not change")
		})
	})

	It("safely handles concurrent read-modify-write operations", func() {
		var (
			oneClusterAsBytes  []byte
			twoClustersAsBytes []byte
		)
		ChangeKubeconfig()

		var err error

		if oneClusterAsBytes, err = ioutil.ReadFile("testdata/one_cluster.golden"); err != nil {
			GinkgoT().Fatalf("failed reading .golden: %v", err)
		}

		if twoClustersAsBytes, err = ioutil.ReadFile("testdata/two_clusters.golden"); err != nil {
			GinkgoT().Fatalf("failed reading .golden: %v", err)
		}

		var wg sync.WaitGroup
		multiplier := 3
		iters := 100
		for i := 0; i < multiplier; i++ {
			for k := 0; k < iters; k++ {
				wg.Add(2)
				go func() {
					defer wg.Done()
					_, err := configFile.Write(oneClusterAsBytes)
					Expect(err).To(BeNil())
				}()
				go func() {
					defer wg.Done()
					_, err := configFile.Write(twoClustersAsBytes)
					Expect(err).To(BeNil())
				}()
			}
		}

		wg.Wait()

	})
})
