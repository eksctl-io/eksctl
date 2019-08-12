//go:generate go run -tags=dev templates_generate.go

// Package clusterautoscaler contains:
// - the logic related to the cluster autoscaler under cluster_autoscaler.go,
// - the cluster autoscaler's YAML templates under the templates/ directory,
// - logic to embed these YAML templates inside the eksctl binary via vfsgen in
//   templates_dev.go and templates_generate.go.
package clusterautoscaler
