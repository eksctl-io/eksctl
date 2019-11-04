# Design Proposal #005: Upgrades

> **STATUS**: This proposal is a _final_ state, and we expect minimal additional refinements.
> If any non-trivial changes are needed to functionality defined here, in particular the user
> experience, those changes should be suggested via a PR to this proposal document.
> Any other changes to the text of the proposal or technical corrections are also very welcome.

With `eksctl` users should be able to easily upgrade from one version of Kubernetes to another.

Cluster upgrades are inherently a multi-step process, and can be fairly complex. It's important
that users can drive upgrades at their own pace, if they must.

There is a standard set of steps a user needs to take, most of which can be automated as a single
command also. However, there may be additional manual steps with some versions, however most parts
should be automated.

When upgrading a cluster, one needs to call `eks.UpdateClusterVersion`. In some cases other changes
to cluster stack may be required. After that they need to replace nodegroups one by one. 

## Initial phase

- provide command that checks cluster stack for upgradability
  - lets user update cluster stack to cater for any additional resources
  - allows to call `eks.UpdateClusterVersion` out-of-band and wait for completion
- provide instruction on how to iteratively replace nodegroups
- provide helper commands for upgrading add-ons

## Second phase (more automation)

- use CloudFormation instead of calling `eks.UpdateClusterVersion` directly
- provide automated command that replaces all nodegroups by copying configuration,
  or at least one command for each nodegroup

## CLI Design

For initial phase:

Update cluster version and append new resources to the cluster stack (if needed):
```
eksctl update cluster --name=<clusterName>
```

As this can have significant implications and some checks are appropriate before
we try updating, `--dry-run` should be enabled by default.

As control plane version can be incremented only one at a time, it would be a good
idea to have `--version` flag. The default should be `--version=next`. It should be
possible to jump many versions with `--version=<v>`, and go for latest right away
with `--version=latest`. Users will have to make sure that this doesn't take them
too far ahead, so their nodegroup(s) don't drift too far behind.

For each `<currentNodeGroup>`, create `<replacementNodeGroup>`:
```
eksctl create ng --cluster=<clusterName> --name=<replacementNodeGroup>
eksctl delete ng --cluster=<clusterName> --name=<currentNodeGroup>
```
In this case new nodegroup will inherit cluster version automatically (`--version=auto`),
however it's also possible to override that behavior.

Alternatively, one can use `cluster.yaml` config file.
They will need to update `metadata.version` field set in the config file, then run:
```
eksctl update cluster --config-file=cluster.yaml
```
> NOTE: to begin with this will only update cluster version and append any resources
> to the stacks as needed

And having added definition for new nodegroup(s), they can:
```
eksctl create ng --config-file=cluster.yaml --only=<replacementNodeGroup>
eksctl delete ng --cluster=<clusterName> --name=<currentNodeGroup>
```

In the future a flag should be added to `eksctl update cluster` that allows users to
control what gets updated exactly, e.g. `--only=version,vpc,iam`.

For second phase `eksctl update cluster` should update everything with (or without) the
config file, including add-ons (although that may be more convenient to move into another
command - TBD).

## Default Add-ons

When upgrading from 1.10 to 1.11, following add-ons should be upgraded:

- `kube-system:daemonset/kube-proxy` - only image tag has to change, it matches Kubernetes version

- `kube-system:deployment/kube-dns` has to be replaced with `kube-system:deployment/coredns`, the manifest lives in [S3 bucket](https://amazon-eks.s3-us-west-2.amazonaws.com)
- `kube-system:daemonset/aws-node` the manifest lives [in GitHub](https://github.com/aws/amazon-vpc-cni-k8s/tree/master/config)


It should be possible to provide simple command under utils for each of these.

## Downgrades

At the time of writing of this proposal version downgrades were not supported by EKS.
It's trusted that upgrades will work as expected with EKS SLA.

More specifically, an `eks.UpdateClusterVersion` call for cluster running 1.12+ trying
to go back to 1.10 or 1.11 results in the following error:
```
An error occurred (InvalidParameterException) when calling the UpdateClusterVersion operation: unsupported Kubernetes version update from the current version, 1.11, to 1.10
```

## General Config Changes

Eventually `eksctl update cluster` will evolve enough to support any meaningful changes
to configuration, albeit within what is supported by EKS. At present it's only capable
of appending new resources, such as shared SG.
