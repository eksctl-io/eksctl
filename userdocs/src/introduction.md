## Getting started

### Basic cluster creation

To create a basic cluster, but with a different name, run:

```

eksctl create cluster --name=cluster-1 --nodes=4

```

EKS supports versions `1.22`, `1.23`, `1.24`, `1.25` (default), `1.26` and `1.27`.
With `eksctl` you can deploy any of the supported versions by passing `--version`.

```

eksctl create cluster --version=1.24

```

### Listing clusters

To list the details about a cluster or all of the clusters, use:

```

eksctl get cluster [--name=<name>][--region=<region>]

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

#### Caching Credentials

`eksctl` supports caching credentials. This is useful when using MFA and not wanting to continuously enter the MFA
token on each `eksctl` command run.

To enable credential caching set the following environment property `EKSCTL_ENABLE_CREDENTIAL_CACHE` as such:

```
export EKSCTL_ENABLE_CREDENTIAL_CACHE=1
```

By default, this will result in a cache file under `~/.eksctl/cache/credentials.yaml` which will contain creds per profile
that is being used. To clear the cache, delete this file.

It's also possible to configure the location of this cache file using `EKSCTL_CREDENTIAL_CACHE_FILENAME` which should
be the **full path** to a file in which to store the cached credentials. These are credentials, so make sure the access
of this file is restricted to the current user and in a secure location.

### Autoscaling

To use a 3-5 node Auto Scaling Group, run:

```

eksctl create cluster --name=cluster-5 --nodes-min=3 --nodes-max=5

```

???+ note
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

[AWS Systems Manager (SSM)](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-sessions-start.html#sessions-start-cli) is enabled by default, so it can be used to SSH onto nodes.


```

eksctl create cluster --enable-ssm

```

???+ note
    If you are creating managed nodes with a custom launch template, the `--enable-ssm` flag is disallowed.

### Tagging

To add custom tags for all resources, use `--tags`.

```

eksctl create cluster --tags environment=staging --region=us-east-1

```

### Volume size

???+ note
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

???+ note
    Cluster info will be cleaned up in kubernetes config file. Please run `kubectl config get-contexts` to select right context.

## Contributions

Code contributions are very welcome. If you are interested in helping make `eksctl` great then see our [contributing guide](https://github.com/eksctl-io/eksctl/blob/master/CONTRIBUTING.md).

_Need help? Join [Eksctl Slack][slackjoin]._
[slackjoin]: https://slack.k8s.io/


## Installation

`eksctl` is available to install from official releases as described below. We recommend that you install `eksctl` from only the official GitHub releases. You may opt to use a third-party installer, but please be advised that AWS does not maintain nor support these methods of installation. Use them at your own discretion.

### Prerequisite

You will need to have AWS API credentials configured. What works for AWS CLI or any other tools (kops, Terraform, etc.) should be sufficient. You can use [`~/.aws/credentials` file][awsconfig]
or [environment variables][awsenv]. For more information read [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html).

[awsenv]: https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html
[awsconfig]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html

You will also need [AWS IAM Authenticator for Kubernetes](https://github.com/kubernetes-sigs/aws-iam-authenticator) command (either `aws-iam-authenticator` or `aws eks get-token` (available in version 1.16.156 or greater of AWS CLI) in your `PATH`.

The IAM account used for EKS cluster creation should have these minimal access levels. 

| AWS Service      | Access Level                                           |
|------------------|--------------------------------------------------------|
| CloudFormation   | Full Access                                            |
| EC2              | **Full:** Tagging **Limited:** List, Read, Write       |
| EC2 Auto Scaling | **Limited:** List, Write                               |
| EKS              | Full Access                                            |
| IAM              | **Limited:** List, Read, Write, Permissions Management |
| Systems Manager  | **Limited:** List, Read                                |

### For Unix
To download the latest release, run:

```sh
# for ARM systems, set ARCH to: `arm64`, `armv6` or `armv7`
ARCH=amd64
PLATFORM=$(uname -s)_$ARCH

curl -sLO "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_$PLATFORM.tar.gz"

# (Optional) Verify checksum
curl -sL "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_checksums.txt" | grep $PLATFORM | sha256sum --check

tar -xzf eksctl_$PLATFORM.tar.gz -C /tmp && rm eksctl_$PLATFORM.tar.gz

sudo mv /tmp/eksctl /usr/local/bin
```

### For Windows

#### Direct download (latest release): [AMD64/x86_64](https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_windows_amd64.zip) - [ARMv6](https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_windows_armv6.zip) - [ARMv7](https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_windows_armv7.zip) - [ARM64](https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_windows_arm64.zip)
Make sure to unzip the archive to a folder in the `PATH` variable. 

Optionally, verify the checksum: 

1. Download the checksum file: [latest](https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_checksums.txt)
2. Use Command Prompt to manually compare `CertUtil`'s output to the checksum file downloaded. 
  ```cmd
  REM Replace amd64 with armv6, armv7 or arm64
  CertUtil -hashfile eksctl_Windows_amd64.zip SHA256
  ```
3. Using PowerShell to automate the verification using the `-eq` operator to get a `True` or `False` result:
```pwsh
# Replace amd64 with armv6, armv7 or arm64
 (Get-FileHash -Algorithm SHA256 .\eksctl_Windows_amd64.zip).Hash -eq ((Get-Content .\eksctl_checksums.txt) -match 'eksctl_Windows_amd64.zip' -split ' ')[0]
 ```

#### Using Git Bash: 
```sh
# for ARM systems, set ARCH to: `arm64`, `armv6` or `armv7`
ARCH=amd64
PLATFORM=windows_$ARCH

curl -sLO "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_$PLATFORM.zip"

# (Optional) Verify checksum
curl -sL "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_checksums.txt" | grep $PLATFORM | sha256sum --check

unzip eksctl_$PLATFORM.zip -d $HOME/bin

rm eksctl_$PLATFORM.zip
```

The `eksctl` executable is placed in `$HOME/bin`, which is in `$PATH` from Git Bash.

### Docker

For every release and RC a container image is pushed to ECR repository `public.ecr.aws/eksctl/eksctl`. Learn more about the usage on [ECR Public Gallery - eksctl](https://gallery.ecr.aws/eksctl/eksctl). For example, 
```bash
docker run --rm -it public.ecr.aws/eksctl/eksctl version
```

### Third-Party Installers (Not Recommended)
#### For MacOS
[Homebrew](https://brew.sh)

```
brew tap weaveworks/tap
brew install weaveworks/tap/eksctl
```

[MacPorts](https://www.macports.org)

```
port install eksctl
```
#### For Windows
[Chocolatey](https://chocolatey.org)

```
chocolatey install eksctl
```

[Scoop](https://scoop.sh)

```
scoop install eksctl
```

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

Note if you're not running a distribution like oh-my-zsh you may first have to enable autocompletion (and put in `~/.zshrc` to make it persistent):

```
autoload -U compinit
compinit
```

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
