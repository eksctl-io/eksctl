# eksctl :heart:s your contributions

Thank you for taking the time to contribute to `eksctl`.

We gratefully welcome improvements to all areas; from code, to documentation;
from bug reports to feature design.

This guide should cover all aspects of how to interact with the project
and how to get involved in development as smoothly as possible.

If we have missed anything you think should be included, or if anything is not
clear, we also accept contributions to this contribution doc :smile:.

For information on how to get in touch, or how the project is run, please see
the [Community page](COMMUNITY.md).

_(If you are a Maintainer, head over to the [Maintainer's guide](https://github.com/weaveworks/eksctl-handbook).)_

Reading docs is often tedious, so we'll put our most important contributing rule
right at the top: **Always be kind!**

Looking forward to seeing you in the repo! :sparkles:

# Table of Contents

- [Legal bits](#legal-bits)
- [Where can I get involved?](#where-can-i-get-involved)
	- [The eksctl roadmap](#the-eksctl-roadmap)
- [Opening Issues](#opening-issues)
	- [Bug report guide](#bug-report-guide)
	- [Feature request guide](#feature-request-guide)
	- [Help request guide](#help-request-guide)
- [Submitting PRs](#submitting-prs)
	- [Choosing something to work on](#choosing-something-to-work-on)
	- [Developing eksctl](#developing-eksctl)
	- [Asking for help](#asking-for-help)
	- [PR submission guidelines](#pr-submission-guidelines)
- [How the Maintainers process contributions](#how-the-maintainers-process-contributions)
	- [Prioritizing issues](#prioritizing-issues)
	- [Reviewing PRs](#reviewing-prs)
- [Proposals](#proposals)

---

# Legal bits

## License
`eksctl` is [Apache 2.0 licenced](LICENSE).

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. No action from you is required, but it's a good idea to see the
[DCO](DCO) file for details before you start contributing code to eksctl.

---

# Where can I get involved?

We are happy to see people in pretty much all areas of eksctl's development.
Here is a non-exhaustive list of ways you can help out:

1. Open a [PR](#submitting-prs). :woman_technologist:

    Beyond fixing bugs and submitting new features, there are other things you can submit
    which, while less flashy, will be deeply appreciated by all who interact with the codebase:

      - Backfilling tests! (our coverage is super low right now.)
      - Refactoring! (omigod such tech debt.)
      - Reviewing and updating [documentation](https://eksctl.io/)! (seems to be we only notice the documentation that doesn't exist.)

   (See also [Choosing something to work on](#choosing-something-to-work-on) below.)
1. Open an [issue](#opening-issues). :interrobang:

    We have 3 forms of issue: [bug reports](#bug-report-guide), [feature requests](#feature-request-guide) and [help requests](#help-request-guide).
    If you are not sure which category you need, just make the best guess and provide as much information as possible.
1. Help out others in the [community slack channel](https://weave-community.slack.com/messages/CAYBZBWGL/). :sos:
1. Chime in on [bugs](https://github.com/weaveworks/eksctl/issues?q=is%3Aopen+is%3Aissue+label%3Akind%2Fbug+), [feature requests](https://github.com/weaveworks/eksctl/issues?q=is%3Aopen+is%3Aissue+label%3Akind%2Ffeature) and asks for [help](https://github.com/weaveworks/eksctl/issues?q=is%3Aopen+is%3Aissue+label%3Akind%2Fhelp). :thought_balloon:
1. Dig into some [`needs-investigation` issues](https://github.com/weaveworks/eksctl/labels/needs-investigation) :detective:
1. Get involved in some [PR reviews](https://github.com/weaveworks/eksctl/pulls). :nerd_face:
1. Review old issues; poke or suggest closing ones which are stale or follow up those which still look good. :speech_balloon:
1. Think deeply about the future of `eksctl`, then [talk about it](#proposals). :crystal_ball:

## The eksctl roadmap

... can be found [here](https://eksctl.io/community/).

This is the general long game we are playing. Contributions which steer us further towards
these overarching goals will be enthusiastically welcomed, while contributions which are
orthogonal or leading in another direction will likely be rejected.

For bonus points, please highlight how feature requests/PRs fit into the roadmap when you can!

---

# Opening Issues

These guides aim to help you write issues in a way which will ensure that they are processed
as quickly as possible.

_See below for [how issues are prioritized](#prioritizing-issues)_.

**General rules**:

1. Before opening anything, take a good look through existing issues.
1. More is more: give as much information as it is humanly possible to give.
  Highly detailed issues are more likely to be picked up because they can be prioritized and
  scheduled for work faster. They are also more accessible
  to the community, meaning that you may not have to wait for the core team to get to it.
1. Please do not open an issue with a description that is simply a link to another issue,
  a link to a slack conversation, a quote from either one of those, or anything else
  equally opaque. This raises the bar for entry and makes it hard for the community
  to get involved. Take the time to write a proper description and summarise key points.
1. Take care with formatting. Ensure the [markdown is tidy](https://docs.github.com/en/free-pro-team@latest/github/writing-on-github/getting-started-with-writing-and-formatting-on-github),
  use [code blocks](https://docs.github.com/en/free-pro-team@latest/github/writing-on-github/creating-and-highlighting-code-blocks) etc etc.
  The faster something can be read, the faster it can be dealt with.
1. Keep it civil. Yes, it is annoying when things don't work, but it is way more fun helping out
  someone who is not... the worst. Remember that conversing via text exacerbates
  everyone's negativity bias, so throw in some emoji when in doubt :+1: :smiley: :rocket: :tada:.

**Dedicated guides**:
- [Bug report guide](#bug-report-guide)
- [Feature request guide](#feature-request-guide)
- [Help request guide](#help-request-guide)

## Bug report guide

We hope to get to bug reports within a couple of working days, but please wait for at least
7 before nudging. (Unless it is a super critical end-of-the world bug, then by all means
make some noise :loudspeaker:.)

Below are the criteria we like our bug reports to cover in order to gather the bare minimum of
information. Add more that what is asked for if you can :smiley:.

1. **Search existing issues.** If something similar already exists, and is still open, please contribute to the discussion there.

1. **Bump to the latest version of eksctl** to see whether your issue has already been fixed.

1. **Write a concise and descriptive title**, like you would a commit message, something which others can easily
  find when doing step 1 above.

1. **Detail what it was that you were trying to do and what you expected would happen**.
  Give some background information around what you have already done to your cluster, any custom configuration etc.
  With sufficient information you can pre-empt any questions others may have. This should cut out some obvious
  back-and-forth and help get people to the heart of the issue quickly.

1. **Explain what actually happened**. Provide the relevant error message and key logs.

1. **Provide a reproduction**, for example the exact command or a yaml file (sensitive info redacted).
  Please try to reduce your reproduction to the minimal necessary to help whoever is helping you
  get to the broken state without needing to recreate your entire environment.

1. **If possible, reproduce the issue with logging verbosity set to at least 4** (`-v=4`), if you have not already done so. Ensure
  logs are formatted with [code blocks](https://docs.github.com/en/free-pro-team@latest/github/writing-on-github/creating-and-highlighting-code-blocks).
  If they are long (>50 lines) please provide them in a Gist or collapsed using
  [HTML details tags](https://gist.github.com/ericclemmons/b146fe5da72ca1f706b2ef72a20ac39d).
  Take care to redact any sensitive info.

1. Paste in the outputs of `eksctl info` where relevant, as
  well as anything else you have running which you think may be relevant.

1. Detail any workarounds which you tried, it may help others who experience the same problem.

1. If you already have a fix in mind, note that on the report and go ahead and open a
  PR whenever you are ready. A core team-member will assign the issue to you.

## Feature request guide

We hope to respond to and prioritize new feature requests within 7 working days. Please wait for
up to 14 before nudging us.

A feature request is the start of a discussion, so don't be put off if it is not
accepted. Features which either do not contribute to or directly work against
the project [goals](#the-eksctl-vision) will likely be rejected, as will highly
specialised usecases.

Once you have opened the ticket, feel free to post it in the community
[slack channel](https://weave-community.slack.com/messages/CAYBZBWGL/) to get more opinions on it.

Below are the steps we encourage people to take when creating a new feature request:

1. **Search existing issues.** If something similar already exists, and is still open, please contribute to the discussion there.

1. **Explain clearly why you want this feature.**

1. **Describe the behaviour you'd like to see.** As well as an explanation, please
  provide some example commands/config/output. Please ensure everything is formatted
  nicely with [code blocks](https://docs.github.com/en/free-pro-team@latest/github/writing-on-github/creating-and-highlighting-code-blocks).
  If you have strong ideas, be as detailed as you like.

1. Explain how this feature would work towards the [project's vision](#the-eksctl-vision),
  or how it would benefit the community.

1. Note the deliverable of this issue: should the outcome be a simple PR to implement the
  feature? Or does it need a design [proposal](#proposals)?

1. If the feature is small (maybe it is more of an improvement) and you already have
  a solution in mind, explain what you plan to do on the issue and open a PR!
  A core team member will assign the task to you.

## Help request guide

While you can ask for general help with `eksctl` usage in the [slack channel](https://weave-community.slack.com/messages/CAYBZBWGL/),
opening an issue creates a more searchable history for the community.

We hope to respond to requests for help within a couple of working days, but please wait
for a week before nudging us.

Once you have created your issue, we recommend posting it in the slack channel
to get more eyes on it faster.

Please following these steps when opening a new help issue:

1. **Search existing issues.** If something similar already exists, and is still open, please contribute to the discussion there.

1. Write a title with the format "How to x".

1. Explain what you are trying to accomplish, what you have tried, and the behaviour you are seeing.

1. Please include your config (removing any sensitive information) or exact the commands you're using.
   Please ensure everything is formatted nicely with [code blocks](https://docs.github.com/en/free-pro-team@latest/github/writing-on-github/creating-and-highlighting-code-blocks).

1. When providing verbose logs, please use either a Gist or [HTML detail tags](https://gist.github.com/ericclemmons/b146fe5da72ca1f706b2ef72a20ac39d).

---

# Submitting PRs
## Choosing something to work on

If you are not here to report a bug, ask for help or request some new behaviour, this
is the section for you. We have curated a set of issues for anyone who simply
wants to build up their open-source cred :muscle:.

- Issues labelled [`good first issue`](https://github.com/weaveworks/eksctl/labels/good%20first%20issue)
  should be accessible to folks new to the repo, as well as to open source in general.

  These issues should present a low/non-existent barrier to entry with a thorough description,
  easy-to-follow reproduction (if relevant) and enough context for anyone to pick up.
  The objective should be clear, possibly with a suggested solution or some pseudocode.
  If anything similar has been done, that work should be linked.

  If you have come across an issue tagged with `good first issue` which you think you would
  like to claim but isn't 100% clear, please ask for more info! When people write issues
  there is a _lot_ of assumed knowledge which is very often taken for granted. This is
  something we could all get better at, so don't be shy in asking for what you need
  to do great work :smile:.

  See more on [asking for help](#asking-for-help) below!

- [`help wanted` issues](https://github.com/weaveworks/eksctl/labels/help%20wanted)
  are for those a little more familiar with the code base, but should still be accessible enough
  to newcomers.

- All other issues labelled `kind/<x>` or `priority/<x>` are also up for grabs, but
  are likely to require a fair amount of context.

## Developing eksctl

**Sections:**
- [Getting started](#getting-started)
- [Setting up your Go environment](#setting-up-your-go-environment)
- [Forking the repo](#forking-the-repo)
- [Building eksctl](#building-eksctl)
- [Running the unit tests](#running-the-unit-tests)
- [Running the integration tests](#running-the-integration-tests)
- [Writing your solution](#writing-your-solution)

> WARNING: All commands in this section have only been tested on Linux/Unix systems.
> There is no guarantee that they will work on Windows.

### Getting started

Before you begin writing code, you may want to have a play with `eksctl` to get familiar
with the tool. Check out the [README](README.md) for basic installation and usage,
then head to the [main docs](https://eksctl.io/) for more information.

### Setting up your Go environment

This project is written in Go. To be able to contribute you will need:

1. A working Go installation of Go >= 1.12. You can check the
[official installation guide](https://golang.org/doc/install).

2. Make sure that `$(go env GOPATH)/bin` is in your shell's `PATH`. You can do so by
   running `export PATH="$(go env GOPATH)/bin:$PATH"`

3. (Optional) [User documentation](https://eksctl.io) is built and generated with [mkdocs](https://www.mkdocs.org/).
   Please make sure you have python3 and pip installed on your local system.

### Forking the repo

Make a fork of this repository and clone it by running:

```bash
git clone git@github.com:<yourusername>/eksctl.git
```

It is not recommended to clone under your `GOPATH` (if you define one), otherwise, you will need to set
`GO111MODULE=on` explicitly.

You may also want to add the original repo to your remotes to keep up to date
with changes.

### Building eksctl

> NOTE: If you are working on Windows, you cannot use `make` at the moment,
> as the `Makefile` is currently not portable.
> However, if you have Git and Go installed, you can still build a binary
> and run unit tests.
> ```
> go build .\cmd\eksctl
> go test .\pkg\...
> ```


Once in your cloned repo, you can install the dependencies and build the binary.
The binary will be created in the root of your repo `./eksctl`.

```bash
make install-build-deps
make build
```

To build and serve the user docs, execute the following:

```bash
# Requires python3 and pip3 installed in your local system
make install-site-deps
make build-pages
make serve-pages
```

### Running the unit tests

To run the tests simply run the following after `install-build-deps`:

```bash
make test
```

If you prefer to use Docker, the same way it is used in CI, you can use the
following command:

```
make -f Makefile.docker test
```

> NOTE: It is not the most convenient way of working on the project, as
> binaries are built inside the container and cannot be tested manually,
> also majority of end-users consume binaries and not Docker images.
> It is recommended to use `make build` etc, unless there is an issue in CI
> that need troubleshooting.

### Running the integration tests

> NOTE: Some parts of the integration tests are not configurable and therefore
> cannot be run by folks outside the core team. If you are NOT contributing to the
> gitops functionality, you can run a subset of the tests which cover your change,
> see below.

> NOTE: Integration tests a lot of infrastructure and are therefore quite expensive
  (in both sense of the word) to run. It is therefore not essential for community
  members to run them as the core team does this as part of the release process.

The integration tests are long and unfortunately there are some flakes (help is very
welcome!).

There are several ways to run the tests. Requirements are:
- An AWS account (it is recommended to use [gsts](https://github.com/ruimarinho/gsts) to authenticate)
  which you are logged in to with a session token which won't expire for at least 2 hours.
- An empty repository for testing gitops operations
- A private key to that gitops repository

To run the full suite including cluster creations/teardowns:

```bash
TEST_V=1 make integration-test
```

To run the tests and save output to a file (recommended), run:

```bash
TEST_V=1 make integration-test 2>&1 | tee <some-name>.log
```

> NOTE: The test suites are run in parallel, which means they write to `stdout` simultaneously.
> To parse logs for a specific test output, you can grep the logs based on the node number
> (eg. `grep '\[9\]' int-tests.log > suite9.log`)

To run the tests with a pre-created cluster for a faster turnaround:

```bash
TEST_CLUSTER=<name> make create-integration-test-dev-cluster
TEST_CLUSTER=<name> make integration-test-dev
TEST_CLUSTER=<name> make delete-integration-test-dev-cluster
```

To run a specific suite:

```bash
ginkgo -tags integration -v --progress integration/tests/<suite name>/... -- -test.v -ginkgo.v
```

### Writing your solution

Once you have your environment set up and have completed a clean run of the unit
tests you can get to work :tada: .

1. First create a topic branch from where you want to base your work (this is usually
  from `main`):

      ```bash
      git checkout -b <feature-name>
      ```

1. Write your solution. Try to align with existing patterns and standards.
  _However_, there is significant tech debt, so any refactoring or changes which would
  improve things even a little would be very welcome. (See [#2931](https://github.com/weaveworks/eksctl/issues/2931)
  for our current efforts.)

1. Try to commit in small chunks so that changes are well described
  and not lumped in with something unrelated. This will make debugging easier in
  the future.
  Make sure commit messages are in the [proper format](#commit-message-formatting).

1. Be sure to include at least unit tests to cover your changes. See the [addon](https://github.com/weaveworks/eksctl/blob/main/pkg/actions/addon)
  package for a good example of tests.

      > NOTE: We are trying to move away from using [`mockery`](https://github.com/vektra/mockery)
      > to generate our fakes. Where possible, please use [`counterfeiter`](https://github.com/maxbrunsfeld/counterfeiter)
      > instead.

1. For extra special bonus points, if you see any tests missing from the area you are
  working on, please add them! It will be much appreciated :heart: .

1. Check the documentation and update it to cover your changes,
  either in the [README](README.md) or in the [docs](docs/) folder.
  If you have added a completely new feature please ensure that it is documented
  thoroughly and include an [example](examples/).

1. Before you [open your PR](#pr-submission-guidelines), run all the unit tests and manually
  verify that your solution works.

1. Note that editing certain things (eg. `pkg/apis/eksctl.io/v1alpha5/types.go`) will mean
  you need to ensure autogenerated files are updated. Not doing so will result in
  merge conflicts on your PR. Running `make test` will handle the generation. If you
  get into a state where you have forgotten to generate, and there is a conflict,
  you can resolve this by accepting `HEAD` in the conflict resolution, and then
  running `make test` again. You can then commit the generated files as part of your
  PR.

## Asking for help

If you need help at any stage of your work, please don't hesitate to ask!

- To get more detail on the issue you have chosen, it is a good idea to start by asking
  whoever created it to provide more information.
  If they do not respond, or more help is needed,
  you can then bring in one of the [core maintainers](COMMUNITY.md#core-team).

- If you are struggling with something while working on your PR, or aren't quite
  sure of your approach, you can open a [draft](https://github.blog/2019-02-14-introducing-draft-pull-requests/)
  (prefix the title with `WIP: `) and explain what you are thinking.
  You can tag in one of the core team, or drop the PR link into [slack](https://weave-community.slack.com/messages/eksctl/) and get
  help from the community.

- We are happy to pair with contributors over a slack call to help them fine-tune their
  implementation. You can ping us directly, or head to the [channel](https://weave-community.slack.com/messages/eksctl/)
  to see if anyone in the community is up for being a buddy :smiley: .

## PR submission guidelines

Push your changes to the branch on your fork and submit a pull request to the original repository
against the `main` branch.
Where possible, please squash your commits to ensure a tidy and descriptive history.

```bash
git push <remote-name> <feature-name>
```

If your PR is still a work in progress, please open a [Draft PR](https://github.blog/2019-02-14-introducing-draft-pull-requests/).

Our GitHub Actions integration will run the automated tests and give you feedback in the review section. We will review your
changes and give you feedback as soon as possible. We also encourage people to post
links to their PRs in slack to get more eyes on the work.

We recommend that you regularly rebase from `main` of the original repo to keep your
branch up to date.

Please ensure that `Allow edits and access to secrets by maintainers` is checked.
While the maintainers will of course wait for you to edit your own work, if you are
unresponsive for over a week, they may add corrections or even complete the work for you,
especially if what you are contributing is very cool :metal: .

PRs which adhere to our guidelines are more likely to be accepted
(when opening the PR, please use the checklist in the template):

1. **The description is thorough.** When writing your description, please be as detailed as possible: don't make people
  guess what you did or simply link back to the issue (the issue explains the problem
  you are trying to solve, not how you solved it.)
  Guide your reviewers through your solution by highlighting
  key changes and implementation choices. Try and pre-empt any obvious questions
  they may have. Providing snippets (or screenshots) of output is very helpful to
  demonstrate new behaviour or UX changes. (Snippets are more searchable than screenshots,
  but we wont be mad at a sneak peek at your terminal envs :eyes: .)

1. **The change has been manually tested.** If you are supplying output above
  then that can be your manual test, with proof :clap: .

1. **The PR has a snappy title**. Your PR title will end up in the release notes,
  so make it a good one. Using the same rule as for the title of a [commit message](#commit-message-formatting)
  is generally a good idea. Try to use the imperative and centre it around the behaviour
  or the user value it delivers, rather than any implementation detail.

    eg: `"changed SomeFunc in file.go to also handle clusters in new region X"`
    is not useful for folks stopping by to quickly see what new stuff they can do with
    `eksctl`, save that for commit messages or the PR description.
    The title `"Add support for region X"` conveys the intent clearly.

1. **There are new tests for new code.**

1. **There are new tests for old code!** This will earn you the title of "Best Loved
  and Respected Contributor" :boom: :sunglasses: .

1. **There are well-written commit messages** ([see below](#commit-message-formatting))
  which will make future debugging fun.


In general, we will merge a PR once a maintainer has reviewed and approved it.
Trivial changes (e.g., corrections to spelling) may get waved through.
For substantial changes, more people may become involved, and you might get asked to resubmit the PR or divide the
changes into more than one PR.

### Commit message formatting

_For more on how to write great commit messages, and why you should, check out
[this excellent blog post](https://cbea.ms/git-commit/)._

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why.

The subject line should feature the _what_ and
the body of the commit should describe the _why_.

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

---

# How the Maintainers process contributions

## Prioritizing issues

The core team regularly processes incoming issues. There may be some delay over holiday periods.

Every issue will be assigned a `priority/<x>` label. The levels of priorities are:

- [`critical`](https://github.com/weaveworks/eksctl/labels/priority%2Fcritical): These are time sensitive issues which should be picked up immediately.
  If an issue labelled `critical` is not assigned or being actively worked on,
  someone is expected to drop what they're doing immediately to work on it.
  This usually means the core team, but community members are welcome to claim
  issues at any priority level if they get there first. _However, given the pressing
  timeframe, should a non-core contributor request to be assigned to a `critical` issue,
  they will be paired with a core team-member to manage the tracking, communication and release of any fix
  as well as to assume responsibility of all progess._

- [`important-soon`](https://github.com/weaveworks/eksctl/labels/priority%2Fimportant-soon): Must be assigned as soon as capacity becomes available.
  Ideally something should be delivered in time for the next release.

- [`important-longterm`](https://github.com/weaveworks/eksctl/labels/priority%2Fimportant-longterm): Important over the long term, but may not be currently
  staffed and/or may require multiple releases to complete.

- [`backlog`](https://github.com/weaveworks/eksctl/labels/priority%2Fbacklog): There appears to be general agreement that this would be good to have,
  but we may not have anyone available to work on it right now or in the immediate future.
  PRs are still very welcome, although it might take a while to get them reviewed if
  reviewers are fully occupied with higher priority issues, for example immediately before a release.

- [`needs-investigation`](https://github.com/weaveworks/eksctl/labels/needs-investigation):  There is currently insufficient information to either categorize properly,
  or to understand and implement a solution. This could be because the issue opener did
  not provide enough relevant information, or because more in-depth research is required
  before work can begin.

These priority categories have been inspired by [the Kubernetes contributing guide](https://github.com/kubernetes/community/blob/master/contributors/guide/issue-triage.md).

## Reviewing PRs

The core team aims to clear the PR queue as quickly as possible. Community members
should also feel free to keep an eye on things and provide their own thoughts and expertise.

---

# Proposals

For chunky features which require Serious Thought™, the first step is the submission
of a design proposal to the [docs](docs/) folder through the standard PR process.

### Process

A template can be found in [`docs/proposal-000-template.md`](docs/proposal-000-template.md). Simply create a copy
of the file (replacing the number with the next in the sequence, and 'template'
with your feature name), and fill in the required fields. When ready, open a PR.

For the initial PR, we can try to avoid getting hung up on specific details
and instead aim to get the motivation and the goals/non-goals sections of the
proposal clarified and merged quickly.
The best way to do this is to just start with the high-level sections and
fill out details incrementally in subsequent PRs.

Initial bare-bones merging does not mean that the proposal is approved, the Status
section will convey whether the proposal has been accepted and work has begun.
Any proposals not marked as 'approved' is a working document and subject to change.

When editing proposals, aim for tightly-scoped, single-topic PRs to keep discussions
focused. If you disagree with what is already in a document, open a new PR
with suggested changes.

A proposal which has been accepted should become a living document. Even the best-laid
plans rarely work out, so as things are learned during implementation the doc
should be updated to accurately reflect the state of the world.

### Sections

Each proposal/design doc should cover the following _at a minimum_:

- **The author(s).** So that people know where to direct questions at any point in time.

- **Status.** Proposals serve as documentation on the design of our codebase.
  It is useful to indicate whether what is documented reflects the state of the project.

- **Summary** A TLDR of everything that is discussed in the doc.This should be
  written in a way that anyone can come by and quickly understand what the proposal is for.

- **Motivation** for why we should do this, how it fits into
  the project's goals and how it will help users in the long term. What is the
  problem this proposal aims to solve? (Includes clear-cut Goals and Non-Goals,
  and all possible Context.)

- **Proposal** of the solution to the above problem. This can be high-level, detail
  comes later. (This includes User Stories as well as Risks and Mitigations.)

- **Design details.** How the proposal should be implemented. (Includes Test Plans,
  and Graduation Criteria.)

- **Alternatives.** List the pros and cons of each solution considered,
  illustrating why the final one was chosen.

- **Known unknowns or open questions.** The author should list any questions for things they are unsure about
  or to direct reviewers to particular areas where their expertise is needed.

### Writing tips

1. Write simply and keep your language accessible. The easier it is to understand,
  the more input you will get from a wider range of sources. Bear in mind that your target audience
  is anyone who comes into contact with eksctl: maintainers, contributors and end users.

1. Refer to the [Kubernetes documentation style guide](https://github.com/kubernetes/community/blob/master/contributors/guide/style-guide.md).

1. Don't assume too much knowledge. Make sure terms have explanations or links for
  the same reason as the last point.

1. Use lots of yaml and diagrams (colour-coded, if possible :wink:).

---

# :rocket: :tada: Thanks for reading! :tada: :rocket:
