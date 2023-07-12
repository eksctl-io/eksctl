# EKS Managed Nodegroups

[Amazon EKS managed nodegroups][eks-user-guide] is a feature that automates the provisioning and lifecycle management of nodes (EC2 instances) for Amazon EKS Kubernetes clusters. Customers can provision optimized groups of nodes for their clusters and EKS will keep their nodes up to date with the latest Kubernetes and host OS versions. 

An EKS managed node group is an autoscaling group and associated EC2 instances that are managed by AWS for an Amazon EKS cluster. Each node group uses the Amazon EKS-optimized Amazon Linux 2 AMI. Amazon EKS makes it easy to apply bug fixes and security patches to nodes, as well as update them to the latest Kubernetes versions. Each node group launches an autoscaling group for your cluster, which can span multiple AWS VPC availability zones and subnets for high-availability.

**NEW** [Launch Template support for managed nodegroups](launch-template-support.md)

???+ info
    The term "unmanaged nodegroups" has been used to refer to nodegroups that eksctl has supported since the beginning (represented via the `nodeGroups` field). The `ClusterConfig` file continues to use the `nodeGroups` field for defining unmanaged nodegroups, and managed nodegroups are defined with the `managedNodeGroups` field.

## Creating managed nodegroups

```
$ eksctl create nodegroup
```

### New clusters

To create a new cluster with a managed nodegroup, run

```console
$ eksctl create cluster
```

To create multiple managed nodegroups and have more control over the configuration, a config file can be used.

