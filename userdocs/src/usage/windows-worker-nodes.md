# Windows Worker Nodes

From version 1.14, Amazon EKS supports [Windows Nodes][eks-user-guide] that allow running Windows containers.
In addition to having Windows nodes, a Linux node in the cluster is required to run CoreDNS, as Microsoft doesn't support host-networking mode yet. Thus, a Windows EKS cluster will be a mixture of Windows nodes and at least one Linux node.
The Linux nodes are critical to the functioning of the cluster, and thus, for a production-grade cluster, it's recommended to have at least two `t2.large` Linux nodes for HA.

???+ note
    You no longer need to install the VPC resource controller on Linux worker nodes to run Windows workloads in EKS clusters
    created after October 22, 2021.
    You can enable Windows IP address management on the EKS control plane via a ConﬁgMap setting (see https://docs.aws.amazon.com/eks/latest/userguide/windows-support.html for details).
    eksctl will automatically patch the ConfigMap to enable Windows IP address management when a Windows nodegroup is created.

## Creating a new Windows cluster

The config file syntax allows creating a fully-functioning Windows cluster in a single command:

```yaml
# cluster.yaml
# An example of ClusterConfig containing Windows and Linux node groups to support Windows workloads
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: windows-cluster
  region: us-west-2

nodeGroups:
  - name: windows-ng
    amiFamily: WindowsServer2019FullContainer
    minSize: 2
    maxSize: 3

managedNodeGroups:
  - name: linux-ng
    instanceType: t2.large
    minSize: 2
    maxSize: 3

  - name: windows-managed-ng
    amiFamily: WindowsServer2019FullContainer
    minSize: 2
    maxSize: 3
```

```console
eksctl create cluster -f cluster.yaml
```


To create a new cluster with Windows un-managed nodegroup without using a config file, issue the following commands:

```console
eksctl create cluster --managed=false --name=windows-cluster --node-ami-family=WindowsServer2019CoreContainer
```


## Adding Windows support to an existing Linux cluster
To enable running Windows workloads on an existing cluster with Linux nodes (`AmazonLinux2` AMI family), you need to add a Windows nodegroup.

**NEW** Support for Windows managed nodegroup has been added (--managed=true or omit the flag).
```console
eksctl create nodegroup --managed=false --cluster=existing-cluster --node-ami-family=WindowsServer2019CoreContainer
eksctl create nodegroup --cluster=existing-cluster --node-ami-family=WindowsServer2019CoreContainer
```

To ensure workloads are scheduled on the right OS, they must have a `nodeSelector` targeting the OS it must run on:

```yaml
# Targeting Windows
  nodeSelector:
    kubernetes.io/os: windows
    kubernetes.io/arch: amd64
```

```yaml
# Targeting Linux
  nodeSelector:
    kubernetes.io/os: linux
    kubernetes.io/arch: amd64
```

If you are using a cluster older than `1.19` the `kubernetes.io/os` and `kubernetes.io/arch` labels need to be replaced with `beta.kubernetes.io/os` and `beta.kubernetes.io/arch` respectively.

### Further information

- [EKS Windows Support][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/windows-support.html

