# Cluster upgrades

An _`eksctl`-managed_ cluster can be upgraded in 3 easy steps:

1. upgrade control plane version with `eksctl upgrade cluster`
2. replace each of the nodegroups by creating a new one and deleting the old one
3. update default add-ons (more about this [here](https://eksctl.io/usage/addon-upgrade/)):
    - `kube-proxy`
    - `aws-node`
    - `coredns`

Please make sure to read this section in full before you proceed.

!!!info
    Kubernetes supports version drift of up-to two minor versions during upgrade
    process. So nodes can be up to two minor versions ahead or behind the control plane
    version. You can only upgrade the control plane one minor version at a time, but
    nodes can be upgraded more than one minor version at a time, provided the nodes stay
    within two minor versions of the control plane.

!!!info
    The old `eksctl update cluster` will be deprecated. Use `eksctl upgrade cluster` instead.

## Updating control plane version

Control plane version upgrades must be done for one minor version at a time.

To upgrade control plane to the next available version run:

```
eksctl upgrade cluster --name=<clusterName>
```

This command will not apply any changes right away, you will need to re-run it with
`--approve` to apply the changes.

The target version for the cluster upgrade can be specified both with the CLI flag:

```
eksctl upgrade cluster --name=<clusterName> --version=1.16
```

or with the config file

```
cat cluster1.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-north-1
  version: "1.16"

eksctl upgrade cluster --config-file cluster1.yaml
```

!!!warning
    The only values allowed for the `--version` and `metadata.version` arguments are the current version of the cluster
    or one version higher. Upgrades of more than one Kubernetes version are not supported at the moment.

