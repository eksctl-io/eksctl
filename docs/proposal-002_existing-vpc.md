# Design Proposal #002: Existing VPC

> **STATUS**: This proposal is a _working draft_, it will get refined and augment as needed.
> If any non-trivial changes are need to functionality defined here, in particular the user
> experience, those changes should be suggested via a PR to this proposal document.
> Any other changes to the text of the proposal or technical corrections are also very welcome.

Initial implementation of eksctl intentionally avoided having too many options for user to
customise the VPC configuration. That was in oder to simplify how the tools works.
Since then, many users asked for various features in relation to custom VPC.

There a few specific use-cases:

1. co-location with kops [#50](https://github.com/weaveworks/eksctl/issues/50)
2. set custom CIDR [#279](https://github.com/weaveworks/eksctl/issues/279)
3. private/public subnets [#120](https://github.com/weaveworks/eksctl/issues/120)
4. use any existing VPC [#42](https://github.com/weaveworks/eksctl/issues/42)
5. use same subnet/AZ for a nodegroup [#232](https://github.com/weaveworks/eksctl/issues/232)

Out of the above, we already have support for 1 & 2 (as of 0.1.8) as well as 3 (as of 0.1.9).

The main challenge with these customisations is user experience and how different flags interact.
Ultimately, once Cluster API support is implemented (planned for 0.3.0), some of these flags will not be needed
and there will be more fine-grain parameters available via Cluster API. This proposal sets out flags that should
serve well for most common use-cases until 0.3.0.

## Flags defined in 0.1.9

These are fairly simple, and rather more user-friendly. These are assumed satisfactory for the
purpose of this proposal.

- `--vpc-from-kops-cluster=cluster-1.k8s.local`: import subnets from given kops cluster `cluster-1.k8s.local`
- `--node-private-networking`: use private subnets for initial nodegroup
- `--vpc-cidr=10.9.0.0/16`: use 10.9.0.0/16 as the global CIDR for the VPC

## Complete set of flags required for using existing VPC/subnets

In order to use eksctl with any custom VPC, one has to provide subnet ID. VPC ID can be derived by making an API call.

So the following flags will be needed:

- `--vpc-private-subnets="subnet-1,subnet-2" --vpc-public-subnets="subnet-3,subnet-4"`: use given subnet IDs, detect VPC ID

Also, the following combinations must be valid:

- `--vpc-private-subnets="subnet-1,subnet-2"`: use only given private subnet IDs, detect VPC ID
- `--vpc-public-subnets="subnet-3,subnet-4"`: use given public subnet IDs, detect VPC ID

User will have to make sure IDs are correct as well as their private/public designations.

For some use-cases, specifying VPC ID maybe preferred, and in that case the VPC maybe the default one (in which case it
is ought to be possible to use default subnets or specific subnets). Hence the following combinations ought to be valid:

- `--vpc-id=default`: use default VPC and all of its existing subnets
- `--vpc-id=default --vpc-private-subnets="subnet-1,subnet-2" --vpc-public-subnets="subnet-3,subnet-4"`: use default VPC,
  and given subnet IDs
- `--vpc-id="vpc-xxx"`: use given VPC, create dedicated subnets
- `--vpc-id="vpc-xxx" --vpc-private-subnets="subnet-1,subnet-2" --vpc-public-subnets="subnet-3,subnet-4"`: use given VPC ID,
  and given subnets (same as if `--vpc-id` was unset, but less explicit)

Needless to say that when none of `--vpc-*` flags given, eksctl will create a dedicated VPC and subnets (public as well as
private). As more flags are being added, grouping them in `--help` output would be very helpful; especially because most of
these flags are for advanced usage.

## Security Groups

It must be noted that security groups are managed by eksctl only, as certain configuration is required to ensure cluster is
is fully functional. Unlike with subnets and VPC, security groups can be update after they were created and any ingress and/or
egress rules can be added or removed. If more advanced usage is commonly required, a separate design proposal will be required.
It is also expected that advanced functionality will be available via Cluster API (which is planned for 0.3.0).

## Additional resources/commands

It would be plausible to also provide utility for managing VPC, i.e.:

- `eksctl <create|get|delete> vpc --name=<vpcName>`: manage VPC
- `eksctl create cluster --vpc-name=<vpcName>`: create cluster in the given VPC

This would allow users to create VPCs with recommended configurations without having to worry about the requirements, yet they
will be able to do VPC resource management at a different level of access. It should be noted that similar functionality is likely
to be required for IAM resources (subject to separate proposal).

## Notes on AWS CNI driver

It can be safely assumed (based on simple tests) that the driver is meant to coexist very happily with anything you have in your VPC,
from plain EC2 instances, through load balancers and gateways to Lambda functions. It is thereby ought to be possible to share a VPC
between two or more cluster running separate instances of the driver. It is expected that driver will try and use what IP addresses
are available, which is only a subject to IP address space usage, of course. Allocations and other properties of the driver can be
monitored via CloudWatch and/or Prometheus.

However, when very small subnets are used in a VPC, or there are too many resources using up the IP address space, one will ultimately
run into issue sooner than later. If that's the case, VPC peering or some other alternative should be consider. Peering should already
be possible, and custom CIDR flag is available as of 0.1.8.

It must be noted, that there exists a plan to provide Weave Net as an overlay option. It maybe possible to use IPv6 also, but more
experiments and research will need to be done (it not clear if EKS supports IPv6 properly yet).

## Requirements of custom VPC

When user decides to use any of `--vpc-*` flags, it is up to them to ensure that:

- suffiecent IP addresses are available
- suffiencent number of subnets (minimum 2)
- public/private subnets are configured correctly (i.e. routing and internet/NAT gateways are configured as intended)
- tagging of subnets
  - `kubernetes.io/cluster/<name>` tag set to either `shared` or `owned`
  - `kubernetes.io/role/internal-elb` tag set to `1` for private subnets
  - any other tags

There maybe other requirements imposed by EKS or Kubernetes, and it is entirely up to the user to stay up-to-date on any requirements
and/or recommendations, and find ways in which those work well with the practices at their own organisations and any requirements imposed
by it. It is not our goal to keep this section of the proposal up to date. However, `eksctl create cluster` aims to always provide an
up-to-date configuration with a dedicated VPC.
