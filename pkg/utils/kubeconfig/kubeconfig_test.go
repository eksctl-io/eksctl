package kubeconfig_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	eksctlapi "github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestCreateNewKubeConfig(t *testing.T) {
	configFile, _ := ioutil.TempFile("", "")
	defer os.Remove(configFile.Name())

	testConfig := api.Config{
		AuthInfos: map[string]*api.AuthInfo{
			"test-user": {Token: "test-token"}},
		Clusters: map[string]*api.Cluster{
			"test-cluster": {Server: "https://127.0.0.1:8443"}},
		Contexts: map[string]*api.Context{
			"test-context": {AuthInfo: "test-user", Cluster: "test-cluster", Namespace: "test-ns"}},
	}

	filename, err := kubeconfig.Write(configFile.Name(), &testConfig, false)
	if err != nil {
		t.Fatalf("Error writing configuration: %v", err)
	}

	readConfig, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	if len(readConfig.Clusters) != 1 || readConfig.Clusters["test-cluster"].Server != testConfig.Clusters["test-cluster"].Server {
		t.Fatalf("Cluster contents not the same")
	}

	if len(readConfig.AuthInfos) != 1 || readConfig.AuthInfos["test-user"].Token != testConfig.AuthInfos["test-user"].Token {
		t.Fatalf("AuthInfos contents not the same")
	}

	if len(readConfig.Contexts) != 1 || readConfig.Contexts["test-context"].Namespace != testConfig.Contexts["test-context"].Namespace {
		t.Fatalf("Context contents not the same")
	}
}

func TestNewConfigSetsContext(t *testing.T) {
	const expectedContext = "test-context"

	configFile, _ := ioutil.TempFile("", "")
	defer os.Remove(configFile.Name())

	testConfig := api.Config{
		AuthInfos: map[string]*api.AuthInfo{
			"test-user": {Token: "test-token"}},
		Clusters: map[string]*api.Cluster{
			"test-cluster": {Server: "https://127.0.0.1:8443"}},
		Contexts: map[string]*api.Context{
			expectedContext: {AuthInfo: "test-user", Cluster: "test-cluster", Namespace: "test-ns"}},
		CurrentContext: expectedContext,
	}

	filename, err := kubeconfig.Write(configFile.Name(), &testConfig, true)
	if err != nil {
		t.Fatalf("Error writing configuration: %v", err)
	}

	readConfig, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	if readConfig.CurrentContext != expectedContext {
		t.Fatalf("Current context is %s but expected %s", readConfig.CurrentContext, expectedContext)
	}
}

func TestMergeKubeConfig(t *testing.T) {
	configFile, _ := ioutil.TempFile("", "")
	defer os.Remove(configFile.Name())

	err := writeConfig(configFile.Name())
	if err != nil {
		t.Errorf("Error writing initial configuration file: %v", err)
	}

	testConfig := api.Config{
		AuthInfos: map[string]*api.AuthInfo{
			"test-user": {Token: "test-token"}},
		Clusters: map[string]*api.Cluster{
			"test-cluster": {Server: "https://127.0.0.1:8443"}},
		Contexts: map[string]*api.Context{
			"test-context": {AuthInfo: "test-user", Cluster: "test-cluster", Namespace: "test-ns"}},
	}

	filename, err := kubeconfig.Write(configFile.Name(), &testConfig, false)
	if err != nil {
		t.Fatalf("Error writing configuration: %v", err)
	}

	readConfig, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	if len(readConfig.Clusters) != 2 || readConfig.Clusters["test-cluster"].Server != testConfig.Clusters["test-cluster"].Server {
		t.Fatalf("Cluster contents not the same")
	}
	if readConfig.Clusters["minikube"].Server != "https://192.168.64.19:8443" {
		t.Fatalf("Error in merging as existing cluster configuration not the same")
	}

	if len(readConfig.AuthInfos) != 2 || readConfig.AuthInfos["test-user"].Token != testConfig.AuthInfos["test-user"].Token {
		t.Fatalf("AuthInfos contents not the same")
	}
	if readConfig.AuthInfos["minikube"].ClientCertificate != "/Users/bob/.minikube/client.crt" {
		t.Fatalf("Error in merging as existing AuthInfos configuration not the same")
	}

	if len(readConfig.Contexts) != 2 || readConfig.Contexts["test-context"].Namespace != testConfig.Contexts["test-context"].Namespace {
		t.Fatalf("Context contents not the same")
	}
	if readConfig.Contexts["minikube"].Cluster != "minikube" {
		t.Fatalf("Error in merging as existing Contexts configuration not the same")
	}
}

