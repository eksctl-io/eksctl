package connector

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

const (
	connectorManifestsURL        = "https://amazon-eks.s3.us-west-2.amazonaws.com/eks-connector/manifests/eks-connector/latest/eks-connector.yaml"
	connectorBindingManifestsURL = "https://amazon-eks.s3.us-west-2.amazonaws.com/eks-connector/manifests/eks-connector-roles-example/latest/eks-connector-roles-example.yaml"
)

// ManifestTemplate holds the manifest templates for EKS Connector
type ManifestTemplate struct {
	Connector   []byte
	RoleBinding []byte
}

// GetManifestTemplate returns the resources for EKS Connector.
func GetManifestTemplate() (ManifestTemplate, error) {
	client := &http.Client{
		Timeout: 45 * time.Second,
	}

	connectorManifests, err := getResource(client, connectorManifestsURL)
	if err != nil {
		return ManifestTemplate{}, err
	}

	connectorBindingManifests, err := getResource(client, connectorBindingManifestsURL)
	if err != nil {
		return ManifestTemplate{}, err
	}
	return ManifestTemplate{
		Connector:   connectorManifests,
		RoleBinding: connectorBindingManifests,
	}, nil
}

func getResource(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("expected status code %d; got %d", http.StatusOK, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

// WriteResources writes the EKS Connector resources to the current directory.
func WriteResources(manifestList *ManifestList) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "error getting current directory")
	}

	writeFile := func(filename string, data []byte) error {
		if err := os.WriteFile(path.Join(wd, filename), data, 0664); err != nil {
			return err
		}
		logger.Info("wrote file %s to %s", filename, wd)
		return nil
	}

	if err := writeFile("eks-connector.yaml", manifestList.ConnectorResources); err != nil {
		return err
	}

	if err := writeFile("eks-connector-binding.yaml", manifestList.ClusterRoleResources); err != nil {
		return err
	}
	// TODO blog link
	logger.Warning(`note: ClusterRoleBinding in "eks-connector-binding.yaml" gives cluster-admin permissions to IAM identity %q, edit if required; read %s for more info`, manifestList.IAMIdentityARN, "https://eksctl.io/usage/eks-connector")

	logger.Info("run `kubectl apply -f eks-connector.yaml,eks-connector-binding.yaml` before %s to connect the cluster", manifestList.Expiry.Format(time.RFC822))
	return nil
}
