# VPC Networking

By default `eksctl create cluster` will create a dedicated VPC for the cluster.
This is done in order to avoid interference with existing resources for a
variety of reasons, including security, but also because it is challenging to detect all settings in an existing VPC.

The default VPC CIDR used by `eksctl` is `192.168.0.0/16`. It is divided into 8 (`/19`) subnets (3 private, 3 public & 2 reserved).
The initial nodegroup is created in public subnets, with SSH access disabled unless `--allow-ssh` is specified.
The nodegroup by default allows inbound traffic from the control plane security group on ports 1025 - 65535.

!!! note
    In `us-east-1` eksctl only creates 2 public and 2 private subnets by default.

!!! important
    From `eksctl` version `0.17.0` and onwards public subnets will have the property `MapPublicIpOnLaunch` enabled, and
    the property `AssociatePublicIpAddress` disabled in the Auto Scaling Group for the nodegroups. This means that when
    creating a **new nodegroup** on a **cluster made with an earlier version** of `eksctl`, the nodegroup must **either** be private
    **or** have `MapPublicIpOnLaunch` enabled in its public subnets. Without one of these, the new nodes won't have access to
    the internet and won't be able to download the basic add-ons (CNI plugin, kube-proxy, etc). To help set up
    subnets correctly for old clusters you can use the new command `eksctl utils update-legacy-subnet-settings`.

If the default functionality doesn't suit you, the following options are currently available.

## Change VPC CIDR

If you need to setup peering with another VPC, or simply need a larger or smaller range of IPs, you can use `--vpc-cidr` flag to
change it. Please refer to [the AWS docs][vpcsizing] for guides on choosing CIDR blocks which are permitted for use in an AWS VPC.

[vpcsizing]: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html#VPC_Sizing

## Use private subnets for initial nodegroup

If you prefer to isolate the initial nodegroup from the public internet, you can use the `--node-private-networking` flag.
When used in conjunction with the `--ssh-access` flag, the SSH port can only be accessed from inside the VPC.

!!! note
    Using the `--node-private-networking` flag will result in outgoing traffic to go through the NAT gateway using its
    Elastic IP. On the other hand, if the nodes are in a public subnet, the outgoing traffic won't go through the
    NAT gateway and hence the outgoing traffic has the IP of each individual node.

## Use an existing VPC: shared with kops

