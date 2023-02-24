# Installation

## Initial setup
Before you get started, you need to configure your AWS credentials on your local machine. What works for the AWS CLI or other tools, such as kops, Terraform, etc., is sufficient. 

### AWS CLI
To configure AWS CLI credentials on your local machine, see the following resources. 
- [Environment variables to configure the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html)
- [Configuration and credential file settings](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html)

### AWS IAM Authenticator
You also need AWS CLI version 1.16.156 or above installed on your local machine. This version contains the [AWS IAM Authenticator for Kubernetes](https://github.com/kubernetes-sigs/aws-iam-authenticator).
- To update your AWS CLI version, see [Installing or updating the latest version of the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html).
- If you are not able to upgrade to the latest version of the CLI, the [AWS IAM Authenticator for Kubernetes](https://github.com/kubernetes-sigs/aws-iam-authenticator) needs to be in your `PATH`.

## Check the version
To check the version of `eksctl` running on your local machine, run the following command.
```bash
eksctl version
```
You should see the version in the response output. For example:
```bash
0.130.0
```

## curl
1. To download the latest release using [curl](https://curl.se/), run the following command. 
```bash
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin
```

## macOS 
### Homebrew
1. To install or upgrade `eksctl` using [Homebrew](https://brew.sh), run the following command. 
```bash
brew upgrade eksctl && { brew link --overwrite eksctl; } || { brew tap weaveworks/tap; brew install weaveworks/tap/eksctl; }
```
This Homebrew recipe also installs the [aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html), which is required.

### MacPorts
1. To install `eksctl` using [MacPorts](https://www.macports.org), run the following command. 
```bash
port install eksctl
```

## Linux
1. To download and extract the latest release of `eksctl`, run the following command.
```bash
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
```
2. Move the extracted binary to `/usr/local/bin`.
```bash
sudo mv /tmp/eksctl /usr/local/bin
```

## Windows 
### Install with Chocolatey
1. To install `eksctl` using [Chocolatey](https://chocolatey.org/install), run the following command. 
```bash
choco install -y eksctl 
```
### Upgrade with Chocolatey
1. To upgrade `eksctl` using [Chocolatey](https://chocolatey.org/install), run the following command. 
```bash
choco upgrade -y eksctl 
```

### Scoop
1. To install `eksctl` using [scoop](https://scoop.sh), run the following command. 
```bash
scoop install eksctl
```

### Docker
For every release and RC, a docker image is pushed to [weaveworks/eksctl](https://hub.docker.com/r/weaveworks/eksctl).


## Enable autocompletion
### Bash
1. To enable autocompletion in bash, run the following command, or put it in `~/.bashrc` or `~/.profile`.
```
. <(eksctl completion bash)
```

### Zsh
1. To enable autocompletion in zsh, run the following command.
```
mkdir -p ~/.zsh/completion/
eksctl completion zsh > ~/.zsh/completion/_eksctl
```
2. Add the following to `~/.zshrc`.
```
fpath=($fpath ~/.zsh/completion)
```

Note if you're not running a distribution, like oh-my-zsh, you may first need to enable autocompletion:
```
autoload -U compinit
compinit
```

To make the above persistent, run the first two lines, and put the above in `~/.zshrc`.

### Fish
1. To enable autocompletion in fish, run the following command.
```
mkdir -p ~/.config/fish/completions
eksctl completion fish > ~/.config/fish/completions/eksctl.fish
```

### Powershell
1. To enable autocompletion in PowerShell, run the following command.
```
eksctl completion powershell > C:\Users\Documents\WindowsPowerShell\Scripts\eksctl.ps1
```
???+ note
    You may need to specify a different path, depending on your system settings.
