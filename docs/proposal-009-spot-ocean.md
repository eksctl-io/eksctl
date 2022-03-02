# Spot ocean integration

## Authors

Spot By NetApp Ocean (@spotinst/sig-developers)

## Status

In process.

## Table of Contents
<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
    - [Goals](#goals)
    - [Non-Goals](#non-goals)
    - [Linked Docs](#linked-docs)
- [Proposal](#proposal)
- [Design Details](#design-details)
    - [Test Plan](#test-plan)
- [Alternatives](#alternatives)
<!-- /toc -->

## Summary

We implemented Spot Ocean structures that are based on the eksctl Cluster and NodeGroup structures from release `0.146.0`. This implementation
allows spot-ocean users to utilize eksctl in various ways on their clusters and node groups.
We note that no dependencies exist between the spot-ocean and eksctl structures that could create problematic issues in the future.

The value in integrating Spot Ocean with `eksctl` is simply to bring existing and future AWS customers a way of:

a) Creating new clusters and/or node groups with spot ocean integration using a
single command.

b) Modifying clusters and/or node groups with spot ocean integration using a
single command.

Spot by Netapp pledges to fully maintain this integration.
This includes:
- Monthly updates with new features
- Code reviews and feature assessment from the direct EKSCTL community
- Feature parity with our direct API and UI enabling EKSCTL all the latest features
- Spot by Netapp fully managing Support and maintenance of this integration
  - Bug fixes directly from the EKSCTL community
  - Urgent 24/7 support available on our platform
  - Ensuring full compatibility with the newest versions of Kubernetes and EKS

## Motivation

The overall motivation of this proposal is to solve 2 problems:

- There are many AWS customers with eks clusters, with a demand for spot ocean integration.
- AWS Customers want to integrate their eks clusters and nodegroups with spot ocean via eksctl's configuration.

### Goals

- Enable AWS users to create spot ocean clusters and nodegroups using eksctl.
- Enable AWS users to modify their spot ocean cluster configs and nodegroups using eksctl.
- Enable AWS users to perform utility actions on their spot ocean clusters and nodegroups using eksctl.

### Linked Docs

[Original PR](https://github.com/weaveworks/eksctl/pull/6731).
[Spot Ocean docs](../userdocs/src/usage/spot).
[Expansion issue](https://github.com/weaveworks/eksctl/issues/6694).

## Proposal

This design proposes adding a new field `spotOcean` to both cluster and nodegroup level,
and creates cluster with spot ocean managed nodegroups.

for example:

```bash
eksctl create cluster \
 --name example \
 --spot-ocean
 --managed=false
```

will result in a new spot ocean cluster.
In addition, the design proposes 2 new utils options `update-spot-ocean-cluster` and `update-spot-ocean-credentials`.

for example:
```bash
eksctl utils update-spot-ocean-cluster -v 4 -f ./cluster.yaml
```

while the `cluster.yaml` contains the new updated cluster definition.

## Design Details

The new arg option `--spot-ocean` will be added to `eksctl create cluster` and `eksctl create nodegroup`. That option will also be supported in the ClusterConfig file for self-managed nodegroups.
In addition, we have added 2 new options for utils actions of eksctl, `update-spot-ocean-cluster` and `update-spot-ocean-credentials`, both require a configuration file, mainly meant for update action regarding the cluster.
- For more details feel free to browse our [spot ocean guides](../userdocs/src/usage/spot/ocean/spot-ocean-cluster.md)

### Test Plan

Following maintenance or the release of a new feature, we check the following:

- Running all the existing unit tests to make sure nothing broke from our changes.
- Creation of new eks clusters on the various up to date k8s versions.
- Creation and modification of nodegroups inside said clusters.
- Verification of utility actions concerning ocean cluster management within eksctl.

## Alternatives

The current alternative is use of our own branch forked from the main eksctl branch [repo](https://github.com/spotinst/weaveworks-eksctl/releases/tag/v0.146.0) for customer purposes.
