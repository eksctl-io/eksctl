# Support for Node Repair Config in EKS Managed Nodegroups

EKS Managed Nodegroups now supports Node Repair, where the health of managed nodes are monitored,
and unhealthy worker nodes are replaced or rebooted in response.

## Creating a cluster a managed nodegroup with node repair enabled

To create a cluster with a managed nodegroup using node repair, pass the `--enable-node-repair` flag:

```shell
$ eksctl create cluster --enable-node-repair
```

To create a managed nodegroup using node repair on an existing cluster:

```shell
$ eksctl create nodegroup --cluster=<clusterName> --enable-node-repair
```

To create a cluster with a managed nodegroup using node repair via a config file:

```yaml
# node-repair-nodegroup-cluster.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-44
  region: us-west-2

managedNodeGroups:
- name: ng-1
  nodeRepairConfig:
    enabled: true

```

```shell
$ eksctl create cluster -f node-repair-nodegroup-cluster.yaml
```

## Further information

- [EKS Managed Nodegroup Node Health][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/node-health.html
