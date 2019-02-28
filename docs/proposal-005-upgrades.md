# Design Proposal #005: Upgrades

> **STATUS**: This proposal is a _final_ state, and we expect minimal additional refinements.
> If any non-trivial changes are needed to functionality defined here, in particular the user
> experience, those changes should be suggested via a PR to this proposal document.
> Any other changes to the text of the proposal or technical corrections are also very welcome.

Cluster upgrades are inherintly a multi-step process, and can be fairly complex.

With `eksctl` users should be able to easily upgrade from one version of Kubernetes to another.
There maybe additional manual steps with some versions, however most parts should be automated.

When upgrading a cluster, one needs to call `eks.UpdateClusterVersion`. After that they need
to replace nodegroups one by one.

## Initial phase

- provide command that checks cluster stack for upgradability
  - let's user update cluster stack to cater for any additional reources
  - allows to call `eks.UpdateClusterVersion` out-of-band and wait for completion
- provide instruction on how to iterate and replace nodegoups
- provide instruction on how to 

## Final phase

- use CloudFormation instead of calling `eks.UpdateClusterVersion` directly
- provide automated command
