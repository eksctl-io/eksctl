# `eksctl` - a CLI for Amazon EKS

[![Circle CI](https://circleci.com/gh/weaveworks/eksctl/tree/master.svg?style=shield)](https://circleci.com/gh/weaveworks/eksctl/tree/master) [![Coverage Status](https://coveralls.io/repos/github/weaveworks/eksctl/badge.svg?branch=master)](https://coveralls.io/github/weaveworks/eksctl?branch=master)[![Go Report Card](https://goreportcard.com/badge/github.com/weaveworks/eksctl)](https://goreportcard.com/report/github.com/weaveworks/eksctl)

`eksctl` is a simple CLI tool for creating clusters on EKS - Amazon's new managed Kubernetes service for EC2. It is written in Go, and uses CloudFormation.

You can create a cluster in minutes with just one command – **`eksctl create cluster`**!

![Gophers: E, K, S, C, T, & L](logo/eksctl.png)

## Usage

### Install

To download the latest release, run:

```
curl --silent --location "https://github.com/weaveworks/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin
```

Alternatively, macOS users can use [Homebrew](https://brew.sh):
```
brew install weaveworks/tap/eksctl
```

You will need to have AWS API credentials configured. What works for AWS CLI or any other tools (kops, Terraform etc), should be sufficient. You can use [`~/.aws/credentials` file][awsconfig]
or [environment variables][awsenv]. For more information read [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html).

[awsenv]: https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html
[awsconfig]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html

#### Stock kubectl vs EKS-vended kubectl

The recommended way to install kubectl on macOS is `brew install kubernetes-cli`. For Kubernetes v1.10, the stock kubectl client works with EKS after installing AWS IAM Authenticator for Kubernetes. Alternatively, you can install Amazon's patched [EKS-vended kubectl binary](https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html) which works with EKS out of the box but lacks automatic upgrades.

#### AWS IAM Authenticator for Kubernetes

Upon creating a cluster, you'll find the cluster credentials added to `~/.kube/config`. If you have `kubectl` v1.10+ as well as the [AWS IAM Authenticator for Kubernetes](https://github.com/kubernetes-sigs/aws-iam-authenticator) command (either `aws-iam-authenticator` or `heptio-authenticator-aws`) in your PATH, you can use `kubectl` with EKS. Make sure to use the same AWS API credentials for this also. Check [EKS docs][ekskubectl] for instructions. If you install `eksctl` via Homebrew, it will automatically install `heptio-authenticator-aws`.

##### aws-iam-authenticator vs heptio-authenticator-aws

The AWS IAM authenticator for EKS was initially developed by Heptio. It was originally named `heptio-authenticator-aws` and was renamed to `aws-iam-authenticator`. Both names refer to the same package.

[ekskubectl]: https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html

## Commands

To create a basic cluster, run:

```
eksctl create cluster
```

A cluster will be created with default parameters
- exciting auto-generated name, e.g. "fabulous-mushroom-1527688624"
- 2x `m5.large` nodes (this instance type suits most common use-cases, and is good value for money)
- use official AWS EKS AMI
- `us-west-2` region
- dedicated VPC (check your quotas)
- using static AMI resolver

Example output:
```
$ eksctl create cluster
2018-08-06T16:32:59+01:00 [ℹ]  setting availability zones to [us-west-2c us-west-2b us-west-2a]
2018-08-06T16:32:59+01:00 [ℹ]  creating EKS cluster "adorable-painting-1533569578" in "us-west-2" region
2018-08-06T16:32:59+01:00 [ℹ]  will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup
2018-08-06T16:32:59+01:00 [ℹ]  if you encounter any issues, check CloudFormation console first
2018-08-06T16:32:59+01:00 [ℹ]  creating cluster stack "eksctl-adorable-painting-1533569578-cluster"
2018-08-06T16:43:43+01:00 [ℹ]  creating nodegroup stack "eksctl-adorable-painting-1533569578-nodegroup-0"
2018-08-06T16:47:14+01:00 [✔]  all EKS cluster resource for "adorable-painting-1533569578" had been created
2018-08-06T16:47:14+01:00 [✔]  saved kubeconfig as "/Users/ilya/.kube/config"
2018-08-06T16:47:20+01:00 [ℹ]  the cluster has 0 nodes
2018-08-06T16:47:20+01:00 [ℹ]  waiting for at least 2 nodes to become ready
2018-08-06T16:47:57+01:00 [ℹ]  the cluster has 2 nodes
2018-08-06T16:47:57+01:00 [ℹ]  node "ip-192-168-115-52.us-west-2.compute.internal" is ready
2018-08-06T16:47:57+01:00 [ℹ]  node "ip-192-168-217-205.us-west-2.compute.internal" is ready
2018-08-06T16:48:00+01:00 [ℹ]  kubectl command should work with "~/.kube/config", try 'kubectl get nodes'
2018-08-06T16:48:00+01:00 [✔]  EKS cluster "adorable-painting-1533569578" in "us-west-2" region is ready
```

To list the details about a cluster or all of the clusters, use:

```
eksctl get cluster [--name=<name>] [--region=<region>]
```

To create the same kind of basic cluster, but with a different name, run:

```
eksctl create cluster --name=cluster-1 --nodes=4
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
eksctl utils write-kubeconfig --name=<name> [--kubeconfig=<path>] [--set-kubeconfig-context=<bool>]
```

To use a 3-5 node Auto Scaling Group, run:

```
eksctl create cluster --name=cluster-5 --nodes-min=3 --nodes-max=5
```

To use 30 `c4.xlarge` nodes and prevent updating current context in `~/.kube/config`, run:

```
eksctl create cluster --name=cluster-6 --nodes=30 --node-type=c4.xlarge --set-kubeconfig-context=false
```


In order to allow SSH access to nodes, `eksctl` imports `~/.ssh/id_rsa.pub` by default, to use a different SSH public key, e.g. `my_eks_node_id.pub`, run:

```
eksctl create cluster --ssh-public-key=my_eks_node_id.pub
```

To use a pre-existing EC2 key pair in `us-east-1` region, you can specify key pair name (which must not resolve to a local file path), e.g. to use `my_kubernetes_key` run:

```
eksctl create cluster --ssh-public-key=my_kubernetes_key --region=us-east-1
```

To add custom tags for all resources, use `--tags`. Note that until
https://github.com/weaveworks/eksctl/issues/25 is resolved, tags will
apply to CloudFormation stacks but not EKS clusters.

```
eksctl create cluster --tags environment=staging --region=us-east-1
```

> NOTE: In `us-east-1` you are likely to get `UnsupportedAvailabilityZoneException`. If you do, copy the suggested zones and pass `--zones` flag, e.g. `eksctl create cluster --region=us-east-1 --zones=us-east-1a,us-east-1b,us-east-1d`. This may occur in other regions, but less likely. You shouldn't need to use `--zone` flag otherwise.

To delete a cluster, run:

```
eksctl delete cluster --name=<name> [--region=<region>]
```
### GPU Support

If you'd like to use GPU instance types (i.e. [p2](https://aws.amazon.com/ec2/instance-types/p2/) or [p3](https://aws.amazon.com/ec2/instance-types/p3/) ) then the first thing you need to do is subscribe to the [EKS-optimized AMI with GPU Support](https://aws.amazon.com/marketplace/pp/B07GRHFXGM). If you don't do this then node creation will fail.

After subscribing to the AMI you can create a cluster specifying the GPU instance type you'd like to use for the nodes. For example:

```
eksctl create cluster --node-type=p2.xlarge
```

The AMI resolvers (both static and auto) will see that you want to use a GPU instance type (p2 or p3 only) and they will select the correct AMI.

Once the cluster is created you will need to install the [NVIDIA Kubernetes device plugin](https://github.com/NVIDIA/k8s-device-plugin). Check the repo for the most up to date instructions but you should be able to run this:

```
kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/v1.11/nvidia-device-plugin.yml
```

> Once `addon` support has been added as part of 0.2.0 its envisioned that there will be a addon to install the NVIDIA Kubernetes Device Plugin.  This addon could potentially be installed automatically as we know an GPU instance type is being used.

### Latest & Custom AMI Support
With the the 0.1.2 release we have introduced the `--node-ami` flag for use when creating a cluster. This enables a number of advanced use cases such as using a custom AMI or querying AWS in realtime to determine which AMI to use (non-GPU and GPU instances).

The `--node-ami` can take the AMI image id for an image to explicitly use. It also can take the following 'special' keywords:

| Keyword | Description |
| ------------ | -------------- |
| static       | Indicates that the AMI images ids embedded into eksctl should be used. This relates to the static resolvers. |
| auto        | Indicates that the AMI to use for the nodes should be found by querying AWS. This relates to the auto resolver. | 

If, for example, AWS release a new version of the EKS node AMIs and a new version of eksctl hasn't been released you can use the latest AMI by doing the following:

```
eksctl create cluster --node-ami=auto
```

<!-- TODO for 0.3.0
To use more advanced configuration options, [Cluster API](https://github.com/kubernetes-sigs/cluster-api):

```
eksctl apply --cluster-config advanced-cluster.yaml
```
-->

## Project Roadmap

### Developer use-case (0.2.0)

It should suffice to install a cluster for development with just a single command. Here are some examples:

To create a cluster with default configuration (2 `m4.large` nodes), run:

```
eksctl create cluster
```

The developer may choose to pre-configure popular addons, e.g.:

- Weave Net: `eksctl create cluster --networking weave`
- Helm: `eksctl create cluster --addons helm`
- AWS CI tools (CodeCommit, CodeBuild, ECR): `eksctl create cluster --addons aws-ci`
- Jenkins X: `eksctl create cluster --addons jenkins-x`
- AWS CodeStar: `eksctl create cluster --addons aws-codestar`
- Weave Scope and Flux: `eksctl create cluster --addons weave-scope,weave-flux`

It should be possible to combine any or all of these addons.

It would also be possible to add any of the addons after cluster was created with `eksctl create addons`.

### Manage EKS the GitOps way (0.3.0)

Just like `kubectl`, `eksctl` aims to be compliant with GitOps model, and can be used as part of a GitOps toolkit!

For example, `eksctl apply --cluster-config prod-cluster.yaml` will manage cluster state declaratively.

And `eksctld` will be a controller inside of one cluster that can manage multiple other clusters based on Kubernetes Cluster API definitions (CRDs).

## Contributions

Code contributions are very welcome. If you are interested in helping make eksctl great then see our [contributing guide](CONTRIBUTING.md).

## Get in touch

[Create an issue](https://github.com/weaveworks/eksctl/issues/new), or login to [Weave Community Slack (#eksctl)](https://weave-community.slack.com/messages/CAYBZBWGL/) ([signup](https://slack.weave.works/)).

> ***Logo Credits***
>
> *Original Gophers drawn by [Ashley McNamara](https://twitter.com/ashleymcnamara), unique E, K, S, C, T & L Gopher identities had been produced with [Gopherize.me](https://gopherize.me/).*
