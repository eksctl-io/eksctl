# How to Contribute

*eksctl* is [Apache 2.0 licenced](LICENSE) and accepts contributions via GitHub
pull requests. This document outlines some of the conventions on the development
workflow, commit message formatting, contact points and other resources to make
it easier to get your contribution accepted.

We gratefully welcome improvements to documentation as well as to code.

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. No action from you is required, but it's a good idea to see the
[DCO](DCO) file for details before you start contributing code to eksctl.

## Chat

The project uses Slack. If you get stuck or just have a question then you are encouraged to join the
[Weave Community](https://weaveworks.github.io/community-slack/) Slack workspace and use the
[#eksctl](https://weave-community.slack.com/messages/eksctl/) channel.

Regular contributor meetings are held on Slack, see [`docs/contributor-meetings.md`](docs/contributor-meetings.md) for
the latest information.

## Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for getting started as a user and learn how/where to ask for help
- If you want to contribute as a developer, continue reading this document for further instructions
- Play with the project, submit bugs, submit pull requests!

### Contribution workflow

#### 1. Set up your Go environment

This project is written in Go. To be able to contribute you will need:

1. A working Go installation of Go >= 1.12. You can check the
[official installation guide](https://golang.org/doc/install).

2. Make sure that `$(go env GOPATH)/bin` is in your shell's `PATH`. You can do so by
   running `export PATH="$(go env GOPATH)/bin:$PATH"`

#### 2. Fork and clone the repo

Make a fork of this repository and clone it by running:

```bash
git clone git@github.com:<yourusername>/eksctl.git
```

It is not recommended to clone under your `GOPATH` (if you define one). Otherwise, you will need to set
`GO111MODULE=on` explicitly.

#### 3. Run the tests and build eksctl

Make sure you can run the tests and build the binary.

```bash
make install-build-deps
make test
make build
```

> NOTE: Windows users should install Docker for Windows and run `make eksctl-image` to build their code.

> TODO: Improve Windows instructions, ensure `go build` works.

There are integration tests for *eksctl* being developed and more details of
how to run them will be included here. You can follow the progress [here](https://github.com/weaveworks/eksctl/issues/151).

#### 4. Write your feature

- Find an [issue](https://github.com/weaveworks/eksctl/issues) to work on or
  create your own. If you are a new contributor take a look at issues marked
  with [good first issue](https://github.com/weaveworks/eksctl/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

- Then create a topic branch from where you want to base your work (usually branched from master):

    ```bash
    git checkout -b <feature-name>
    ```

- Write your feature. Make commits of logical units and make sure your
  commit messages are in the [proper format](#format-of-the-commit-message).

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

## Acceptance policy

These things will make a PR more likely to be accepted:

- a well-described requirement
- tests for new code
- tests for old code!
- new code and tests follow the conventions in old code and tests
- a good commit message (see below)

In general, we will merge a PR once a maintainer has reviewed and approved it.
Trivial changes (e.g., corrections to spelling) may get waved through.
For substantial changes, more people may become involved, and you might get asked to resubmit the PR or divide the
changes into more than one PR.

### Format of the Commit Message

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```text
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

```text
<short title for what changed>
<BLANK LINE>
<why this change was made and what changed>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various git tools.

## Release Process

1. Ensure integration tests pass (ETA: 45 minutes ; more details below).
2. Determine next release tag (e.g. `0.1.35`).
3. Create release notes file for the given tag – `docs/release_notes/<tag>.md` (e.g. `docs/release_notes/0.1.35.md`).
4. Run `./tag-release-candidate.sh <tag>-rc.<N>` or `./tag-release.sh <tag>` (e.g. `./tag-release-candidate.sh 0.1.35-rc.0` or `./tag-release.sh 0.1.35`).
5. Ensure release jobs succeeded in [CircleCI](https://circleci.com/gh/weaveworks/eksctl).
6. Ensure the release was successfully [published in Github](https://github.com/weaveworks/eksctl/releases).
7. Download the binary just released, verify its checksum, and perform any relevant manual testing.

### Notes on Integration Tests

It's recommended to run containerised tests with `make integration-test-container TEST_V=1 AWS_PROFILE="<AWS profile name>"`. The tests require access to an AWS account. If there is an issue with access (e.g. expired MFA token), you will see all tests failing (albeit the error message may be slightly unclear).

At present we ignore flaky tests, so if you see output like show below, you don't need to worry about this for the purpose of the release. However, you might consider reviewing the issues in question after you made the release.

```console
$ make integration-test-container TEST_V=1 AWS_PROFILE="default-mfa"
[...]
Summarizing 2 Failures:

[Fail] (Integration) Create, Get, Scale & Delete when creating a cluster with 1 node and add the second nodegroup and delete the second nodegroup [It] {FLAKY: https://github.com/weaveworks/eksctl/issues/717} should make it 4 nodes total
/go/src/github.com/weaveworks/eksctl/integration/creategetdelete_test.go:376

[Fail] (Integration) Create, Get, Scale & Delete when creating a cluster with 1 node and scale the initial nodegroup back to 1 node [It] {FLAKY: https://github.com/weaveworks/eksctl/issues/717} should make it 1 nodes total
/go/src/github.com/weaveworks/eksctl/integration/creategetdelete_test.go:403

Ran 26 of 26 Specs in 2556.238 seconds
FAIL! -- 24 Passed | 2 Failed | 0 Pending | 0 Skipped
--- FAIL: TestSuite (2556.25s)
```

### Notes on Automation

When you run `./tag-release.sh <tag>` it will push a commit to master and a tag, which will trigger [release workflow](https://github.com/weaveworks/eksctl/blob/38364943776230bcc9ad57a9f8a423c7ec3fb7fe/.circleci/config.yml#L28-L42) in Circle CI. This runs `make eksctl-image` followed by `make release`. Most of the logic is defined in [`do-release.sh`](https://github.com/weaveworks/eksctl/blob/master/do-release.sh).

You want to keep an eye on Circle CI for the progress of the release ([0.1.35 example logs](https://circleci.com/workflow-run/3553542c-88ad-4a77-bd42-441da4c87fa1)). It normally takes 15-20 minutes.

### Notes on Artefacts

We use `latest_release` floating tag, in order to enable static URLs for release artefacts, i.e. `latest_release` gets shifted on every release.

That means you will see two entries on [the release page](https://github.com/weaveworks/eksctl/releases):

- [**eksctl 0.1.35 (permalink)**](https://github.com/weaveworks/eksctl/releases/tag/0.1.35)
- [**eksctl 0.1.35**](https://github.com/weaveworks/eksctl/releases/tag/latest_release)
