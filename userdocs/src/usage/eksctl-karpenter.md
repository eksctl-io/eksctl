# Karpenter Support

`eksctl` provides support for adding [Karpenter](https://karpenter.sh/) to a newly created cluster. It will create all the necessary
prerequisites outlined in Karpenter's [Getting Started](https://karpenter.sh/docs/getting-started/) section including installing
Karpenter itself using Helm. We currently support installing versions `0.28.0+`. See the [Karpenter compatibility](https://karpenter.sh/docs/upgrading/compatibility/) section for further details.

The following cluster configuration outlines a typical Karpenter installation:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-with-karpenter
  region: us-west-2
  version: '1.32' # requires a version of Kubernetes compatible with Karpenter
  tags:
    karpenter.sh/discovery: cluster-with-karpenter # here, it is set to the cluster name
iam:
  withOIDC: true # required

karpenter:
  version: '1.2.1' # Exact version should be specified according to the Karpenter compatibility matrix

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
  version: '1.2.1'
  createServiceAccount: true # default is false
  defaultInstanceProfile: 'KarpenterNodeInstanceProfile' # default is to use the IAM instance profile created by eksctl
  withSpotInterruptionQueue: true # adds all required policies and rules for supporting Spot Interruption Queue, default is false
```

OIDC must be defined in order to install Karpenter.

Once Karpenter is successfully installed, add [NodePool(s)](https://karpenter.sh/docs/concepts/nodepools/) and [NodeClass(es)](https://karpenter.sh/docs/concepts/nodeclasses/) to allow Karpenter
to start adding nodes to the cluster.

The NodePool's `nodeClassRef` section must match the name of an `EC2NodeClass`. For example:

```yaml
apiVersion: karpenter.sh/v1
kind: NodePool
metadata:
  name: example
  annotations:
    kubernetes.io/description: "Example NodePool"
spec:
  template:
    spec:
      requirements:
        - key: kubernetes.io/arch
          operator: In
          values: ["amd64"]
        - key: kubernetes.io/os
          operator: In
          values: ["linux"]
        - key: karpenter.sh/capacity-type
          operator: In
          values: ["on-demand"]
        - key: karpenter.k8s.aws/instance-category
          operator: In
          values: ["c", "m", "r"]
        - key: karpenter.k8s.aws/instance-generation
          operator: Gt
          values: ["2"]
      nodeClassRef:
        group: karpenter.k8s.aws
        kind: EC2NodeClass
        name: example # must match the name of an EC2NodeClass
```

```yaml
apiVersion: karpenter.k8s.aws/v1
kind: EC2NodeClass
metadata:
  name: example
  annotations:
    kubernetes.io/description: "Example EC2NodeClass"
spec:
  role: "eksctl-KarpenterNodeRole-${CLUSTER_NAME}" # replace with your cluster name
  subnetSelectorTerms:
    - tags:
        karpenter.sh/discovery: "${CLUSTER_NAME}" # replace with your cluster name
  securityGroupSelectorTerms:
    - tags:
        karpenter.sh/discovery: "${CLUSTER_NAME}" # replace with your cluster name
  amiSelectorTerms:
    - alias: al2023@latest # Amazon Linux 2023
```

Note that you must specify one of `role` or `instanceProfile` for lauch nodes. If you choose to use `instanceProfile`
the name of the profile created by `eksctl` follows the pattern: `eksctl-KarpenterNodeInstanceProfile-<cluster-name>`.

## Automatic Security Group Tagging

`eksctl` automatically tags the cluster's shared node security group with `karpenter.sh/discovery` when both Karpenter is enabled (`karpenter.version` specified) and the `karpenter.sh/discovery` tag exists in `metadata.tags`. This enables AWS Load Balancer Controller compatibility.
