# VPC Configuration

## Change VPC CIDR

If you need to set up peering with another VPC, or simply need a larger or smaller range of IPs, you can use `--vpc-cidr` flag to
change it. Please refer to [the AWS docs][vpcsizing] for guides on choosing CIDR blocks which are permitted for use in an AWS VPC.

[vpcsizing]: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html#VPC_Sizing

If you are creating an IPv6 cluster you can also bring your own IPv6 pool by configuring `VPC.IPv6Cidr` and `VPC.IPv6Pool`.
See [AWS docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html) on how to import your own pool.

## Use an existing VPC: shared with kops

You can use the VPC of an existing Kubernetes cluster managed by [kops](https://github.com/kubernetes/kops). This feature is provided to facilitate migration and/or cluster peering.

If you have previously created a cluster with kops, e.g. using commands similar to this:

```
export KOPS_STATE_STORE=s3://kops
kops create cluster cluster-1.k8s.local --zones=us-west-2c,us-west-2b,us-west-2a --networking=weave --yes
```

You can create an EKS cluster in the same AZs using the same VPC subnets (NOTE: at least 2 AZs/subnets are required):

```
eksctl create cluster --name=cluster-2 --region=us-west-2 --vpc-from-kops-cluster=cluster-1.k8s.local
```

## Use existing VPC: other custom configuration

`eksctl` provides some, but not complete, flexibility for custom VPC and subnet topologies.

You can use an existing VPC by supplying private and/or public subnets using the `--vpc-private-subnets` and `--vpc-public-subnets` flags.
It is up to you to ensure the subnets you use are categorised correctly, as there is no simple way to verify whether a subnet is actually private or
public, because configurations vary.

Given these flags, `eksctl create cluster` will determine the VPC ID automatically, but it will not create any routing tables or other
resources, such as internet/NAT gateways. It will, however, create dedicated security groups for the initial nodegroup and the control
plane.

You must ensure to provide **at least 2 subnets in different AZs**. There are other requirements that you will need to follow (listed below), but it's
entirely up to you to address those. (For example, tagging is not strictly necessary, tests have shown that it is possible to create
a functional cluster without any tags set on the subnets, however there is no guarantee that this will always hold and tagging is
recommended.)

Standard requirements:

- all given subnets must be in the same VPC, within the same block of IPs
- a sufficient number IP addresses are available, based on needs
- sufficient number of subnets (minimum 2), based on needs
- subnets are tagged with at least the following:
    - `kubernetes.io/cluster/<name>` tag set to either `shared` or `owned`
    - `kubernetes.io/role/internal-elb` tag set to `1` for _private_ subnets
    - `kubernetes.io/role/elb` tag set to `1` for _public_ subnets
- correctly configured internet and/or NAT gateways
- routing tables have correct entries and the network is functional
- **NEW**: all public subnets should have the property `MapPublicIpOnLaunch` enabled (i.e. `Auto-assign public IPv4 address` in the AWS console)

There may be other requirements imposed by EKS or Kubernetes, and it is entirely up to you to stay up-to-date on any requirements and/or
recommendations, and implement those as needed/possible.

Default security group settings applied by `eksctl` may or may not be sufficient for sharing access with resources in other security
groups. If you wish to modify the ingress/egress rules of the security groups, you might need to use another tool to automate
changes, or do it via EC2 console.

When in doubt, don't use a custom VPC. Using `eksctl create cluster` without any `--vpc-*` flags will always configure the cluster
with a fully-functional dedicated VPC.

**Examples**

Create a cluster using a custom VPC with 2x private and 2x public subnets:

```
eksctl create cluster \
  --vpc-private-subnets=subnet-0ff156e0c4a6d300c,subnet-0426fb4a607393184 \
  --vpc-public-subnets=subnet-0153e560b3129a696,subnet-009fa0199ec203c37
```

or use the following equivalent config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-test
  region: us-west-2

vpc:
  id: "vpc-11111"
  subnets:
    private:
      us-west-2a:
          id: "subnet-0ff156e0c4a6d300c"
      us-west-2c:
          id: "subnet-0426fb4a607393184"
    public:
      us-west-2a:
          id: "subnet-0153e560b3129a696"
      us-west-2c:
          id: "subnet-009fa0199ec203c37"

nodeGroups:
  - name: ng-1
```

Create a cluster using a custom VPC with 3x private subnets and make initial nodegroup use those subnets:

```
eksctl create cluster \
  --vpc-private-subnets=subnet-0ff156e0c4a6d300c,subnet-0549cdab573695c03,subnet-0426fb4a607393184 \
  --node-private-networking
```

or use the following equivalent config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-test
  region: us-west-2

vpc:
  id: "vpc-11111"
  subnets:
    private:
      us-west-2d:
          id: "subnet-0ff156e0c4a6d300c"
      us-west-2c:
          id: "subnet-0549cdab573695c03"
      us-west-2a:
          id: "subnet-0426fb4a607393184"

nodeGroups:
  - name: ng-1
    privateNetworking: true

```

Create a cluster using a custom VPC 4x public subnets:

```
eksctl create cluster \
  --vpc-public-subnets=subnet-0153e560b3129a696,subnet-0cc9c5aebe75083fd,subnet-009fa0199ec203c37,subnet-018fa0176ba320e45
```


```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-test
  region: us-west-2

vpc:
  id: "vpc-11111"
  subnets:
    public:
      us-west-2d:
          id: "subnet-0153e560b3129a696"
      us-west-2c:
          id: "subnet-0cc9c5aebe75083fd"
      us-west-2a:
          id: "subnet-009fa0199ec203c37"
      us-west-2b:
          id: "subnet-018fa0176ba320e45"

nodeGroups:
  - name: ng-1
```

More examples can be found in the repo's `examples` folder:

- [using an existing VPC](https://github.com/weaveworks/eksctl/blob/master/examples/04-existing-vpc.yaml)
- [using a custom VPC CIDR](https://github.com/weaveworks/eksctl/blob/master/examples/02-custom-vpc-cidr-no-nodes.yaml)

## Custom Shared Node Security Group

`eksctl` will create and manage a shared node security group that allows communication between
unmanaged nodes and the cluster control plane and managed nodes.

If you wish to provide your own custom security group instead, you may override the `sharedNodeSecurityGroup`
field in the config file:


```yaml
vpc:
  sharedNodeSecurityGroup: sg-0123456789
```

By default, when creating the cluster, `eksctl` will add rules to this security group to allow communication to and
from the default cluster security group that EKS creates. The default cluster security group is used by both
the EKS control plane and managed node groups.

If you wish to manage the security group rules yourself, you may prevent `eksctl` from creating the rules
by setting `manageSharedNodeSecurityGroupRules` to `false` in the config file:

```yaml
vpc:
  sharedNodeSecurityGroup: sg-0123456789
  manageSharedNodeSecurityGroupRules: false
```

## NAT Gateway

The NAT Gateway for a cluster can be configured to be `Disabled`, `Single` (default) or `HighlyAvailable`.
The `HighlyAvailable` option will deploy a NAT Gateway in each Availability Zone of the Region, so that if
an AZ is down, nodes in the other AZs will still be able to communicate to the Internet.

It can be specified through the `--vpc-nat-mode` CLI flag or in the cluster config file like the example below:


```yaml
vpc:
  nat:
    gateway: HighlyAvailable # other options: Disable, Single (default)
```

See the complete example [here](https://github.com/weaveworks/eksctl/blob/master/examples/09-nat-gateways.yaml).

**Note**: Specifying the NAT Gateway is only supported during cluster creation. It isn't touched during a cluster
upgrade. There are plans to support changing between different modes on cluster update in the future.
