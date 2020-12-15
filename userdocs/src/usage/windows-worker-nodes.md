# Windows Worker Nodes

From version 1.14, Amazon EKS supports [Windows Nodes][eks-user-guide] that allow running Windows containers.
In addition to having Windows nodes, a Linux node in the cluster is required to run the VPC resource controller and CoreDNS, as Microsoft doesn't support host-networking mode yet. Thus, a Windows EKS cluster will be a mixed-mode cluster containing Windows nodes and at least one Linux node.
The Linux nodes are critical to the functioning of the cluster, and thus, for a production-grade cluster, it's recommended to have at least two `t2.large` Linux nodes for HA.

`eksctl` provides a flag to install the VPC resource controller as part of cluster creation, and a command to install it after a cluster has been created.

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
  - name: linux-ng
    instanceType: t2.large
    minSize: 2
    maxSize: 3
```

```console
eksctl create cluster -f cluster.yaml --install-vpc-controllers
```


To create a new cluster without using a config file, issue the following commands:

```console
eksctl create cluster --name=windows-cluster --node-ami-family=WindowsServer2019CoreContainer
eksctl create nodegroup --cluster=windows-cluster --node-ami-family=AmazonLinux2 --nodes-min=2 --node-type=t2.large
eksctl utils install-vpc-controllers --cluster=windows-cluster --approve
```


## Adding Windows support to an existing Linux cluster
To enable running Windows workloads on an existing cluster with Linux nodes (`AmazonLinux2` AMI family), you need to add a Windows node group and install the Windows VPC controller:

```console
eksctl create nodegroup --cluster=existing-cluster --node-ami-family=WindowsServer2019CoreContainer
eksctl utils install-vpc-controllers --cluster=windows-cluster --approve
```

To ensure workloads are scheduled on the right OS, they must have a `nodeSelector` targeting the OS it must run on:

```yaml
# Targeting Windows
  nodeSelector:
    beta.kubernetes.io/os: windows
    beta.kubernetes.io/arch: amd64

# Targeting Linux
  nodeSelector:
    beta.kubernetes.io/os: linux
    beta.kubernetes.io/arch: amd64
```
### Further information

- [EKS Windows Support][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/windows-support.html