You can use the VPC of an existing Kubernetes cluster managed by [kops](https://github.com/kubernetes/kops). This feature is provided to facilitate migration and/or
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

Further examples can be found in the repo's `examples` dir:

- [using an existing VPC](https://github.com/weaveworks/eksctl/blob/master/examples/04-existing-vpc.yaml)
- [using a custom VPC CIDR](https://github.com/weaveworks/eksctl/blob/master/examples/02-custom-vpc-cidr-no-nodes.yaml)


### Custom subnet topology

`eksctl` version `0.32.0` introduced further subnet topology customisation with the ability to:

- List multiple subnets per AZ in VPC configuration
- Specify subnets in nodegroup configuration

In earlier versions custom subnets had to be provided by availability zone, meaning just one subnet per AZ could be listed.
From `0.32.0` the identifying keys can be arbitrary.

```yaml
vpc:
  id: "vpc-11111"
  subnets:
    public:
      public-one:                           # arbitrary key
          id: "subnet-0153e560b3129a696"
      public-two:
          id: "subnet-0cc9c5aebe75083fd"
      us-west-2b:                           # or list by AZ
          id: "subnet-018fa0176ba320e45"
    private:
      private-one:
          id: "subnet-0153e560b3129a696"
      private-two:
          id: "subnet-0cc9c5aebe75083fd"
```

!!! important
    If using the AZ as the identifying key, the `az` value can be omitted.

    If using an arbitrary string as the identifying key, like above, either:

	* `id` must be set (`az` and `cidr` optional)
	* or `az` must be set (`cidr` optional)

	If a user specifies a subnet by AZ without specifying CIDR and ID, a subnet
	in that AZ will be chosen from the VPC, arbitrarily if multiple such subnets
	exist.

!!! note
    A complete subnet spec must be provided, ie. both `public` and `private` configurations
    declared in the VPC spec.

Nodegroups can be restricted to named subnets via the configuration.
When specifying subnets on nodegroup configuration, use the identifying key as given in the VPC spec **not** the subnet id.
For example:

```yaml
vpc:
  id: "vpc-11111"
  subnets:
    public:
      public-one:
          id: "subnet-0153e560b3129a696"
    ... # subnet spec continued

nodeGroups:
  - name: ng-1
    instanceType: m5.xlarge
    desiredCapacity: 2
    subnets:
      - public-one
```

!!! note
    Only one of `subnets` or `availabilityZones` can be provided in nodegroup configuration.

When placing nodegroups inside a private subnet, `privateNetworking` must be set to `true`
on the nodegroup:

```yaml
vpc:
  id: "vpc-11111"
  subnets:
    public:
      private-one:
          id: "subnet-0153e560b3129a696"
    ... # subnet spec continued

nodeGroups:
  - name: ng-1
    instanceType: m5.xlarge
    desiredCapacity: 2
    privateNetworking: true
    subnets:
      - private-one
```

See [here](https://github.com/weaveworks/eksctl/blob/master/examples/24-nodegroup-subnets.yaml) for a full
configuration example.

## Custom Cluster DNS address

There are two ways of overwriting the DNS server IP address used for all the internal and external DNS lookups. This
is the equivalent of the `--cluster-dns` flag for the `kubelet`.

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
[`kubeletExtraConfig` parameter](../customizing-the-kubelet):

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

**Note**: Specifying the NAT Gateway is only supported during cluster creation and it is not touched during a cluster
upgrade. There are plans to support changing between different modes on cluster update in the future.

## Managing Access to the Kubernetes API Server Endpoints

The default creation of an EKS cluster exposes the Kubernetes API server publicly but not directly from within the
VPC subnets (public=true, private=false). Traffic destined for the API server from within the VPC must first exit the
VPC networks (but not Amazon's network) and then re-enter to reach the API server.

The Kubernetes API server endpoint access for a cluster can be configured for public and private access when creating
the cluster using the cluster config file. Example below:

```yaml
vpc:
  clusterEndpoints:
    publicAccess:  <true|false>
    privateAccess: <true|false>
```

There are some additional caveats when configuring Kubernetes API endpoint access:

1. EKS doesn't allow one to create or update a cluster without at least one of private or public access being
   enabled.
1. EKS does allow creating a configuration which allows only private access to be enabled, but eksctl doesn't
   support it during cluster creation as it prevents eksctl from being able to join the worker nodes to the cluster.
1. Updating a cluster to have private only Kubernetes API endpoint access means that Kubernetes commands, by default,
   (e.g. `kubectl`) as well as `eksctl delete cluster`, `eksctl utils write-kubeconfig`, and possibly the command
   `eksctl utils update-kube-proxy` must be run within the cluster VPC.  This requires some changes to various AWS
   resources.  See:
   [EKS user guide](https://docs.aws.amazon.com/en_pv/eks/latest/userguide/cluster-endpoint#private-access)
   A user can elect to supply vpc.extraCIDRs which will append additional CIDR ranges to the ControlPlaneSecurityGroup,
   allowing subnets outside the VPC to reach the kubernetes API endpoint.

The following is an example of how one could configure the Kubernetes API endpoint access using the `utils` sub-command:

```
eksctl utils update-cluster-endpoints --name=<clustername> --private-access=true --public-access=false
```

To update the setting using a `ClusterConfig` file, use:

```console
eksctl utils update-cluster-endpoints -f config.yaml --approve
```

Note that if you don't pass a flag in it will keep the current value. Once you are satisfied with the proposed changes,
add the `approve` flag to make the change to the running cluster.

## Restricting Access to the EKS Kubernetes Public API endpoint

The default creation of an EKS cluster exposes the Kubernetes API server publicly. To restrict access to the public API
endpoint to a set of CIDRs when creating a cluster, set the `publicAccessCIDRs` field:

```yaml
vpc:
  publicAccessCIDRs: ["1.1.1.1/32", "2.2.2.0/24"]
```

To update the restrictions on an existing cluster, use:

```console
eksctl utils set-public-access-cidrs --cluster=<cluster> 1.1.1.1/32,2.2.2.0/24
```

To update the restrictions using a `ClusterConfig` file, set the new CIDRs in `vpc.publicAccessCIDRs` and run:

```console
eksctl utils set-public-access-cidrs -f config.yaml
```

!!!note
    This feature only applies to the public endpoint. The
    [API server endpoint access configuration options](https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html)
    won't change, and you will still have the option to disable the public endpoint so your cluster is not accessible from
    the internet. (Source: https://github.com/aws/containers-roadmap/issues/108#issuecomment-552766489)

    Implementation notes: https://github.com/aws/containers-roadmap/issues/108#issuecomment-552698875
