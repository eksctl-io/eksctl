Authors: @weaveworks/eksctl-reviewers
### Status: Design ongoing

# v1.0 changes

This proposal outlines a list of features and breaking changes to be included in the v1.0 release of eksctl. This
proposal can be used as a place to reference other proposal for changes and should serve as an overarching document to
see what will change in v1.0.

### Draft list of changes
This is a draft list of changes that could be made in a v1.0 release. Some of these may required their own proposals.

- Introduction of [eksctl apply](https://github.com/weaveworks/eksctl/pull/3037)

- Removal of deprecated functionality
  - [`eksctl update` command](https://github.com/weaveworks/eksctl/blob/c7a2c6632c4ccf53dc4bdb69b28f465a8281eb5d/pkg/ctl/update/cluster.go#L21-L22)

  - [`--name` flag for cluster name](https://github.com/weaveworks/eksctl/blob/c7a2c6632c4ccf53dc4bdb69b28f465a8281eb5d/pkg/ctl/cmdutils/cmdutils.go#L139-L146)

  - [`--role` flag for IAMIdentityMapping](https://github.com/weaveworks/eksctl/blob/c7a2c6632c4ccf53dc4bdb69b28f465a8281eb5d/pkg/ctl/cmdutils/iam_flags.go#L16)

  - [`static` AMI value](https://github.com/weaveworks/eksctl/blob/c7a2c6632c4ccf53dc4bdb69b28f465a8281eb5d/pkg/ctl/cmdutils/nodegroup_flags.go#L42)

  - [`update-cluster-stack` command](https://github.com/weaveworks/eksctl/blob/2e3ce5b9d1fcc0bed960d88bfb4251da8479ab9f/pkg/ctl/utils/update_cluster_stack.go#L13-L18)

  - [`wait-nodes` command](https://github.com/weaveworks/eksctl/blob/2e3ce5b9d1fcc0bed960d88bfb4251da8479ab9f/pkg/ctl/utils/wait_nodes.go#L14-L16)

  - [Fix typo in asset filename](https://github.com/weaveworks/eksctl/pull/3020)

- New config spec. Needs its own proposal and will be heavily linked to [eksctl apply](https://github.com/weaveworks/eksctl/pull/3037)

- Refactoring to consistently use or remove the `--approve` flag. This is used in some commands such as `eksctl upgrade cluster`
  but not in others. Deciding a consistent, documented usage of this will help improve UX

- Removal of the `--include` `--exclude` `--only-missing` flags in the `eksctl delete nodegroup` command. With the
  introduction of `eksctl apply` these are no longer needed.


## Goals
- TODO

## Non-Goals
- TODO
## Risks
- Adoption of the new CLI version will be slow

- Users will dislike migrating config file format

- Community will ask us to maintain the older version