func TestMergeSetsContext(t *testing.T) {
	const expectedContext = "test-context"

	configFile, _ := ioutil.TempFile("", "")
	defer os.Remove(configFile.Name())

	err := writeConfig(configFile.Name())
	if err != nil {
		t.Errorf("Error writing initial configuration file: %v", err)
	}

	testConfig := api.Config{
		AuthInfos: map[string]*api.AuthInfo{
			"test-user": {Token: "test-token"}},
		Clusters: map[string]*api.Cluster{
			"test-cluster": {Server: "https://127.0.0.1:8443"}},
		Contexts: map[string]*api.Context{
			expectedContext: {AuthInfo: "test-user", Cluster: "test-cluster", Namespace: "test-ns"}},
		CurrentContext: expectedContext,
	}

	filename, err := kubeconfig.Write(configFile.Name(), &testConfig, true)
	if err != nil {
		t.Fatalf("Error writing configuration: %v", err)
	}

	readConfig, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	if readConfig.CurrentContext != expectedContext {
		t.Fatalf("Current context is %s but expected %s", readConfig.CurrentContext, expectedContext)
	}
}

func TestMergeDoesNotSetContext(t *testing.T) {
	expectedContext := "minikube"
	configFile, _ := ioutil.TempFile("", "")
	defer os.Remove(configFile.Name())

	err := writeConfig(configFile.Name())
	if err != nil {
		t.Errorf("Error writing initial configuration file: %v", err)
	}

	testConfig := api.Config{
		AuthInfos: map[string]*api.AuthInfo{
			"test-user": {Token: "test-token"}},
		Clusters: map[string]*api.Cluster{
			"test-cluster": {Server: "https://127.0.0.1:8443"}},
		Contexts: map[string]*api.Context{
			"test-context": {AuthInfo: "test-user", Cluster: "test-cluster", Namespace: "test-ns"}},
		CurrentContext: "test-context",
	}

	filename, err := kubeconfig.Write(configFile.Name(), &testConfig, false)
	if err != nil {
		t.Fatalf("Error writing configuration: %v", err)
	}

	readConfig, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	if readConfig.CurrentContext != expectedContext {
		t.Fatalf("Current context is %s but expected %s", readConfig.CurrentContext, expectedContext)
	}
}

// Checks to see if MaybeDeleteConfig 1) removes a cluster from the kubeconfig if the kubeconfig
// file includes the cluster and 2) does not change the kubeconfig if the kubeconfig does not
// include the cluster
func TestMaybeDeleteConfigInfo(t *testing.T) {
	configFile, _ := ioutil.TempFile("", "")

	// NOTE: Currently, when the program calls clientcmd.ModifyConfig() from MaybeDeleteConfig(),
	// ModifyConfig() uses the path assigned to the KUBECONFIG env variable instead of DefaultPath.
	// Changing the KUBECONFIG to the configFile path allows the test to simulate normal behavior.
	resetKubeconfig, err := manageConfigFile(configFile.Name())
	if err != nil {
		t.Fatalf("Error setting kubeconfig path: %v", err)
	}
	defer resetKubeconfig()

	_, err = configFile.Write(twoClustersAsBytes)
	if err != nil {
		t.Fatalf("Error writing configuration: %v", err)
	}

	// 'cluster-two' is the name of a cluster written to configFile from twoClustersAsBytes
	existingClusterConfig := getClusterConfig("cluster-two")
	kubeconfig.MaybeDeleteConfig(existingClusterConfig)

	configFileAsBytes, err := ioutil.ReadFile(configFile.Name())
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	// Checks if an existing cluster is removed.
	if !bytes.Equal(configFileAsBytes, oneClusterAsBytes){
		t.Fatalf("Failed to delete cluster from config.\n\nEXPECTED:\n%v\nGOT:\n%v", string(oneClusterAsBytes), string(configFileAsBytes))
	}

	nonExistentClusterConfig := getClusterConfig("not-a-cluster")
	kubeconfig.MaybeDeleteConfig(nonExistentClusterConfig)

	configFileAsBytes, err = ioutil.ReadFile(configFile.Name())
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
	// Checks if no changes are made for a cluster that does not exist.
	if !bytes.Equal(configFileAsBytes, oneClusterAsBytes){
		t.Fatalf("Failed to delete cluster from config.\n\nEXPECTED:\n%v\nGOT:\n%v", string(oneClusterAsBytes), string(configFileAsBytes))
	}
}

