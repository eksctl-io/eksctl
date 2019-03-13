package defaultaddons

import (
	"bytes"
	"compress/gzip"
	"strings"
)

func init() {
	// we must patch the manifest until we can use fake
	// clientset that supports CRDs

	awsNode := string(MustAsset("aws-node.yaml"))

	awsNodeParts := strings.Split(awsNode, "---\n")
	awsNodeParts = awsNodeParts[0 : len(awsNodeParts)-1]

	awsNode = strings.Join(awsNodeParts, "---\n")

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if _, err := zw.Write([]byte(awsNode)); err != nil {
		panic(err)
	}

	if err := zw.Close(); err != nil {
		panic(err)
	}

	_awsNodeYaml = buf.Bytes()
}
