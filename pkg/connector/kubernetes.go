package connector

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

const (
	connectorManifestsURL     = "https://amazon-eks.s3.us-west-2.amazonaws.com/eks-connector/manifests/eks-connector/latest/eks-connector.yaml"
	clusterRoleManifestsURL   = "https://amazon-eks.s3.us-west-2.amazonaws.com/eks-connector/manifests/eks-connector-console-roles/eks-connector-clusterrole.yaml"
	consoleAccessManifestsURL = "https://amazon-eks.s3.us-west-2.amazonaws.com/eks-connector/manifests/eks-connector-console-roles/eks-connector-console-dashboard-full-access-group.yaml"
)

// ManifestTemplate holds the manifest templates for EKS Connector.
type ManifestTemplate struct {
	Connector     ManifestFile
	ClusterRole   ManifestFile
	ConsoleAccess ManifestFile
}

type ManifestFile struct {
	Data     []byte
	Filename string
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

	clusterRoleManifests, err := getResource(client, clusterRoleManifestsURL)
	if err != nil {
		return ManifestTemplate{}, err
	}

	consoleAccessManifests, err := getResource(client, consoleAccessManifestsURL)
	if err != nil {
		return ManifestTemplate{}, err
	}

	return ManifestTemplate{
		Connector:     connectorManifests,
		ClusterRole:   clusterRoleManifests,
		ConsoleAccess: consoleAccessManifests,
	}, nil
}

func getResource(client *http.Client, url string) (ManifestFile, error) {
	resp, err := client.Get(url)
	if err != nil {
		return ManifestFile{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ManifestFile{}, errors.Errorf("expected status code %d; got %d", http.StatusOK, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ManifestFile{}, err
	}

	_, filename := filepath.Split(resp.Request.URL.Path)

	return ManifestFile{
		Data:     data,
		Filename: filename,
	}, nil
}

// GetManifestFilenames gets the filenames for EKS Connector manifests
func GetManifestFilenames() ([]string, error) {
	var filenames []string
	for _, u := range []string{connectorManifestsURL, clusterRoleManifestsURL, consoleAccessManifestsURL} {
		filename, err := filenameFromURL(u)
		if err != nil {
			return nil, err
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func filenameFromURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrapf(err, "unexpected error getting filename for URL %q", u)
	}
	_, filename := filepath.Split(parsed.Path)
	return filename, nil
}

// WriteResources writes the EKS Connector resources to the current directory.
func WriteResources(fs afero.Fs, manifestList *ManifestList) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "error getting current directory")
	}

	writeFile := func(filename string, data []byte) error {
		if err := afero.WriteFile(fs, path.Join(wd, filename), data, 0664); err != nil {
			return err
		}
		logger.Info("wrote file %s to %s", filename, wd)
		return nil
	}

	var filenames []string
	for _, m := range []ManifestFile{manifestList.ConnectorResources, manifestList.ClusterRoleResources, manifestList.ConsoleAccessResources} {
		if err := writeFile(m.Filename, m.Data); err != nil {
			return errors.Wrapf(err, "error writing file %s", m.Filename)
		}
		filenames = append(filenames, m.Filename)
	}

	logger.Warning(`note: %q and %q give full EKS Console access to IAM identity %q, edit if required; read %s for more info`,
		manifestList.ClusterRoleResources.Filename, manifestList.ConsoleAccessResources.Filename, manifestList.IAMIdentityARN,
		"https://docs.aws.amazon.com/eks/latest/userguide/connector-grant-access.html")

	logger.Info("run `kubectl apply -f %s` before %s to connect the cluster", strings.Join(filenames, ","), manifestList.Expiry.Format(time.RFC822))
	return nil
}
