# Subnet Settings

## Use private subnets for initial nodegroup

If you prefer to isolate the initial nodegroup from the public internet, you can use the `--node-private-networking` flag.
When used in conjunction with the `--ssh-access` flag, the SSH port can only be accessed from inside the VPC.

???+ note
    Using the `--node-private-networking` flag will result in outgoing traffic to go through the NAT gateway using its
    Elastic IP. On the other hand, if the nodes are in a public subnet, the outgoing traffic won't go through the
    NAT gateway and hence the outgoing traffic has the IP of each individual node.

## Custom subnet topology

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

???+ note
    A complete subnet spec must be provided, i.e. both `public` and `private` configurations
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

???+ note
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
