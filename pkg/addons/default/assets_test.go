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
	nonCRDs := []string{}
	for _, part := range awsNodeParts {
		if strings.Contains(part, ": \"CustomResourceDefinition\"") {
			continue
		}
		nonCRDs = append(nonCRDs, part)
	}

	awsNode = strings.Join(nonCRDs, "---\n")

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
