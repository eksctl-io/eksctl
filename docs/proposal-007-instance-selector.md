# EC2 Instance Selector Support for Nodegroups

## Table of Contents

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
- [Proposal](#proposal)
- [Design](#design)
  - [eksctl create cluster](#eksctl-create-cluster)
  - [eksctl create nodegroup](#eksctl-create-nodegroup)
  - [API Changes](#api-changes)
  - [Caveats](#caveats)
<!-- /toc -->

## Summary

This proposal adds support for the EC2 instance selector for both managed and unmanaged nodegroups.


## Motivation

With support for multiple instance types and Spot instances for EKS Managed Nodegroups, users can create Spot instances where EKS will configure and launch an Autoscaling Group of Spot instances following  Spot best practices. Users can set `managedNodeGroup.spot` and supply a list of instance types in `managedNodeGroup.instanceTypes` to instruct eksctl to create a managed nodegroup using Spot instances.

It is a Spot best practice to use multiple instance types, but with over 270 EC2 instance types, users have to spend time figuring out which instance types would be well suited for their nodegroup.

The [EC2 Instance Selector](https://github.com/aws/amazon-ec2-instance-selector) tool addresses this problem by generating a list of instance types based on resource criteria such as vCPUs, memory, etc. This proposal aims to integrate the EC2 instance selector with eksctl to enable creating nodegroups with multiple instance types by passing the resource criteria.

### Problem: a simple integration
We could add the instance selector options—as CLI options (`--instance-selector-vcpus` and `--instance-selector-memory`) and in the config file—to `create cluster` and `create nodegroup` that would select the instance types matching the resource criteria for the user and proceed to creating the nodegroup. However, this gives users no control over selecting the instance types. Users may want to inspect and change the instances matching the instance selector criteria before proceeding to creating a nodegroup, and also to document the instance types being used for their nodegroups in the ClusterConfig file.

## Proposal

This design proposes adding a new general-purpose `--dry-run` option to `eksctl create cluster` and `eksctl create nodegroup` that skips cluster and nodegroup creation and instead outputs a ClusterConfig file representing the CLI options passed. This functionality extends to the instance selector integration where eksctl will set the instance types for a nodegroup when the instance selector options are supplied. This gives users the opportunity to inspect and, if required, change the instance types before proceeding with cluster and nodegroup creation. Users can omit the `--dry-run` option if they do not care about the instance types selected by eksctl.

There are certain one-off options that cannot be represented in the ClusterConfig file, e.g., `--install-vpc-controllers`.
Users would expect `eksctl create cluster --<options...> --dry-run > config.yaml` followed by `eksctl create cluster -f config.yaml` to be equivalent to running the first command without `--dry-run`. eksctl would therefore disallow passing options that cannot be represented in the config file when `--dry-run` is passed.

## Design
Along with the `--dry-run` option, two new options `--instance-selector-vcpus` and `--instance-selector-memory` will be added to `eksctl create cluster` and `eksctl create nodegroup`. These options will also be supported in the ClusterConfig file for both managed and self-managed nodegroups with the following schema:

```yaml
instanceSelector:
  vCPUs: <integer>
  memory: <integer>
```

At launch only the `vCPUs` and `memory` options will be supported, additional options may be added later based on community feedback.

**Note:** Because fields are case-sensitive, `vCPUs` may be difficult to get right, especially for users unfamiliar with Go naming conventions. However, `vCPUs` with the same case is standard enough that it's widely used by all major cloud providers. `vcpus` is easier to remember but it's inconsistent with the other field names.

### eksctl create cluster

When `eksctl create cluster` is called with the instance selector options and `--dry-run`, eksctl will output a ClusterConfig file containing a nodegroup representing the CLI options and the instance types set to the instances matched by the instance selector resource criteria.

```shell
$ eksctl create cluster --name development --managed --instance-selector-vcpus=2 --instance-selector-memory=4 --dry-run

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: development
  region: us-west-2

managedNodeGroups:
- name: <name>
  instanceTypes: ["c5.large", "c5a.large", "c5ad.large", "c5d.large", "t2.medium", "t3.medium", "t3a.medium"]
  instanceSelector:
    vCPUs: 2
    memory: 4

```

The `instanceSelector` field representing the CLI options will also be added to the ClusterConfig file for visibility and documentation purposes. When `--dry-run` is omitted, this field will be ignored and the `instanceTypes` field will be used, otherwise any changes to `instanceTypes` would get overridden by eksctl.

For self-managed nodegroups, `instancesDistribution.instanceTypes` will be set:

```shell
$ eksctl create cluster --name development --instance-selector-vcpus=2 --instance-selector-memory=4 --dry-run

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: development
  region: us-west-2

nodeGroups:
- name: <name>
  instancesDistribution:
    instanceTypes: ["c5.large", "c5a.large", "c5ad.large", "c5d.large", "t2.medium", "t3.medium", "t3a.medium"]
  instanceSelector:
    vCPUs: 2
    memory: 4


```


If `--dry-run` is omitted, eksctl will proceed to creating the cluster and nodegroup with the instance types matching the instance selector resource criteria. 

---

When a ClusterConfig file is passed with `--dry-run`, eksctl will output a ClusterConfig file containing the same set of nodegroups after expanding each nodegroup's instance selector resource criteria.


```shell
$ cat cluster.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster
  region: us-west-2

managedNodeGroups:
- name: linux
  instanceSelector:
    vCPUs: 2
    memory: 4 
  
nodeGroups:
- name: windows
  amiFamily: WindowsServer2019CoreContainer
  instanceSelector:
    vCPUs: 2
    memory: 4
  instancesDistribution:
    maxPrice: 0.017
    onDemandBaseCapacity: 0
    onDemandPercentageAboveBaseCapacity: 50
    spotInstancePools: 2  
```

```shell
$ eksctl create cluster -f cluster.yaml --dry-run 

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster
  region: us-west-2

managedNodeGroups:
- name: linux
  instanceTypes: ["c3.large", "c4.large", "c5.large", "c5d.large", "c5n.large", "c5a.large"]
  instanceSelector:
    vCPUs: 2
    memory: 4 
  
nodeGroups:
- name: windows
  amiFamily: WindowsServer2019CoreContainer
  instanceSelector:
    vCPUs: 1
    memory: 2
  instancesDistribution:
    maxPrice: 0.017
    instanceTypes: ["c3.large", "c4.large", "c5.large", "c5d.large", "c5n.large", "c5a.large"]
    onDemandBaseCapacity: 0
    onDemandPercentageAboveBaseCapacity: 50
    spotInstancePools: 2

```

**Note:** It is an error to pass both a ClusterConfig file (`-f`) and the instance selector CLI options. When using a config file, the instance selector options must be set in the `instanceSelector` field.


### eksctl create nodegroup
`eksctl create nodegroup` is different from `eksctl create cluster` in that it operates on an existing cluster. When a ClusterConfig file is not passed, eksctl will not attempt to reconcile the cluster state and will instead output a ClusterConfig file containing a single nodegroup representing the CLI options. The cluster may have more than one nodegroup, so the generated ClusterConfig file will not represent the actual state of the cluster. It is the user's responsibility to copy and add the generated nodegroup config to their existing ClusterConfig file.

```shell
$ eksctl create nodegroup --cluster=dev --managed --spot --instance-selector-vcpus=2 --instance-selector-memory=4 --dry-run

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev
  region: us-west-2

managedNodeGroups:
- name: <name>
  spot: true
  instanceTypes: ["c3.large", "c4.large", "c5.large", "c5d.large", "c5n.large", "c5a.large"]
  instanceSelector:
    vCPUs: 2
    memory: 4 

```

By outputting a well-formed ClusterConfig file instead of just the nodegroup section, users can subsequently run `eksctl create nodegroup -f generated-config.yaml` to create a new nodegroup.
 

**Note:** If a ClusterConfig file exists, it's recommended to add your nodegroups to the file and pass it to `eksctl create nodegroup`.

---

When a ClusterConfig file is passed with `--dry-run`, eksctl will output the same set of nodegroups after expanding each nodegroup's instance selector resource criteria.

```shell
$ cat cluster.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev
  region: us-west-2

managedNodeGroups:
- name: m1
  instanceSelector:
    vCPUs: 2
    memory: 4

- name: m2
  desiredCapacity: 5

$ eksctl create nodegroup -f cluster.yaml --dry-run

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev
  region: us-west-2

managedNodeGroups:
- name: m1
  instanceTypes: ["c3.large", "c4.large", "c5.large", "c5d.large", "c5n.large", "c5a.large"]
  instanceSelector:
    vCPUs: 2
    memory: 4

- name: m2
  desiredCapacity: 5

```


This consistent behaviour for both the config file version and the non-config file version ensures that subsequently running `eksctl create nodegroup -f cluster.yaml` will be identical to running the initial command without `--dry-run`. 



When the `--dry-run` option is omitted in either case, eksctl will proceed to nodegroup creation after applying the instance selector criteria.

### API Changes
A new type `InstanceSelector` will be added to represent the instance selector options:

```go
type InstanceSelector struct {
	VCPUs  int `json:"vCPUs"`
	Memory int `json:"memory"`
}
```

A new field `InstanceSelector` will be added to `NodeGroupBase`, which will make the field available for both `NodeGroup` and `ManagedNodeGroup`:

```go
type NodeGroupBase struct {
	InstanceSelector *InstanceSelector `json:"instanceSelector"`
}
```


### Caveats
When a ClusterConfig file is passed with `--dry-run`, the formatting and order of fields in the YAML output may not match the supplied file's formatting.
