# Getting started

!!! tip "New for 2024"
    `eksctl` now supports new region Kuala Lumpur (`ap-southeast-5`)

    EKS Add-ons now support receiving IAM permissions via [EKS Pod Identity Associations](/usage/pod-identity-associations/#eks-add-ons-support-for-pod-identity-associations)

    `eksctl` now supports AMIs based on AmazonLinux2023

!!! tip "eksctl main features in 2023"
    `eksctl` now supports configuring cluster access management via [AWS EKS Access Entries](/usage/access-entries).

    `eksctl` now supports configuring fine-grained permissions to EKS running apps via [EKS Pod Identity Associations](/usage/pod-identity-associations)

    `eksctl` now supports [updating the subnets and security groups](/usage/cluster-subnets-security-groups) associated with the EKS control plane.

    `eksctl` now supports creating fully private clusters on [AWS Outposts](/usage/outposts).

    `eksctl` now supports new ISO regions `us-iso-east-1` and `us-isob-east-1`.

    `eksctl` now supports new regions - Calgary (`ca-west-1`), Tel Aviv (`il-central-1`), Melbourne (`ap-southeast-4`), Hyderabad (`ap-south-2`), Spain (`eu-south-2`) and Zurich (`eu-central-2`).

`eksctl` is a simple CLI tool for creating and managing clusters on EKS - Amazon's managed Kubernetes service for EC2.
It is written in Go, uses CloudFormation, was created by [Weaveworks](https://www.weave.works/) and it welcomes
contributions from the community.

!!! example "Create a basic cluster in minutes with just one command"
    ```
    eksctl create cluster
    ```
    ![eksctl create cluster](img/eksctl-gopher.png){ align=right width=25% }

    A cluster will be created with default parameters:

    - exciting auto-generated name, e.g., `fabulous-mushroom-1527688624`
    - two `m5.large` worker nodes (this instance type suits most common use-cases, and is good value for money)
    - use the official AWS [EKS AMI](https://github.com/awslabs/amazon-eks-ami)
    - `us-west-2` region
    - a dedicated VPC (check your quotas)

    ???- info "Example output"
        ```sh
          $ eksctl create cluster
          [ℹ]  using region us-west-2
          [ℹ]  setting availability zones to [us-west-2a us-west-2c us-west-2b]
          [ℹ]  subnets for us-west-2a - public:192.168.0.0/19 private:192.168.96.0/19
          [ℹ]  subnets for us-west-2c - public:192.168.32.0/19 private:192.168.128.0/19
          [ℹ]  subnets for us-west-2b - public:192.168.64.0/19 private:192.168.160.0/19
          [ℹ]  nodegroup "ng-98b3b83a" will use "ami-05ecac759c81e0b0c" [AmazonLinux2/1.11]
          [ℹ]  creating EKS cluster "floral-unicorn-1540567338" in "us-west-2" region
          [ℹ]  will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup
          [ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=us-west-2 --cluster=floral-unicorn-1540567338'
          [ℹ]  2 sequential tasks: { create cluster control plane "floral-unicorn-1540567338", create nodegroup "ng-98b3b83a" }
          [ℹ]  building cluster stack "eksctl-floral-unicorn-1540567338-cluster"
          [ℹ]  deploying stack "eksctl-floral-unicorn-1540567338-cluster"
          [ℹ]  building nodegroup stack "eksctl-floral-unicorn-1540567338-nodegroup-ng-98b3b83a"
          [ℹ]  --nodes-min=2 was set automatically for nodegroup ng-98b3b83a
          [ℹ]  --nodes-max=2 was set automatically for nodegroup ng-98b3b83a
          [ℹ]  deploying stack "eksctl-floral-unicorn-1540567338-nodegroup-ng-98b3b83a"
          [✔]  all EKS cluster resource for "floral-unicorn-1540567338" had been created
          [✔]  saved kubeconfig as "~/.kube/config"
          [ℹ]  adding role "arn:aws:iam::376248598259:role/eksctl-ridiculous-sculpture-15547-NodeInstanceRole-1F3IHNVD03Z74" to auth ConfigMap
          [ℹ]  nodegroup "ng-98b3b83a" has 1 node(s)
          [ℹ]  node "ip-192-168-64-220.us-west-2.compute.internal" is not ready
          [ℹ]  waiting for at least 2 node(s) to become ready in "ng-98b3b83a"
          [ℹ]  nodegroup "ng-98b3b83a" has 2 node(s)
          [ℹ]  node "ip-192-168-64-220.us-west-2.compute.internal" is ready
          [ℹ]  node "ip-192-168-8-135.us-west-2.compute.internal" is ready
          [ℹ]  kubectl command should work with "~/.kube/config", try 'kubectl get nodes'
          [✔]  EKS cluster "floral-unicorn-1540567338" in "us-west-2" region is ready
        ```

Customize your cluster by using a config file. Just run

```sh
eksctl create cluster -f cluster.yaml
```

to apply a `cluster.yaml` file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: basic-cluster
  region: eu-north-1

nodeGroups:
  - name: ng-1
    instanceType: m5.large
    desiredCapacity: 10
  - name: ng-2
    instanceType: m5.xlarge
    desiredCapacity: 2
```

Once you have created a cluster, you will find that cluster credentials were added in `~/.kube/config`. If you have
`kubectl` v1.10.x as well as `aws-iam-authenticator` commands in your PATH, you should be
able to use `kubectl`. You will need to make sure to use the same AWS API credentials for this also. Check
[EKS docs][ekskubectl] for instructions. If you installed `eksctl` via Homebrew, you should have all of these
dependencies installed already.

To learn more about how to create clusters and other features continue reading the
[Creating and Managing Clusters section](usage/creating-and-managing-clusters).

[ekskubectl]: https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html

## Listing clusters

To list the details about a cluster or all of the clusters, use:

```sh
eksctl get cluster [--name=<name>] [--region=<region>]
```

## Basic cluster creation

To create a basic cluster, but with a different name, run:

```sh
eksctl create cluster --name=cluster-1 --nodes=4
```

### Supported versions

EKS supports versions `1.23` (extended), `1.24` (extended), `1.25`, `1.26`, `1.27`, `1.28`, `1.29`, **`1.30`** (default) and `1.31`.
With `eksctl` you can deploy any of the supported versions by passing `--version`.

```sh
eksctl create cluster --version=1.28
```

### Config-based creation

You can also create a cluster passing all configuration information in a file
using `--config-file`:

```sh
eksctl create cluster --config-file=<path>
```

To create a cluster using a configuration file and skip creating
nodegroups until later:

```sh
eksctl create cluster --config-file=<path> --without-nodegroup
```

### Cluster credentials

To write cluster credentials to a file other than default, run:

```sh
eksctl create cluster --name=cluster-2 --nodes=4 --kubeconfig=./kubeconfig.cluster-2.yaml
```

To prevent storing cluster credentials locally, run:

```sh
eksctl create cluster --name=cluster-3 --nodes=4 --write-kubeconfig=false
```

To let `eksctl` manage cluster credentials under `~/.kube/eksctl/clusters` directory, run:

```sh
eksctl create cluster --name=cluster-3 --nodes=4 --auto-kubeconfig
```

To obtain cluster credentials at any point in time, run:

```sh
eksctl utils write-kubeconfig --cluster=<name> [--kubeconfig=<path>] [--set-kubeconfig-context=<bool>]
```

### Caching Credentials

`eksctl` supports caching credentials. This is useful when using MFA and not wanting to continuously enter the MFA
token on each `eksctl` command run.

To enable credential caching set the following environment property `EKSCTL_ENABLE_CREDENTIAL_CACHE` as such:

```sh
export EKSCTL_ENABLE_CREDENTIAL_CACHE=1
```

By default, this will result in a cache file under `~/.eksctl/cache/credentials.yaml` which will contain creds per profile
that is being used. To clear the cache, delete this file.

It's also possible to configure the location of this cache file using `EKSCTL_CREDENTIAL_CACHE_FILENAME` which should
be the **full path** to a file in which to store the cached credentials. These are credentials, so make sure the access
of this file is restricted to the current user and in a secure location.

## Autoscaling

To use a 3-5 node Auto Scaling Group, run:

```sh
eksctl create cluster --name=cluster-5 --nodes-min=3 --nodes-max=5
```

You will still need to install and configure Auto Scaling. See the "Enable Auto Scaling" section. Also
note that depending on your workloads you might need to use a separate nodegroup for each AZ. See [Zone-aware
Auto Scaling](/usage/autoscaling/) for more info.

## SSH access

In order to allow SSH access to nodes, `eksctl` imports `~/.ssh/id_rsa.pub` by default, to use a different SSH public key, e.g. `my_eks_node_id.pub`, run:

```sh
eksctl create cluster --ssh-access --ssh-public-key=my_eks_node_id.pub
```

To use a pre-existing EC2 key pair in `us-east-1` region, you can specify key pair name (which must not resolve to a local file path), e.g. to use `my_kubernetes_key` run:

```sh
eksctl create cluster --ssh-access --ssh-public-key=my_kubernetes_key --region=us-east-1
```

[AWS Systems Manager (SSM)](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-sessions-start.html#sessions-start-cli) is enabled by default, so it can be used to SSH onto nodes.

```sh
eksctl create cluster --enable-ssm
```

If you are creating managed nodes with a custom launch template, the `--enable-ssm` flag is disallowed.

## Tagging

To add custom tags for all resources, use `--tags`.

```sh
eksctl create cluster --tags environment=staging --region=us-east-1
```

## Volume size

To configure node root volume, use the `--node-volume-size` (and optionally `--node-volume-type`), e.g.:

```sh
eksctl create cluster --node-volume-size=50 --node-volume-type=io1
```

???+ note
    The default volume size is 80G.

## Deletion

To delete a cluster, run:

```sh
eksctl delete cluster --name=<name> [--region=<region>]
```

???+ note
    Cluster info will be cleaned up in kubernetes config file. Please run `kubectl config get-contexts` to select right context.
