# `eksctl` - a CLI for Amazon EKS

[![Circle CI](https://circleci.com/gh/polyverse/eksctl/tree/master.svg?style=shield)](https://circleci.com/gh/polyverse/eksctl/tree/master)

`eksctl` is a simple CLI tool for creating clusters on EKS - Amazon's new managed Kubernetes service for EC2. It is written in Go, and based on Amazon's official CloudFormation templates.

You can create a cluster in minutes with just one command – **`eksctl create cluster`**!

![Gophers: E, K, S, C, T, & L](logo/eksctl.png)

## Polyverse fork - intended for Polyverse-internal Production-safety

The Polyverse fork of eksctl is intended to have code that is dirty, ugly, logs verbosely, doesn't "think through" go primitives, in order to prioritize above all else production-safety first. All changes are posted through PRs upstream.

This fork provides:
1. Idempotent operations - you can run create/delete repeatedly for the same cluster. You can CTRL+C an operation and re-run it to pick it up where it left off. You can repair a deleted stack, by re-running create.

2. Getting-it-working-first priority - rather than "develop for sandbox; then figure out how to make it work", we focus mainly on ensuring launching in existing VPCs is easy, adding user-data to an existing AMI is possible, etc. etc. How people *should* do things is not enforced by limiting this tool, but rather through education and encouragement. We might have features that are anti-patterns if that's what getting a scenario up-and-running requires. EKS and Kubernetes are going to work when they can be launched right now.

3. Non-ideal world; Production-safety - verbose logging that assumes things won't be how you expect. Rather than error out and leave it as an exercise to the user, these situations might be verbosely logged, explained, etc. Logs may contain remediation commands inline. Information is power.

## Usage

To download the latest release, run:

```
curl --silent --location "https://github.com/polyverse/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin
```

You will need to have AWS API credentials configured. What works for AWS CLI or any other tools (kops, Terraform etc), should be sufficient. You can use [`~/.aws/credentials` file][awsconfig]
or [environment variables][awsenv]. For more information read [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html).

[awsenv]: https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html
[awsconfig]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html

To create a basic cluster, run:

```
eksctl create cluster
```

A cluster will be created with default parameters
- exciting auto-generated name, e.g. "fabulous-mushroom-1527688624"
- 2x `m5.large` nodes (this instance type suits most common use-cases, and is good value for money)
- default EKS AMI
- `us-west-2` region
- dedicated VPC (check your quotas)

Once you have created a cluster, you will find `kubeconfig` in your current working directory. If you have `kubectl` v1.10.x as well as `heptio-authenticator-aws` commands in your PATH, you should be
able to use `kubectl`. You will need to make sure to use the same AWS API credentials for this also. Check [EKS docs][ekskubectl] for instructions.

[ekskubectl]: https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html

Example output:
```
$ eksctl create cluster
2018-06-06T16:40:58+01:00 [ℹ]  importing SSH public key "~/.ssh/id_rsa.pub" as "EKS-extravagant-sculpture-1528299658"
2018-06-06T16:40:58+01:00 [ℹ]  creating EKS cluster "extravagant-sculpture-1528299658" in "us-west-2" region
2018-06-06T16:40:58+01:00 [ℹ]  creating VPC stack "EKS-extravagant-sculpture-1528299658-VPC"
2018-06-06T16:40:58+01:00 [ℹ]  creating ServiceRole stack "EKS-extravagant-sculpture-1528299658-ServiceRole"
2018-06-06T16:41:19+01:00 [✔]  created ServiceRole stack "EKS-extravagant-sculpture-1528299658-ServiceRole"
2018-06-06T16:42:19+01:00 [✔]  created VPC stack "EKS-extravagant-sculpture-1528299658-VPC"
2018-06-06T16:42:19+01:00 [ℹ]  creating control plane "extravagant-sculpture-1528299658"
2018-06-06T16:50:41+01:00 [✔]  created control plane "extravagant-sculpture-1528299658"
2018-06-06T16:50:41+01:00 [ℹ]  creating DefaultNodeGroup stack "EKS-extravagant-sculpture-1528299658-DefaultNodeGroup"
2018-06-06T16:54:22+01:00 [✔]  created DefaultNodeGroup stack "EKS-extravagant-sculpture-1528299658-DefaultNodeGroup"
2018-06-06T16:54:22+01:00 [✔]  all EKS cluster "extravagant-sculpture-1528299658" resources has been created
2018-06-06T16:54:22+01:00 [ℹ]  wrote "kubeconfig2"
2018-06-06T16:54:23+01:00 [ℹ]  the cluster has 0 nodes
2018-06-06T16:54:23+01:00 [ℹ]  waiting for at least 2 nodes to become ready
2018-06-06T16:54:49+01:00 [ℹ]  the cluster has 2 nodes
2018-06-06T16:54:49+01:00 [ℹ]  node "ip-192-168-185-142.ec2.internal" is ready
2018-06-06T16:54:49+01:00 [ℹ]  node "ip-192-168-221-172.ec2.internal" is ready
2018-06-06T16:54:49+01:00 [ℹ]  EKS cluster "extravagant-sculpture-1528299658" is ready in "us-west-2" region
```

To list the details about a cluster or all of the clusters, use:

```
eksctl get cluster [--name <name>] [--region <region>]
```

To create the same kind of basic cluster, but with a different name, run:

```
eksctl create cluster --name cluster-1 --nodes 4
```

To write cluster credentials to a file other than default, run:

```
eksctl create cluster --name cluster-2 --nodes 4 --kubeconfig ./kubeconfig.cluster-2.yaml
```

To prevent storing cluster credentials locally, run:

```
eksctl create cluster --name cluster-3 --nodes 4 --write-kubeconfig=false
```

To let `eksctl` manage cluster credentials under `~/.kube/eksctl/clusters` directory, run:

```
eksctl create cluster --name cluster-3 --nodes 4 --auto-kubeconfig
```

To obtain cluster credentials at any point in time, run:

```
eksctl utils write-kubeconfig --name <name> [--kubeconfig <path>]
```

To use a 3-5 node Auto Scaling Group, run:

```
eksctl create cluster --name cluster-5 --nodes-min 3 --nodes-max 5
```

To use 30 `c4.xlarge` nodes, run:

```
eksctl create cluster --name cluster-6 --nodes 30 --node-type c4.xlarge
```

To delete a cluster, run:

```
eksctl delete cluster --name <name> [--region <region>]
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

Code contributions are very welcome, however until a 0.1.0 release testing and bug reports are the contributions that authors will appreciate the most.

## Get in touch

[Create an issue](https://github.com/weaveworks/eksctl/issues/new), or login to [Weave Community Slack (#eksctl)](https://weave-community.slack.com/messages/CAYBZBWGL/) ([signup](https://slack.weave.works/)).

> ***Logo Credits***
>
> *Original Gophers drawn by [Ashley McNamara](https://twitter.com/ashleymcnamara), unique E, K, S, C, T & L Gopher identities had been produced with [Gopherize.me](https://gopherize.me/).*
