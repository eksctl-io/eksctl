# `eksctl apply`

This proposal addresses a GitOps-compatible subcommand of `eksctl` called
_tentatively_ `apply`. Its purpose is to subsume the various imperative commands
`eksctl` provides now into one `apply` command that reconciles the actual state
of an EKS cluster with the intended state as specified in the `eksctl` config file.

This proposal is only concerned with the `apply` subcommand. Existing behavior
and commands are unaffected. Continuing to support them shouldn't complicate
the implementation too much as they are more or less subsets of
operations performed by `apply`.

## Questions to be answered

-   Does the config represent the entire state of the config vs a partial
    spec?
    -   e.g. if I `apply` with a list of 1 nodegroups, are any other existing
        nodegroups deleted?
    -   do we also delete non-eksctl managed/created resources? i.e. does `eksctl`
        own the cluster by default? (potential for configurability here)
-   Surfacing which fields are immutable (i.e. only settable during
    cluster/nodegroup creation)
    -   One possibility might be structuring the config in such a way that
        immutable fields are separated out
-   How to introduce gradual, partial support?
-   How do we handle the config
    1. Leave config alone
    2. Introduce a _new_ type which only specifies
       those parts of the config currently supported by `apply`
       and put it under a new top level field
    3. Restructure the config gradually

## Proposal

1. `apply` is put behind the explicit "experimental" flag (cf `enable profile/repo`)
2. We take full ownership of `eksctl`-managed resources (i.e. delete missing resources by default)
3. Options for non-`eksctl`-managed resources as well as resources where
   ownership is unknowable are covered below
4. Immutability is surfaced by erroring when discrepancies are detected
5. We gradually expand `apply` support to different parts of the cluster
6. As we expand `apply` support, we reevaluate the config structure

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
deleted from this file.

#### Non-taggable resources

Some resources may not be taggable as owned by `eksctl`.
Imagine there weren't a `tags` field on addons, for example. `eksctl`-created
addons would then be indistinguishable from addons created in the AWS console.
There are at least 2 ways to deal with this case:

##### `eksctl` owns everything

With this option, we assume that `eksctl` owns the set of such resources and
delete any not listed in the config.

##### `terraform`-style state

Storing some kind of state/config "somewhere" is another option and has been discussed in
https://github.com/weaveworks/eksctl/issues/642
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
    -   `terraform`'s answer here is `resource` vs `data`.