var minikubeSample = `
kind: Config
apiVersion: v1
clusters:
- cluster:
    certificate-authority: /Users/bob/.minikube/ca.crt
    server: https://192.168.64.19:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    as-user-extra: {}
    client-certificate: /Users/bob/.minikube/client.crt
    client-key: /Users/bob/.minikube/client.key		
`

func writeConfig(filename string) error {
	return ioutil.WriteFile(filename, []byte(minikubeSample), os.FileMode(0755))
}

// If tempKubeconfigPath is a valid file path, sets the KUBECONFIG env variable to the file path
// provided by tempKubeconfigPath and returns a function that restores KUBECONFIG to its previous
// value and deletes the file at tempKubeconfigPath.
// Returns an error if tempKubeconfigPath is not a valid file path.
func manageConfigFile(tempKubeconfigPath string) (func(), error) {
	if _, err := os.Stat(tempKubeconfigPath); os.IsNotExist(err) {
		return nil, err
	}
	kubeconfigPathToRestoreAtEnd, hasKubeconfigPathToRestoreAtEnd := os.LookupEnv("KUBECONFIG")
	os.Setenv("KUBECONFIG", tempKubeconfigPath)

	return func() {
		if hasKubeconfigPathToRestoreAtEnd {
			os.Setenv("KUBECONFIG", kubeconfigPathToRestoreAtEnd)
		} else {
			os.Unsetenv("KUBECONFIG")
		}
		os.Remove(tempKubeconfigPath)
	}, nil
}

// Cluster names are 'cluster-one.us-west-2.eksctl.io' and 'cluster-two.us-west-2.eksctl.io'.
// All the information for cluster-one.us-west-2.eksctl.io is identical to one-cluster.yaml. If all
// the information for cluster-two.us-west-2.eksctl.io is deleted, the file should be identical to
// one-cluster.yaml and oneClusterAsBytes.
var twoClustersAsBytes, _ = ioutil.ReadFile("./two-clusters.yaml")

// Cluster name is 'cluster-one.us-west-2.eksctl.io'.
// All the information is identical to cluster cluster-one.us-west-2.eksctl.io in two-clusters.yaml.
var oneClusterAsBytes, _ = ioutil.ReadFile("./one-cluster.yaml")

// Default cluster name is 'foo' and region is 'us-west-2'
var apiClusterConfigSample = eksctlapi.ClusterConfig{
	Region: "us-west-2",
	Profile: "",
	Tags: map[string]string{},
	ClusterName: "foo",
	NodeAMI: "",
	NodeType: "m5.large",
	Nodes: 2,
	MinNodes: 0,
	MaxNodes: 0,
	MaxPodsPerNode: 0,
	NodePolicyARNs: []string(nil),
	NodeSSH: false,
	SSHPublicKeyPath: "~/.ssh/id_rsa.pub",
	SSHPublicKey: []uint8(nil),
	SSHPublicKeyName: "",
	WaitTimeout: 1200000000000,
	SecurityGroup: "",
	Subnets: []string(nil),
	VPC: "",
	Endpoint: "",
	CertificateAuthorityData: []uint8(nil),
	ARN: "",
	ClusterStackName: "",
	NodeInstanceRoleARN: "",
	AvailabilityZones: []string{"us-west-2b", "us-west-2a", "us-west-2c"},
	Addons: eksctlapi.ClusterAddons{},
}

// Returns an ClusterConfig with a cluster name equal to the provided clusterName.
func getClusterConfig(clusterName string) *eksctlapi.ClusterConfig {
	apiClusterConfig := apiClusterConfigSample
	apiClusterConfig.ClusterName = clusterName
	return &apiClusterConfig
}
