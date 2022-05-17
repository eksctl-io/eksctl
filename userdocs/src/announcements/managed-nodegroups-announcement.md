# Managed Nodegroups Default

As of [eksctl v0.58.0](https://github.com/weaveworks/eksctl/releases/tag/0.58.0), eksctl creates managed nodegroups by
default when a `ClusterConfig` file isn't specified for `eksctl create cluster` and `eksctl create nodegroup`.
To create a self-managed nodegroup, pass `--managed=false`. This may break scripts not using a config file if a feature
not supported in managed nodegroups, e.g., Windows nodegroups, is being used.
To fix this, pass `--managed=false`, or specify your nodegroup config in a `ClusterConfig` file using the
`nodeGroups` field which creates a self-managed nodegroup.
