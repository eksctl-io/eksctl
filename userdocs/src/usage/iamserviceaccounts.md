# IAM Roles for Service Accounts

## Introduction

Amazon EKS supports [IAM Roles for Service Accounts (IRSA)][eks-user-guide] that allows cluster operators to map AWS IAM Roles to Kubernetes Service Accounts.

This provides fine-grained permission management for apps that run on EKS and use other AWS services. These could be apps that use S3,
any other data services (RDS, MQ, STS, DynamoDB), or Kubernetes components like AWS Load Balancer controller or ExternalDNS.

You can easily create IAM Role and Service Account pairs with `eksctl`.

!!!note
    If you used [instance roles](/usage/iam-policies), and are considering to use IRSA instead, you shouldn't mix the two.

## How it works

It works via IAM OpenID Connect Provider (OIDC) that EKS exposes, and IAM Roles must be constructed with reference to the IAM OIDC Provider (specific to a given EKS cluster), and a reference to the Kubernetes Service Account it will be bound to.
Once an IAM Role is created, a service account should include the ARN of that role as an annotation (`eks.amazonaws.com/role-arn`).
By default the service account will be created or updated to include the role annotation, this can be disabled using the flag `--role-only`.

Inside EKS, there is an [admission controller](https://github.com/aws/amazon-eks-pod-identity-webhook/) that injects AWS session credentials into pods respectively of the roles based on the annotation on the Service Account used by the pod. The credentials will get exposed by `AWS_ROLE_ARN` & `AWS_WEB_IDENTITY_TOKEN_FILE` environment variables. Given a recent version of AWS SDK is used (see [AWS documentation][eks-user-guide-sdk] for details of exact version), the application will use these credentials.

In `eksctl` the name of the resource is _iamserviceaccount_, which represents an IAM Role and Service Account pair.

### Usage without config files

!!!note
    IAM Roles for Service Accounts require Kubernetes version 1.13 or above.

The IAM OIDC Provider is not enabled by default, you can use the following command to enable it, or use config file (see below):

```console
eksctl utils associate-iam-oidc-provider --cluster=<clusterName>
```

Once you have the IAM OIDC Provider associated with the cluster, to create a IAM role bound to a service account, run:

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=<serviceAccountName> --namespace=<serviceAccountNamespace> --attach-policy-arn=<policyARN>
```

!!!note
    You can specify `--attach-policy-arn` multiple times to use more than one policy.

More specifically, you can create a service account with read-only access to S3 by running:

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=s3-read-only --attach-policy-arn=arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
```

By default, it will be created in `default` namespace, but you can specify any other namespace, e.g.:
```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=s3-read-only --namespace=s3-app --attach-policy-arn=arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
```

!!!note
    If the namespace doesn't exist already, it will be created.

If you have service account already created in the cluster (without an IAM Role), you will need to use `--override-existing-serviceaccounts` flag.

Custom tagging may also be applied to the IAM Role by specifying `--tags`:

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=<serviceAccountName> --tags "Owner=John Doe,Team=Some Team"
```

CloudFormation will generate a role name that includes a random string. If you prefer a predetermined role name you can specify `--role-name`:

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=<serviceAccountName> --role-name "custom-role-name"
```

When the service account is created and managed by some other tool, such as helm, use `--role-only` to prevent conflicts.
The other tool is then responsible for maintaining the role ARN annotation. Note that `--override-existing-serviceaccounts` has no effect on `roleOnly`/`--role-only` service accounts, the role will always be created.

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=<serviceAccountName> --role-only --role-name=<customRoleName>
```

When you have an existing role which you want to use with a service account, you can provide the `--attach-role-arn` flag instead of providing the policies. To ensure the role can only be assumed by the specified service account, you should set a [trust relationship policy document](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts-technical-overview.html#iam-role-configuration).

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=<serviceAccountName> --attach-role-arn=<customRoleARN>
```

To update a service accounts roles permissions you can run `eksctl update iamserviceaccount`.

!!!note
    `eksctl delete iamserviceaccount` deletes Kubernetes `ServiceAccounts` even if they were not created by `eksctl`.

### Usage with config files

To manage `iamserviceaccounts` using config file, you will be looking to set `iam.withOIDC: true` and list account you want under `iam.serviceAccount`.

All of the commands support `--config-file`, you can manage _iamserviceaccounts_ the same way as _nodegroups_.
The `eksctl create iamserviceaccount` command supports `--include` and `--exclude` flags (see
[this section](/usage/managing-nodegroups#include-and-exclude-rules) for more details about how these work).
And the `eksctl delete iamserviceaccount` command supports `--only-missing` as well, so you can perform deletions the same way as nodegroups.

The option to enable `wellKnownPolicies` is included for using IRSA with well-known
use cases like `cluster-autoscaler` and `cert-manager`, as a shorthand for lists
of policies.

Supported well-known policies and other properties of `serviceAccounts` are documented at
[the config schema](https://eksctl.io/usage/schema/#iam-serviceAccounts).

You use the following config example with `eksctl create cluster`:

```YAML
# An example of ClusterConfig with IAMServiceAccounts:
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-13
  region: us-west-2

iam:
  withOIDC: true
  serviceAccounts:
  - metadata:
      name: s3-reader
      # if no namespace is set, "default" will be used;
      # the namespace will be created if it doesn't exist already
      namespace: backend-apps
      labels: {aws-usage: "application"}
    attachPolicyARNs:
    - "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
    tags:
      Owner: "John Doe"
      Team: "Some Team"
  - metadata:
      name: cache-access
      namespace: backend-apps
      labels: {aws-usage: "application"}
    attachPolicyARNs:
    - "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess"
    - "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess"
  - metadata:
      name: cluster-autoscaler
      namespace: kube-system
      labels: {aws-usage: "cluster-ops"}
    wellKnownPolicies:
      autoScaler: true
    roleName: eksctl-cluster-autoscaler-role
    roleOnly: true
  - metadata:
      name: some-app
      namespace: default
    attachRoleARN: arn:aws:iam::123:role/already-created-role-for-app
nodeGroups:
  - name: "ng-1"
    tags:
      # EC2 tags required for cluster-autoscaler auto-discovery
      k8s.io/cluster-autoscaler/enabled: "true"
      k8s.io/cluster-autoscaler/cluster-13: "owned"
    desiredCapacity: 1
```

If you create a cluster without these fields set, you can use the following commands to enable all you need:

```console
eksctl utils associate-iam-oidc-provider --config-file=<path>
eksctl create iamserviceaccount --config-file=<path>
```

### Further information

- [Introducing Fine-grained IAM Roles For Service Accounts](https://aws.amazon.com/blogs/opensource/introducing-fine-grained-iam-roles-service-accounts/)
- [AWS EKS User Guide - IAM Roles For Service Accounts][eks-user-guide]
- [Mapping IAM users and role to Kubernetes RBAC roles](https://eksctl.io/usage/iam-identity-mappings/)

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html
[eks-user-guide-sdk]: https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts-minimum-sdk.html
