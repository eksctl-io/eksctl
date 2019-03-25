# How to Contribute

*eksctl* is [Apache 2.0 licenced](LICENSE) and accepts contributions via GitHub
pull requests. This document outlines some of the conventions on the development
workflow, commit message formatting, contact points and other resources to make
it easier to get your contribution accepted.

We gratefully welcome improvements to documentation as well as to code.

# Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. No action from you is required, but it's a good idea to see the
[DCO](DCO) file for details before you start contributing code to eksctl.

# Chat

The project uses Slack. If you get stuck or just have a question then you are encouraged to join the 
[Weave Community](https://weaveworks.github.io/community-slack/) Slack workspace and use the 
[#eksctl](https://weave-community.slack.com/messages/eksctl/) channel.

Regular contributor meetings are held on Slack, see [`docs/contributor-meetings.md`](docs/contributor-meetings.md) for 
the latest information.

# Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for getting started as a user and learn how/where to ask for help
- If you want to contribute as a developer, continue reading this document for further instructions
- Play with the project, submit bugs, submit pull requests!

## Contribution workflow


#### 1. Set up your Go environment

This project is written in Go. To be able to contribute you will need a working Go installation. You can check the 
[official installation guide](https://golang.org/doc/install). Make sure your `GOPATH` and `GOBIN` 
[environment variables](https://github.com/golang/go/wiki/SettingGOPATH) are set correctly.

Here are some quick start steps you can follow:

Select a folder path for your go projects, for example `~/eksctl`
```bash
mkdir ~/eksctl
cd ~/eksctl
export GOPATH="$(pwd)"
export GOBIN="${GOPATH}/bin"
```

> NOTE: Windows users should install Docker for Windows and run `make eksctl-image` to build their code.

> TODO: Improve Windows instructions, ensure `go build` works.

#### 2. Fork and clone the repo

Clone the repo using go:

```bash
go get -d github.com/weaveworks/eksctl/...
```

Make a fork of this repository and add it as a remote:

```bash
cd src/github.com/weaveworks/eksctl
git remote add <username> git@github.com:<username>/eksctl.git
```

#### 3. Run the tests and build eksctl

Make sure you can run the tests and build the binary. For this project you need go version 1.12 or higher.

```bash 
make install-build-deps
make test
make build
```


There are integration tests for *eksctl* being developed and more details of how to run them will be included here. You 
can follow the progress [here](https://github.com/weaveworks/eksctl/issues/151).

#### 4. Write your feature

- Find an [issue](https://github.com/weaveworks/eksctl/issues) to work on or create your own. If you are a new 
contributor take a look at issues marked with 
[good first issue](https://github.com/weaveworks/eksctl/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

- Then create a topic branch from where you want to base your work (usually branched from master):

    ```bash
    git checkout -b <feature-name>
    ```

- Write your feature. Make commits of logical units and make sure your commit messages are in the 
[proper format](#format-of-the-commit-message).

- Add automated tests to cover your changes. See the [az](https://github.com/weaveworks/eksctl/tree/master/pkg/az) 
package for a good example of tests.

- If needed, update the documentation, either in the [README](README.md) or in the [docs](docs/) folder.

- Make sure the tests are running successfully.

#### 5. Submit a pull request

Push your changes to your fork and submit a pull request to the original repository. If your PR is a work in progress 
then make sure you prefix the title with `WIP: `. This lets everyone know that this is still being worked on. Once its 
ready remove the `WIP: ` title prefix and where possible squash your commits.

```bash
git push <username> <feature-name>
```

Our CircleCI integration will run the automated tests and give you feedback in the review section. We will review your 
changes and give you feedback as soon as possible.

# Acceptance policy

These things will make a PR more likely to be accepted:

 * a well-described requirement
 * tests for new code
 * tests for old code!
 * new code and tests follow the conventions in old code and tests
 * a good commit message (see below)

In general, we will merge a PR once a maintainer has reviewed and approved it.
Trivial changes (e.g., corrections to spelling) may get waved through.
For substantial changes, more people may become involved, and you might get asked to resubmit the PR or divide the 
changes into more than one PR.

### Format of the Commit Message

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```
Added AWS Profile Support

Changes to ensure that AWS profiles are supported. This involved making
sure that the AWS config file is loaded (SharedConfigEnabled) and
also making sure we have a TokenProvider set.

Added an explicit --profile flag that can be used to explicity specify
which AWS profile you would like to use. This will override any profile
that you have specified via AWS_PROFILE.

If endpoints are being overriden then the credentials from the initial
session creation are shared with any subsequent session creation to
ensure that the tokens are shared (otherwise you may get multiple MFA
prompts).

Issue #57
```

The format can be described more formally as follows:

```
<short title for what changed>
<BLANK LINE>
<why this change was made and what changed>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various git tools.

