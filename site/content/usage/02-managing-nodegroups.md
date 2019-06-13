---
title: "Managing nodegroups"
weight: 20
---

# Managing nodegroups

You can add one or more nodegroups in addition to the initial nodegroup created along with the cluster.

To create an additional nodegroup, use:

```
eksctl create nodegroup --cluster=<clusterName> [--name=<nodegroupName>]
```

> NOTE: By default, new nodegroups inherit the version from the control plane (`--version=auto`), but you can specify a different
> version e.g. `--version=1.10`, you can also use `--version=latest` to force use of whichever is the latest version.

Additionally, you can use the same config file used for `eksctl create cluster`:

```
eksctl create nodegroup --config-file=<path>
```

If there are multiple nodegroups specified in the file, you can select
a subset via `--include=<glob,glob,...>` and `--exclude=<glob,glob,...>`:

```
eksctl create nodegroup --config-file=<path> --include='ng-prod-*-??' --exclude='ng-test-1-ml-a,ng-test-2-?'
```

## Creating a nodegroup from a config file

Nodegroups can also be created through a cluster definition or config file. Given the following example config file 
and an existing cluster called ``dev-cluster:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev-cluster
  region: eu-north-1

nodeGroups:
  - name: ng-1-workers
    labels: {role: workers}
    instanceType: m5.xlarge
    desiredCapacity: 10
    privateNetworking: true
  - name: ng-2-builders
    labels: {role: builders}
    instanceType: m5.2xlarge
    desiredCapacity: 2
    privateNetworking: true
```

The nodegroups `ng-1-workers` and `ng-2-builders` can be created with this command:

```bash
eksctl create nodegroup --config-file=dev-cluster.yaml
```

## Listing nodegroups

To list the details about a nodegroup or all of the nodegroups, use:

```
eksctl get nodegroup --cluster=<clusterName> [--name=<nodegroupName>]
```

## Nodegroup immutability

By design, nodegroups are immutable. This means that if you need to change something (other than scaling) like the 
AMI or the instance type of a nodegroup, you would need to create a new nodegroup with the desired changes, move the 
load and delete the old one. Check [Deleting and draining](#deleting-and-draining).

## Scaling

A nodegroup can be scaled by using the `eksctl scale nodegroup` command:

```
eksctl scale nodegroup --cluster=<clusterName> --nodes=<desiredCount> --name=<nodegroupName>
```

For example, to scale nodegroup `ng-a345f4e1` in `cluster-1` to 5 nodes, run:

```
eksctl scale nodegroup --cluster=cluster-1 --nodes=5 ng-a345f4e1
```

If the desired number of nodes is greater than the current maximum set on the ASG then the maximum value will be increased to match the number of requested nodes. And likewise for the minimum.

Scaling a nodegroup works by modifying the nodegroup CloudFormation stack via a ChangeSet.

> NOTE: Scaling a nodegroup down/in (i.e. reducing the number of nodes) may result in errors as we rely purely on changes to the ASG. This means that the node(s) being removed/terminated aren't explicitly drained. This may be an area for improvement in the future.

You can also enable SSH, ASG access and other feature for each particular nodegroup, e.g.:

```
eksctl create nodegroup --cluster=cluster-1 --node-labels="autoscaling=enabled,purpose=ci-worker" --asg-access --full-ecr-access --ssh-access
```

## Update labels

There are no specific commands in `eksctl`to update the labels of a nodegroup but that can easily be achieved using 
`kubectl`:

```bash
kubectl label nodes -l alpha.eksctl.io/nodegroup-name=ng-1 new-label=foo
```

## Deleting and draining

To delete a nodegroup, run:
```
eksctl delete nodegroup --cluster=<clusterName> --name=<nodegroupName>
```
> NOTE: this will drain all pods from that nodegroup before the instances are deleted.

All nodes are cordoned and all pods are evicted from a nodegroup on deletion,
but if you need to drain a nodegroup without deleting it, run:
```
eksctl drain nodegroup --cluster=<clusterName> --name=<nodegroupName>
```

To uncordon a nodegroup, run:
```
eksctl drain nodegroup --cluster=<clusterName> --name=<nodegroupName> --undo
```

## Nodegroup selection in config files

To perform a create or delete operation on only a subset of the nodegroups specified in a config file, there are two 
CLI flags: `include` and `exclude`. These accept a list of globs such as `ng-dev-*`, for example.

Using the example config file above, one can create all the workers nodegroup except the workers one with the following 
command:

```bash
eksctl create nodegroup --config-file=dev-cluster.yaml --exclude=ng-1-workers
```

Or one could delete the builders nodegroup with:

```bash
eksctl delete nodegroup --config-file=dev-cluster.yaml --include=ng-2-builders --approve
```

In this case, we also need to supply the `--approve` command to actually delete the nodegroup.

