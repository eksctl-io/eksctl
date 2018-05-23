# get-eks

> ***THIS IS FOR INTERNAL USE ONLY***
> #### What is the purpsose of this?
> 
> We may release some version of this to the public, but most of the thing you see here will be reviewed,
> so when you read this please don't try to review every detail of the presentation.
> 
> #### What do we have here?
>
> Right now we have a vary naive CLI solution, it's implemented in bash and I only spent one evening doing
> it, but it kind of does the job (in the most naive sense).
> We had to do this quickly, before EKS preview access is shut down. So we know the steps to make it work.
> 
> It exists in oder to:
> a) document (as code) current steps required to create an EKS cluster
> b) help us discuss CLI solution show and find the best route for a GA implementation
> c) help us consider how CLI solution shown can be translated to a Weave Cloud feature

Create an EKS cluster with one command, not dozens:

```
./create-cluster.sh
```

This tool wraps multiple CloudFormation steps into one, offering user a few most useful parameters, without exposing some of the less useful parameters.

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

### Improved design

To create a basic cluster run:

```
get-eks create --cluster-name cluster-1
```

It will be created in `us-west-2`, using default EKS AMI and 2 `m4.large` nodes. Name will be `cluster-1`.

To create the same kind of basic cluster, but with a different name run:

```
get-eks create --cluster-name cluster-2 --nodes 4
```

To write cluster credentials to a file other then default, run:

```
get-eks create --cluster-name cluster-2 --nodes 4 --kubeconfig ./kubeconfig.yaml
```

To prevent storing cluster credentials localy, run:

```
get-eks create --cluster-name cluster-2 --nodes 4 --write-kubeconfig=false
```

To use 3-5 node ASG, run:

```
get-eks create --cluster-name cluster-2 --nodes-min 3 --nodes-max 5
```

To use 30 `c4.xlarge` nodes, run:

```
get-eks create --cluster-name cluster-2 --nodes 30 --node-type c4.xlarge
```

To use more advanced configuration options, use [Cluster API](https://github.com/kubernetes-sigs/cluster-api)

```
get-eks apply --cluster-config advanced-cluster.yaml
```
