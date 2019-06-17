---
title: "Introduction"
weight: 10
---

# `eksctl` - a CLI for Amazon EKS

[![Circle CI](https://circleci.com/gh/weaveworks/eksctl/tree/master.svg?style=shield)](https://circleci.com/gh/weaveworks/eksctl/tree/master) [![Coverage Status](https://coveralls.io/repos/github/weaveworks/eksctl/badge.svg?branch=master)](https://coveralls.io/github/weaveworks/eksctl?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/weaveworks/eksctl)](https://goreportcard.com/report/github.com/weaveworks/eksctl)

`eksctl` is a simple CLI tool for creating clusters on EKS - Amazon's new managed Kubernetes service for EC2. It is written in Go, and uses CloudFormation.

You can create a cluster in minutes with just one command – **`eksctl create cluster`**!

![Gophers: E, K, S, C, T, & L](../images/eksctl.png)

*Need help? Join [Weave Community Slack][slackjoin].*

## Usage


To create a basic cluster, run:

```
eksctl create cluster
```

A cluster will be created with default parameters
- exciting auto-generated name, e.g. "fabulous-mushroom-1527688624"
- 2x `m5.large` nodes (this instance type suits most common use-cases, and is good value for money)
- use official AWS EKS AMI
- `us-west-2` region
- dedicated VPC (check your quotas)
- using static AMI resolver

Once you have created a cluster, you will find that cluster credentials were added in `~/.kube/config`. If you have `kubectl` v1.10.x as well as `aws-iam-authenticator` commands in your PATH, you should be
able to use `kubectl`. You will need to make sure to use the same AWS API credentials for this also. Check [EKS docs][ekskubectl] for instructions. If you installed `eksctl` via Homebrew, you should have all of these dependencies installed already.

[ekskubectl]: https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html

Example output:
```
$ eksctl create cluster
[ℹ]  using region us-west-2
[ℹ]  setting availability zones to [us-west-2a us-west-2c us-west-2b]
[ℹ]  subnets for us-west-2a - public:192.168.0.0/19 private:192.168.96.0/19
[ℹ]  subnets for us-west-2c - public:192.168.32.0/19 private:192.168.128.0/19
[ℹ]  subnets for us-west-2b - public:192.168.64.0/19 private:192.168.160.0/19
[ℹ]  nodegroup "ng-98b3b83a" will use "ami-05ecac759c81e0b0c" [AmazonLinux2/1.11]
[ℹ]  creating EKS cluster "floral-unicorn-1540567338" in "us-west-2" region
[ℹ]  will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup
[ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=us-west-2 --name=floral-unicorn-1540567338'
[ℹ]  2 sequential tasks: { create cluster control plane "floral-unicorn-1540567338", create nodegroup "ng-98b3b83a" }
[ℹ]  building cluster stack "eksctl-floral-unicorn-1540567338-cluster"
[ℹ]  deploying stack "eksctl-floral-unicorn-1540567338-cluster"
[ℹ]  building nodegroup stack "eksctl-floral-unicorn-1540567338-nodegroup-ng-98b3b83a"
[ℹ]  --nodes-min=2 was set automatically for nodegroup ng-98b3b83a
[ℹ]  --nodes-max=2 was set automatically for nodegroup ng-98b3b83a
[ℹ]  deploying stack "eksctl-floral-unicorn-1540567338-nodegroup-ng-98b3b83a"
[✔]  all EKS cluster resource for "floral-unicorn-1540567338" had been created
[✔]  saved kubeconfig as "~/.kube/config"
[ℹ]  adding role "arn:aws:iam::376248598259:role/eksctl-ridiculous-sculpture-15547-NodeInstanceRole-1F3IHNVD03Z74" to auth ConfigMap
[ℹ]  nodegroup "ng-98b3b83a" has 1 node(s)
[ℹ]  node "ip-192-168-64-220.us-west-2.compute.internal" is not ready
[ℹ]  waiting for at least 2 node(s) to become ready in "ng-98b3b83a"
[ℹ]  nodegroup "ng-98b3b83a" has 2 node(s)
[ℹ]  node "ip-192-168-64-220.us-west-2.compute.internal" is ready
[ℹ]  node "ip-192-168-8-135.us-west-2.compute.internal" is ready
[ℹ]  kubectl command should work with "~/.kube/config", try 'kubectl get nodes'
[✔]  EKS cluster "floral-unicorn-1540567338" in "us-west-2" region is ready
$
```

To list the details about a cluster or all of the clusters, use:

```
eksctl get cluster [--name=<name>] [--region=<region>]
```

To create the same kind of basic cluster, but with a different name, run:

```
eksctl create cluster --name=cluster-1 --nodes=4
```

EKS supports versions `1.10`, `1.11` and `1.12` (default), with `eksctl` you can deploy either version by passing `--version`.

```
eksctl create cluster --version=1.10
```

A default [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) (gp2 volume type provisioned by EBS) will be added automatically when creating a cluster.  If you want to prevent this, use the `--storage-class` flag.  For example:

```
eksctl create cluster --storage-class=false
```

To write cluster credentials to a file other than default, run:

```
eksctl create cluster --name=cluster-2 --nodes=4 --kubeconfig=./kubeconfig.cluster-2.yaml
```

To prevent storing cluster credentials locally, run:

```
eksctl create cluster --name=cluster-3 --nodes=4 --write-kubeconfig=false
```

To let `eksctl` manage cluster credentials under `~/.kube/eksctl/clusters` directory, run:

```
eksctl create cluster --name=cluster-3 --nodes=4 --auto-kubeconfig
```

To obtain cluster credentials at any point in time, run:

```
eksctl utils write-kubeconfig --name=<name> [--kubeconfig=<path>] [--set-kubeconfig-context=<bool>]
```

To use a 3-5 node Auto Scaling Group, run:

```
eksctl create cluster --name=cluster-5 --nodes-min=3 --nodes-max=5
```

> NOTE: You will still need to install and configure autoscaling. See the "Enable Autoscaling" section below. Also note that depending on your workloads you might need to use a separate nodegroup for each AZ. See [Zone-aware Autoscaling](#zone-aware-autoscaling) below for more info.

To use 30 `c4.xlarge` nodes and prevent updating current context in `~/.kube/config`, run:

```
eksctl create cluster --name=cluster-6 --nodes=30 --node-type=c4.xlarge --set-kubeconfig-context=false
```


In order to allow SSH access to nodes, `eksctl` imports `~/.ssh/id_rsa.pub` by default, to use a different SSH public key, e.g. `my_eks_node_id.pub`, run:

```
eksctl create cluster --ssh-access --ssh-public-key=my_eks_node_id.pub
```

To use a pre-existing EC2 key pair in `us-east-1` region, you can specify key pair name (which must not resolve to a local file path), e.g. to use `my_kubernetes_key` run:

```
eksctl create cluster --ssh-access  --ssh-public-key=my_kubernetes_key --region=us-east-1
```

To add custom tags for all resources, use `--tags`.

> NOTE: Until [#25](https://github.com/weaveworks/eksctl/issues/25) is resolved, tags cannot be applied to EKS cluster itself, but most of other resources (e.g. EC2 nodes).

```
eksctl create cluster --tags environment=staging --region=us-east-1
```

To configure node root volume, use the `--node-volume-size` (and optionally `--node-volume-type`), e.g.:

```
eksctl create cluster --node-volume-size=50 --node-volume-type=io1
```

> NOTE: In `us-east-1` you are likely to get `UnsupportedAvailabilityZoneException`. If you do, copy the suggested zones and pass `--zones` flag, e.g. `eksctl create cluster --region=us-east-1 --zones=us-east-1a,us-east-1b,us-east-1d`. This may occur in other regions, but less likely. You shouldn't need to use `--zone` flag otherwise.

You can also create a cluster passing all configuration information in a file
using `--config-file`:

```
eksctl create cluster --config-file=<path>
```

To create a cluster using a configuration file and skip creating
nodegroups until later:

```
eksctl create cluster --config-file=<path> --without-nodegroup
```

To delete a cluster, run:

```
eksctl delete cluster --name=<name> [--region=<region>]
```
> NOTE: Cluster info will be cleaned up in kubernetes config file. Please run `kubectl config get-contexts` to select right context.


## Contributions

Code contributions are very welcome. If you are interested in helping make `eksctl` great then see our [contributing guide](CONTRIBUTING.md).

> ***Logo Credits***
>
> *Original Gophers drawn by [Ashley McNamara](https://twitter.com/ashleymcnamara), unique E, K, S, C, T & L Gopher identities had been produced with [Gopherize.me](https://gopherize.me/).*
