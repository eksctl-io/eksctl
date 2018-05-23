# eksctl

> ***CURRENT STATE: EAERLY PROTOTYPE***

What is `eksctl`? It's a simple CLI tool for creating EKS clusters, for most common use-cases.

It's inspired by `kubectl`, and one of the goals of the project is to implement a Cluster API controller for EKS.

It is not inteded to be a like-for-like alternative to well-established community tools (`kops`, `kubicorn`, `kubeadm`).
However, the intention is to work well with other popular tools.

## Developer use-case

It should suffice to install a cluster for development with just a single command, here are some examples.

To create a cluster with default configurations (2 `m4.large` nodes), run:
```
eksctl create cluster dev-cluster
```

It supposrts many popular addons, including:

* Weave Net: `eksctl create cluster dev-cluster --networking=weave`
* Helm: `eksctl create cluster dev-cluster --addons=helm`
* AWS CI tools (CodeCommit, CodeBuild, ECR): `eksctl create cluster dev-cluster --addons=aws-ci`
* AWS CodeStar: `eksctl create cluster dev-cluster --addons=aws-codestar`
* Weave Scope and Flux: `eksctl create cluster dev-cluster --addons=weave-scope,weave-flux`

You can combine any or all of these.

You can also add any of these addons after you create a cluster with `eksctl addons install <addon>...`.

## Manage EKS the GitOps way

Just like `kubectl`, `eksclt` is aimed to be compliant with GitOps model, and can be used as part GitOps toolkit!

For example, you can use `eksctl apply --cluster-config prod-cluster.yaml`.

You can also use `eksctld`, which you'd normaly run aa controller inside of another
cluster, you can manage multiple clusters this way.

## Current Design (prototype)

Usage: ***`./create-cluster.sh [<clusterName> [<numberOfNodes> [<nodeType>]]]`***

So to create a basic cluster run:

```
./create-cluster.sh
```

It will be created in `us-west-2`, using default EKS AMI and 2 `m4.large` nodes. Name will be `cluster-1`.

To create the same kind of basic cluster, but with a different name run:

```
./create-cluster.sh cluster-2
```

To use 3 nodes, run:

```
./create-cluster.sh cluster-2 3
```

To use 3 `c4.xlarge` nodes, run:

```
./create-cluster.sh cluster-2 3 c4.xlarge
```

Example output:

```console
 [0] >> ./create-cluster.sh cluster-2
Creating EKS-cluster-2-ServiceRole and EKS-cluster-2-VPC stacks we need first
{
    "StackId": "arn:aws:cloudformation:us-west-2:376248598259:stack/EKS-cluster-2-ServiceRole/909e04b0-5e5b-11e8-a5a3-50a68a0bca9a"
}
{
    "StackId": "arn:aws:cloudformation:us-west-2:376248598259:stack/EKS-cluster-2-VPC/918186e0-5e5b-11e8-80c5-503aca41a0fd"
}
Waiting until the EKS-cluster-2-ServiceRole and EKS-cluster-2-VPC stacks are ready
Collect outputs from the EKS-cluster-2-ServiceRole and EKS-cluster-2-VPC stacks
Creating cluster cluster-2
{
    "cluster": {
        "clusterName": "cluster-2",
        "clusterArn": "arn:aws:eks:us-west-2:376248598259:cluster/cluster-2",
        "createdAt": 1527060875149000,
        "desiredMasterVersion": "1.10",
        "roleArn": "arn:aws:iam::376248598259:role/EKS-cluster-2-ServiceRole-AWSServiceRoleForAmazonE-7NS9V7ERKDXO",
        "subnets": [
            "subnet-f3b009b8",
            "subnet-9f3aa6e6"
        ],
        "securityGroups": [
            "sg-2976a258"
        ],
        "status": "NEW",
        "certificateAuthority": {}
    }
}
Creating EKS-cluster-2-DefaultNodeGroup stack
{
    "StackId": "arn:aws:cloudformation:us-west-2:376248598259:stack/EKS-cluster-2-DefaultNodeGroup/bece5bf0-5e5b-11e8-9b25-50a68d01a68d"
}
Waiting until cluster is ready
Saving cluster credentials in /Users/ilya/Code/eks-preview/get-eks/cluster-2.us-west-2.yaml
Waiting until EKS-cluster-2-DefaultNodeGroup stack is ready
configmap "aws-auth" created
Cluster is ready, nodes will be added soon
Use the following command to monitor the nodes
$ kubectl --kubeconfig='/Users/ilya/Code/eks-preview/get-eks/cluster-2.us-west-2.yaml' get nodes --watch
 [0] >>
```

## Limitations

- Written in bash
- kubectl and heptio-authenticator-aws binaries are vendored in the repo
- Doesn't handle most errors
- Doesn't offer parameters for important things (like region, AMI, node SSH key)
- Cannot use custom VPC or customise networking in any way
- Manual deletion

## Various notes

- Rewrite in Go (or maybe Python, as AWS CLI extension)
- Use named flags instead of postional arguments
- Use Cluster API for the sake of GitOps etc (initially CLI only, later offer a controller)
- Single CloudFormation template (nested stack)
- Call home (and mention in the readme) - time, cluster type, regions, IP (or hash of) [no need to count deletions]
- Add short-cuts for Weave Net (most certainly) and Weave Cloud (maybe)
- Consider repurposing kops (or even kubicorn), or some of its code (it may be easier to use the AWS API the way kops does, instead of CloudFormation - TBD, but kops node bootstrap code may not be very useful)
- On EKS GA date Terraform module for EKS will be available – perhaps try it
- Find partners and contributors (e.g. Jenkins X and/or Heptio)
- Could persuade Docker to work on LinuxKit node AMIs
- Node upgrade controller
- Consider kubeadm join

### Improved design – MVP

To create a basic cluster run:
```
eksctl create cluster --cluster-name cluster-1
```
It will be created in `us-west-2`, using default EKS AMI and 2 `m4.large` nodes. Name will be `cluster-1`.

To create the same kind of basic cluster, but with a different name run:
```
eksctl create cluster --cluster-name cluster-2 --nodes 4
```

To write cluster credentials to a file other then default, run:
```
eksctl create cluster --cluster-name cluster-3 --nodes 4 --kubeconfig ./kubeconfig.yaml
```

To prevent storing cluster credentials localy, run:
```
eksctl create cluster --cluster-name cluster-4 --nodes 4 --write-kubeconfig=false
```

To use 3-5 node ASG, run:
```
eksctl create cluster --cluster-name cluster-5 --nodes-min 3 --nodes-max 5
```

To use 30 `c4.xlarge` nodes, run:
```
eksctl create cluster --cluster-name cluster-6 --nodes 30 --node-type c4.xlarge
```

To delete a cluster, run:
```
eksctl delete cluster <name>
```

To use more advanced configuration options, use [Cluster API](https://github.com/kubernetes-sigs/cluster-api):
```
eksctl apply --cluster-config=advanced-cluster.yaml
```
