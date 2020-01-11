---
title: "EKS Managed Nodegroups"
weight: 150
url: usage/eks-managed-nodegroups
---

## EKS Managed Nodegroups

[Amazon EKS managed nodegroups][eks-user-guide] is a feature that automates the provisioning and lifecycle management of nodes (EC2 instances) for Amazon EKS Kubernetes clusters. Customers can provision optimized groups of nodes for their clusters and EKS will keep their nodes up to date with the latest Kubernetes and host OS versions. 

An EKS managed node group is an autoscaling group and associated EC2 instances that are managed by AWS for an Amazon EKS cluster. Each node group uses the Amazon EKS-optimized Amazon Linux 2 AMI. Amazon EKS makes it easy to apply bug fixes and security patches to nodes, as well as update them to the latest Kubernetes versions. Each node group launches an autoscaling group for your cluster, which can span multiple AWS VPC availability zones and subnets for high-availability.


>NOTE: The term "unmanaged nodegroups" has been used to refer to nodegroups that eksctl has supported since the
beginning and uses by default. The `ClusterConfig` file continues to use the `nodeGroups` field for defining unmanaged
nodegroups, and a new field `managedNodeGroups` has been added for defining managed nodegroups.


### Creating a cluster with Managed Nodegroups

You can add a managed node group to new or existing clusters. To create a new cluster with a managed nodegroup, run

```console
eksctl create cluster --managed
```

To create multiple managed nodegroups and have more control over the configuration, a config file can be used.
> NOTE: Managed nodegroups do not have complete feature parity with unmanaged nodegroups.

```yaml
# cluster.yaml
# A cluster with two managed nodegroups
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: managed-cluster
  region: us-west-2

managedNodeGroups:
  - name: managed-ng-1
    minSize: 2
    maxSize: 4
    desiredCapacity: 3
    volumeSize: 20
    ssh:
      allow: true
      publicKeyPath: ~/.ssh/ec2_id_rsa.pub
      # new feature for restricting SSH access to certain AWS security group IDs
      sourceSecurityGroupIds: ["sg-00241fbb12c607007"]
    labels: {role: worker}
    tags:
      nodegroup-role: worker
    iam:
      withAddonPolicies:
        externalDNS: true
        certManager: true

  - name: managed-ng-2
    instanceType: t2.large
    minSize: 2
    maxSize: 3
```

It's possible to have a cluster with both managed and unmanaged nodegroups. Unmanaged nodegroups do not show up in
the AWS EKS console but `eksctl get nodegroup` will list both types of nodegroups.


```yaml
# cluster.yaml
# A cluster with an unmanaged nodegroup and two managed nodegroups.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: managed-cluster
  region: us-west-2

nodeGroups:
  - name: ng-1
    minSize: 2

managedNodeGroups:
  - name: managed-ng-1
    minSize: 2
    maxSize: 4
    desiredCapacity: 3
    volumeSize: 20
    ssh:
      allow: true
      publicKeyPath: ~/.ssh/ec2_id_rsa.pub
      # new feature for restricting SSH access to certain AWS security group IDs
      sourceSecurityGroupIds: ["sg-00241fbb12c607007"]
    labels: {role: worker}
    tags:
      nodegroup-role: worker
    iam:
      withAddonPolicies:
        externalDNS: true
        certManager: true

  - name: managed-ng-2
    instanceType: t2.large
    minSize: 2
    maxSize: 3
```


### Upgrading managed nodegroups
You can update a nodegroup to the latest EKS-optimized AMI release version for the AMI type you are using at any time.

If your nodegroup is the same Kubernetes version as the cluster, you can update to the latest AMI release version
for that Kubernetes version of the AMI type you are using. If your nodegroup is the previous Kubernetes version from
the cluster’s Kubernetes version, you can update the nodegroup to the latest AMI release version that matches the
nodegroup’s Kubernetes version, or update to the latest AMI release version that matches the clusters Kubernetes
version. You cannot roll back a nodegroup to an earlier Kubernetes version.

To upgrade a managed nodegroup to the latest AMI release version:

```console
eksctl upgrade nodegroup --name=managed-ng-1 --cluster=managed-cluster
```

If a nodegroup is on Kubernetes 1.13, and the cluster's Kubernetes version is 1.14, the nodegroup can be upgraded to
the latest AMI release for Kubernetes 1.14 using:

```console
eksctl upgrade nodegroup --name=managed-ng-1 --cluster=managed-cluster --kubernetes-version=1.14
```

### Nodegroup Health issues
EKS Managed Nodegroups automatically checks the configuration of your nodegroup and nodes for health issues and reports
them through the EKS API and console.
To view health issues for a nodegroup:

```console
eksctl utils nodegroup-health --name=managed-ng-1 --cluster=managed-cluster
```

### Managing Labels
EKS Managed Nodegroups supports attaching labels that are applied to the Kubernetes nodes in the nodegroup. This is
specified via the `labels` field in eksctl during cluster or nodegroup creation.

To set new labels or updating existing labels on a nodegroup:

```console
eksctl set labels --cluster managed-cluster --nodegroup managed-ng-1 --labels kubernetes.io/managed-by=eks,kubernetes.io/role=worker
```


To unset or remove labels from a nodegroup:

```console
eksctl unset labels --cluster managed-cluster --nodegroup managed-ng-1 --labels kubernetes.io/managed-by,kubernetes.io/role
```

To view all labels set on a nodegroup:

```console
eksctl get labels --cluster managed-cluster --nodegroup managed-ng-1
```

### Scaling Managed Nodegroups
`eksctl scale nodegroup` also supports managed nodegroups. The syntax for scaling a managed or unmanaged nodegroup is
the same.

```console
eksctl scale nodegroup --name=managed-ng-1 --cluster=managed-cluster --nodes=4
```


### Feature parity with unmanaged nodegroups
EKS Managed Nodegroups are managed by AWS EKS and do not offer the same level of configuration as unmanaged nodegroups.
The unsupported options are noted below.

- No support for private networking (`nodeGroups[*].privateNetworking`).
- Tags (`managedNodeGroups[*].tags`) in managed nodegroups apply to the EKS Nodegroup resource and do not propagate to
the provisioned Autoscaling Group like in unmanaged nodegroups.
- `iam.instanceProfileARN` is not supported for managed nodegroups.
- The `amiFamily` field supports only `AmazonLinux2`
- `instancesDistribution` field is not supported
- `volumeSize` is the only field supported for configuring volumes
- Control over the node bootstrapping process and customization of the kubelet are not supported. This includes the
following fields: `loadBalancerNames`, `maxPodsPerNode`, `taints`, `targetGroupARNs`, `preBootstrapCommands`, `overrideBootstrapCommand`,
`clusterDNS` and `kubeletExtraConfig`.

### Note for eksctl versions below 0.12.0
- For clusters upgraded from EKS 1.13 to EKS 1.14, managed nodegroups will not be able to communicate with unmanaged
nodegroups. As a result, pods in a managed nodegroup will be unable to reach pods in an unmanaged
nodegroup, and vice versa.
To fix this, use eksctl 0.12.0 or above and run `eksctl update cluster`.
To fix this manually, add ingress rules to the shared security group and the default cluster
security group to allow traffic from each other. The shared security group and the default cluster security groups have
the naming convention `eksctl-<cluster>-cluster-ClusterSharedNodeSecurityGroup-<id>` and
`eks-cluster-sg-<cluster>-<id>-<id>` respectively.


### Further information

- [EKS Managed Nodegroups][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html

