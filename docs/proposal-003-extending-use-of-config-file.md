# Design Proposal #003: Extending Use of `--config-file`

> **STATUS**: This proposal is a _final_ state, and we expect minimal additional refinements.
> If any non-trivial changes are needed to functionality defined here, in particular the user
> experience, those changes should be suggested via a PR to this proposal document.
> Any other changes to the text of the proposal or technical corrections are also very welcome.

Support for config files had been added in eksctl 0.1.17 as an experimental alpha feature, it has
become popular amongst users.
Users currently can specify `ClusterConfig` object in JSON or YAML, which allows them to set more
fields then they currently can with CLI flags. They can also specify several nodegroups,
and create a cluster in one go. Besides many missing features in `v1alpha4` incarnation (where there
is little of  management functionality beyond initial cluster creation), one of commonly requested
features is management of nodegroups through config files. 

This proposal aims to define how existing `ClusterConfig` object can be enhanced to support new type
of objects that allow users to define and manage nodegroups via a separate step.

Let's consider a simple config:

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

There are two nodegroups defined inline. When this file is passed to `eksctl create cluster --config-file=cluster.yaml` the cluster will be created with the nodegroups as defined.  

To create additional nodegroups for this cluster, we will enhance the `eksctl create nodegroup` cli to selectively, using glob patterns, pull nodegroups from this configuration file.  This enables the entire cluster definition to live in a single file.

For example, if we add 2 nodegroups to our original cluster config:

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
  - name: ng3-dev
    instanceType: m4.large
    desiredCapacity: 2
    privateNetworking: true
  - name: ng2-test
    instanceType: m5.large
    desiredCapacity: 4
    privateNetworking: true
```
`eksctl create nodegroup` will be updated to support glob pattern matching for nodegroup names contained within the cluster config.  Therefore, a user can create these additional nodegroups using `eksctl create nodegroup --config-file=cluster.yaml --only="*dev,*test"`.  Eksctl will create two nodegroups in cluster named cluster-5.

Here are some additional examples:

```
eksctl create ng --config-file=cluster.yaml ## all nodegroups will be created
eksctl create ng --config-file=cluster.yaml --only=ng1-public  ## only one nodegroup will be created
eksctl create ng --config-file=cluster.yaml --only=ng1-public,ng2-private ## both nodegroups will be created
```

The `eksctl create cluster` cli will be updated to optionally ignore all nodegroups in the cluster config. https://github.com/weaveworks/eksctl/issues/555

By using a single configuration file, users can keep their cluster definitions together and most importantly checked into source control.

## Possible future enhancements

- Update `eksctl create cluster` to allow the user to create selected nodegroups defined in the cluster config
> `eksctl create cluster --config-file=cluster.yaml --only='.*-private' ## only nodegroups matching glob pattern will be created`
- Support storing nodegroups within their own configuration files.

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

Two ways that these config files can be used to create cluster:

- `eksctl create cluster --config-file=cluster.yaml && eksctl create ng --config-file=extra-nodegroup.yaml`
- `cat cluster.yaml extra-nodegroup.yaml | eksctl create cluster --config-file -`

### Questions

- Is it appropriate to also allow referencing nodegroup that is defined externally within `ClusterConfig`?
   - Yes.  That is the initial implementation we landed on.
- Where in standalone `NodeGroup` config should lables live? (If we copy the same struct, we will get `spec.labels`,
  while it might make more sense to have them as `metadata.labels`).
   - We will address the definition of the actual nodegroup config file later.
- How will `ClusterConfig` and `NodeGroup` comapare when ispected in a running cluster? Namely, are all nodegroups supposed to be inline in `ClusterConfig` as well as represented as distinct object? (This is currently hard to say,  but in the future it will become importnat that there is coherent representation).
   - For the initial version, all `nodeGroup` definitions will live in the `ClusterConfig` 

