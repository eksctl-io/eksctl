# IAM permissions boundary

A [permissions boundary][permissions-boundary] is an advanced AWS IAM feature in which the maximum permissions that an identity-based policy can grant to an IAM entity have been set; where those entities are either users or roles. When a permissions boundary is set for an entity, that entity can only perform the actions that are allowed by both its identity-based policies and its permissions boundaries.

You can provide your permissions boundary so that all identity-based entities created by eksctl are created within that boundary. This example demonstrates how a permissions boundary can be provided to the various identity-based entities that are created by eksctl:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-17
  region: us-west-2

iam:
  withOIDC: true
  serviceRolePermissionsBoundary: "arn:aws:iam::11111:policy/entity/boundary"
  fargatePodExecutionRolePermissionsBoundary: "arn:aws:iam::11111:policy/entity/boundary"
  serviceAccounts:
    - metadata:
        name: s3-reader
      attachPolicyARNs:
      - "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
      permissionsBoundary: "arn:aws:iam::11111:policy/entity/boundary"

nodeGroups:
  - name: "ng-1"
    desiredCapacity: 1
    iam:
      instanceRolePermissionsBoundary: "arn:aws:iam::11111:policy/entity/boundary"
```

!!! warning
    It is not possible to provide both a role ARN and a permissions boundary!

[permissions-boundary]: https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_boundaries.html

## Setting the VPC CNI Permission Boundary
Please note that when you create a cluster with OIDC enabled eksctl will automatically create an `iamserviceaccount` for the VPC-CNI for [security reasons](security.md). If
you would like to add a permission boundary to it then you must specify the `iamserviceaccount` in your config file manually:
```yaml
iam:
  serviceAccounts:
    - metadata:
        name: aws-node
        namespace: kube-system
      attachPolicyARNs:
      - "arn:aws:iam::<arn>:policy/AmazonEKS_CNI_Policy"
      permissionsBoundary: "arn:aws:iam::11111:policy/entity/boundary"
```
