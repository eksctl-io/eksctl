---
hide:
  - navigation
---
# eksctl Community :fontawesome-solid-people-group:

We welcome contributions from the community. :octicons-heart-fill-24:{ .heart }
For more information, please head to our [Community][community] and [Contributing][contributing] docs in Github.

[community]: https://github.com/eksctl-io/eksctl/blob/main/COMMUNITY.md
[contributing]: https://github.com/eksctl-io/eksctl/blob/main/CONTRIBUTING.md

## Get in touch :simple-wechat:

[Create an issue](https://github.com/eksctl-io/eksctl/issues/new), or login to [Eksctl Slack (#eksctl)][slackchan] ([signup][slackjoin]).

[slackjoin]: https://slack.k8s.io/
[slackchan]: https://slack.k8s.io/messages/eksctl/

## Release Cadence :material-clipboard-check-multiple-outline:

Minor releases of `eksctl` are loosely scheduled for weekly on Fridays. Patch
releases will be made available as needed.

One or more release candidate(s) (RC) builds will be made available prior to
each minor release. RC builds are intended only for testing purposes.

## Eksctl Roadmap :octicons-project-roadmap-16:

The EKS section of the AWS Containers Roadmap contains the overall roadmap for EKS. All the upcoming features for `eksctl` built in partnership with AWS can be found [here](https://github.com/aws/containers-roadmap/projects/1?card_filter_query=label%3Aeks).

## 2021 Roadmap

The following are the features/epics we will focus on and hope to ship this year.
We will take their completion as a marker for graduation to v1.
General maintenance of `eksctl` is still implied alongside this work,
but all subsequent features which are suggested during the year will be weighed
in relation to the core targets.

Progress on the roadmap can be tracked [here](https://github.com/eksctl-io/eksctl/projects/2).

### Technical Debt

Not a feature, but a vital pre-requisite to making actual feature work straightforward.

Key aims within this goal include, but are not limited to:

- [Refactoring/simplifying the Provider](https://github.com/eksctl-io/eksctl/issues/2931)
- [Expose core `eksctl` workflows through a library/SDK](https://github.com/eksctl-io/eksctl/issues/813)
- Greater integration test coverage and resilience
- Greater unit test coverage (this will either be dependent on, or help drive out,
  better internal interface boundaries)

### Declarative configuration and cluster reconciliation

This has been on the TODO list for quite a while, and we are very excited to bring
it into scope for 2021

Current interaction with `eksctl` is imperative, we hope to add support for declarative
configuration and cluster reconciliation via a new `eksctl apply -f config.yaml`
command.  This model will additionally allow users to manage a cluster via a git repo.

A [WIP proposal](https://github.com/eksctl-io/eksctl/blob/main/docs/proposal-007-apply.md)
is already under consideration, to participate in the development of this feature
please refer to the [tracking issue](https://github.com/eksctl-io/eksctl/issues/2774)
and our [proposal contributing guide](https://github.com/eksctl-io/eksctl/blob/main/CONTRIBUTING.md#proposals).

### Flux v2 integration (GitOps Toolkit)

In 2019 `eksctl` gave users a way to easily create a Gitops-ready ([Flux v1](https://docs.fluxcd.io/en/1.21.1/))
cluster and to declare a set of pre-installed applications Quickstart profiles which can be managed via a git repo.

Since then, the practice of GitOps has matured, therefore `eksctl`'s support of
GitOps has changed to keep up with current standards. From version 0.76.0 Flux v1 support was removed after an 11
month deprecation period. In its place support for [Flux v2](https://fluxcd.io/) can be used via
`eksctl enable flux`


