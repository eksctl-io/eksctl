## Getting started

### Basic cluster creation

To create a basic cluster, but with a different name, run:

```

eksctl create cluster --name=cluster-1 --nodes=4

```

EKS supports versions `1.22`, `1.23`, `1.24` (default) and `1.25`.
With `eksctl` you can deploy any of the supported versions by passing `--version`.

```

eksctl create cluster --version=1.24

```

### Listing clusters

To list the details about a cluster or all of the clusters, use:

```

eksctl get cluster [--name=<name>][--region=<region>]

```

#### Config-based creation

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

#### Cluster credentials

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

#### Caching Credentials

`eksctl` supports caching credentials. This is useful when using MFA and not wanting to continuously enter the MFA
token on each `eksctl` command run.

To enable credential caching set the following environment property `EKSCTL_ENABLE_CREDENTIAL_CACHE` as such:

```
export EKSCTL_ENABLE_CREDENTIAL_CACHE=1
```

By default, this will result in a cache file under `~/.eksctl/cache/credentials.yaml` which will contain creds per profile
that is being used. To clear the cache, delete this file.

It's also possible to configure the location of this cache file using `EKSCTL_CREDENTIAL_CACHE_FILENAME` which should
be the **full path** to a file in which to store the cached credentials. These are credentials, so make sure the access
of this file is restricted to the current user and in a secure location.

### Autoscaling

To use a 3-5 node Auto Scaling Group, run:

```

eksctl create cluster --name=cluster-5 --nodes-min=3 --nodes-max=5

```

???+ note
    You will still need to install and configure Auto Scaling. See the "Enable Auto Scaling" section. Also
    note that depending on your workloads you might need to use a separate nodegroup for each AZ. See [Zone-aware
    Auto Scaling](/usage/autoscaling/) for more info.

### SSH access

In order to allow SSH access to nodes, `eksctl` imports `~/.ssh/id_rsa.pub` by default, to use a different SSH public key, e.g. `my_eks_node_id.pub`, run:

```

eksctl create cluster --ssh-access --ssh-public-key=my_eks_node_id.pub

```

To use a pre-existing EC2 key pair in `us-east-1` region, you can specify key pair name (which must not resolve to a local file path), e.g. to use `my_kubernetes_key` run:

```

eksctl create cluster --ssh-access --ssh-public-key=my_kubernetes_key --region=us-east-1

```

[AWS Systems Manager (SSM)](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-sessions-start.html#sessions-start-cli) is enabled by default, so it can be used to SSH onto nodes.


```

eksctl create cluster --enable-ssm

```

???+ note
    If you are creating managed nodes with a custom launch template, the `--enable-ssm` flag is disallowed.

### Tagging

To add custom tags for all resources, use `--tags`.

```

eksctl create cluster --tags environment=staging --region=us-east-1

```

### Volume size

???+ note
    The default volume size is 80G.

To configure node root volume, use the `--node-volume-size` (and optionally `--node-volume-type`), e.g.:

```

eksctl create cluster --node-volume-size=50 --node-volume-type=io1

```

### Deletion

To delete a cluster, run:

```

eksctl delete cluster --name=<name> [--region=<region>]

```

???+ note
    Cluster info will be cleaned up in kubernetes config file. Please run `kubectl config get-contexts` to select right context.

## Contributions

Code contributions are very welcome. If you are interested in helping make `eksctl` great then see our [contributing guide](https://github.com/weaveworks/eksctl/blob/master/CONTRIBUTING.md).

_Need help? Join [Weave Community Slack][slackjoin]._
[slackjoin]: https://slack.weave.works/

