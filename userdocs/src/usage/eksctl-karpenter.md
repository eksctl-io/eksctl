# Karpenter Support

`eksctl` provides adding [Karpenter](https://karpenter.sh/) to a newly created cluster. It will create all the necessary
prerequisites outlined in Karpenter's [Getting Started](https://karpenter.sh/docs/getting-started/) section including installing
Karpenter itself using Helm. We currently support installing versions starting `0.28.0` and above. See in The [Compatibility](https://karpenter.sh/docs/upgrading/compatibility/) section.

???+ info
    With [v0.17.0](https://karpenter.sh/docs/upgrading/upgrade-guide/) Karpenter’s Helm chart package is now stored in Karpenter’s OCI (Open Container Initiative) registry.
    Clusters created on previous versions shouldn't be affected by this change. If you wish to upgrade your current installation of Karpenter please refer to the [upgrade guide](https://karpenter.sh/docs/upgrading/upgrade-guide/)
    You have to be logged out of ECR repositories to be able to pull the OCI artifact by running `helm registry logout public.ecr.aws` or `docker logout public.ecr.aws`, failure to do so will result in a 403 error when trying to pull the chart.

To that end, a new configuration value has been introduced into `eksctl` cluster config called `karpenter`. The following
yaml outlines a typical installation configuration:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-with-karpenter
  region: us-west-2
  version: '1.32' # requires a version of kubernetes compatible with karpenter
  tags:
    karpenter.sh/discovery: cluster-with-karpenter # here, it is set to the cluster name
iam:
  withOIDC: true # required

karpenter:
  version: '1.2.1' # Exact version should be specified according to the karpenter compatibility matrix

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

Once Karpenter is successfully installed, add a [NodePools](https://karpenter.sh/docs/concepts/nodepools/) and [NodeClasses](https://karpenter.sh/docs/concepts/nodeclasses/) so Karpenter
can start adding the right nodes to the cluster.

The NodePool's `nodeClassRef` section must match the created `EC2NodeClass` metadata name. For example:

```yaml
apiVersion: karpenter.sh/v1
kind: NodePool
metadata:
  name: general-purpose
  annotations:
    kubernetes.io/description: "General purpose NodePool for generic workloads"
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
        name: default

```

```yaml
apiVersion: karpenter.k8s.aws/v1
kind: EC2NodeClass
metadata:
  name: default
  annotations:
    kubernetes.io/description: "General purpose EC2NodeClass for running Amazon Linux 2 nodes"
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

Note that you must specify one of `role` or `instanceProfile` for lauch nodes. If you choose to use `instanceProfile` the name is
`eksctl-KarpenterNodeInstanceProfile-<cluster-name>`.
