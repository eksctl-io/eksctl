# Karpenter Support

`eksctl` provides adding [Karpenter](https://karpenter.sh/) to a newly created cluster. It will create all the necessary
prerequisites outlined in Karpenter's [Getting Started](https://karpenter.sh/docs/getting-started/) section including installing
Karpenter itself using Helm. We currently support installing versions up to `0.15.0`

To that end, a new configuration value has been introduced into `eksctl` cluster config called `karpenter`. The following
yaml outlines a typical installation configuration:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-with-karpenter
  region: us-west-2
  version: '1.20'
  tags:
    karpenter.sh/discovery: cluster-with-karpenter # here, it is set to the cluster name
iam:
  withOIDC: true # required

karpenter:
  version: '0.15.0'

managedNodeGroups:
  - name: managed-ng-1
    minSize: 1
    maxSize: 2
    desiredCapacity: 1
```

The version is Karpenter's version as it can be found in their Helm Repository. The following options are also available
to be set: 

```yaml
karpenter:
  version: '0.15.0'
  createServiceAccount: true # default is false
  defaultInstanceProfile: 'KarpenterNodeInstanceProfile' # default is to use the IAM instance profile created by eksctl
```

OIDC must be defined in order to install Karpenter.

Once Karpenter is successfully installed, add a [Provisioner](https://karpenter.sh/docs/provisioner/) so Karpenter
can start adding the right nodes to the cluster.

The provisioner's `instanceProfile` section must match the created `NodeInstanceProfile` role's name. For example:

```yaml
apiVersion: karpenter.sh/v1alpha5
kind: Provisioner
metadata:
  name: default
spec:
  requirements:
    - key: karpenter.sh/capacity-type
      operator: In
      values: ["on-demand"]
  limits:
    resources:
      cpu: 1000
  provider:
    instanceProfile: eksctl-KarpenterNodeInstanceProfile-${CLUSTER_NAME}
    subnetSelector:
      karpenter.sh/discovery: cluster-with-karpenter # must match the tag set in the config file
    securityGroupSelector:
      karpenter.sh/discovery: cluster-with-karpenter # must match the tag set in the config file
  ttlSecondsAfterEmpty: 30
```

Note that unless `defaultInstanceProfile` is defined, the name used for `instanceProfile` is
`eksctl-KarpenterNodeInstanceProfile-<cluster-name>`.
