# Design Proposal #003: ClusterConfig and NodeGroupConfig

> **STATUS**: This proposal is a _working draft_, it will get refined and augment as needed.
> If any non-trivial changes are need to functionality defined here, in particular the user
> experience, those changes should be suggested via a PR to this proposal document.
> Any other changes to the text of the proposal or technical corrections are also very welcome.

Support for config files had been added in eksctl 0.1.17 as an experimental alpha feature, it has
become popular amongs users.
Users currently can specify `ClusterConfig` object in JSON or YAML, which allows them to set more
different fields then they currently can with CLI flags. They can alos specify several nodegroups,
and create a cluster in one go. Besides many missing features in `v1alpha4` incornation (where there
is little of  management functionality beyond initial cluster creation), one of commonly requested
features is management of nodegroups through config files. 

This proposal aims to define how existing `ClusterConfig` object can be enhanced to support new type
of objects that allow users to define and manage nodegroups via a separate step.

Let's conside a simple config:

```YAML
# cluster.yaml
---
apiVersion: eksctl.io/v1alpha4
kind: ClusterConfig

metadata:
  name: cluster-5
  region: eu-north-1

nodeGroups:
  - name: ng1-public
    instanceType: m5.xlarge
    desiredCapacity: 4
  - name: ng2-private
    instanceType: m5.large
    desiredCapacity: 10
    privateNetworking: true
```

There are two nodegroups are defined inline.

It should be possible to create this cluster and add more nodegroups based on a separate configs later.

For example:

```YAML
# extra-nodegroup.yaml
---
apiVersion: eksctl.io/v1alpha4
kind: NodeGroupConfig

metadata:
  cluster: cluster-5
  region: eu-north-1
  name: ng3-extra

spec:
  instanceType: c4.xlarge
  minSize: 10
  maxSize: 20
  privateNetworking: true
```

Two way that these config files can be use to create cluster:

- `eksctl create cluster --config-file=cluster.yaml && eksctl create ng --config-file=extra-nodegroup.yaml`
- `cat cluster.yaml extra-nodegroup.yaml | eksctl create cluster --config-file -`

### Questions

- Is it appropriate to also allow referencing nodegroup that is defined externally withing `ClusterConfig`?
- Where in standalone `NodeGroup` config should lables live? (If we copy the same struct, we will get `spec.labels`,
  while it might make more sense to have them as `metadata.labels`).
- How will `ClusterConfig` and `NodeGroup` comapare when ispected in a running cluster? Namely, are all nodegroups
  supposed to be inline in `ClusterConfig` as well as represented as distinct object? (This is currently hard to
  say,  but in the future it will become importnat that there is coherent representation).

### Alternative Approach

One alternative could be to avoid creating a separate object, and keep objects and relationships simpler.

Given the `cluster.yaml` above, we could allow creating cluster with an empty list of nodegroups, as well as
a flag to force creation of node or some of the nodegroups.

Examples:

```
eksctl create cluster --config-file=cluster.yaml --exclude-nodegroups='.*-private' ## only nodegroups matching regex will be created
eksctl create cluster --config-file=cluster.yaml --exclude-nodegroups='.*' ## no nodegroups will be created all
eksctl create ng --config-file=cluster.yaml ## all nodegroups will be created
eksctl create ng --config-file=cluster.yaml --only=ng1-public ## only one nodegroup will be created
eksctl create ng --config-file=cluster.yaml --only=ng1-public,ng2-private ## both nodegroups will be created
```