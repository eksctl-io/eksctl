# `eksctl apply`

This proposal addresses a GitOps-compatible subcommand of `eksctl` called
`apply`.

**Status**: WIP

* [Goals](#goals)
* [Non-goals](#non-goals)
* [Motivation](#motivation)
   * [Questions for discussion](#questions-for-discussion)
      * [Deletion of resources](#deletion-of-resources)
      * [Immutability](#immutability)
      * [Gradual support and config](#gradual-support-and-config)
* [Proposal](#proposal)
   * [Experimental flag](#experimental-flag)
   * [Ownership](#ownership)
      * [eksctl-managed resources](#eksctl-managed-resources)
         * [Non-authoritative resources](#non-authoritative-resources)
      * [Non-taggable resources](#non-taggable-resources)
         * [eksctl owns everything](#eksctl-owns-everything)
         * [terraform-style state](#terraform-style-state)
      * [Ignore resource option](#ignore-resource-option)
      * [non-eksctl-managed resources](#non-eksctl-managed-resources)
   * [New config](#new-config)
      * [Goals/justification for config refactors](#goalsjustification-for-config-refactors)
* [Open questions](#open-questions)

## Goals

The purpose of `apply` is to subsume the various imperative commands
`eksctl` provides now into one `apply` command that reconciles the actual state
of an EKS cluster with the intended state as specified in the `eksctl` config file.
It should be GitOpsable aka GitOps-compatible, i.e. usuable as part of a GitOps
pipeline.

## Non-goals

This proposal is only concerned with the `apply` subcommand. Existing behavior
and commands are unaffected. Continuing to support them shouldn't complicate
the implementation too much as they are more or less subsets of
operations performed by `apply`.

## Motivation

`eksctl` supplies users with imperative commands for clusters and cluster
resources that allow them to take explicit actions, like `create` or `delete`.
We also encourage the use of a config file in which the desired state of the
cluster is described. Still, users are required to figure out which imperative
steps to take to reconcile the cluster with the desired state.
The missing capability is to give `eksctl` the ability to do this.
Conceptually, `eksctl` would gather all of the information it can about the
current state of the cluster in order to build a `ClusterConfig` describing
the real-world cluster and then `diff` this config with the user-provided
`ClusterConfig`. The `diff` is used to setup a plan of modification to the
cluster and resources.

### Questions for discussion

This proposal also opens the floor for discussion around non-straightforward
questions around reconciliation.

#### Deletion of resources

When looking at a list of nodegroups in a cluster, it's not clear how to handle
deletion.

-   Is the given list intended to be a complete list of nodegroups and
    no others should exist? Are they automatically deleted?
-   How do we handle non-eksctl managed resources?

#### Immutability

In order to change some properties, entire resource will have to be recreated.
This should be explicit to the user.

#### Gradual support and config

Because it's impractical to introduce complete support for reconciliation at
once, we need to instead gradually support properties and resources of a
cluster.

This leads to the question of whether the current `ClusterConfig` is a good
structure for a config meant for reconciliation.

## Proposal

Introduce a subcommand `apply` which has only a config file argument.
`apply` would reconcile desired with real-world state
[as described in the motivation](#motivation).
The exact details per resource will potentially differ depending on
the resources and won't be enumerated here.

1. `apply` is put behind the explicit "experimental" flag (cf `enable profile/repo`)
2. Take full ownership of `eksctl`-managed resources (i.e. delete missing resources by default)
3. Options for non-`eksctl`-managed resources as well as resources where
   ownership is unknowable are covered below
4. Changing immutable fields through recreation is the default
5. We gradually expand `apply` support to different parts of the cluster
6. As we expand `apply` support, we reevaluate and update the config structure as
   necessary in a new API version.

Some further discussion/justification follows.

### Experimental flag

We make it explicit that `apply` is experimental, very open to feedback and
subject to change.

### Ownership

Ownership comes into play most importantly with deletion.

#### `eksctl`-managed resources

Given the config:

```
nodeGroups:
  - name: ng-2
    desiredCapacity: 1
```

we discover an existing nodegroup `ng-1` with the `eksctl.io/owned: true` tag (for example),
do we delete it?

The proposed answer is yes, it should be deleted. With `apply` we assume that all
previous changes initiated by `eksctl` were initiated through the same config
file/reconciliation process and that this nodegroup must have been
deleted from this file with the intention of deleting it from the cluster.

##### Non-authoritative resources

There is precedent for adding the ability to specify resources in a
non-authoritative list. For example, the [`google` terraform provider supports IAM with 3 different levels of authoritativeness](
https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_iam#google_storage_bucket_iam_policy#iam-policy-for-cloud-storage-bucket).

Initial `apply` support will focus on authoritative resource specification.

#### Non-taggable resources

Some resources may not be taggable as owned by `eksctl`.
If there weren't a `tags` field on addons, for example. `eksctl`-created
addons would then be indistinguishable from addons created in the AWS console.
There are at least 2 ways to deal with this case:

##### `eksctl` owns everything

With this option, we assume that `eksctl` owns the set of such resources and
delete any not listed in the config.

##### `terraform`-style state

Storing some kind of state/config "somewhere" is another option and has been discussed in
https://github.com/eksctl-io/eksctl/issues/642

#### Ignore resource option

An `ignore` flag (cf `terraform`'s `prevent_destroy`) in the schema for all objects
could be added to tell `eksctl` about such objects and prevent deletions
and potentially skip reconciliation.

The user would prefix their list:

```
nodeGroups:
  - name: ng-1
    ignore: true
  - name: ng-2
    desiredCapacity: 1
```

#### non-`eksctl`-managed resources

The behavior for resources created outside of `eksctl` could also
be toggleable by settings in the config/flags on `apply`.

### New config

Proposed is creating a new API version, starting with `v1alpha6`, where, as we add
`apply` support for cluster/nodegroup properties, we allow for aggressive reevaluation
of the existing config structure that would hold these newly supported properties/aspects
of the cluster/nodegroup.

If we decide, _for example_, to add support for reconciling tags, we might
decide it fits better as a top level key. We would then, with `v1alpha6`, remove
`tags` under `metadata` and make it a top level field. It should be possible to
automatically rewrite the previous version into the current version.

#### Goals/justification for config refactors

-   Ensure that the config makes sense in a world centered around reconciliation.
-   Auto upgrade makes it easier for users to use an existing config
-   **"Parse, don't validate"**: the code that uses the config should be insulated from how we
    serialize/deserialize the config.

## Open questions

-   What, if anything, about the VPC is owned by `eksctl` and under what
    circumstances?
    -   What does it mean if
        the user specifies a VPC not created by `eksctl` but `eksctl` finds a
        discrepancy between the actual state and the config
    -   `terraform`'s answer here is explicitness: `resource` vs `data`.
