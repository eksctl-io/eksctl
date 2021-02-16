## Getting started

_Need help? Join [Weave Community Slack][slackjoin]._
[slackjoin]: https://slack.weave.works/

### Listing clusters

To list the details about a cluster or all of the clusters, use:

```

eksctl get cluster [--name=<name>][--region=<region>]

```

### Basic cluster creation

To create a basic cluster, but with a different name, run:

```

eksctl create cluster --name=cluster-1 --nodes=4

```

EKS supports versions `1.15`, `1.16`, `1.17`, `1.18` (default) and `1.19`.
With `eksctl` you can deploy any of the supported versions by passing `--version`.

```

eksctl create cluster --version=1.18

```

#### Config-based creation

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

#### Cluster credentials

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

eksctl utils write-kubeconfig --cluster=<name> [--kubeconfig=<path>][--set-kubeconfig-context=<bool>]

```

### Autoscaling

To use a 3-5 node Auto Scaling Group, run:

```

eksctl create cluster --name=cluster-5 --nodes-min=3 --nodes-max=5

```

!!! note
    You will still need to install and configure Auto Scaling. See the "Enable Auto Scaling" section. Also
    note that depending on your workloads you might need to use a separate nodegroup for each AZ. See [Zone-aware
    Auto Scaling](/usage/autoscaling/) for more info.

### SSH access

In order to allow SSH access to nodes, `eksctl` imports `~/.ssh/id_rsa.pub` by default, to use a different SSH public key, e.g. `my_eks_node_id.pub`, run:

```

eksctl create cluster --ssh-access --ssh-public-key=my_eks_node_id.pub

```

To use a pre-existing EC2 key pair in `us-east-1` region, you can specify key pair name (which must not resolve to a local file path), e.g. to use `my_kubernetes_key` run:

```

eksctl create cluster --ssh-access --ssh-public-key=my_kubernetes_key --region=us-east-1

```

To use [AWS Systems Manager (SSM)](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-sessions-start.html#sessions-start-cli) to SSH onto nodes, you can specify the `--enable-ssm` flag:

```

eksctl create cluster --enable-ssm

```

!!! note
    If you are creating managed nodes with a custom launch template, the `--enable-ssm` flag is disallowed.

### Tagging

To add custom tags for all resources, use `--tags`.

```

eksctl create cluster --tags environment=staging --region=us-east-1

```

### Volume size

!!! note
    The default volume size is 80G.

To configure node root volume, use the `--node-volume-size` (and optionally `--node-volume-type`), e.g.:

```

eksctl create cluster --node-volume-size=50 --node-volume-type=io1

```

### Deletion

To delete a cluster, run:

```

eksctl delete cluster --name=<name> [--region=<region>]

```

!!! note
    Cluster info will be cleaned up in kubernetes config file. Please run `kubectl config get-contexts` to select right context.

## Contributions

Code contributions are very welcome. If you are interested in helping make `eksctl` great then see our [contributing guide](https://github.com/weaveworks/eksctl/blob/master/CONTRIBUTING.md).


## Installation

To download the latest release, run:

```
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin
```

Alternatively, macOS users can use [Homebrew](https://brew.sh):

```
brew tap weaveworks/tap
brew install weaveworks/tap/eksctl
```

or [MacPorts](https://www.macports.org):

```
port install eksctl
```

and Windows users can use [chocolatey](https://chocolatey.org):

```
chocolatey install eksctl
```

or [scoop](https://scoop.sh):

```
scoop install eksctl
```

You will need to have AWS API credentials configured. What works for AWS CLI or any other tools (kops, Terraform etc), should be sufficient. You can use [`~/.aws/credentials` file][awsconfig]
or [environment variables][awsenv]. For more information read [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html).

[awsenv]: https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html
[awsconfig]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html

You will also need [AWS IAM Authenticator for Kubernetes](https://github.com/kubernetes-sigs/aws-iam-authenticator) command (either `aws-iam-authenticator` or `aws eks get-token` (available in version 1.16.156 or greater of AWS CLI) in your `PATH`.

### Docker

For every release and RC a docker image is pushed to [weaveworks/eksctl](https://hub.docker.com/r/weaveworks/eksctl).

### Shell Completion

#### Bash
To enable bash completion, run the following, or put it in `~/.bashrc` or `~/.profile`:

```
. <(eksctl completion bash)
```

#### Zsh
For zsh completion, please run:

```
mkdir -p ~/.zsh/completion/
eksctl completion zsh > ~/.zsh/completion/_eksctl
```

and put the following in `~/.zshrc`:

```
fpath=($fpath ~/.zsh/completion)
```

Note if you're not running a distribution like oh-my-zsh you may first have to enable autocompletion:

```
autoload -U compinit
compinit
```

To make the above persistent, run the first two lines, and put the above in `~/.zshrc`.

#### Fish
The below commands can be used for fish auto completion:

```
mkdir -p ~/.config/fish/completions
eksctl completion fish > ~/.config/fish/completions/eksctl.fish
```

#### Powershell

The below command can be referred for setting it up. Please note that the path might be different depending on your
system settings.

```
eksctl completion powershell > C:\Users\Documents\WindowsPowerShell\Scripts\eksctl.ps1
```

## Features

The features that are currently implemented are:

- Create, get, list and delete clusters
- Create, drain and delete nodegroups
- Scale a nodegroup
- Update a cluster
- Use custom AMIs
- Configure VPC Networking
- Configure access to API endpoints
- Support for GPU nodegroups
- Spot instances and mixed instances
- IAM Management and Add-on Policies
- List cluster Cloudformation stacks
- Install coredns
- Write kubeconfig file for a cluster
