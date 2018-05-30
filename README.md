# `eksctl` - CLI tool to create Amazon EKS clusters

[![Circle CI](https://circleci.com/gh/weaveworks/eksctl/tree/master.svg?style=shield)](https://circleci.com/gh/weaveworks/eksctl/tree/master)

EKS is Amazon's new managed Kubernetes service for EC2.
What is `eksctl`? It's a simple CLI tool for creating EKS clusters, for most common use-cases.

It's inspired by `kubectl`. It provides an easy way to create and manage clusters, and aims to implement a [Cluster API](https://github.com/kubernetes-sigs/cluster-api) controller for EKS also (`eksctld`).

It is not intended to be a like-for-like alternative to well-established community tools (`kops`, `kubicorn`, `kubeadm`).
However, the intention is to work well with most popular tools, and collaborate very closely, so that Kubernetes makes the
cloud-native world even more amazing to live in!

> **Download Today**
>
> Linux, macOS and Windows binaries for 0.1.0-alpha1 release are [available for download](https://github.com/weaveworks/eksctl/releases/tag/0.1.0-alpha1).
>
> **Roadmap**
>
> Stable 0.1.0 release will made available based on user-feedback.
> Release 0.2.0 will add support for addons, and 0.3.0 is planned to support Cluster API.
>
> **Contributions**
>
> Code contributions are very welcome, however until 0.1.0 release testing and bug reports are the contributions that authors will appreciate the most.
> 
> **Get in touch**
>
> [Create and issue](https://github.com/weaveworks/eksctl/issues/new), or login to [Weave Community Slack (#eksctl)](https://weave-community.slack.com/messages/CAYBZBWGL/) ([signup](https://slack.weave.works/)).

## Developer use-case

It should suffice to install a cluster for development with just a single command, here are some examples.

To create a cluster with default configurations (2 `m4.large` nodes), run:
```
eksctl create cluster
```

In 0.2.0, it will support many popular addons, e.g.:

* Weave Net: `eksctl create cluster --networking weave`
* Helm: `eksctl create cluster --addons helm`
* AWS CI tools (CodeCommit, CodeBuild, ECR): `eksctl create cluster --addons aws-ci`
* Jenkins X: `eksctl create cluster --addons jenkins-x`
* AWS CodeStar: `eksctl create cluster --addons aws-codestar`
* Weave Scope and Flux: `eksctl create cluster --addons weave-scope,weave-flux`

<!-- TODO
You can combine any or all of these.

You can also add any of these addons after you create a cluster with `eksctl addons install <addon>...`.
-->

## Manage EKS the GitOps way (0.3.0)

Just like `kubectl`, `eksclt` is aimed to be compliant with GitOps model, and can be used as part GitOps toolkit!

For example, you can use `eksctl apply --cluster-config prod-cluster.yaml`.

You can also use `eksctld`, which you'd normally run as a controller inside of another
cluster, you can manage multiple clusters this way.

## Usage

To create a basic cluster run:
```
eksctl create cluster
```
A cluster will be created with default parameters
- exciting auto-generated name, e.g. "fabulous-mushroom-1527688624"
- 2x `m5.large` nodes (this instance type suits most common use-cases, and is good value for money)
- default EKS AMI
- `us-west-2` region

To create the same kind of basic cluster, but with a different name run:
```
eksctl create cluster --cluster-name cluster-1 --nodes 4
```

To write cluster credentials to a file other then default, run:
```
eksctl create cluster --cluster-name cluster-2 --nodes 4 --kubeconfig ./kubeconfig.cluster-2.yaml
```

To prevent storing cluster credentials locally, run:
```
eksctl create cluster --cluster-name cluster-3 --nodes 4 --write-kubeconfig=false
```

To use 3-5 node ASG, run:
```
eksctl create cluster --cluster-name cluster-4 --nodes-min 3 --nodes-max 5
```

To use 30 `c4.xlarge` nodes, run:
```
eksctl create cluster --cluster-name cluster-5 --nodes 30 --node-type c4.xlarge
```

To delete a cluster, run:
```
eksctl delete cluster --cluster-name <name> [--region <region>]
```

<!-- TODO for 0.3.0
To use more advanced configuration options, [Cluster API](https://github.com/kubernetes-sigs/cluster-api):
```
eksctl apply --cluster-config advanced-cluster.yaml
```
-->
