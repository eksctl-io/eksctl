# EKS Auto Mode

## Introduction

eksctl supports [EKS Auto Mode][eks-user-guide], a feature that extends AWS management of Kubernetes clusters beyond the cluster itself,
to allow AWS to also set up and manage the infrastructure that enables the smooth operation of your workloads.
This allows you to delegate key infrastructure decisions and leverage the expertise of AWS for day-to-day operations.
Cluster infrastructure managed by AWS includes many Kubernetes capabilities as core components, as opposed to add-ons,
such as compute autoscaling, pod and service networking, application load balancing, cluster DNS, block storage, and GPU support.

## Creating an EKS cluster with Auto Mode enabled

`eksctl` has added a new `autoModeConfig` field to enable and configure Auto Mode. The shape of the `autoModeConfig` field is

```yaml
autoModeConfig:
    # defaults to false
    enabled: boolean
    # optional, defaults to [general-purpose, system].
    # To disable creation of nodePools, set it to the empty array ([]).
    nodePools: []string
    # optional, eksctl creates a new role if this is not supplied
    # and nodePools are present.
    nodeRoleARN: string
```

If `autoModeConfig.enabled` is true, eksctl creates an EKS cluster by passing `computeConfig.enabled: true`,
`kubernetesNetworkConfig.elasticLoadBalancing.enabled: true`, and `storageConfig.blockStorage.enabled: true` to the EKS API,
enabling management of data plane components like compute, storage and networking.

To create an EKS cluster with Auto Mode enabled, set `autoModeConfig.enabled: true`, as in

```yaml
# auto-mode-cluster.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
    name: auto-mode-cluster
    region: us-west-2

autoModeConfig:
    enabled: true
```

```shell
$ eksctl create cluster -f auto-mode-cluster.yaml
```

eksctl creates a node role to use for nodes launched by Auto Mode. eksctl also creates the `general-purpose` and `system` node pools.
To disable creation of the default node pools, e.g., to configure your own node pools that use a different set of subnets, set `nodePools: []`, as in

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
    name: auto-mode-cluster
    region: us-west-2

autoModeConfig:
    enabled: true
    nodePools: [] # disables creation of default node pools.
```

## Updating an EKS cluster to use Auto Mode
To update an existing EKS cluster to use Auto Mode, run

```yaml
# cluster.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
    name: cluster
    region: us-west-2

autoModeConfig:
    enabled: true
```

```shell
$ eksctl update auto-mode-config -f cluster.yaml
```

???+ note
    If the cluster was created by eksctl, and it uses public subnets as cluster subnets, Auto Mode will launch nodes in public subnets.
    To use private subnets for worker nodes launched by Auto Mode, [update the cluster to use private subnets](https://eksctl.io/usage/cluster-subnets-security-groups/).


## Disabling Auto Mode
To disable Auto Mode, set `autoModeConfig.enabled: false` and run

```yaml
# cluster.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
    name: auto-mode-cluster
    region: us-west-2

autoModeConfig:
    enabled: false
```

```shell
$ eksctl update auto-mode-config -f cluster.yaml
```

## Further information

- [EKS Auto Mode][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/automode.html
