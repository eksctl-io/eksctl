# `eksctl` - a CLI for Amazon EKS

[![Circle CI](https://circleci.com/gh/weaveworks/eksctl/tree/master.svg?style=shield)](https://circleci.com/gh/weaveworks/eksctl/tree/master) [![Coverage Status](https://coveralls.io/repos/github/weaveworks/eksctl/badge.svg?branch=master)](https://coveralls.io/github/weaveworks/eksctl?branch=master)[![Go Report Card](https://goreportcard.com/badge/github.com/weaveworks/eksctl)](https://goreportcard.com/report/github.com/weaveworks/eksctl)

`eksctl` is a simple CLI tool for creating clusters on EKS - Amazon's new managed Kubernetes service for EC2. It is written in Go, and uses CloudFormation.

You can create a cluster in minutes with just one command – **`eksctl create cluster`**!

![Gophers: E, K, S, C, T, & L](logo/eksctl.png)

## Usage

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

Once you have created a cluster, you will find that cluster credentials were added in `~/.kube/config`. If you have `kubectl` v1.10.x as well as `aws-iam-authenticator` commands in your PATH, you should be
able to use `kubectl`. You will need to make sure to use the same AWS API credentials for this also. Check [EKS docs][ekskubectl] for instructions. If you installed `eksctl` via Homebrew, you should have all of these dependencies installed already.

[ekskubectl]: https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html

Example output:
```
$ eksctl create cluster
2018-10-26T16:22:17+01:00 [ℹ]  using region us-west-2
2018-10-26T16:22:19+01:00 [ℹ]  setting availability zones to [us-west-2a us-west-2b us-west-2c]
2018-10-26T16:22:19+01:00 [ℹ]  subnets for us-west-2a - public:192.168.0.0/19 private:192.168.96.0/19
2018-10-26T16:22:19+01:00 [ℹ]  subnets for us-west-2b - public:192.168.32.0/19 private:192.168.128.0/19
2018-10-26T16:22:19+01:00 [ℹ]  subnets for us-west-2c - public:192.168.64.0/19 private:192.168.160.0/19
2018-10-26T16:22:19+01:00 [ℹ]  using "ami-0a54c984b9f908c81" for nodes
2018-10-26T16:22:19+01:00 [ℹ]  creating EKS cluster "floral-unicorn-1540567338" in "us-west-2" region
2018-10-26T16:22:19+01:00 [ℹ]  will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup
2018-10-26T16:22:19+01:00 [ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=us-west-2 --name=floral-unicorn-1540567338'
2018-10-26T16:22:19+01:00 [ℹ]  creating cluster stack "eksctl-floral-unicorn-1540567338-cluster"
2018-10-26T16:33:03+01:00 [ℹ]  creating nodegroup stack "eksctl-floral-unicorn-1540567338-nodegroup-0"
2018-10-26T16:36:44+01:00 [✔]  all EKS cluster resource for "floral-unicorn-1540567338" had been created
2018-10-26T16:36:44+01:00 [✔]  saved kubeconfig as "/Users/ilya/.kube/config"
2018-10-26T16:36:46+01:00 [ℹ]  the cluster has 0 nodes
2018-10-26T16:36:46+01:00 [ℹ]  waiting for at least 2 nodes to become ready
2018-10-26T16:37:22+01:00 [ℹ]  the cluster has 2 nodes
2018-10-26T16:37:22+01:00 [ℹ]  node "ip-192-168-25-215.us-west-2.compute.internal" is ready
2018-10-26T16:37:22+01:00 [ℹ]  node "ip-192-168-83-60.us-west-2.compute.internal" is ready
2018-10-26T16:37:23+01:00 [ℹ]  kubectl command should work with "~/.kube/config", try 'kubectl get nodes'
2018-10-26T16:37:23+01:00 [✔]  EKS cluster "floral-unicorn-1540567338" in "us-west-2" region is ready
$
```

To list the details about a cluster or all of the clusters, use:

```
eksctl get cluster [--name=<name>] [--region=<region>]
```

To create the same kind of basic cluster, but with a different name, run:

```
eksctl create cluster --name=cluster-1 --nodes=4
```

EKS supports two versions `1.10` and `1.11` (default), with `eksctl` you can deploy either version by passing `--version`.

```
eksctl create cluster --version=1.10
```

A default [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) (gp2 volume type provisioned by EBS) will be added automatically when creating a cluster.  If you want to prevent this, use the `--storage-class` flag.  For example:

```
eksctl create cluster --storage-class=false
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
eksctl create cluster --ssh-access --ssh-public-key=my_eks_node_id.pub
```

To use a pre-existing EC2 key pair in `us-east-1` region, you can specify key pair name (which must not resolve to a local file path), e.g. to use `my_kubernetes_key` run:

```
eksctl create cluster --ssh-access  --ssh-public-key=my_kubernetes_key --region=us-east-1
```

To add custom tags for all resources, use `--tags`.

> NOTE: Until [https://github.com/weaveworks/eksctl/issues/25] is resolved, tags cannot be applied to EKS cluster itself, but most of other resources (e.g. EC2 nodes).

```
eksctl create cluster --tags environment=staging --region=us-east-1
```

To configure node volume size, use the `--node-volume-size` flag.

```
eksctl create cluster --node-volume-size=50
```

> NOTE: In `us-east-1` you are likely to get `UnsupportedAvailabilityZoneException`. If you do, copy the suggested zones and pass `--zones` flag, e.g. `eksctl create cluster --region=us-east-1 --zones=us-east-1a,us-east-1b,us-east-1d`. This may occur in other regions, but less likely. You shouldn't need to use `--zone` flag otherwise.

To delete a cluster, run:

```
eksctl delete cluster --name=<name> [--region=<region>]
```

### Scaling nodegroup

The initial nodegroup can be scaled by using the `eksctl scale nodegroup` command. For example, to scale to 5 nodes:

```
eksctl scale nodegroup --name=<name> --nodes=5
```

If the desired number of nodes is greater than the current maximum set on the ASG then the maximum value will be increased to match the number of requested nodes. And likewise for the minimum.

Scaling a nodegroup works by modifying the nodegroup CloudFormation stack via a ChangeSet.

> NOTE: Scaling a nodegroup down/in (i.e. reducing the number of nodes) may result in errors as we rely purely on changes to the ASG. This means that the node(s) being removed/terminated aren't explicitly drained. This may be an area for improvement in the future.

### VPC Networking

By default, `eksctl create cluster` instatiates a dedicated VPC, in order to avoid interference with any existing resources for a
variety of reasons, including security, but also because it's challenging to detect all the settings in an existing VPC. 
Default VPC CIDR used by `eksctl` is `192.168.0.0/16`, it is divided into 8 (`/19`) subnets (3 private, 3 public & 2 reserved).
Initial nodegroup is create in public subnets, with SSH access disabled unless `--allow-ssh` is specified. However, this implies
that each of the EC2 instances in the initial nodegroup gets a public IP and can be accessed on ports 1025 - 65535, which is
not insecure in principle, but some compromised workload could risk an access violation.

If that functionality doesn't suit you, the following options are currently available.

#### change VPC CIDR

If you need to setup peering with another VPC, or simply need larger or smaller range of IPs, you can use `--vpc-cidr` flag to
change it. You cannot use just any sort of CIDR, there only certain ranges that can be used in [AWS VPC][vpcsizing].

[vpcsizing]: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html#VPC_Sizing

#### use private subnets for initial nodegroup

If you prefer to isolate initial nodegroup from the public internet, you can use `--node-private-networking` flag.
When used in conjunction with `--ssh-access` flag, SSH port can only be accessed inside the VPC.

#### use existing VPC: shared with kops

You can use a VPC of an existing Kubernetes cluster managed by kops. This feature is provided to facilitate migration and/or
cluster peering.

If you have previously created a cluster with kops, e.g. using commands similar to this:

```
export KOPS_STATE_STORE=s3://kops
kops create cluster cluster-1.k8s.local --zones=us-west-2c,us-west-2b,us-west-2a --networking=weave --yes
```

You can create an EKS cluster in the same AZs using the same VPC subnets (NOTE: at least 2 AZs/subnets are required):

```
eksctl create cluster --name=cluster-2 --region=us-west-2 --vpc-from-kops-cluster=cluster-1.k8s.local
```

#### use existing VPC: any custom configuration

Use this feature if you must configure a VPC in a way that's different to how dedicated VPC is configured by `eksctl`, or have to
use a VPC that already exists so your EKS cluster gets shared access to some resources inside that existing VPC, or you have any
other use-case that requires you to manage VPCs separately.

You can use an existing VPC by supplying private and/or public subnets using `--vpc-private-subnets` and `--vpc-public-subnets` flags.
It is up to you to ensure which subnets you use, as there is no simple way to determine automatically whether a subnets is private or
public, because configurations vary.
Given these flags, `eksctl create cluster` will determine the VPC ID automatically, but it will not create any routing tables or other
resources, such as internet/NAT gateways. It will, however, create dedicated security groups for the initial nodegroup and the control
plane.

You must ensure to provide at least 2 subnets in different AZs. There are other requirements that you will need to follow, but it's
entirely up to you to address those. For example, tagging is not strictly necessary, tests have shown that its possible to create
a functional cluster without any tags set on the subnets, however there is no guarantee that this will always hold and tagging is
recommended.

- all subnets in the same VPC, within the same block of IPs
- sufficient IP addresses are available
- sufficient number of subnets (minimum 2)
- internet and/or NAT gateways are configured correctly
- routing tables have correct entries and the network is functional
- tagging of subnets
  - `kubernetes.io/cluster/<name>` tag set to either `shared` or `owned`
  - `kubernetes.io/role/internal-elb` tag set to `1` for private subnets

There maybe other requirements imposed by EKS or Kubernetes, and it is entirely up to you to stay up-to-date on any requirements and/or
recommendations, and implement those as needed/possible.

Default security group settings applied by `eksctl` may or may not be sufficient for sharing access with resources in other security
groups. If you wish to modify the ingress/egress rules of the either of security groups, you might need to use another tool to automate
changes, or do it via EC2 console.

If you are in doubt, don't use a custom VPC. Using `eksctl create cluster` without any `--vpc-*` flags will always configure the cluster
with a fully-functional dedicated VPC.

To create a cluster using 2x private and 2x public subnets, run:

```
eksctl create cluster \
  --vpc-private-subnets=subnet-0ff156e0c4a6d300c,subnet-0426fb4a607393184 \
  --vpc-public-subnets=subnet-0153e560b3129a696,subnet-009fa0199ec203c37
```

To create a cluster using 3x private subnets and make initial nodegroup use those subnets, run:

```
eksctl create cluster \
  --vpc-private-subnets=subnet-0ff156e0c4a6d300c,subnet-0549cdab573695c03,subnet-0426fb4a607393184 \
  --node-private-networking
```

To create a cluster using 4x public subnets, run:

```
eksctl create cluster \
  --vpc-public-subnets=subnet-0153e560b3129a696,subnet-0cc9c5aebe75083fd,subnet-009fa0199ec203c37,subnet-018fa0176ba320e45
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

> NOTE: Once `addon` support has been added as part of 0.2.0 it is envisioned that there will be a addon to install the NVIDIA Kubernetes Device Plugin.  This addon could potentially be installed automatically as we know an GPU instance type is being used.

### Latest & Custom AMI Support

With the 0.1.2 release we have introduced the `--node-ami` flag for use when creating a cluster. This enables a number of advanced use cases such as using a custom AMI or querying AWS in realtime to determine which AMI to use (non-GPU and GPU instances).

The `--node-ami` can take the AMI image id for an image to explicitly use. It also can take the following 'special' keywords:

| Keyword | Description |
| ------------ | -------------- |
| static       | Indicates that the AMI images ids embedded into eksctl should be used. This relates to the static resolvers. |
| auto        | Indicates that the AMI to use for the nodes should be found by querying AWS. This relates to the auto resolver. |

If, for example, AWS release a new version of the EKS node AMIs and a new version of eksctl hasn't been released you can use the latest AMI by doing the following:

```
eksctl create cluster --node-ami=auto
```

With the 0.1.9 release we have introduced the `--node-ami-family` flag for use when creating the cluster. This makes it possible to choose between different offically supported EKS AMI families.

The `--node-ami-family` can take following keywords:

| Keyword | Description |
| --- | --- |
| AmazonLinux2 | Indicates that the EKS AMI image based on Amazon Linux 2 should be used. (default)|
| Ubuntu1804 | Indicates that the EKS AMI image based on Ubuntu 18.04 should be used. |

<!-- TODO for 0.3.0
To use more advanced configuration options, [Cluster API](https://github.com/kubernetes-sigs/cluster-api):

```
eksctl apply --cluster-config advanced-cluster.yaml
```
-->

### Shell Completion

To enable bash completion, run the following, or put it in `~/.bashrc` or `~/.profile`:
```
. <(eksctl completion bash)
```

Or for zsh, run:
```
mkdir -p ~/.zsh/completion/
eksctl completion zsh > ~/.zsh/completion/_eksctl
```
and put the following in `~/.zshrc`:
```
fpath=($fpath ~/.zsh/completion)
```
Note if you're not running a distribution like oh-my-zsh you may first have to enable autocompletion:
```
autoload -U compinit
compinit
```

To make the above persistent, run the first two lines, and put the


## Project Roadmap

### Developer use-case (0.2.0)

It should suffice to install a cluster for development with just a single command. Here are some examples:

To create a cluster with default configuration (2 `m5.large` nodes), run:

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
