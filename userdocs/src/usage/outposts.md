# AWS Outposts Support

[AWS Outposts][eks-outposts] support in eksctl lets you create local clusters with the entire Kubernetes cluster, including the EKS control plane and worker nodes, running locally on AWS Outposts.
Customers can either create a local cluster with both the EKS control plane and worker nodes running locally on AWS Outposts, or they can extend an existing EKS cluster running in an AWS region
to AWS Outposts by creating worker nodes on Outposts.

!!! warning
    EKS Managed Nodegroups are not supported on Outposts.

???+ info
    Local clusters support Outpost racks only.


## Creating a local cluster on AWS Outposts

To create the EKS control plane and nodegroups on AWS Outposts, set `outpost.controlPlaneOutpostARN` to the Outpost ARN, as in:

```yaml
# outpost.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: outpost
  region: us-west-2

outpost:
  # Required.
  controlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
  # Optional, defaults to the smallest available instance type on the Outpost.
  controlPlaneInstanceType: m5d.large
```

```shell
$ eksctl create cluster -f outpost.yaml
```

This instructs eksctl to create the EKS control plane and subnets on the specified Outpost. Since an Outposts rack exists in a single availability zone, eksctl creates only one public and private subnet.
eksctl does not associate the created VPC with a [local gateway](https://docs.aws.amazon.com/outposts/latest/userguide/outposts-local-gateways.html) and, as such, eksctl will lack connectivity to the API server
and will be unable to create nodegroups. Therefore, if the `ClusterConfig` contains any nodegroups during cluster creation, the command must be run with `--without-nodegroup`,
as in:

```shell
eksctl create cluster -f outpost.yaml --without-nodegroup
```

It is the customer's responsibility to associate the eksctl-created VPC with the local gateway after cluster creation to enable connectivity to the API server. After this step, nodegroups can be created using `eksctl create nodegroup`.

You can optionally specify the instance type for the control plane nodes in `outpost.controlPlaneInstanceType` or for the nodegroups in `nodeGroup.instanceType`, but the instance type must exist on Outpost or eksctl will return an error.
By default, eksctl attempts to choose the smallest available instance type on Outpost for the control plane nodes and nodegroups.

When the control plane is on Outposts, nodegroups are created on that Outpost. You can optionally specify the Outpost ARN for the nodegroup in `nodeGroup.outpostARN` but it must match the control plane's Outpost ARN.


!!! abstract "New"
    [Fully-private cluster support](eks-private-cluster.md)

    eksctl now supports creating fully-private clusters on AWS Outposts.

```yaml
# outpost-fully-private.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: outpost-fully-private
  region: us-west-2

privateCluster:
  enabled: true

outpost:
  # Required.
  controlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
  # Optional, defaults to the smallest available instance type on the Outpost.
  controlPlaneInstanceType: m5d.large
```

!!! abstract "New"
    **Placement group support**

    A placement group name can be specified in `controlPlanePlacement.groupName` to satisfy high-availability requirements according to your Outpost deployment topology. If a placement group is not specified, the default EC2 placement is used.

```yaml
# outpost.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: outpost
  region: us-west-2

outpost:
  # Required.
  controlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
  # Optional, defaults to the smallest available instance type on the Outpost.
  controlPlaneInstanceType: m5d.large

  controlPlanePlacement:
    groupName: placement-group-name
```

### Existing VPC
Customers with an existing VPC can create local clusters on AWS Outposts by specifying the subnet configuration in `vpc.subnets`, as in:

```yaml
# outpost-existing-vpc.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: outpost
  region: us-west-2

vpc:
  id: vpc-1234
  subnets:
    private:
      outpost-subnet-1:
        id: subnet-1234

nodeGroups:
  - name: outpost-ng
    privateNetworking: true

outpost:
    # Required.
    controlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
    # Optional, defaults to the smallest available instance type on the Outpost.
    controlPlaneInstanceType: m5d.large
```

```shell
$ eksctl create cluster -f outpost-existing-vpc.yaml
```

The subnets must exist on the Outpost specified in `outpost.controlPlaneOutpostARN` or eksctl will return an error. You can also specify nodegroups during cluster creation if you have access
to the local gateway for the subnet, or have connectivity to VPC resources.

## Extending existing clusters to AWS Outposts
Customers can extend an existing EKS cluster running in an AWS region to AWS Outposts by setting `nodeGroup.outpostARN` for new nodegroups to create nodegroups on Outposts, as in:

```yaml
# extended-cluster.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: existing-cluster
  region: us-west-2

nodeGroups:
  # Nodegroup will be created in an AWS region.
  - name: ng

  # Nodegroup will be created on the specified Outpost.
  - name: outpost-ng
    privateNetworking: true
    outpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234"

```

```shell
$ eksctl create nodegroup -f extended-cluster.yaml
```

In this setup, the EKS control plane runs in an AWS region while nodegroups with `outpostARN` set run on the specified Outpost.
When a nodegroup is being created on Outposts for the first time, eksctl extends the VPC by creating subnets on the specified Outpost. These subnets are used to create nodegroups that have `outpostARN` set.

Customers with a pre-existing VPC are required to create the subnets on Outposts and pass them in `nodeGroup.subnets`, as in:

```yaml
# extended-cluster-vpc.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: extended-cluster-vpc
  region: us-west-2

vpc:
  id: vpc-1234
    subnets:
      private:
        outpost-subnet-1:
          id: subnet-1234

nodeGroups:
  # Nodegroup will be created in an AWS region.
  - name: ng

  # Nodegroup will be created on the specified Outpost.
  - name: outpost-ng
    privateNetworking: true
    # Subnet IDs for subnets created on Outpost.
    subnets: [subnet-5678]
    outpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
```

???+ note
    - Only Amazon Linux 2 is supported for nodegroups when the control plane is on Outposts.
    - Only EBS gp2 volume types are supported for nodegroups on Outposts.


## Features unsupported on local clusters
* [Fully-private clusters](/usage/eks-private-cluster)
* [Addons](/usage/addons)
* [IAM Roles for Service Accounts](/usage/iamserviceaccounts)
* [IPv6](/usage/vpc-ip-family)
* [Identity Providers](https://github.com/weaveworks/eksctl/blob/main/examples/27-oidc-provider.yaml)
* [Fargate](/usage/fargate-support)
* [KMS Encryption](/usage/kms-encryption)
* [Local Zones](https://github.com/weaveworks/eksctl/blob/main/examples/33-local-zones.yaml)
* [Karpenter](/usage/eksctl-karpenter)
* [GitOps](/usage/gitops-v2)
* [Instance Selector](/usage/instance-selector)
* Availability Zones cannot be specified as it defaults to the Outpost availability zone.
* `vpc.publicAccessCIDRs` and `vpc.autoAllocateIPv6` are not supported.
* Public endpoint access to the API server is not supported as a local cluster can only be created with private-only endpoint access.


## Further information

- [Amazon EKS on AWS Outposts][eks-outposts]
- [Local clusters for Amazon EKS on AWS Outposts](https://docs.aws.amazon.com/eks/latest/userguide/eks-outposts-local-cluster-overview.html)
- [Creating local clusters](https://docs.aws.amazon.com/eks/latest/userguide/eks-outposts-local-cluster-create.html)
- [Launching self-managed Amazon Linux nodes on an Outpost](https://docs.aws.amazon.com/eks/latest/userguide/eks-outposts-self-managed-nodes.html)


[eks-outposts]: https://docs.aws.amazon.com/eks/latest/userguide/eks-outposts.html
