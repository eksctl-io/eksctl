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

1. Ensure integration tests pass (ETA: 45 minutes ; more details below).
2. Determine the next release tag, e.g.:

   - for a release candidate, `0.4.0-rc.0`, or
   - for a release, `0.4.0`.

3. Create a `docs/release_notes/<tag>.md` release notes file for the given tag, e.g.:

    ```console
    touch docs/release_notes/0.4.0.md
    ```

4. Check out the latest `master`:

    ```console
    git checkout master
    git fetch origin master
    git merge --ff-only origin/master
    ```

5. Run:

   - for a release candidate: `./tag-release-candidate.sh <tag>-rc.<N>`, e.g.:

     ```console
     ./tag-release-candidate.sh 0.4.0-rc.0
     ```

   - for a release: `./tag-release.sh <tag>`, e.g.:

     ```console
     ./tag-release.sh 0.4.0
     ```

6. Ensure release jobs succeeded in [CircleCI](https://circleci.com/gh/weaveworks/eksctl).
7. Ensure the release was successfully [published in Github](https://github.com/weaveworks/eksctl/releases).
8. Download the binary just released, verify its checksum, and perform any relevant manual testing.

### Notes on Integration Tests

It's recommended to run containerised tests with `make integration-test-container TEST_V=1 AWS_PROFILE="<AWS profile name>"`. The tests require:

- Access to an AWS account. If there is an issue with access (e.g. expired MFA token), you will see all tests failing (albeit the error message may be slightly unclear).
- Access to the private SSH key for the Git repository to use for testing gitops-related operations. It is recommended to extract the private SSH key available [here](https://weaveworks.1password.com/vaults/all/allitems/kuxa5ujn7424jzkqqk7qtngovi) into `~/.ssh/eksctl-bot_id_rsa`, and then let the integration tests mount this path and use this key.

### Notes on Automation

When you run `./tag-release.sh <tag>` it will push a commit to master and a tag, which will trigger [release workflow](https://github.com/weaveworks/eksctl/blob/38364943776230bcc9ad57a9f8a423c7ec3fb7fe/.circleci/config.yml#L28-L42) in Circle CI. This runs `make eksctl-image` followed by `make release`. Most of the logic is defined in [`do-release.sh`](https://github.com/weaveworks/eksctl/blob/master/do-release.sh).

You want to keep an eye on Circle CI for the progress of the release ([0.3.1 example logs](https://circleci.com/workflow-run/02d8b5fb-bc7f-404c-9051-68307c124649)). It normally takes around 30 minutes.

### Notes on Artefacts

We use `latest_release` floating tag, in order to enable static URLs for release artefacts, i.e. `latest_release` gets shifted on every release.

That means you will see two entries on [the release page](https://github.com/weaveworks/eksctl/releases):

- [**eksctl 0.4.0 (permalink)**](https://github.com/weaveworks/eksctl/releases/tag/0.4.0)
- [**eksctl 0.4.0**](https://github.com/weaveworks/eksctl/releases/tag/latest_release)
