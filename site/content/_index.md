---
title: Introduction
type: docs
---

<div id="github"><script async defer src="https://buttons.github.io/buttons.js"></script><a class="github-button" href="https://github.com/weaveworks/eksctl" data-icon="octicon-star" data-show-count="true" aria-label="Star weaveworks/eksctl on GitHub">Star</a><div id="github-section"><a class="github-button" href="https://github.com/weaveworks/eksctl/fork" data-icon="octicon-repo-forked" data-show-count="true" aria-label="Fork weaveworks/eksctl on GitHub">Fork</a></div></div>

# **`eksctl` - The official CLI for Amazon EKS**

### sponsored by [![Weaveworks](introduction/images/weaveworks.svg#inline-ww)](https://www.weave.works/) and built by [![Contributors](introduction/images/gophers.png#inline)](https://github.com/weaveworks/eksctl/graphs/contributors) on [![Github](introduction/images/octocat.svg#inline)](https://github.com/weaveworks/eksctl)

`eksctl` is a simple CLI tool for creating clusters on EKS - Amazon's new managed Kubernetes service for EC2. It is
written in Go, uses CloudFormation, was created by [Weaveworks](https://www.weave.works/) and it welcomes
contributions from the community. Create a basic cluster in minutes with just one command :

 <h2 id="intro-code">**`eksctl create cluster`**</h2>

A cluster will be created with default parameters

> ∙ exciting auto-generated name, e.g. "fabulous-mushroom-1527688624"

> ∙ 2x `m5.large` nodes (this instance type suits most common use-cases, and is good value for money)

> ∙ use official AWS EKS AMI

> ∙ `us-west-2` region

> ∙ dedicated VPC (check your quotas)

> ∙ using static AMI resolver

Example output:

```
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

_Need help? Join [Weave Community Slack][slackjoin]._
[slackjoin]: https://slack.weave.works/

Customize your cluster by using a config file. Just run

```
eksctl create cluster -f cluster.yaml
```

to apply a `cluster.yaml` file:

 <div id="styled-code">

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

</div>

Once you have created a cluster, you will find that cluster credentials were added in `~/.kube/config`. If you have
`kubectl` v1.10.x as well as `aws-iam-authenticator` commands in your PATH, you should be
able to use `kubectl`. You will need to make sure to use the same AWS API credentials for this also. Check
[EKS docs][ekskubectl] for instructions. If you installed `eksctl` via Homebrew, you should have all of these
dependencies installed already.

To learn more about how to create clusters and other features continue reading the
[usage section](usage/creating-and-managing-clusters/).

[ekskubectl]: https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html

[![Circle CI](https://circleci.com/gh/weaveworks/eksctl/tree/master.svg?style=shield)](https://circleci.com/gh/weaveworks/eksctl/tree/master) [![Coverage Status](https://coveralls.io/repos/github/weaveworks/eksctl/badge.svg?branch=master)](https://coveralls.io/github/weaveworks/eksctl?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/weaveworks/eksctl)](https://goreportcard.com/report/github.com/weaveworks/eksctl)
![Gophers: E, K, S, C, T, & L](introduction/images/eksctl.png)

> **_Credits_**
>
> _Original Gophers drawn by [Ashley McNamara](https://twitter.com/ashleymcnamara), unique E, K, S, C, T & L Gopher identities had been produced with [Gopherize.me](https://gopherize.me/)._
>
> Site powered by [Netlify ©](https://netlify.com)
