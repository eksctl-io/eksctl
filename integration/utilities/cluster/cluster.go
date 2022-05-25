package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	. "github.com/onsi/gomega"

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
