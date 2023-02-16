# Spot instances

## Managed Nodegroups

`eksctl` supports [Spot worker nodes using EKS Managed Nodegroups][eks-user-guide], a feature that allows EKS customers with
fault-tolerant applications to easily provision and manage EC2 Spot Instances for their EKS clusters.
EKS Managed Nodegroup will configure and launch an EC2 Autoscaling group of Spot Instances following Spot best
practices and draining Spot worker nodes automatically before the instances are interrupted by AWS. There is no
incremental charge to use this feature and customers pay only for using the AWS resources, such as EC2 Spot Instances
and EBS volumes.

To create a cluster with a managed nodegroup using Spot instances, pass the `--spot` flag and an optional list of instance types:

```console
$ eksctl create cluster --spot --instance-types=c3.large,c4.large,c5.large
```

To create a managed nodegroup using Spot instances on an existing cluster:

```console
$ eksctl create nodegroup --cluster=<clusterName> --spot --instance-types=c3.large,c4.large,c5.large
```

To create Spot instances using managed nodegroups via a config file:

```yaml
# spot-cluster.yaml

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: spot-cluster
  region: us-west-2

managedNodeGroups:
- name: spot
  instanceTypes: ["c3.large","c4.large","c5.large","c5d.large","c5n.large","c5a.large"]
  spot: true


# `instanceTypes` defaults to [`m5.large`]
- name: spot-2
  spot: true

# On-Demand instances
- name: on-demand
  instanceTypes: ["c3.large", "c4.large", "c5.large"]

```

```console
$ eksctl create cluster -f spot-cluster.yaml
```

???+ note
    Unmanaged nodegroups do not support the `spot` and `instanceTypes` fields, instead the `instancesDistribution` field
    is used to configure Spot instances. [See below](spot-instances.md#unmanaged-nodegroups)


### Further information

- [EKS Spot Nodegroups][eks-user-guide]
- [EKS Managed Nodegroup Capacity Types](https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html#managed-node-group-capacity-types)

[eks-user-guide]: https://aws.amazon.com/blogs/containers/amazon-eks-now-supports-provisioning-and-managing-ec2-spot-instances-in-managed-node-groups/



## Unmanaged Nodegroups
`eksctl` has support for spot instances through the MixedInstancesPolicy for Auto Scaling Groups.

Here is an example of a nodegroup that uses 50% spot instances and 50% on demand instances:

```yaml
nodeGroups:
  - name: ng-1
    minSize: 2
    maxSize: 5
    instancesDistribution:
      maxPrice: 0.017
      instanceTypes: ["t3.small", "t3.medium"] # At least one instance type should be specified
      onDemandBaseCapacity: 0
      onDemandPercentageAboveBaseCapacity: 50
      spotInstancePools: 2
```

Note that the `nodeGroups.X.instanceType` field shouldn't be set when using the `instancesDistribution` field.

This example uses GPU instances:

```yaml
nodeGroups:
  - name: ng-gpu
    instanceType: mixed
    desiredCapacity: 1
    instancesDistribution:
      instanceTypes:
        - p2.xlarge
        - p2.8xlarge
        - p2.16xlarge
      maxPrice: 0.50
```

This example uses the capacity-optimized spot allocation strategy:

```yaml
nodeGroups:
  - name: ng-capacity-optimized
    minSize: 2
    maxSize: 5
    instancesDistribution:
      maxPrice: 0.017
      instanceTypes: ["t3.small", "t3.medium"] # At least one instance type should be specified
      onDemandBaseCapacity: 0
      onDemandPercentageAboveBaseCapacity: 50
      spotAllocationStrategy: "capacity-optimized"
```

This example uses the capacity-optimized-prioritized spot allocation strategy:

```yaml
nodeGroups:
  - name: ng-capacity-optimized-prioritized
    minSize: 2
    maxSize: 5
    instancesDistribution:
      maxPrice: 0.017
      instanceTypes: ["t3a.small", "t3.small"] # At least two instance types should be specified
      onDemandBaseCapacity: 0
      onDemandPercentageAboveBaseCapacity: 0
      spotAllocationStrategy: "capacity-optimized-prioritized"
```

[Use the `capacity-optimized-prioritized` allocation strategy and then set the order of instance types in the list of launch template overrides from highest to lowest priority (first to last in the list). Amazon EC2 Auto Scaling honors the instance type priorities on a best-effort basis but optimizes for capacity first. This is a good option for workloads where the possibility of disruption must be minimized, but also the preference for certain instance types matters.](https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-purchase-options.html#asg-spot-strategy)

Note that the `spotInstancePools` field shouldn't be set when using the `spotAllocationStrategy` field. If the `spotAllocationStrategy` is not specified, EC2 will default to use the `lowest-price` strategy.

Here is a minimal example:

```yaml
nodeGroups:
  - name: ng-1
    instancesDistribution:
      instanceTypes: ["t3.small", "t3.medium"] # At least one instance type should be specified
```

To distinguish nodes between spot or on-demand instances you can use the kubernetes label `node-lifecycle` which will have the value `spot` or `on-demand` depending on its type.

### Parameters in instancesDistribution

Please see [the config parameters](/usage/schema/#nodeGroups-instancesDistribution) for details.
