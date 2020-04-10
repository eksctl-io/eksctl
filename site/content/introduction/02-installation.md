---
title: "Installation"
weight: 20
url: introduction/installation
---

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

You will need to have AWS API credentials configured. What works for AWS CLI or any other tools (kops, Terraform etc), should be sufficient. You can use [`~/.aws/credentials` file][awsconfig]
or [environment variables][awsenv]. For more information read [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html).

[awsenv]: https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html
[awsconfig]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html

You will also need [AWS IAM Authenticator for Kubernetes](https://github.com/kubernetes-sigs/aws-iam-authenticator) command (either `aws-iam-authenticator` or `aws eks get-token` (available in version 1.16.156 or greater of AWS CLI) in your `PATH`.

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
