# `eksctl` - a CLI for Amazon EKS

[![Circle CI](https://circleci.com/gh/weaveworks/eksctl/tree/master.svg?style=shield)](https://circleci.com/gh/weaveworks/eksctl/tree/master) [![Coverage Status](https://coveralls.io/repos/github/weaveworks/eksctl/badge.svg?branch=master)](https://coveralls.io/github/weaveworks/eksctl?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/weaveworks/eksctl)](https://goreportcard.com/report/github.com/weaveworks/eksctl)

`eksctl` is a simple CLI tool for creating clusters on EKS - Amazon's new managed Kubernetes service for EC2. It is written in Go, and uses CloudFormation.

You can create a cluster in minutes with just one command – **`eksctl create cluster`**!

![Gophers: E, K, S, C, T, & L](logo/eksctl.png)

*Need help? Join [Weave Community Slack][slackjoin].*

## Usage

To download the latest release, run:

```
curl --silent --location "https://github.com/weaveworks/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin
```

Alternatively, macOS users can use [Homebrew](https://brew.sh):
```
brew tap weaveworks/tap
brew install weaveworks/tap/eksctl
```

and Windows users can use [chocolatey](https://chocolatey.org):
```
chocolatey install eksctl
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
[ℹ]  using region us-west-2
[ℹ]  setting availability zones to [us-west-2a us-west-2c us-west-2b]
[ℹ]  subnets for us-west-2a - public:192.168.0.0/19 private:192.168.96.0/19
[ℹ]  subnets for us-west-2c - public:192.168.32.0/19 private:192.168.128.0/19
[ℹ]  subnets for us-west-2b - public:192.168.64.0/19 private:192.168.160.0/19
[ℹ]  nodegroup "ng-98b3b83a" will use "ami-05ecac759c81e0b0c" [AmazonLinux2/1.11]
[ℹ]  creating EKS cluster "floral-unicorn-1540567338" in "us-west-2" region
[ℹ]  will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup
[ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=us-west-2 --name=floral-unicorn-1540567338'
[ℹ]  2 sequential tasks: { create cluster control plane "floral-unicorn-1540567338", create nodegroup "ng-98b3b83a" }
[ℹ]  building cluster stack "eksctl-floral-unicorn-1540567338-cluster"
[ℹ]  deploying stack "eksctl-floral-unicorn-1540567338-cluster"
[ℹ]  buildings nodegroup stack "eksctl-floral-unicorn-1540567338-nodegroup-ng-98b3b83a"
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

EKS supports versions `1.10`, `1.11` (default) and `1.12`, with `eksctl` you can deploy either version by passing `--version`.

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

> NOTE: You will still need to install and configure autoscaling. See the "Enable Autoscaling" section below. Also note that depending on your workloads you might need to use a separate nodegroup for each AZ. See [Zone-aware Autoscaling](#zone-aware-autoscaling) below for more info.

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


### Managing nodegroups

You can add one or more nodegroups in addition to the initial nodegroup created along with the cluster.

To create an additional nodegroup, use:

```
eksctl create nodegroup --cluster=<clusterName> [--name=<nodegroupName>]
```

> NOTE: By default, new nodegroups inherit the version from the control plane (`--version=auto`), but you can specify a different
> version e.g. `--version=1.10`, you can also use `--version=latest` to force use of whichever is the latest version.

Additionally, you can use the same config file used for `eksctl create cluster`:

```
eksctl create nodegroup --config-file=<path>
```

If there are multiple nodegroups specified in the file, you can select
a subset via `--only=<glob,glob,...>`:

```
eksctl create nodegroup --config-file=<path> --only='ng-prod-*-??'
```

To list the details about a nodegroup or all of the nodegroups, use:

```
eksctl get nodegroup --cluster=<clusterName> [--name=<nodegroupName>]
```

A nodegroup can be scaled by using the `eksctl scale nodegroup` command:

```
eksctl scale nodegroup --cluster=<clusterName> --nodes=<desiredCount> --name=<nodegroupName>
```

For example, to scale nodegroup `ng-a345f4e1` in `cluster-1` to 5 nodes, run:

```
eksctl scale nodegroup --cluster=cluster-1 --nodes=5 ng-a345f4e1
```

If the desired number of nodes is greater than the current maximum set on the ASG then the maximum value will be increased to match the number of requested nodes. And likewise for the minimum.

Scaling a nodegroup works by modifying the nodegroup CloudFormation stack via a ChangeSet.

> NOTE: Scaling a nodegroup down/in (i.e. reducing the number of nodes) may result in errors as we rely purely on changes to the ASG. This means that the node(s) being removed/terminated aren't explicitly drained. This may be an area for improvement in the future.

You can also enable SSH, ASG access and other feature for each particular nodegroup, e.g.:

```
eksctl create nodegroup --cluster=cluster-1 --node-labels="autoscaling=enabled,purpose=ci-worker" --asg-access --full-ecr-access --ssh-access
```

To delete a nodegroup, run:
```
eksctl delete nodegroup --cluster=<clusterName> --name=<nodegroupName>
```
> NOTE: this will drain all pods from that nodegroup before the instances are deleted.

All nodes are cordoned and all pods are evicted from a nodegroup on deletion,
but if you need to drain a nodegroup without deleting it, run:
```
eksctl drain nodegroup --cluster=<clusterName> --name=<nodegroupName>
```

To uncordon a nodegroup, run:
```
eksctl drain nodegroup --cluster=<clusterName> --name=<nodegroupName> --undo
```

### Cluster Upgrades

An _`eksctl`-managed_ cluster can be upgraded in 3 easy steps:

1. update control plane version with `eksctl update cluster`
2. replace each of the nodegroups by creating a new one and deleting the old one
3. update default add-ons:
  - `kube-proxy`
  - `aws-node`
  - `kube-dns` or `coredns`

Please make sure to read this section in full before you proceed.

> NOTE: Kubernetes supports version drift of up-to 2 minor versions during upgrade process.

### Updating control plane version

Control plane version updates must be done for one minor version at a time.

To update control plane to the next available version run:
```
eksctl update cluster --name=<clusterName>
```

This command will not apply any changes right away, you will need to re-run it with
`--dry-run=false` to apply the changes.

#### Updating nodegroups

You should update nodegroups only after you ran `eksctl update cluster`.

If you have a simple cluster with just an initial nodegroup (i.e. created with
`eksctl create cluster`), the process is very simple.

Get the name of old nodegroup:
```
eksctl get nodegroups --cluster=<clusterName>
```
> NOTE: you should see only one nodegroup here, if you see more - read the next section

Create new nodegroup:
```
eksctl create nodegroup --cluster=<clusterName>
```

Delete old nodegroup:
```
eksctl delete nodegroup --cluster=<clusterName> --name=<oldNodeGroupName>
```
> NOTE: this will drain all pods from that nodegroup before the instances are deleted.

##### Updating multiple nodegroups

If you have multiple nodegroups, it's your responsibility to track how each one was configured.
You can do this by using config files, but if you haven't used it already, you will need to inspect
your cluster to find out how each nodegroup was configured.

In general terms, you are looking to:
- review nodegroups you have and which ones can be deleted or must be replaced for the new version
- note down configuration of each nodegroup, consider using config file to ease upgrades next time

To create a new nodegroup:
```
eksctl create nodegroup --cluster=<clusterName> --name=<newNodeGroupName>
```

To delete old nodegroup:
```
eksctl delete nodegroup --cluster=<clusterName> --name=<oldNodeGroupName>
```

##### Updating multiple nodegroups with config file

If you are using config file, you will need to do the following.

Get the list of old nodegroups:
```
old_nodegroups="$(eksctl get ng --cluster=<clusterName> --output=json | jq -r '.[].Name')"
```

Edit config file to add new nodegroups. You can remove old nodegroups from the config file now
or do it later.

To create all of new nodegroups defined in the config file, run:
```
eksctl create nodegroup --config-file=<path>
```

Once you have new nodegroups in place, you can delete old ones:
```
for ng in $old_nodegroups ; do eksctl delete nodegroup --cluster=<clusterName> --name=$ng ; done
```
> NOTE: all pods will be drained.

#### Updating default add-ons

There are 3 default add-ons that get included in each EKS cluster, the process for updating each of them is different, hence
there are 3 distinct commands that you will need to run.

> NOTE: all of the following commands accept `--config-file`.

To update `kube-proxy`, run:
```
eksctl utils update-kube-proxy
```

To update `aws-node`, run:
```
eksctl utils update-aws-node
```

If you have upgraded from 1.10 to 1.11, you will need to replace `kube-dns` with `coredns`.
To do that, run:
```
eksctl utils install-coredns
```

If you have upgraded from 1.11 to 1.12, run:
```
eksctl utils update-coredns
```

Once upgraded, be sure to run `kubectl get pods -n kube-system` and check if all addon pods are in ready state, you should see
something like this:
```
NAME                       READY   STATUS    RESTARTS   AGE
aws-node-g5ghn             1/1     Running   0          2m
aws-node-zfc9s             1/1     Running   0          2m
coredns-7bcbfc4774-g6gg8   1/1     Running   0          1m
coredns-7bcbfc4774-hftng   1/1     Running   0          1m
kube-proxy-djkp7           1/1     Running   0          3m
kube-proxy-mpdsp           1/1     Running   0          3m
```

### Enable Autoscaling

You can create a cluster (or nodegroup in an existing cluster) with IAM role that will allow use of [cluster autoscaler][]:

```
eksctl create cluster --asg-access
```

Once cluster is running, you will need to install [cluster autoscaler][] itself. This flag also sets `k8s.io/cluster-autoscaler/enabled`
and `k8s.io/cluster-autoscaler/<clusterName>` tags, so nodegroup discovery should work.

[cluster autoscaler]: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md

#### Zone-aware Autoscaling

If your workloads are zone-specific you'll need to create separate nodegroups for each zone.  This is because the `cluster-autoscaler` assumes that all nodes in a group are exactly equivalent.  So, for example, if a scale-up event is triggered by a pod which needs a zone-specific PVC (e.g. an EBS volume), the new node might get scheduled in the wrong AZ and the pod will fail to start.

You won't need a separate nodegroup for each AZ if your environment meets the following criteria:

- No zone-specific storage requirements.
- No required podAffinity with topology other than host.
- No required nodeAffinity on zone label.
- No nodeSelector on a zone label.

(Read more [here](https://github.com/kubernetes/autoscaler/pull/1802#issuecomment-474295002) and [here](https://github.com/weaveworks/eksctl/pull/647#issuecomment-474698054).)

If you meet all of the above requirements (and possibly others) then you should be safe with a single nodegroup which spans multiple AZs.  Otherwise you'll want to create separate, single-AZ nodegroups:

BEFORE:

```yaml
nodeGroups:
  - name: ng1-public
    instanceType: m5.xlarge
    # availabilityZones: ["eu-west-2a", "eu-west-2b"]
```

AFTER:

```yaml
nodeGroups:
  - name: ng1-public-2a
    instanceType: m5.xlarge
    availabilityZones: ["eu-west-2a"]
  - name: ng1-public-2b
    instanceType: m5.xlarge
    availabilityZones: ["eu-west-2b"]
```

### VPC Networking

By default, `eksctl create cluster` will build a dedicated VPC, in order to avoid interference with any existing resources for a
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

### Using Config Files

You can create a cluster using a config file instead of flags.

First, create `cluster.yaml` file:
```YAML
apiVersion: eksctl.io/v1alpha4
kind: ClusterConfig

metadata:
  name: basic-cluster
  region: eu-north-1

nodeGroups:
  - name: ng-1
    instanceType: m5.large
    desiredCapacity: 10
```

Next, run this command:
```
eksctl create cluster -f cluster.yaml
```

This will create a cluster as described.

If you needed to use an existing VPC, you can use a config file like this:
```YAML
apiVersion: eksctl.io/v1alpha4
kind: ClusterConfig

metadata:
  name: cluster-in-existing-vpc
  region: eu-north-1

vpc:
  subnets:
    private:
      eu-north-1a: {id: subnet-0ff156e0c4a6d300c}
      eu-north-1b: {id: subnet-0549cdab573695c03}
      eu-north-1c: {id: subnet-0426fb4a607393184}

nodeGroups:
  - name: ng-1-workers
    labels: {role: workers}
    instanceType: m5.xlarge
    desiredCapacity: 10
    privateNetworking: true
  - name: ng-2-builders
    labels: {role: builders}
    instanceType: m5.2xlarge
    desiredCapacity: 2
    privateNetworking: true
    iam:
      withAddonPolicies:
        imageBuilder: true
```

To delete this cluster, run:
```
eksctl delete cluster -f cluster.yaml
```

See [`examples/`](https://github.com/weaveworks/eksctl/tree/master/examples) directory for more sample config files.

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
| static       | Indicates that the AMI images ids embedded into `eksctl` should be used. This relates to the static resolvers. |
| auto        | Indicates that the AMI to use for the nodes should be found by querying AWS. This relates to the auto resolver. |

If, for example, AWS release a new version of the EKS node AMIs and a new version of `eksctl` hasn't been released you can use the latest AMI by doing the following:

```
eksctl create cluster --node-ami=auto
```

With the 0.1.9 release we have introduced the `--node-ami-family` flag for use when creating the cluster. This makes it possible to choose between different officially supported EKS AMI families.

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

To make the above persistent, run the first two lines, and put the above in `~/.zshrc`.


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

Code contributions are very welcome. If you are interested in helping make `eksctl` great then see our [contributing guide](CONTRIBUTING.md).

## Get in touch

[Create an issue](https://github.com/weaveworks/eksctl/issues/new), or login to [Weave Community Slack (#eksctl)][slackchan] ([signup][slackjoin]).

[slackjoin]: https://slack.weave.works/
[slackchan]: https://weave-community.slack.com/messages/CAYBZBWGL/

> ***Logo Credits***
>
> *Original Gophers drawn by [Ashley McNamara](https://twitter.com/ashleymcnamara), unique E, K, S, C, T & L Gopher identities had been produced with [Gopherize.me](https://gopherize.me/).*
