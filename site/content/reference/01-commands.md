---
title: "Commands reference"
weight: 10
---

# Commands reference

a CLI for Amazon EKS

Usage: `eksctl [command] [flags]`

Commands:

`eksctl completion` | Generates shell completion scripts
------------- | -------------
`eksctl create` | Create resource(s)
`eksctl delete` | Delete resource(s)
`eksctl drain` | drain resources(s)
`eksctl get` | Get resource(s)
`eksctl help` | Help about any command
`eksctl scale` | Scale resources(s)
`eksctl update` | Update resource(s)
`eksctl utils` | Various utils
`eksctl version` | Output the version of eksctl

Use `eksctl [command] --help` for more information about a command.


## common flags

These flags can be applied to all eksctl commands and behave consistently between them.


[`-C, --color`](../02-flags#color)  | [`-h, --help`](../02-flags#help) | [`-v, --verbose`](../02-flags#verbose)
------------- | ------------- | -------------
 

## completion

Generates shell completion scripts for various Unix command-line shells.

Usage: `eksctl completion [command]`

The `completion` commands do not support any additional flags.

Without an invoking a sub-command, the completion command simply prints the help screen.

### bash

Generates shell completion scripts for the bash shell.

To load completion run

```
. <(eksctl completion bash)
```

To configure your bash shell to load completions for each session add to your bashrc

```
# ~/.bashrc or ~/.profile
. <(eksctl completion bash)
```

If you are stuck on Bash 3 (macOS) use

```
source /dev/stdin <<<"$(eksctl completion bash)"
```

### zsh

Generates shell completion scripts for the zsh shell.

To configure your zsh shell, run:

```
mkdir -p ~/.zsh/completion/
eksctl completion zsh > ~/.zsh/completion/_eksctl
```

and put the following in ~/.zshrc:

```
fpath=($fpath ~/.zsh/completion)
```

## create

Creates cluster, nodegroup, and IAM idenitity mapping resources.

Usage: `eksctl create [command] [flags]`

### common create flags <a name="create-common-flags"></a>

These flags can be applied to all eksctl commands and behave consistently between them.

[`-f, --config-file`](../02-flags#config-file) | [`-p, --profile`](../02-flags#profile) | [`-r, --region`](../02-flags#region) | [`--timeout`](../02-flags#timeout)
------------- | ------------- | ------------- | -------------

### cluster <a name="create-cluster"></a>

Create a cluster.

`Usage: eksctl create cluster [flags]`

#### flags

[`--alb-ingress-access`](../02-flags#alb-ingress-access) | [`--appmesh-access`](../02-flags#appmesh-access)
  ------------- | -------------
[`--asg-access`](../02-flags#asg-access) | [`--authenticator-role-arn`](../02-flags#authenticator-role-arn)
[`--auto-kubeconfig`](../02-flags#auto-kubeconfig) | [`--cfn-role-arn`](../02-flags#cfn-role-arn)
[`--external-dns-access`](../02-flags#external-dns-access) | [`--full-ecr-access`](../02-flags#full-ecr-access)
[`--kubeconfig`](../02-flags#kubeconfig) | [`--max-pods-per-node`](../02-flags#max-pods-per-node)
[`-n, --name`](../02-flags#name) | [`--node-ami`](../02-flags#node-ami)
[`--node-ami-family`](../02-flags#node-ami-family) | [`--node-labels`](../02-flags#node-labels)
[`--node-security-groups`](../02-flags#node-security-groups) | [`--node-volume-size`](../02-flags#node-volume-size)
[`--node-volume-type`](../02-flags#node-volume-type) | [`--node-zones`](../02-flags#node-zones)
[`--nodegroup-name`](../02-flags#nodegroup-name) |[`-N, --nodes`](../02-flags#nodes)
[`-M, --nodes-max`](../02-flags#nodes-max) | [`-m, --nodes-min`](../02-flags#nodes-min)
[`-P, --node-private-networking`](../02-flags#node-private-networking) | [`-t, --node-type`](../02-flags#node-type)
[`--set-kubeconfig-context`](../02-flags#set-kubeconfig-context) | [`--ssh-access`](../02-flags#ssh-access)
[`--ssh-public-key`](../02-flags#ssh-public-key) | [`-tags`](../02-flags#tags)
[`-version`](../02-flags#version) | [`--vpc-cidr`](../02-flags#vpc-cidr)
[`--vpc-from-kops-cluster`](../02-flags#vpc-from-kops-cluster) | [`--vpc-nat-mode`](../02-flags#vpc-nat-mode)
[`--vpc-private-subnets`](../02-flags#vpc-private-subnets) | [`--vpc-public-subnets`](../02-flags#vpc-public-subnets)
[`--without-nodegroup`](../02-flags#without-nodegroup) | [`--write-kubeconfig`](../02-flags#write-kubeconfig)
[`-zones`](../02-flags#zones) |

### iamidentitymapping <a name="create-iamidentitymapping"></a>

Creates a mapping from IAM role to Kubernetes user and groups.

*Note:* aws-iam-authenticator only considers the last entry for any given
role. If you create a duplicate entry it will shadow all the previous
username and groups mappings.

Usage: `eksctl create iamidentitymapping [flags]`

[`--role`](../02-flags#role) | [`--username`](../02-flags#username) | [`--group`](../02-flags#group) | [`-n, --name`](../02-flags#name)
  ------------- | ------------- | ------------- | -------------

***Note***: the string value of the name flag denotes the name of the EKS cluster for which the `create identitymapping` command will create the identity mapping and **NOT** the name of the identity mapping resource.

### nodegroup, ng <a name="create-nodegroup"></a>

Create a nodegroup

Usage: `eksctl create nodegroup [flags]`

#### flags

[`--alb-ingress-access`](../02-flags#alb-ingress-access) | [`--appmesh-access`](../02-flags#appmesh-access)
  ------------- | -------------
[`--asg-access`](../02-flags#asg-access) | [`--cfn-role-arn`](../02-flags#cfn-role-arn)
[`--cluster`](../02-flags#cluster) | [`--exclude`](../02-flags#exclude)
[`--external-dns-access`](../02-flags#external-dns-access) | [`--full-ecr-access`](../02-flags#full-ecr-access)
[`--include`](../02-flags#include) | [`--max-pods-per-node`](../02-flags#max-pods-per-node)
[`-n, --name`](../02-flags#name) | [`--node-ami`](../02-flags#node-ami)
[`--node-ami-family`](../02-flags#node-ami-family) | [`--node-labels`](../02-flags#node-labels)
[`-P, --node-private-networking`](../02-flags#node-private-networking) | [`--node-security-groups`](../02-flags#node-security-groups)
[`-t, --node-type`](../02-flags#node-type) | [`--node-volume-size`](../02-flags#node-volume-size)
[`--node-volume-type`](../02-flags#node-volume-type) | [`--node-zones`](../02-flags#node-zones)
[`-N, --nodes`](../02-flags#nodes) | [`-M, --nodes-max`](../02-flags#nodes-max)
[`-m, --nodes-min`](../02-flags#nodes-min) | [`--ssh-access`](../02-flags#ssh-access)
[`--ssh-public-key`](../02-flags#ssh-public-key) | [`--update-auth-configmap`](../02-flags#update-auth-configmap)
[`--version`](../02-flags#version) |



## delete

Deletes resource(s).

## drain

Drains resources(s).

## get

Get resource(s).

## help

Prints a help screen for any `eksctl` command.

## scale

Scales resources(s).

## update

Updates resource(s).

## utils

Various utils (?)

### update-cluster-logging

## version

Prints the installed `eksctl` version.
