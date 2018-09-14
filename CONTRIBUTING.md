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

The project uses Slack. If you get stuck or just have a question then you are encouraged to join the [Weave Community](https://weaveworks.github.io/community-slack/) Slack workspace and use the [#eksctl](https://weave-community.slack.com/messages/eksctl/) channel.

## Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for getting started as a user and learn how/where to ask for help
- If you want to contribute as a developer, continue reading this document for further instructions
- Play with the project, submit bugs, submit pull requests!

## Contribution workflow

This is a rough outline of how to prepare a contribution:

- Find an [issue](https://github.com/weaveworks/eksctl/issues) to work on. If you are a new contributor
take a look at issues marked with [good first issue](https://github.com/weaveworks/eksctl/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).
- Create a topic branch from where you want to base your work (usually branched from master).
- Make commits of logical units.
- Make sure your commit messages are in the proper format (see below).
- Push your changes to a topic branch in your fork of the repository.
- If you changed code:
   - add automated tests to cover your changes. See the [az](https://github.com/weaveworks/eksctl/tree/master/pkg/az) package for a good example of tests.
- Submit a pull request to the original repository.

If your PR is a work in progress then make sure you prefix the title with `WIP: `. This lets everyone know that this is still being worked on. Once its ready
remove the `WIP: ` title prefix and where possible squash your commits.

## How to build and run the project

```bash
make build
./eksctl get clusters
```

## How to run the test suite

You can run the unit tests by simply doing

```bash
make test
```

There are integration tests for *eksctl* being developed and more details of how to run them will be included here. You can follow the progress [here](https://github.com/weaveworks/eksctl/issues/151).

# Acceptance policy

These things will make a PR more likely to be accepted:

 * a well-described requirement
 * tests for new code
 * tests for old code!
 * new code and tests follow the conventions in old code and tests
 * a good commit message (see below)

In general, we will merge a PR once a maintainer has reviewed and approved it.
Trivial changes (e.g., corrections to spelling) may get waved through.
For substantial changes, more people may become involved, and you might get asked to resubmit the PR or divide the changes into more than one PR.

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

