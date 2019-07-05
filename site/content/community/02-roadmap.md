---
title: "Project Roadmap"
weight: 20
url: community/roadmap
---

## Release Cadence

Starting with 0.2.0 onwards, minor releases of `eksctl` should be expected every two weeks and patch releases will be made available as needed.

One or more release candidate(s) (RC) builds will be made available prior to each minor release. RC builds are intended only for testing purposes.


## Project Roadmap Themes

### Cluster quickstart with GitOps

It should be easy to create a cluster with various applications pre-installed, e.g. Weave Flux, Helm 2 (Tiller), ALB Ingress controller. It should also be easy to manage these applications in a declarative way using config files in a git repo (with [GitOps](https://www.weave.works/blog/what-is-gitops-really)).

### Declarative configuration management for clusters

One should be able to make EKS cluster configuration through declarative config files (`eksctl apply`). Additionally, they should be able to manage a cluster via a git repo.

### Cluster addons

Understanding how the [add-ons spec by SIG Cluster Lifecycle](https://github.com/kubernetes/enhancements/pull/746) will evolve and how we can implement management of some cluster applications via the addon spec.
