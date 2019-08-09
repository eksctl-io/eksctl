---
title: "Command Reference"
weight: 10
---

## eksctl commands

```
a CLI for Amazon EKS

Usage: eksctl [command] [flags]

Commands:
  eksctl completion  Generates shell completion scripts
  eksctl create      Create resource(s)
  eksctl delete      Delete resource(s)
  eksctl drain       drain resources(s)
  eksctl get         Get resource(s)
  eksctl help        Help about any command
  eksctl scale       Scale resources(s)
  eksctl update      Update resource(s)
  eksctl utils       Various utils
  eksctl version     Output the version of eksctl

Use 'eksctl [command] --help' for more information about a command.
```

`Usage: eksctl [command] [flags]`

##### common flags<a name="common-flags"></a>

These flags can be applied to all eksctl commands and behave consistently between them.

- [C, color](02-flags.md#color)
- [h, help](02-flags.md#help)
- [v, verbose](02-flags.md#verbose)

### completion <a name="completion"></a>

Generates shell completion scripts for various Unix command-line shells.

`Usage: eksctl completion [command]`

The `completion` commands do not support any additional flags.

Without an invoking a sub-command, the completion command simply prints the help screen.

#### bash <a name="completion-bash"></a>

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

#### zsh <a name="completion-zsh"></a>

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

### create <a name="create"></a>

Creates cluster, nodegroup, and IAM idenitity mapping resources.

`Usage: eksctl create [command] [flags]`

##### common create flags <a name="create-common-flags"></a>

These flags can be applied to all eksctl commands and behave consistently between them.

- [f, config-file](02-flags#config-file)
- [p, profile](02-flags.md#profile)
- [r, region](02-flags.md#region)
- [timeout](02-flags.md#timeout)

#### cluster <a name="create-cluster"></a>

Create a cluster.

`Usage: eksctl create cluster [flags]`

##### flags
- [n, name](02-flags.md#name)
- [tags](02-flags.md#tags)
- [zones](02-flags.md#zones)
- [version](02-flags.md#version)
- [nodegroup-name](02-flags.md#nodegroup-name)
- [without-nodegroup](02-flags.md#without-nodegroup)
- [t, node-type](02-flags.md#node-type)
- [N, nodes](02-flags.md#nodes)
- [m, nodes-min](02-flags.md#nodes-min)
- [M, nodes-max](02-flags.md#nodes-max)
- [node-volume-size](02-flags.md#node-volume-size)
- [node-volume-type](02-flags.md#node-volume-type)
- [max-pods-per-node](02-flags.md#max-pods-per-node)
- [ssh-access](02-flags.md#ssh-access)
- [ssh-public-key](02-flags.md#ssh-public-key)
- [node-ami](02-flags.md#node-ami)
- [node-ami-family](02-flags.md#node-ami-family)
- [P, node-private-networking](02-flags.md#node-private-networking)
- [node-security-groups](02-flags.md#node-security-groups)
- [node-labels](02-flags.md#node-labels)
- [node-zones](02-flags.md#node-zones)
- [asg-access](02-flags.md#asg-access)
- [external-dns-access](02-flags.md#external-dns-access)
- [full-ecr-access](02-flags.md#full-ecr-access)
- [appmesh-access](02-flags.md#appmesh-access)
- [alb-ingress-access](02-flags.md#alb-ingress-access)
- [vpc-cidr](02-flags.md#vpc-cidr)
- [vpc-private-subnets](02-flags.md#vpc-private-subnets)
- [vpc-public-subnets](02-flags.md#vpc-public-subnets)
- [vpc-from-kops-cluster](02-flags.md#vpc-from-kops-cluster)
- [vpc-nat-mode](02-flags.md#vpc-nat-mode)
- [cfn-role-arn](02-flags.md#cfn-role-arn)
- [kubeconfig](02-flags.md#kubeconfig)
- [authenticator-role-arn](02-flags.md#authenticator-role-arn)
- [set-kubeconfig-context](02-flags.md#set-kubeconfig-context)
- [auto-kubeconfig](02-flags.md#auto-kubeconfig)
- [write-kubeconfig](02-flags.md#write-kubeconfig)

#### iamidentitymapping <a name="create-iamidentitymapping"></a>

Creates a mapping from IAM role to Kubernetes user and groups.

*Note:* aws-iam-authenticator only considers the last entry for any given
role. If you create a duplicate entry it will shadow all the previous
username and groups mappings.

`Usage: eksctl create iamidentitymapping [flags]`

##### flags
- [role](02-flags.md#role)
- [username](02-flags.md#username)
- [group](02-flags.md#group)
- [n, name](02-flags.md#name)
  - ***Note***: the string value of the name flag denotes the name of the EKS cluster for which the `create identitymapping` command will create the identity mapping and **NOT** the name of the identity mapping resource.

#### nodegroup <a name="create-nodegroup"></a>

Create a nodegroup

`Usage: eksctl create nodegroup [flags]`

##### aliases
nodegroup, ng

##### flags
- [cluster](02-flags.md#cluster)
- [version](02-flags.md#version)
- [include](02-flags.md#include)
- [exclude](02-flags.md#exclude)
- [update-auth-configmap](02-flags.md#update-auth-configmap)
- [n, name](02-flags.md#name)
- [t, node-type](02-flags.md#node-type)
- [N, nodes](02-flags.md#nodes)
- [m, nodes-min](02-flags.md#nodes-min)
- [M, nodes-max](02-flags.md#nodes-max)
- [node-volume-size](02-flags.md#node-volume-size)
- [node-volume-type](02-flags.md#node-volume-type)
- [max-pods-per-node](02-flags.md#max-pods-per-node)
- [ssh-access](02-flags.md#ssh-access)
- [ssh-public-key](02-flags.md#ssh-public-key)
- [node-ami](02-flags.md#node-ami)
- [node-ami-family](02-flags.md#node-ami-family)
- [P, node-private-networking](02-flags.md#node-private-networking)
- [node-security-groups](02-flags.md#node-security-groups)
- [node-labels](02-flags.md#node-labels)
- [node-zones](02-flags.md#node-zones)
- [asg-access](02-flags.md#asg-access)
- [external-dns-access](02-flags.md#external-dns-access)
- [full-ecr-access](02-flags.md#full-ecr-access)
- [appmesh-access](02-flags.md#appmesh-access)
- [alb-ingress-access](02-flags.md#alb-ingress-access)
- [cfn-role-arn](02-flags.md#cfn-role-arn)

### delete <a name="delete"></a>

Deletes resource(s).

### drain <a name="drain"></a>

Drains resources(s).

### get <a name="get"></a>

Get resource(s).

### help <a name="help"></a>

Prints a help screen for any `eksctl` command.

### scale <a name="scale"></a>

Scales resources(s).

### update <a name="update"></a>

Updates resource(s).

### utils <a name="utils"></a>

Various utils (?)

### version <a name="version"></a>

Prints the installed `eksctl` version.
