# Instance Selector

eksctl supports specifying multiple instance types for managed and self-managed nodegroups, but with over 270 EC2 instance types,
users have to spend time figuring out which instance types would be well suited for their nodegroup. It's even harder
when using Spot instances because you need to choose a set of instances that works together well with the Cluster Autoscaler.

eksctl now integrates with the [EC2 instance selector](https://github.com/aws/amazon-ec2-instance-selector),
which addresses this problem by generating a list of instance types based on resource criteria such as vCPUs, memory, etc.
When the instance selector criteria is passed, eksctl creates a nodegroup with the instance types set to the instance types
matching the supplied criteria.


## Create cluster and nodegroups
To create a cluster with a single nodegroup that uses instance types matched by the instance selector resource
criteria passed to eksctl, run

```console
$ eksctl create cluster --instance-selector-vcpus=2 --instance-selector-memory=4
```

This will create a cluster and nodegroup with the `instancesDistribution.instanceTypes` field set to
`[c5.large, c5a.large, c5ad.large, c5d.large, t2.medium, t3.medium, t3a.medium]` (the set of instance types returned may change).


For managed nodegroups, the `instanceTypes` field will be set:

```console
$ eksctl create cluster --managed --instance-selector-vcpus=2 --instance-selector-memory=4
```

The instance selector criteria can also be specified in ClusterConfig:

```yaml
# instance-selector-cluster.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster
  region: us-west-2

nodeGroups:
- name: ng
  instanceSelector:
    vCPUs: 2
    memory: "4" # 4 GiB, unit defaults to GiB

managedNodeGroups:
- name: mng
  instanceSelector:
    vCPUs: 2
    memory: 2GiB #
    cpuArchitecture: x86_64 # default value
```

```console
$ eksctl create cluster -f instance-selector-cluster.yaml
```

The following instance selector CLI options are supported by `eksctl create cluster` and `eksctl create nodegroup`:

`--instance-selector-vcpus`, `--instance-selector-memory`, `--instance-selector-gpus` and `instance-selector-cpu-architecture`


### Dry Run
The dry-run feature allows you to inspect and change the instances matched by the instance selector before proceeding
to creating a nodegroup.

When `eksctl create cluster` is called with the instance selector options and `--dry-run`, eksctl will output a
ClusterConfig file containing a nodegroup representing the CLI options and the instance types set to the instances
matched by the instance selector resource criteria.

```shell
$ eksctl create cluster --name development --managed --instance-selector-vcpus=2 --instance-selector-memory=4 --dry-run


apiVersion: eksctl.io/v1alpha5
cloudWatch:
  clusterLogging: {}
iam:
  vpcResourceControllerPolicy: true
  withOIDC: false
kind: ClusterConfig
managedNodeGroups:
- amiFamily: AmazonLinux2
  desiredCapacity: 2
  disableIMDSv1: false
  disablePodIMDS: false
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: false
      appMeshPreview: false
      autoScaler: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  instanceSelector:
    memory: "4"
    vCPUs: 2
  instanceTypes:
  - c5.large
  - c5a.large
  - c5ad.large
  - c5d.large
  - t2.medium
  - t3.medium
  - t3a.medium
  labels:
    alpha.eksctl.io/cluster-name: development
    alpha.eksctl.io/nodegroup-name: ng-7bdfb1fb
  maxSize: 2
  minSize: 2
  name: ng-7bdfb1fb
  privateNetworking: false
  securityGroups:
    withLocal: null
    withShared: null
  ssh:
    allow: false
    enableSsm: false
    publicKeyPath: ""
  tags:
    alpha.eksctl.io/nodegroup-name: ng-7bdfb1fb
    alpha.eksctl.io/nodegroup-type: managed
  volumeIOPS: 3000
  volumeSize: 80
  volumeThroughput: 125
  volumeType: gp3
metadata:
  name: development
  region: us-west-2
  version: "1.18"
privateCluster:
  enabled: false
vpc:
  autoAllocateIPv6: false
  cidr: 192.168.0.0/16
  clusterEndpoints:
    privateAccess: false
    publicAccess: true
  manageSharedNodeSecurityGroupRules: true
  nat:
    gateway: Single


```

The generated ClusterConfig can then be passed to `eksctl create cluster`:

```console
$ eksctl create cluster -f generated-cluster.yaml
```

The `instanceSelector` field representing the CLI options will also be added to the ClusterConfig file for visibility and documentation purposes.
When `--dry-run` is omitted, this field will be ignored and the `instanceTypes` field will be used, otherwise any
changes to `instanceTypes` would get overridden by eksctl.


When a ClusterConfig file is passed with `--dry-run`, eksctl will output a ClusterConfig file containing the same set of nodegroups after expanding each nodegroup's instance selector resource criteria.


```yaml
# instance-selector-cluster.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster
  region: us-west-2

nodeGroups:
- name: ng
  instanceSelector:
    vCPUs: 2
    memory: 4 # 4 GiB, unit defaults to GiB

managedNodeGroups:
- name: mng
  instanceSelector:
    vCPUs: 2
    memory: 2GiB #
    cpuArchitecture: x86_64 # default value
```

```console
$ eksctl create cluster -f instance-selector-cluster.yaml --dry-run

apiVersion: eksctl.io/v1alpha5
iam:
  vpcResourceControllerPolicy: true
  withOIDC: false
kind: ClusterConfig
managedNodeGroups:
- amiFamily: AmazonLinux2
  desiredCapacity: 2
  disableIMDSv1: false
  disablePodIMDS: false
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: null
      appMeshPreview: null
      autoScaler: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  instanceSelector:
    cpuArchitecture: x86_64
    memory: 2GiB
    vCPUs: 2
  instanceTypes:
  - t3.small
  - t3a.small
  labels:
    alpha.eksctl.io/cluster-name: cluster
    alpha.eksctl.io/nodegroup-name: mng
  maxSize: 2
  minSize: 2
  name: mng
  privateNetworking: false
  securityGroups:
    withLocal: null
    withShared: null
  ssh:
    allow: false
  tags:
    alpha.eksctl.io/nodegroup-name: mng
    alpha.eksctl.io/nodegroup-type: managed
  volumeIOPS: 3000
  volumeSize: 80
  volumeThroughput: 125
  volumeType: gp3
metadata:
  name: cluster
  region: us-west-2
  version: "1.18"
nodeGroups:
- amiFamily: AmazonLinux2
  disableIMDSv1: false
  disablePodIMDS: false
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: null
      appMeshPreview: null
      autoScaler: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  instanceSelector:
    memory: "4"
    vCPUs: 2
  instanceType: mixed
  instancesDistribution:
    capacityRebalance: false
    instanceTypes:
    - c5.large
    - c5a.large
    - c5ad.large
    - c5d.large
    - t2.medium
    - t3.medium
    - t3a.medium
  labels:
    alpha.eksctl.io/cluster-name: cluster
    alpha.eksctl.io/nodegroup-name: ng
  name: ng
  privateNetworking: false
  securityGroups:
    withLocal: true
    withShared: true
  ssh:
    allow: false
  volumeIOPS: 3000
  volumeSize: 80
  volumeThroughput: 125
  volumeType: gp3
privateCluster:
  enabled: false
vpc:
  autoAllocateIPv6: false
  cidr: 192.168.0.0/16
  clusterEndpoints:
    privateAccess: false
    publicAccess: true
  manageSharedNodeSecurityGroupRules: true
  nat:
    gateway: Single

```


!!!note
    There are certain one-off options that cannot be represented in the ClusterConfig file, e.g., --install-vpc-controllers. It is expected that `eksctl create cluster --<options...> --dry-run` > config.yaml followed by `eksctl create cluster -f config.yaml` would be equivalent to running the first command without `--dry-run`. eksctl therefore disallows passing options that cannot be represented in the config file when `--dry-run` is passed.

