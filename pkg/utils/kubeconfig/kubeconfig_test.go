package kubeconfig_test

import (
	"io/ioutil"
	"os"
	"testing"

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
