package testutils

import (
	"bytes"
	"encoding/json"
	"io"

	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func ClusterConfigReader(clusterConfig *v1alpha5.ClusterConfig) io.Reader {
	data, err := json.Marshal(clusterConfig)
	Expect(err).ToNot(HaveOccurred())
	return bytes.NewReader(data)
}
