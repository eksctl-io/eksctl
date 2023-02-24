# Installation

## Check the version
To check the version of `eksctl` running on your local machine, run the following command.
```bash
eksctl version
```
You should see the version in the response output. For example:
```bash
0.130.0
```

## macOS (Homebrew)
1. To install or upgrade `eksctl` using [Homebrew](https://formulae.brew.sh/formula/eksctl), run the following command. 
```bash
brew install eksctl
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

## Windows (Chocolatey)
1. To install `eksctl` using [Chocolatey](https://chocolatey.org/install), run the following command. 
```bash
choco install -y eksctl 
```
2. To upgrade `eksctl` using [Chocolatey](https://chocolatey.org/install), run the following command. 
```bash
choco upgrade -y eksctl 
```

