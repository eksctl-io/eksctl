# Nodegroups with custom subnet(s)

It's possible to extend an existing VPC with a new subnet and add a Nodegroup to that subnet.

## Why

Should the cluster run out of pre-configured IPs, it's possible to resize the existing VPC with
a new CIDR to add a new subnet to it. To see how to do that, read this guide on AWS [Extending VPCs](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html#vpc-resize).

### TL;DR

Go to the VPC's configuration and add click on Actions->Edit CIDRs and add a new range.
For example:

```diff
192.168.0.0/19 -> existing CIDR
+ 192.169.0.0/19 -> new CIDR
```

Now you need to add a new Subnet. Depending on if it's a new Private or a Public subnet, you will have
to copy the routing information from a private or a public subnet respectively.

Once the subnet is created, add routing, and copy either the NAT gateway ID or the Internet Gateway
from another subnet in the VPC. Take care that if it's a public subnet Enable Automatic IP Assignment.
Actions->Modify auto-assign IP settings->Enable auto-assign public IPv4 address.

Don't forget to also copy the TAGS of the existing subnets depending on Public or Private subnet configuration.
This is important, otherwise the subnet will not be part of the cluster and instances in the subnet
will be unable to join.

When finished, copy the new subnet's ID. Repeat as often as necessary.

## How

To create a nodegroup in the created subnet(s) run the following command:

```bash
eksctl create nodegroup --cluster <cluster-name> --name my-new-subnet --subnet-ids subnet-0edeb3a04bec27141,subnet-0edeb3a04bec27142,subnet-0edeb3a04bec27143
# or for a single subnet id
eksctl create nodegroup --cluster <cluster-name> --name my-new-subnet --subnet-ids subnet-0edeb3a04bec27141
```

Or, use the configuration as such:

```
eksctl create nodegroup -f cluster-managed.yaml
```

With a configuration like this:

```yaml
# A simple example of ClusterConfig object with two nodegroups:
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-3
  region: eu-north-1

nodeGroups:
  - name: new-subnet-nodegroup
    instanceType: m5.large
    desiredCapacity: 1
    subnets:
      - subnet-id1
      - subnet-id2
```

Wait for the nodegroup to be created and the new instances should have the new IP ranges of the subnet(s).

## Deleting the cluster

Since the new addition modified the existing VPC by adding a dependency outside of the CloudFormation stack, CloudFormation
can no longer remove the cluster.

Before deleting the cluster, remove all created extra subnets by hand, then proceed by calling `eksctl`:

```
eksctl delete cluster -n <cluster-name> --wait
```
