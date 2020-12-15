# ARM Support

EKS supports 64 bit ARM architecture with its [Graviton processors](https://aws.amazon.com/ec2/graviton/). To create a cluster,
select one of the Graviton-based instance types (`a1`, `t4g`, `m6g`, `m6gd`, `c6g`, `c6gd`, `r6g`, `r6gd`) and run:


```
eksctl create cluster --node-type=a1.large
```

or use a config file:

```
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-arm-1
  region: us-west-2


nodeGroups:
  - name: ng-arm-1
    instanceType: m6g.medium
    desiredCapacity: 1
```

```
eksctl create cluster -f cluster-arm-1.yaml
```

ARM is also supported in managed nodegroups:

```
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-arm-2
  region: us-west-2

managedNodeGroups:
  - name: mng-arm-1
    instanceType: m6g.medium
    desiredCapacity: 1
```

```
eksctl create cluster -f cluster-arm-2.yaml
```

The AMI resolvers, `auto` and `auto-ssm`, will see that you want to use an ARM instance type and they will select the correct AMI.

!!!note
    Note that currently there are only AmazonLinux2 EKS optimized AMIs for ARM.

!!!note
    ARM is supported for clusters with version 1.15 and higher.

