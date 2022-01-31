# IPv6 Support

## Define IP Family

When `eksctl` creates a vpc, you can define the IP version that will be used. The following options are available to be configured:

- IPv4
- IPv6

To define it, use the following example:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-test
  region: us-west-2
  version: "1.21"

kubernetesNetworkConfig:
  ipFamily: IPv6 # or IPv4

addons:
  - name: vpc-cni
  - name: coredns
  - name: kube-proxy

iam:
  withOIDC: true
```

This is an in config file setting only. When IPv6 is set, the following restriction must be followed:

- OIDC is enabled
- managed addons are defined as shows above
- cluster version must be => 1.21
- vpc-cni addon version must be => 1.10.0
- unmanaged nodegroups are not yet supported with IPv6 clusters
- managed nodegroup creation is not supported with un-owned IPv6 clusters
- `vpc.NAT` and `serviceIPv4CIDR` fields are created by eksctl for ipv6 clusters and thus, are not supported configuration options
- AutoAllocateIPv6 is not supported together with IPv6

The default value is `IPv4`.

Private networking can be done with IPv6 IP family as well. Please follow the instruction outlined under [EKS Private Cluster](/usage/eks-private-cluster).
