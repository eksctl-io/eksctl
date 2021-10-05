package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func Reader(clusterConfig *v1alpha5.ClusterConfig) io.Reader {
	data, err := json.Marshal(clusterConfig)
	Expect(err).ToNot(HaveOccurred())
	return bytes.NewReader(data)
}

func ReaderFromFile(clusterName, region, filename string) io.Reader {
	data, err := os.ReadFile(filename)
	Expect(err).ToNot(HaveOccurred())
	clusterConfig, err := eks.ParseConfig(data)
	Expect(err).ToNot(HaveOccurred())
	clusterConfig.Metadata.Name = clusterName
	clusterConfig.Metadata.Region = region

	data, err = json.Marshal(clusterConfig)
	Expect(err).ToNot(HaveOccurred())
	return bytes.NewReader(data)
}
