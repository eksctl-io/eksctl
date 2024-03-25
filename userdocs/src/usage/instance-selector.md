# Instance Selector

eksctl supports specifying multiple instance types for managed and self-managed nodegroups, but with over 270 EC2 instance types,
users have to spend time figuring out which instance types would be well suited for their nodegroup. It's even harder
when using Spot instances because you need to choose a set of instances that works together well with the Cluster Autoscaler.

eksctl now integrates with the [EC2 instance selector](https://github.com/aws/amazon-ec2-instance-selector),
which addresses this problem by generating a list of instance types based on resource criteria: vCPUs, memory, # of GPUs and CPU architecture.
When the instance selector criteria is passed, eksctl creates a nodegroup with the instance types set to the instance types
matching the supplied criteria.


## Create cluster and nodegroups
To create a cluster with a single nodegroup that uses instance types matched by the instance selector resource
criteria passed to eksctl, run

```console
$ eksctl create cluster --instance-selector-vcpus=2 --instance-selector-memory=4
```

This will create a cluster and a managed nodegroup with the `instanceTypes` field set to
`[c5.large, c5a.large, c5ad.large, c5d.large, t2.medium, t3.medium, t3a.medium]` (the set of instance types returned may change).


For unmanaged nodegroups, the `instancesDistribution.instanceTypes` field will be set:

```console
$ eksctl create cluster --managed=false --instance-selector-vcpus=2 --instance-selector-memory=4
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

???+ note
  By default, GPU instance types are not filtered out. If you wish to do so (e.g. for cost effectiveness, when your applications don't particularly benefit from GPU-accelerated workloads), please explicitly set `gpus: 0` (via config file) or `--instance-selector-gpus=0` (via CLI flag).

An example file can be found [here](https://github.com/eksctl-io/eksctl/blob/main/examples/28-instance-selector.yaml).

### Dry Run
The [dry-run](/usage/dry-run) feature allows you to inspect and change the instances matched by the instance selector before proceeding
to creating a nodegroup.

```console
$ eksctl create cluster --name development --instance-selector-vcpus=2 --instance-selector-memory=4 --dry-run

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
# ...
managedNodeGroups:
- amiFamily: AmazonLinux2
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
...
# other config
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
kind: ClusterConfig
# ...
managedNodeGroups:
- amiFamily: AmazonLinux2
  # ...
  instanceSelector:
    cpuArchitecture: x86_64
    memory: 2GiB
    vCPUs: 2
  instanceTypes:
  - t3.small
  - t3a.small
nodeGroups:
- amiFamily: AmazonLinux2
  # ...
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
# ...
```
