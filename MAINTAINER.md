# Maintainer Guide

This is the ["contribution doc"](CONTRIBUTING.md) for people who are also
maintainers. It is part handbook, part responsibility guide.

All new maintainers should go through this as part of onboarding and [PR](CONTRIBUTING.md/#submitting-prs)
any updates.

For the part of the guide which is secret (credentials and whatnot) go [here](#TODO create private issue with link to AWS only handbook doc).

# Contents

- [Responsibilities](#responsibilites)
- [On call](#on-call)
- [How to...](#how-to)
  - [Triage issues](#triage-issues)
  - [Plan for a sprint](#plan-for-a-sprint)
  - [Create a release](#create-a-release)
  - [Release a snap](#release-a-snap)
- [Updating maintainer guides](#updating-maintainer-guides)

---

# Responsibilities

Maintainers are responsible for:
- Enabling and promoting community values
- Building the eksctl community; encouraging new people to get involved and current
members to stay
- Engaging with end users in slack and through github
- Promoting transparent communication and development
- Weighing the needs of the individual against the needs of the community
- Seeking input and consensus on project decisions
- Making final decisions, ensuring they are documented and followed through

For more on how to maintain an open source project with a community, see Github's
[open source guides](https://opensource.guide/).

---

# "On call"

AKA Community Service :construction_worker_woman: .

While maintainers should keep up with their community and project responsibilities
at all times, there is a weekly rota which designates a single person who is in
charge of handling new issues, as well as day-to-day repo chores.

Duties are as follows:

- Responding to questions in [`#eksctl`](https://weave-community.slack.com/messages/eksctl/).
  - If you don't know the answer, encourage the asker to open a [help issue](CONTRIBUTING.md#help-request-guide).
    This is a useful way to build up a troubleshooting collection, since we have no history on slack.
    Also many times requests for help turn into bugs so :shrug: .

- Keeping an eye on and triaging new issues. See [How to triage new issues](#triage-issues) below.

- Triage integration test failures. The integration tests run every night, and the result
  is sent to the team's private slack channel. If there is a failure overnight, you need
  to determine whether the failure was genuine or a flake.

  In the case of:

    a) a genuine failure (aka a recently introduced bug), try to see whether
      it can be quickly fixed. If yes, submit a PR. If it can't, create a
      [bug report](CONTRIBUTING.md#bug-report-guide) and triage it.

    b) a flake, create a [bug report](CONTRIBUTING.md#bug-report-guide), note
      that it is a flake, and triage as normal.

- Checking documentation is up to date. :tired_face:

- Preparing the next batch of work for the core team. See [How to prep the backlog](#plan-for-a-sprint) below.

- Releasing `eksctl`. See [How to create a release](#create-a-release) below.


---

# How to... :thinking:

## Triage issues

Triaging issues will get easier with more time spent in the project. Triaging means
sorting out which things need to be handled immediately from those which can wait a little
(with a few more levels in-between). At first all issues will seem equally important
or unimportant, but soon it will become easy enough to do by gut instinct.

For each new issue which comes in, do the following:

1. Check for a duplicate. I know we ask people to check before opening, but it often
  doesn't happen. If you find something already exists, link it in a comment and close
  the new issue.

1. Check that they have chosen the right category. When creating new issues, users
  can choose between `bug`, `feature` and `help` and the correct label is added to the
  issue. Most times they get it right, but sometimes a `bug` is confused with a `help` etc.
  Reset the labels to what they should be.

1. Have they provided enough information? We have pretty thorough [guides](CONTRIBUTING.md#opening-issues)
  which ask for all sorts of things. If the detail is so sparse, that **most people** (note that I didn't say "you")
  couldn't get a clear picture of what is happening, stick an `awaiting more information`
  label on it and ask the OP to give a little more detail.

1. If they have provided tons of information and yet you still have no idea what could
  be going wrong, add a `needs-investigation` label.

1. Have they formatted things nicely? Reading a whole wall of terminal text is hard-going,
  and some folks either don't know how to use code blocks or just don't bother.
  Edit the description to make it readable and help out whoever picks it up.

1. Add a `priority/<x>` label. Once you can read it and understand what is happening,
  decide when the work should be scheduled. We have the following categories:

    - [`critical`](https://github.com/weaveworks/eksctl/labels/priority%2Fcritical):
      Examples include user-visible bugs in core features, broken builds
      or tests and critical security issues. If one of these appears, flag it in the
      team channel, add it to the backlog and get someone on it immediately. If there are
      no low effort user workarounds a fix needs to go out asap, otherwise it can go
      out as part of the next release.

    - [`important-soon`](https://github.com/weaveworks/eksctl/labels/priority%2Fimportant-soon):
      Examples include bugs which cause a fair amount of user pain, super useful improvements
      desired by all the community and test flakes which routinely hold up development.
      Ideally these things should be handled fairly soon and
      go out in the next release or not too long afterwards.

    - [`important-longterm`](https://github.com/weaveworks/eksctl/labels/priority%2Fimportant-longterm):
      Examples include larger features which need time to design, things which contribute
      to the goals/wellbeing of the project but are not massively urgent, infrequent test flakes
      and other useful things which we really want but can live without for the moment.
      These things should be implemented and released within a couple of months
      of the request.

    - [`backlog`](https://github.com/weaveworks/eksctl/labels/priority%2Fbacklog).
      Things which would be nice to have but will survive without.

1. Consider the "bar of entry". Is the task really simple, with tons of detail, minimal scope
  and maybe even some hints on how to implement? Assign a `good-first-issue` label.
  Does it have potential to be a `good-first-issue`, but a little more context is needed?
  Either provide that yourself or ask the OP to elaborate.

    These issues should be accessible to folks new to the repo, as well as to open source in general.
    They should present a low/non-existent barrier to entry with a thorough description,
    easy-to-follow reproduction (if relevant) and enough context for anyone to pick up.
    The objective should be clear, possibly with a suggested solution or some pseudocode.
    If anything similar has been done, that work should be linked.

    If the bar is a little higher, add a `help wanted` label.

1. Add the issue to the `eksctl-backlog` if it looks like a good candidate
  for the upcoming sprint.

1. **Respond to the issue**! One of the biggest reasons that people lose interest
  in open-source projects is lack of interaction. It also does not help the project's
  reputation when people feel ignored or do not feel welcome.
  Even if you have absolutely nothing useful to say, there should always
  be _at least_ one comment on every issue (that was not written by the OP).

    Here is some boilerplate if you are lost for words:

    - "Hi @\<op-name\>! Thanks for reporting this bug :+1: If you have a fix in mind,
      we are happy to accept PRs, otherwise the core team will get to it soon."

    - "Hi @\<op-name\>! Thanks for opening this feature request :+1: The team will review
      and see whether this functionality fits in with the project. In the meantime,
      if you would like to demonstrate a solution, please open a PR :smiley: ."

    - "Hi @\<op-name\>! Thanks for creating a help ticket. I will do some research and get
      back to you asap!"


**Notes:**

Try not to forget about issues. If you don't know what they are asking for,
don't understand the bug, or don't even know if it is a thing we should care about,
tag in the rest of the team on slack after you have written your boilerplate "Thanks for opening"
message.


## Plan for a sprint

The core team works in weekly sprints and gathers every Wednesday to decide
what their focus will be for the coming week.

Whoever is [on call](#on-call) both decides what the core team will be working on and runs the
planning meeting. _The power..._ :astonished: :star_struck: .

### Pre-planning

This involves moving tickets around in the `eksctl-backlog` and ensuring there
is a prioritized set in `Ready for work`.

Depending on the complexity of tickets, we want to end up with from 5-10 issues
teed up in `Ready for work`.

- Check in with AWS to see if there are any new features coming out, check the [container
  roadmap](https://github.com/aws/containers-roadmap) for new things coming out and create mirror issues in eksctl.
- Check the `Backlog` column in `eksctl-backlog` to see if anyone has added anything
  good. Move tickets you want into `Ready for work`.
- Check the issues. `priority/critical` issues should already be in the `Sprint In Progess`
  column (or at least in `Ready for work`), but if there are some which haven't been scheduled
  move them into ready.
- Go through other `priority/<x>` issues and move some which look ready for
  work (ie. have sufficient context) into that column.
- Ping the core team to ask if they have anything they would like to prioritize.
- Prioritize the issues into the order you think they should be tackled.

### Running the meeting

Go through each ticket in `Ready for work`, including the ones which are still there
from last sprint. For each one:
- Remind everyone of the context
- Explain why you prioritized it for this week
- Explain the user value it would deliver (if applicable)
- Ask if the scope is clear, if it is ready for work, or if folks feel they need more
- Give space for discussion and questions

## Create a release

1. Check out to master, make sure the branch is up to date.
1. Ensure [integration tests](CONTRIBUTING.md#running-the-integration-tests) pass (est: 3 hours).
1. Make a PR to change the API version (`v1alpha5`) if necessary.
1. Determine the next release tag, e.g.:

   - for a release candidate, `0.13.0-rc.0`, or
   - for a release, `0.13.0`.

1. Create a `docs/release_notes/<version>.md` release notes file for the given tag using the contents of the release
   draft (generated by [`release-drafter`](https://github.com/release-drafter/release-drafter)), e.g.:

    ```console
    touch docs/release_notes/0.13.0.md

    # copy-paste in the contents of the release draft.
    ```
1.
    a) For the first release candidate (`rc.0`) create a new branch after the major and minor numbers of the release (`release-X.Y`):

     ```console
     git checkout master
     git pull --ff-only origin master
     git checkout -b release-0.13
     ```

    b) If this is a subsequent release candidate or the release after an RC check out the existing release branch

     ```console
     git checkout release-0.13
     ```

1. Create the tags so that circleci can start the release process:

   - for a release candidate: `make prepare-release-candidate`. If there was an existing RC it will increase it's number, e.g.: rc.0 -> rc.1

   - for a release: `make prepare-release`. Regardless of whether there was a previous RC or a development version it will create a normal release

1. Ensure release jobs succeeded in [CircleCI](https://circleci.com/gh/weaveworks/eksctl).
1. Ensure the release was successfully [published in Github](https://github.com/weaveworks/eksctl/releases).
1. Download the binary just released, verify its checksum, and perform any relevant manual testing.
1. Merge in the branch you created with the release notes.

### Notes on Automation

When you run `make prepare-release` it will push a commit to master and a tag, which will trigger [release workflow](https://github.com/weaveworks/eksctl/blob/38364943776230bcc9ad57a9f8a423c7ec3fb7fe/.circleci/config.yml#L28-L42) in Circle CI. This runs `make eksctl-image` followed by `make release`. Most of the logic is defined in [`do-release.sh`](https://github.com/weaveworks/eksctl/blob/master/build/scripts/do-release.sh).

You want to keep an eye on Circle CI for the progress of the release ([0.3.1 example logs](https://circleci.com/workflow-run/02d8b5fb-bc7f-404c-9051-68307c124649)). It normally takes around 30 minutes.

### Latest release

To get the latest release you can use the link <https://github.com/weaveworks/eksctl/releases/latest>.

**Note** Previously, eksctl used a floating tag called `latest_release`. This was deprecated after release `0.14.0`.

## Release a snap

Snaps are software packages which run on a variety of Linux flavours.

*Note:* This is in somewhat of a [weird state](https://github.com/weaveworks/eksctl/issues/215).

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

---

# Updating maintainer guides

Where possible, add all new tips and tricks to this doc. Add things to the secret handbook very
sparingly.
