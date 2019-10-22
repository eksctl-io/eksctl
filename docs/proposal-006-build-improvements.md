# Design Proposal #006: Build Improvements

## Current State

At the start of the project we've picked CircleCI as it appears to be a default choice at Weaveworks.

The original solution involved:
- running `docker build` to create a build image
- use build image to run tests produce final image

Publishing build and final images hasn't been automated, but some users requested it.

When we switched [to modules](https://github.com/weaveworks/eksctl/pull/917), we have increased the
compexity of how build image worked. It's hard for a new person to understand it, e.g. for a new
contributor.

We currently build images in CI and run tests as part of that, we should only run tests as a primary
task.

Downside of current setup:

- CircleCI requires delegation of access to GitHub
  - for releases we use a token
  - we don't have a way to automate workflows that require pushing to the repo
- CircleCI cache 
  - it is not very fast, it turns out pulling image with cached dependencies is as fast
  - it is very specific to CircleCI, not portable
  - we do not have a way to control it
- the task of building images should not be mixed with the task of running test
- images are not being automatically pushed to a registry
- update of build image is manual

## Proposed Improvements

### Summary

- streamline the CircleCI setup
- enable GitHub Actions as an experiment

Present approach:

- restore cache from previous build
- mount cache as docker volume
- run tests and build new binary

Proposed improvement:

- pull docker image
- run tests and build new binary

So, `docker pull` for a large image that contains all cached objects takes as long (or less time) then
using CircleCI cache feature. There are other benefits to this, as we can tell exactly what's being cached
and can be reproduced any other CI environment.

### Details

This proposal [#1200](https://github.com/weaveworks/eksctl/pull/1200) aims to streamline the following
aspects:

- automate build image versioning through deterministic git object hashes (see `Makefile.docker`)
- use build image as cache
- initial cleanup of build scripts
    - separate makefiles
    - single static `Dockerfile`
- stop using CircleCI cache and pull build image instead
- switch to Docker executor using our build image in CircleCI
    - stop having to download Go toolchain
    - clear separation of concerns - creating build image vs running tests
- enable GitHub Actions as an experiment, so we can evaluate it against CircleCI
- based on initial tests (in GitHub Actions) - 40%-50% speed-up in on-commit execution (lint+test+build)

### Timing

I have managed to bypass the need for upgrading CircleCI setup (at least for the time being). It's possible upgrading to native Docker executor will provide a speed-up, we might see shorter build queue times and it's possible that layer caching that they advertise will be beneficial, but they don't provide enough details for me to judge if it will. I'm going to evaluate the pricing and discuss it internally.

Current timings on master branch (prior to this change) are between 10 & 11 minutes in total, e.g.:

- restore cache - 02:00
- make image - 04:18
- save cache - 03:20

> NOTE: With current approach `make test` relies on test cache, which means tests
> are a little faster, however that prevents flakes from surfacing as not all tests
> run every time. _Proposed approach doesn't rely on test cache._

New timings on CircleCI are:

- image pull - 01:17
- make test - 06:11
- total - 07:44

With `make build` added, here are two samples:

- image pull - 01:19 / 01:20
- make test - 06:32 / 06:22
- make build - 02:08 / 02:05
- total - 10:16 / 10:04

With `make build` and `make test` combined in one step we save a couple of minutes:

- image pull - 01:19 / 01:17
- make test & build - 06:30 / 06:39
- total - 08:13 / 08:12

Overall, there is a definite 2m speed-up.

## Simplifications

We currently have rather complicated [CircleCI config](https://github.com/weaveworks/eksctl/blob/950f9bc695234107725234b9e1a9c9d2ee54e51f/.circleci/config.yml#L3-L38),
and [`Makefile`](https://github.com/weaveworks/eksctl/blob/950f9bc695234107725234b9e1a9c9d2ee54e51f/Makefile#L188-L204).


New configuration is much simpler, in [CircleCI config](https://github.com/weaveworks/eksctl/blob/d3b6988562b14c9d91f4e1bf7dd2c086e06c2383/.circleci/config.yml#L3-L24)
it comes down to running [`make test`](https://github.com/weaveworks/eksctl/blob/1be26b314333467d1b67a44c77d1cd27460eaa70/Makefile#L68-L73) followed by `make build`.

## Further Improvements

- stop using `docker run` and `docker commit`, use another `Dockerfile` instead of `eksctl-image-builder.sh`
- push images to a registry (GitHub offering could be a good fit, otherwise consider ECR)
- automate other workflows with bots and GitHub Actions (e.g. cherry-picking [#1284](https://github.com/weaveworks/eksctl/issues/1284), AMI updates [#314](https://github.com/weaveworks/eksctl/issues/314))
- automate integration tests
- allow running integration tests on PRs from contributors upon PR approval
