package v1alpha5

import (
	"os"
	"strings"
)

func (c *ClusterConfig) IsCustomEksEndpoint() bool {
	eksEndpoint := os.Getenv("AWS_EKS_ENDPOINT")
	if eksEndpoint == "" {
		eksEndpoint = os.Getenv("AWS_ENDPOINT_URL_EKS")
	}
	return strings.Contains(eksEndpoint, "beta") || strings.Contains(eksEndpoint, "gamma")
}
