# Non eksctl-created clusters

From `eksctl` version `0.40.0` users can run `eksctl` commands against clusters which were
not created by `eksctl`.

???+ note
    Eksctl can only support unowned clusters with names which comply with the guidelines mentioned [here](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-using-console-create-stack-parameters.html). Any cluster names which do not match this will fail CloudFormation API validation check.

## Supported commands

The following commands can be used against clusters created by any means other than `eksctl`.
The commands, flags and config file options can be used in exactly the same way.

If we have missed some functionality, please [let us know](https://github.com/eksctl-io/eksctl/issues).

- [x] Create:
    - [x] `eksctl create nodegroup` ([see note below](#creating-nodegroups))
    - [x] `eksctl create fargateprofile`
    - [x] `eksctl create iamserviceaccount`
    - [x] `eksctl create iamidentitymapping`
- [x] Get:
    - [x] `eksctl get clusters/cluster`
    - [x] `eksctl get fargateprofile`
    - [x] `eksctl get nodegroup`
    - [x] `eksctl get labels`
- [x] Delete:
    - [x] `eksctl delete cluster`
    - [x] `eksctl delete nodegroup`
    - [x] `eksctl delete fargateprofile`
    - [x] `eksctl delete iamserviceaccount`
    - [x] `eksctl delete iamidentitymapping`
- [x] Upgrade:
    - [x] `eksctl upgrade cluster`
    - [x] `eksctl upgrade nodegroup`
- [x] Set/Unset:
    - [x] `eksctl set labels`
    - [x] `eksctl unset labels`
- [x] Scale:
    - [x] `eksctl scale nodegroup`
- [x] Drain:
    - [x] `eksctl drain nodegroup`
- [x] Enable:
    - [x] `eksctl enable profile`
    - [x] `eksctl enable repo`
- [x] Utils:
    - [x] `eksctl utils associate-iam-oidc-provider`
    - [x] `eksctl utils describe-stacks`
    - [x] `eksctl utils install-vpc-controllers`
    - [x] `eksctl utils nodegroup-health`
    - [x] `eksctl utils set-public-access-cidrs`
    - [x] `eksctl utils update-cluster-endpoints`
    - [x] `eksctl utils update-cluster-logging`
    - [x] `eksctl utils write-kubeconfig`
    - [x] `eksctl utils update-coredns`
    - [x] `eksctl utils update-aws-node`
    - [x] `eksctl utils update-kube-proxy`

## Creating nodegroups

`eksctl create nodegroup` is the only command which requires specific input from the user.

Since users can create their clusters with any networking configuration they like,
for the time-being, `eksctl` will not attempt to retrieve or guess these values. This
may change in the future as we learn more about how people are using this command on non eksctl-created clusters.

This means that in order to create nodegroups or managed nodegroups on a cluster which was
not created by `eksctl`, a config file containing VPC details must be provided. At a minimum:

```yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: non-eksctl-created-cluster
  region: us-west-2

vpc:
  id: "vpc-12345"
  securityGroup: "sg-12345"    # this is the ControlPlaneSecurityGroup
  subnets:
    private:
      private1:
          id: "subnet-12345"
      private2:
          id: "subnet-67890"
    public:
      public1:
          id: "subnet-12345"
      public2:
          id: "subnet-67890"

...
```

Further information on VPC configuration options can be found [here](/usage/vpc-networking).
