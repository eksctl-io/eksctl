---
title: "VPC Networking"
weight: 60
url: usage/vpc-networking
---

## VPC Networking

By default, `eksctl create cluster` will build a dedicated VPC, in order to avoid interference with any existing resources for a
variety of reasons, including security, but also because it's challenging to detect all the settings in an existing VPC.
Default VPC CIDR used by `eksctl` is `192.168.0.0/16`, it is divided into 8 (`/19`) subnets (3 private, 3 public & 2 reserved).
Initial nodegroup is create in public subnets, with SSH access disabled unless `--allow-ssh` is specified. However, this implies
that each of the EC2 instances in the initial nodegroup gets a public IP and can be accessed on ports 1025 - 65535, which is
not insecure in principle, but some compromised workload could risk an access violation.

If that functionality doesn't suit you, the following options are currently available.

### Change VPC CIDR

If you need to setup peering with another VPC, or simply need larger or smaller range of IPs, you can use `--vpc-cidr` flag to
change it. You cannot use just any sort of CIDR, there only certain ranges that can be used in [AWS VPC][vpcsizing].

[vpcsizing]: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html#VPC_Sizing

### Use private subnets for initial nodegroup

If you prefer to isolate initial nodegroup from the public internet, you can use `--node-private-networking` flag.
When used in conjunction with `--ssh-access` flag, SSH port can only be accessed inside the VPC.

### Use existing VPC: shared with kops

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

### Use existing VPC: any custom configuration

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

### Custom Cluster DNS address

There are two ways of overwriting the DNS server IP address used for all the internal and external DNs lookups (this
is, the equivalent of the `--cluster-dns` flag for the `kubelet`).

The first, is through the `clusterDNS` field. [Config files](../schema) accept a `string` field called
`clusterDNS` with the IP address of the DNS server to use.
This will be passed to the `kubelet` that in turn will pass it to the pods through the `/etc/resolv.conf` file.

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-north-1

nodeGroups:
- name: ng-1
  clusterDNS: 169.254.20.10
```

Note that this configuration only accepts one IP address. To specify more than one address, use the
[`extraKubeletConfig` parameter](../customizing-the-kubelet):

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-north-1

nodeGroups:
  - name: ng-1
    kubeletExtraConfig:
        clusterDNS: ["169.254.20.10","172.20.0.10"]
```

### NAT Gateway

The NAT Gateway for a cluster can be configured to be `Disabled`, `Single` (default) or `HighlyAvailable`. It can be
specified through the `--vpc-nat-mode` CLI flag or in the cluster config file like the example below:


```yaml
vpc:
  nat:
    gateway: HighlyAvailable # other options: Disable, Single (default)
```

See the complete example [here](https://github.com/weaveworks/eksctl/blob/master/examples/09-nat-gateways.yaml).

**Note**: Specifying the NAT Gateway is only supported during cluster creation and it is not touched during a cluster
upgrade. There are plans to support changing between different modes on cluster update in the future.

### Managing Access to the Kubernetes API Server Endpoints

The default creation of an EKS cluster exposes the Kubernetes API server publicly but not directly
from within the
VPC subnets (Public=true, Private=false). Traffic destined for the API server from within the VPC
must first exit
the VPC networks (but not Amazon's network) and then re-enter to reach the API server. Allowing direct access from
within the VPC subnets is possible but not without some caveats. See [Amazon EKS Cluster
Endpoint Access Control](https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html) for
further information
on additional steps needed when public=false, private=true.

The Kubernetes API server endpoint access for a cluster can be cofigured for public and private
access using the cluster config file. Example below:

```yaml
vpc:
  endpointAccess:
    public:  <true|false>
    private: <true|false>
```
