package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slices"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func Reader(clusterConfig *api.ClusterConfig) io.Reader {
	data, err := json.Marshal(clusterConfig)
	Expect(err).NotTo(HaveOccurred())
	return bytes.NewReader(data)
}

func ReaderFromFile(clusterName, region, filename string) io.Reader {
	clusterConfig := ParseClusterConfig(clusterName, region, filename)
	return Reader(clusterConfig)
}

func ParseClusterConfig(clusterName, region, filename string) *api.ClusterConfig {
	data, err := os.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())
	clusterConfig, err := eks.ParseConfig(data)
	Expect(err).NotTo(HaveOccurred())
	clusterConfig.Metadata.Name = clusterName
	clusterConfig.Metadata.Region = region
	return clusterConfig
}

func GetCurrentAndNextVersionsForUpgrade(testVersion string) (currentVersion, nextVersion string) {
	supportedVersions := api.SupportedVersions()
	if len(supportedVersions) < 2 {
		Fail("Upgrade test requires at least two supported EKS versions")
	}

	// if latest version is used, fetch previous version to upgrade from
	if testVersion == api.LatestVersion {
		previousVersionIndex := slices.Index(supportedVersions, testVersion) - 1
		currentVersion = supportedVersions[previousVersionIndex]
		nextVersion = testVersion
		return
	}

	// otherwise fetch next version to upgrade to
	nextVersionIndex := slices.Index(supportedVersions, testVersion) + 1
	currentVersion = testVersion
	nextVersion = supportedVersions[nextVersionIndex]
	return
}
