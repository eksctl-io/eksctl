# New: GitOps with Flux v2

!!! warning "Flux v2 `eksctl` integration is **experimental**"
    Partial experimental support for [Flux v2 (part of the GitOps Toolkit)](https://toolkit.fluxcd.io/)
    was made available in `eksctl` version `0.38.0`.

    Complete, but still **experimental**, support for Flux v2 is available in `eksctl` from version `0.53.0`.
    There are significant API changes from the `0.38.0-0.52.0` support with this release, please refer to this documentation.

    Disruptions are still likely until development has stabilised.
    All changes will be noted in [releases](https://github.com/eksctl-io/eksctl/releases).

[Gitops][gitops] is a way to do Kubernetes application delivery. It
works by using Git as a single source of truth for Kubernetes resources
and everything else. With Git at the center of your delivery pipelines,
you and your team can make pull requests to accelerate and simplify
application deployments and operations tasks to Kubernetes.

[gitops]: https://www.weave.works/technologies/gitops/

## Installing Flux v2 (GitOps Toolkit)

_Initial Flux v2 support was released in `0.38.0`, but the API has undergone significant changes. An older
example can be found [here](https://github.com/eksctl-io/eksctl/blob/a854bdbd3e47059d860a78600c2c1813281f4ebf/examples/12-gitops-toolkit.yaml)._

From version `0.53.0`, `eksctl` provides an option to alternatively bootstrap [Flux v2](https://toolkit.fluxcd.io/)
components into an EKS cluster, with the `enable flux` subcommand.

```console
eksctl enable flux --config-file <config-file>
```

The `enable flux` command will shell out to the `flux` binary and run the `flux bootstrap` command
against the cluster.

In order to allow users to specify whichever `bootstrap` flags they like, the `eksctl`
API exposes an arbitrary `map[string]string` of `flags`. To find out which flags you need
to bootstrap your cluster, refer to the [Flux docs](https://fluxcd.io/docs/),
or run `flux bootstrap --help`.

Example:
```YAML
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-12
  region: eu-north-1

# other cluster config ...

gitops:
  flux:
    gitProvider: github      # required. options are github, gitlab or git
    flags:                   # required. arbitrary map[string]string for all flux args.
    # these args are not controlled by eksctl. see https://fluxcd.io/docs/get-started/ for all available flags
      owner: "dr-who"
      repository: "our-org-gitops-repo"
      private: "true"
      branch: "main"
      namespace: "flux-system"
      path: "clusters/cluster-12"
      team: "team1,team2"
```

This example configuration can be found [here](https://github.com/eksctl-io/eksctl/blob/main/examples/12-gitops-toolkit.yaml).

???+ note
    Flux v2 configuration can **only** be provided via configuration file; no flags
    are exposed on this subcommand other than `--config-file`.

Flux will install default toolkit components to the cluster, unless told otherwise by your configuration:

```console
kubectl get pods --namespace flux-system
NAME                                       READY   STATUS    RESTARTS   AGE
helm-controller-7cfb98d895-zmmfc           1/1     Running   0          3m30s
kustomize-controller-557986cf44-2jwjh      1/1     Running   0          3m35s
notification-controller-65694dc94d-rhbxk   1/1     Running   0          3m20s
source-controller-7f856877cf-jgwdk         1/1     Running   0          3m39s
```

For instructions on how to use your newly installed Gitops Toolkit,
refer to the [official docs](https://fluxcd.io/docs/get-started/).

### Bootstrap after cluster create

You can have your cluster bootstrapped immediately following a cluster create
by including your Flux configuration in your config file and running:

```console
eksctl create cluster --config-file <config-file>
```

### Requirements

#### Environment variables

Before running `eksctl enable flux`, ensure that you have read the [Flux bootstrap docs](https://fluxcd.io/docs/get-started/).
If you are using Github or Gitlab as your git provider, either `GITHUB_TOKEN` or `GITLAB_TOKEN`
must be exported with your Personal Access Token in your session. Please refer to the Flux docs
for any other requirements.

#### Flux version

Eksctl requires a minimum Flux version of `0.13.3`.

## Quickstart profiles

Quickstart profiles will **not** be supported with Flux v2.

## Further reading

To learn more about gitops and Flux, check the Flux [documentation](https://fluxcd.io/).
