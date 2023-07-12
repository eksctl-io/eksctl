# Flux v2 integration plan: next steps

## Authors

Claudia Beresford (@Callisto13)

## Status

Approved.

## Table of Contents
<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
  - [Linked Docs](#linked-docs)
- [Details](#details)
  - [Removals](#removals)
  - [Git provider](#git-provider)
  - [Args](#args)
  - [&quot;Dry run&quot;](#dry-run)
  - [Documentation improvements](#documentation-improvements)
- [Notes/Constraints/Caveats](#notesconstraintscaveats)
  - [Teams](#teams)
  - [Future work](#future-work)
- [Implementation Stages](#implementation-stages)
  - [Remove all the things](#remove-all-the-things)
  - [Docs](#docs)
  - [Recommended Flux version](#recommended-flux-version)
  - [&quot;Dry run&quot;](#dry-run-1)
- [Test coverage](#test-coverage)
- [API spec](#api-spec)
- [Alternatives](#alternatives)
  - [List of strings](#list-of-strings)
  - [Actual config](#actual-config)
<!-- /toc -->

<!--
install https://github.com/kubernetes-sigs/mdtoc.
tags, and then generate with mdtoc -inplace docs/proposal-008-flux2.md
-->

## Summary

We implemented some very basic Flux 2 integration with release `0.38.0`. This implementation
is marked as experimental, and allows users to set the most common flags for bootstrapping a
cluster with `github` or `gitlab`.

The value in integrating Flux 2 with `eksctl` is simply to give users a way of:

a) setting up their new cluster with flux immediately after it is created, as part of the same
single command

b) letting users save and run their flux bootstrap args via a config file, rather than a list of flags

We implemented the experimental integration with eksctl config values only, and no eksctl flags,
since that would defeat the purpose of value `b`: we intend to keep it this way.


Exposing the rest of the options available to `flux bootstrap` is the next step.
Since Flux development moves quite fast, we want to do this in a way that we don't have
to continuously keep in step; we want to give users a way to provide arbitrary args.
Balancing that we also don't want to lose value `b`; we want to ensure users have a
nice config file UX, and a single string of flags goes against that.
We have had a few issues opened since then from people asking us to expose more options
(which is nice because people are using it), and we don't want to be in a position where
we get an issues like that opened frequently.

The suggestion in this proposal is that we:
- Remove all config options exposed in the initial Flux integration
- Implement a `map[string]string` through which users can pass flags to Flux

For a more rigid way of doing this configuration, see [Alternatives](#alternatives) below.

## Motivation

The overall motivation of this proposal is to solve 2 problems:

- Users want to bootstrap their clusters with flux via eksctl's configuration
- The configuration needs to be flexible enough that users do not have to request
the exposure of new flags/features as Flux develops

The unofficial motivations of this proposal are:

- Sketch out a plan to deliver in stages (multiple PRs) so that Claudia does not open yet another
50 page diff
- Have everything written down before Claudia goes on holiday for 3 weeks

### Goals

- Users can bootstrap to any of the 3 supported git providers (github, gitlab, git)
- Users can benefit from a standard set of configuration options
- Users can provider arbitrary args as they need

### Non-Goals

- Air-gapped bootstrapping
- Terraform bootstrapping

### Linked Docs

[Original PR](https://github.com/eksctl-io/eksctl/pull/3066).
[Current eksctl docs](https://eksctl.io/usage/gitops/#experimental-installing-flux-v2-gitops-toolkit).
[Current Flux api object in eksctl](https://github.com/eksctl-io/eksctl/blob/ab702daa671a87dc926b70348481ce336638f064/pkg/apis/eksctl.io/v1alpha5/types.go#L901-L937).
[Expansion issue](https://github.com/eksctl-io/eksctl/issues/3238).

## Details

### Removals

Pretty much everything which we exposed in the initial implementation will disappear:

```diff
  gitops:
    flux:
      gitProvider: github
-     owner: dr-who
-     repository: my-cluster-gitops
-     personal: true
-     path: "clusters/cluster-27"
-     kubeconfig: /foo/bar
-     branch: main
-     namespace: "flux-system"
```

The addition of an arbitrary list of flags and values means users can immediately use
any Flux feature without waiting for us to expose it.

### Git provider

By "provider" I mean which git thing you want. Flux lets you bootstrap to one of
`github`, `gitlab` and `git` a generic git server.

This is the one explicit bit of config we will leave, simply because it is not a flag,
it is a subcommand and therefore the ordering would matter if we left it to the map.
If is is an explicit config option users do not have to care about this.

```diff
  gitops:
    flux:
      gitProvider: github|gitlab|git
```

### Args

All other Flux bootstrap configuration will be provided via an arbitrary `map[string]string`.
(I have called this key `flags` but open to whatever.)

```diff
  gitops:
    flux:
      gitProvider: github|gitlab|git
+     flags:
+	branch: main
+	namespace: default
+	components: "helm-controller,source-controller"
+	personal: true
```

### "Dry run"

One thing we lose by having these arbitrary fields is the extra validation: are users setting the
correct flags for their chosen git provider? Of course Flux will complain if you get the flags wrong,
and we will be directing users to the Flux docs so that they know what options are available to them,
but I can see people tripping on this and it is really annoying to get to the end of a 25 min cluster
create, only to have your `flux bootstrap` fail.

Ideally we could "validate" the flags are at least applicable to the bootstrap command
before starting the create.

Flux doesn't have a dry-run bootstrap option (AFAIK), so we cannot call the command without
actually running it. But we could hack it by maybe appending
a garbage flag to the end of the list we pass in. Flux will fail on the first flag it
does not recognise, so if it errors on anything other than our nonsense flag, then the
validation is failed.

I am not 100% sold on this, and we don't need it right away, but I think it may be neat
to have something like it.

### Documentation improvements

This should be improved in general and should redirect as much as possible
to the Flux docs. We should be very clear on how much `eksctl` is doing here, and how
users can pass the correct options through to Flux.
The Flux 2 implementation currently shares a page with the Flux v1 eksctl.io docs, they should
be split up. I don't think we are even clear about what Flux commands we run.

## Notes/Constraints/Caveats

### Teams

One important caveat with the `map[string]string` approach is how Flux implements adding teams for github and gitlab
providers. Their current implementation uses `StringArrayVar` which means flags which look like
`--team one --team two` and therefore a yaml which looks like:

```yaml
...
     flags:
       team: one
       team: two
```

This obviously will not work. We intend to ask about changing that flag to a `StringSliceVar`
which would support both `--team one --team two` and `--team one,two` formats (and hence a `team: "one,two"` yaml), but if this is
not possible then we would have to change `flags` to be a `[]string`.

Update: a PR has been opened to do this: [fluxcd/flux2#1390](https://github.com/fluxcd/flux2/pull/1390).

Even if we were able to change that flag in Flux cli, we would have to hope that they didn't
use `StringArrayVar` for bootstrap flags again (at least for a year, see future work [below](#future-work)).

See more on alternatives [below](#alternatives).

### Future work

It is worth noting that changes the future will result in whatever we do now being
removed, maybe a year from now.

## Implementation Stages

_These are broken down with the intention of being done in separate PRs. I will create separate tickets._

### Remove all the things

Remove all current config opts which are not `gitProvider`. Add a `flags` (or whatever) key
which is a `map[string]string`. Ensure documentation is updated and that the API change is
publicised. For bonus points, if they try to use the old structure (or any thing which we do not recognise), explain things have
changed in logs/error and point them to the PR/docs.
If Flux fails after cluster create, suggest running `eksctl enable flux` once config is resolved.

Don't forget to update the example yaml with something thorough.

### Docs

- Document required Env vars and Flux binary
- Make sure to link to https://fluxcd.io/ everywhere possible
- Document discovering all options by running `flux bootstrap <provider> --help`
- Make it clear how flags can be set in eksctl config
- Move Flux 2 docs away from general gitops page in eksctl.io (dedicated page)
- Mark old Flux 1 docs as deprecated
- Note that quickstartprofiles will not be implemented for Flux 2

### Recommended Flux version

Once [our PR](https://github.com/fluxcd/flux2/pull/1390) to Flux cli is merged and released,
we should add a check in `eksctl` for that version or above. This version is not required,
since it only effects those using the `team` flag for more than one team.

We should log a warning that their Flux version is `foo` and in order to use `--team`
they will need at least `bar`. We should not fail. (Or should we? Maybe setting a min version is
just simpler :shrug:?

### "Dry run"

If we think this is a good idea/worthwhile, validate that flags are applicable and  whatever required vars `GITX_TOKEN` etc are set
by calling Flux with a dummy command.

## Test coverage

I don't think we need to get too into testing, given that we are just "passing through"
and don't want to end up testing Flux itself.

The current integration test can be altered pass in arbitrary args,
and the units can handle testing that correct opts end up in the `exec`. This should be sufficient.

<!-- ## Drawbacks -->

<!-- _Being devil's advocate._ -->

## API spec

```go
// GitOps groups all configuration options related to enabling GitOps Toolkit on a
// cluster and linking it to a Git repository.
// Note: this will replace the older Git types
type GitOps struct {
	// [Enable Toolkit](/usage/gitops/#experimental-installing-gitops-toolkit-flux-v2)
	Flux *Flux `json:"flux,omitempty"`
}

// Flux groups all configuration options related to a Git repository used for
// GitOps Toolkit (Flux v2).
type Flux struct {
	// The repository hosting service. Can be either Github, Gitlab or Git. Required.
	GitProvider string `json:"gitProvider,omitempty"`

	// Arbitrary map of key value pairs which will be passed to Flux bootstrap as flags
	Flags map[string]string `json:"namespace,omitempty"`
}
```

## Alternatives

### List of strings

The main alternative for if we cannot cover all Flux flags in a map.

```diff
  gitops:
    flux:
      gitProvider: github
+     args:
+       - "--flag-one=value"
+       - "--flag-two value"
+       - "--bool-flag"
```

### Actual config

In this option we have a mix of 1:1 config:flag mappings and an arbitrary option
for all non-explicit config.

- All provider flags
- Some common/global flags
- Arbitrary list of string

```yaml
gitops:
  flux:
    # gitlab, github and git would not be specified all at once ofc, but you get the idea
    gitlab:
      branch: main
      namespace: "flux-system"
      components:
      - "source-controller"
      - "kustomise-controller"
      fluxExtraArgs:
	- "--flag-one=value"
	- "--flag-two value"
	- "--bool-flag"
      interval: 1m0s
      path: "clusters/cluster-27"
      owner: dr-who
      repository: my-cluster-gitops
      personal: true
      hostname: "gh-enterprise.com"
      interval: 1m0s
      private: true
      readWriteKey: true
      reconcile: true
      teams:
      - "team1"
      - "team2"
    github:
      ... would be same as gitlag
    git:
      branch: main
      namespace: "flux-system"
      components:
      - "source-controller"
      - "kustomise-controller"
      fluxExtraArgs:
	- "--flag-one=value"
	- "--flag-two value"
	- "--bool-flag"
      interval: 1m0s
      path: "clusters/cluster-27"
      username: "dr-who"
      passwordFile: /basic/auth/pass
      url: https://<host>/<org>/<repository>
```

<!-- ## Open Questions / Known Unknowns -->

<!--
List any questions for things you unsure about or to
direct reviewers to particular areas where their expertise is needed.
-->