???+ note
    Managed nodegroups do not have complete feature parity with unmanaged nodegroups.

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
Another example of a config file for creating a managed nodegroup can be found [here](https://github.com/eksctl-io/eksctl/blob/main/examples/15-managed-nodes.yaml).

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
    privateNetworking: true
    minSize: 2
    maxSize: 3
```

**NEW** Support for custom AMI, security groups, `instancePrefix`, `instanceName`, `ebsOptimized`, `volumeType`, `volumeName`,
`volumeEncrypted`, `volumeKmsKeyID`, `volumeIOPS`, `maxPodsPerNode`, `preBootstrapCommands`, `overrideBootstrapCommand`, and `disableIMDSv1`


```yaml
# cluster.yaml
# A cluster with an unmanaged nodegroup and two managed nodegroups.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: managed-cluster
  region: us-west-2

managedNodeGroups:
  - name: custom-ng
    ami: ami-0e124de4755b2734d
    securityGroups:
      attachIDs: ["sg-1234"]
    maxPodsPerNode: 80
    ssh:
      allow: true
    volumeSize: 100
    volumeName: /dev/xvda
    volumeEncrypted: true
    # defaults to true, which enforces the use of IMDSv2 tokens
    disableIMDSv1: false  
    overrideBootstrapCommand: |
      #!/bin/bash
      /etc/eks/bootstrap.sh managed-cluster --kubelet-extra-args '--node-labels=eks.amazonaws.com/nodegroup=custom-ng,eks.amazonaws.com/nodegroup-image=ami-0e124de4755b2734d'
```

If you are requesting an instance type that is only available in one zone (and the eksctl config requires
specification of two) make sure to add the availability zone to your node group request:


```yaml
# cluster.yaml
# A cluster with a managed nodegroup with "availabilityZones"
---

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: flux-cluster
  region: us-east-2
  version: "1.23"
 
availabilityZones: ["us-east-2b", "us-east-2c"]
managedNodeGroups:
  - name: workers
    instanceType: hpc6a.48xlarge
    minSize: 64
    maxSize: 64
    labels: { "fluxoperator": "true" }
    availabilityZones: ["us-east-2b"]   
    efaEnabled: true
    placement:
      groupName: eks-efa-testing
```

This can be true for instance types like [the Hpc6 family](https://aws.amazon.com/ec2/instance-types/hpc6/) that are only available
in one zone.

### Existing clusters

```console
$ eksctl create nodegroup --managed
```

Tip : if you are using a `ClusterConfig` file to describe your whole cluster, describe your new managed node group in the `managedNodeGroups` field and run\:

```console
$ eksctl create nodegroup --config-file=YOUR_CLUSTER.yaml
```

## Upgrading managed nodegroups
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

If a nodegroup is on Kubernetes 1.14, and the cluster's Kubernetes version is 1.15, the nodegroup can be upgraded to
the latest AMI release for Kubernetes 1.15 using:

```console
eksctl upgrade nodegroup --name=managed-ng-1 --cluster=managed-cluster --kubernetes-version=1.15
```

To upgrade to a specific AMI release version instead of the latest version, pass `--release-version`:

```console
eksctl upgrade nodegroup --name=managed-ng-1 --cluster=managed-cluster --release-version=1.19.6-20210310
```

???+ note
    If the managed nodes are deployed using custom AMIs, the following workflow must be followed in order to deploy a new version of the custom AMI.

    - initial deployment of the nodegroup must be done using a launch template. e.g.
      ```yaml
      managedNodeGroups:
        - name: launch-template-ng
          launchTemplate:
            id: lt-1234
            version: "2" #optional (uses the default version of the launch template if unspecified)
      ```
    - create a new version of the custom AMI (using AWS EKS console).
    - create a new launch template version with the new AMI ID (using AWS EKS console).
    - upgrade the nodes to the new version of the launch template. e.g.
      ```
      eksctl upgrade nodegroup --name nodegroup-name --cluster cluster-name --launch-template-version new-template-version
      ```

## Handling parallel upgrades for nodes
Multiple managed nodes can be upgraded simultaneously. To configure parallel upgrades, define the `updateConfig` of a nodegroup when creating the nodegroup. An example `updateConfig` can be found [here](https://github.com/eksctl-io/eksctl/blob/main/examples/15-managed-nodes.yaml).

To avoid any downtime to your workloads due to upgrading multiple nodes at once, you can limit the number of nodes that can become unavailable during an upgrade by specifying this in the `maxUnavailable` field of an `updateConfig`. Alternatively, use `maxUnavailablePercentage`, which defines the maximum number of unavailable nodes as a percentage of the total number of nodes.

Note that `maxUnavailable` cannot be higher than `maxSize`. Also, `maxUnavailable` and `maxUnavailablePercentage` cannot be used simultaneously.

This feature is only available for managed nodes.

## Updating managed nodegroups

`eksctl` allows updating the [UpdateConfig](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-updateconfig.html) section of a managed nodegroup.
This section defines two fields. `MaxUnavailable` and `MaxUnavailablePercentage`. Your nodegroups are unaffected during
the update, thus downtime shouldn't be expected.

The command `update nodegroup` should be used with a config file using the `--config-file` flag. The nodegroup should
contain an `nodeGroup.updateConfig` section. More information can be found [here](https://eksctl.io/usage/schema/#nodeGroups-updateConfig).

## Nodegroup Health issues
EKS Managed Nodegroups automatically checks the configuration of your nodegroup and nodes for health issues and reports
them through the EKS API and console.
To view health issues for a nodegroup:

```console
eksctl utils nodegroup-health --name=managed-ng-1 --cluster=managed-cluster
```

## Managing Labels
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

## Scaling Managed Nodegroups
`eksctl scale nodegroup` also supports managed nodegroups. The syntax for scaling a managed or unmanaged nodegroup is
the same.

```console
eksctl scale nodegroup --name=managed-ng-1 --cluster=managed-cluster --nodes=4 --nodes-min=3 --nodes-max=5
```


## Feature parity with unmanaged nodegroups
EKS Managed Nodegroups are managed by AWS EKS and do not offer the same level of configuration as unmanaged nodegroups.
The unsupported options are noted below.

- `iam.instanceProfileARN` is not supported for managed nodegroups.
- `instancesDistribution` field is not supported
- Full control over the node bootstrapping process and customization of the kubelet are not supported. This includes the
following fields: `classicLoadBalancerNames`, `targetGroupARNs`, `clusterDNS` and `kubeletExtraConfig`.
- No support for enabling metrics on AutoScalingGroups using `asgMetricsCollection`

## Note for eksctl versions below 0.12.0
- For clusters upgraded from EKS 1.13 to EKS 1.14, managed nodegroups will not be able to communicate with unmanaged
nodegroups. As a result, pods in a managed nodegroup will be unable to reach pods in an unmanaged
nodegroup, and vice versa.
To fix this, use eksctl 0.12.0 or above and run `eksctl upgrade cluster`.
To fix this manually, add ingress rules to the shared security group and the default cluster
security group to allow traffic from each other. The shared security group and the default cluster security groups have
the naming convention `eksctl-<cluster>-cluster-ClusterSharedNodeSecurityGroup-<id>` and
`eks-cluster-sg-<cluster>-<id>-<id>` respectively.


## Further information

- [EKS Managed Nodegroups][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html
