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

## Communication

The project uses Slack. If you get stuck or just have a question then you are encouraged to join the
[Weave Community](https://weaveworks.github.io/community-slack/) Slack workspace and use the
[#eksctl](https://weave-community.slack.com/messages/eksctl/) channel and/or the [mailing
list][maillist].

We use the mailing list for some discussion, potentially for sharing documents
and for calendar invites.

Regular contributor meetings are held on Slack, see [`docs/contributor-meetings.md`](docs/contributor-meetings.md) for
the latest information.

[maillist]: https://groups.google.com/forum/#!forum/eksctl

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

*Optional:* the following assets can be updated:
- pkg/nodebootstrap/maxpods.go
- pkg/addons/default/assets/aws-node.yaml
by running:
```bash
make download-assets
```
If any of the files above are missing, then they will attempt to download when building.
Note that when submitting a PR, you should run this so that you get the latest files.

To run integration test you will need an AWS account.
```bash
make integration-test-container TEST_V=1
```

> NOTE: If you are working on Windows, you cannot use `make` at the moment,
> as the `Makefile` is currently not portable.
> However, if you have Git and Go installed, you can still build a binary
> and run unit tests.
> ```
> go build .\cmd\eksctl
> go test .\pkg\...
> ```

If you prefer to use Docker, the same way it is used in CI, you can use the
following comands:

```
make -f Makefile.docker test
make -f Makefile.docker build
```

> NOTE: It is not the most convenient way of working on the project, as
> binaries are built inside the container and cannot be tested manually,
> also majority of end-users consume binaries and not Docker images.
> It is recommended to use `make build` etc, unless there is an issue in CI
> that need troubleshooting.

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
ready remove the `WIP: ` title prefix and where possible squash your commits. Alternatively, you can use `Draft PR`
feature of Github as mentioned [here](https://github.blog/2019-02-14-introducing-draft-pull-requests/)

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

Added an explicit --profile flag that can be used to explicitly specify
which AWS profile you would like to use. This will override any profile
that you have specified via AWS_PROFILE.

If endpoints are being overridden then the credentials from the initial
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

1. Ensure integration tests pass (ETA: 2 hours ; more details below).
2. Determine the next release tag, e.g.:

   - for a release candidate, `0.13.0-rc.0`, or
   - for a release, `0.13.0`.

3. Create a `docs/release_notes/<version>.md` release notes file for the given tag using the contents of the release
   draft (generated by [`release-drafter`](https://github.com/release-drafter/release-drafter)), e.g.:

    ```console
    touch docs/release_notes/0.13.0.md
    ```
4a. For the first release candidate (`rc.0`) create a new branch after the major and minor numbers of the release (`release-X.Y`:

   ```console
   git checkout master
   git pull --ff-only origin master
   git checkout -b release-0.13
   ```

4b. If this is a subsequent release candidate or the release after an RC check out the existing release branch

   ```console
   git checkout release-0.13
   ```

6. Create the tags so that circleci can start the release process:

   - for a release candidate: `make prepare-release-candidate`. If there was an existing RC it will increase it's number, e.g.: rc.0 -> rc.1

   - for a release: `make prepare-release`. Regardless of whether there was a previous RC or a development version it will create a normal release

7. Ensure release jobs succeeded in [CircleCI](https://circleci.com/gh/weaveworks/eksctl).
8. Ensure the release was successfully [published in Github](https://github.com/weaveworks/eksctl/releases).
9. Download the binary just released, verify its checksum, and perform any relevant manual testing.

### Releasing snaps of eksctl

Snaps are software packages which run on a variety of Linux flavours.

*Note:* This is still somewhat [TBD](https://github.com/weaveworks/eksctl/issues/215).

Setting up the environment to build snaps on e.g. Ubuntu:

```sh
sudo apt install snapd
sudo snap install snapcraft --classic
sudo snap install multipass --classic
```

Building the snap can be done by running this command in the top-level directory:

```sh
snapcraft
```

*Note:* By default the `grade` of the snap is automatically set to `devel` and thus cannot be released to the e.g. `stable` channel of the snap. If you want to release a stable version, you want to check out the tag first, so e.g. `git checkout 0.12.0` and then run `snapcraft`.

Testing the resulting package can be done via:

```sh
sudo snap install eksctl_<version>_amd64.snap --dangerous
```

The `--dangerous` flag is required as it's a locally built snap and doesn't come from the store.

Publishing the snap can be done by following these steps.

Login to the Snap Store:

```sh
snapcraft login
```

Or

```sh
snapcraft login --with <login-file>
```

Where `<login-file>` was generated via `snapcraft export-login`.

Then publish via:

```sh
snapcraft push eksctl_<version>_amd64.snap --release [stable,beta,candidate,edge]
```

TODO: Further automate these steps in CircleCI, etc.

### Notes on Integration Tests

It's recommended to run containerised tests with `make integration-test-container TEST_V=1 AWS_PROFILE="<AWS profile name>"`. The tests require:

- Access to an AWS account. If there is an issue with access (e.g. expired MFA token), you will see all tests failing (albeit the error message may be slightly unclear).
- Access to the private SSH key for the Git repository to use for testing gitops-related operations. It is recommended to extract the private SSH key available [here](https://weaveworks.1password.com/vaults/all/allitems/kuxa5ujn7424jzkqqk7qtngovi) into `~/.ssh/eksctl-bot_id_rsa`, and then let the integration tests mount this path and use this key.

### Notes on Automation

When you run `make prepare-release` it will push a commit to master and a tag, which will trigger [release workflow](https://github.com/weaveworks/eksctl/blob/38364943776230bcc9ad57a9f8a423c7ec3fb7fe/.circleci/config.yml#L28-L42) in Circle CI. This runs `make eksctl-image` followed by `make release`. Most of the logic is defined in [`do-release.sh`](https://github.com/weaveworks/eksctl/blob/master/do-release.sh).

You want to keep an eye on Circle CI for the progress of the release ([0.3.1 example logs](https://circleci.com/workflow-run/02d8b5fb-bc7f-404c-9051-68307c124649)). It normally takes around 30 minutes.

### Latest release

To get the latest release you can use the link [https://github.com/weaveworks/eksctl/releases/latest]().

**Note** Previously, eksctl used a floating tag called `latest_release`. This is _deprecated_ and it will stop working after release `0.14.0`.
