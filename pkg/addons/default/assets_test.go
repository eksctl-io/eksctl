package defaultaddons

import (
	"strings"
)

func init() {
	// we must patch the manifest until we can use fake
	// clientset that supports CRDs

	awsNode := string(latestAWSNodeYaml)

	awsNodeParts := strings.Split(awsNode, "---\n")
	nonCRDs := []string{}
	for _, part := range awsNodeParts {
		if strings.Contains(part, ": \"CustomResourceDefinition\"") {
			continue
		}
		nonCRDs = append(nonCRDs, part)
	}

	awsNode = strings.Join(nonCRDs, "---\n")

	latestAWSNodeYaml = []byte(awsNode)
}
