---
title: "Getting started"
weight: 10
url: introduction/getting-started
---

## Getting started

_Need help? Join [Weave Community Slack][slackjoin]._
[slackjoin]: https://slack.weave.works/

To list the details about a cluster or all of the clusters, use:

```

eksctl get cluster [--name=<name>][--region=<region>]

```

To create the same kind of basic cluster, but with a different name, run:

```

eksctl create cluster --name=cluster-1 --nodes=4

```

EKS supports versions `1.10`, `1.11`, `1.12`, `1.13` and `1.14` (default).
With `eksctl` you can deploy either version by passing `--version`.

```

eksctl create cluster --version=1.10

```

To write cluster credentials to a file other than default, run:

```

eksctl create cluster --name=cluster-2 --nodes=4 --kubeconfig=./kubeconfig.cluster-2.yaml

```

To prevent storing cluster credentials locally, run:

```

eksctl create cluster --name=cluster-3 --nodes=4 --write-kubeconfig=false

```

To let `eksctl` manage cluster credentials under `~/.kube/eksctl/clusters` directory, run:

```

eksctl create cluster --name=cluster-3 --nodes=4 --auto-kubeconfig

```

To obtain cluster credentials at any point in time, run:

```

eksctl utils write-kubeconfig --cluster=<name> [--kubeconfig=<path>][--set-kubeconfig-context=<bool>]

```

To use a 3-5 node Auto Scaling Group, run:

```

eksctl create cluster --name=cluster-5 --nodes-min=3 --nodes-max=5

```

> NOTE: You will still need to install and configure Auto Scaling. See the "Enable Auto Scaling" section below. Also
 note that depending on your workloads you might need to use a separate nodegroup for each AZ. See [Zone-aware Auto
 Scaling](#zone-aware-Auto Scaling) below for more info.

To use 30 `c4.xlarge` nodes and prevent updating current context in `~/.kube/config`, run:

```

eksctl create cluster --name=cluster-6 --nodes=30 --node-type=c4.xlarge --set-kubeconfig-context=false

```

In order to allow SSH access to nodes, `eksctl` imports `~/.ssh/id_rsa.pub` by default, to use a different SSH public key, e.g. `my_eks_node_id.pub`, run:

```

eksctl create cluster --ssh-access --ssh-public-key=my_eks_node_id.pub

```

To use a pre-existing EC2 key pair in `us-east-1` region, you can specify key pair name (which must not resolve to a local file path), e.g. to use `my_kubernetes_key` run:

```

eksctl create cluster --ssh-access --ssh-public-key=my_kubernetes_key --region=us-east-1

```

To add custom tags for all resources, use `--tags`.

> NOTE: Until [#25](https://github.com/weaveworks/eksctl/issues/25) is resolved, tags cannot be applied to EKS cluster itself, but most of other resources (e.g. EC2 nodes).

```

eksctl create cluster --tags environment=staging --region=us-east-1

```

To configure node root volume, use the `--node-volume-size` (and optionally `--node-volume-type`), e.g.:

```

eksctl create cluster --node-volume-size=50 --node-volume-type=io1

```

> NOTE: In `us-east-1` you are likely to get `UnsupportedAvailabilityZoneException`. If you do, copy the suggested zones and pass `--zones` flag, e.g. `eksctl create cluster --region=us-east-1 --zones=us-east-1a,us-east-1b,us-east-1d`. This may occur in other regions, but less likely. You shouldn't need to use `--zone` flag otherwise.

You can also create a cluster passing all configuration information in a file
using `--config-file`:

```

eksctl create cluster --config-file=<path>

```

To create a cluster using a configuration file and skip creating
nodegroups until later:

```

eksctl create cluster --config-file=<path> --without-nodegroup

```

To delete a cluster, run:

```

eksctl delete cluster --name=<name> [--region=<region>]

```

> NOTE: Cluster info will be cleaned up in kubernetes config file. Please run `kubectl config get-contexts` to select right context.

### Contributions

Code contributions are very welcome. If you are interested in helping make `eksctl` great then see our [contributing guide](https://github.com/weaveworks/eksctl/blob/master/CONTRIBUTING.md).
