---
title: "IAM Roles for Service Accounts"
weight: 120
url: usage/iamserviceaccounts
---

## Introduction

Amazon EKS supports [IAM Roles for Service Accounts (IRSA)][eks-user-guide] that allows cluster operators to map AWS IAM Roles to Kubernetes Service Accounts.

This provides fine-grained permission management for apps that run on EKS and use other AWS services. These could be apps that use S3,
any other data services (RDS, MQ, STS, DynamoDB), or Kubernetes components like AWS ALB Ingress controller or ExternalDNS.

You can easily create IAM Role and Service Account pairs with `eksctl`.

> NOTE: if you used [instance roles](https://eksctl.io/usage/iam-policies/), and are considering to use IRSA instead, you shouldn't mix the two.

## How it works

It works via IAM OpenID Connect Provider (OIDC) that EKS exposes, and IAM Roles must be constructed with reference to the IAM OIDC Provider (specific to a given EKS cluster), and a reference to the Kubernetes Service Account it will be bound to.
Once an IAM Role is created, a service account should include the ARN of that role as an annotation (`eks.amazonaws.com/role-arn`).

Inside EKS, there is an [admission controller](https://github.com/aws/amazon-eks-pod-identity-webhook/) that injects AWS session credentials into pods respectively of the roles based on the annotation on the Service Account used by the pod. The credentials will get exposed by `AWS_ROLE_ARN` & `AWS_WEB_IDENTITY_TOKEN_FILE` environment variables. Given a recent version of AWS SDK is used (see [AWS documentation][eks-user-guide-sdk] for details of exact version), the application will use these credentials.

In `eksctl` the name of the resource is _iamserviceaccount_, which represents an IAM Role and Service Account pair.

### Usage without config files

> NOTE: IAM Roles for Service Accounts require Kubernetes version 1.13 or above.

The IAM OIDC Provider is not enabled by default, you can use the following command to enable it, or use config file (see below):

```console
eksctl utils associate-iam-oidc-provider --cluster=<clusterName>
```

Once you have the IAM OIDC Provider associated with the cluster, to create a IAM role bound to a service account, run:

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=<serviceAccountName> --namespace=<serviceAccountNamespace> --attach-policy-arn=<policyARN>
```

> NOTE: you can specify `--attach-policy-arn` multiple times to use more then one policy.

More specifically, you can create a service account with read-only access to S3 by running:

```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=s3-read-only --attach-policy-arn=arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
```

By default, it will be created in `default` namespace, but you can specify any other namespace, e.g.:
```console
eksctl create iamserviceaccount --cluster=<clusterName> --name=s3-read-only --namespace=s3-app --attach-policy-arn=arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
```

> NOTE: If the namespace doesn't exist already, it will be created.

If you have service account already created in the cluster (without an IAM Role), you will need to use `--override-existing-serviceaccounts` flag.

Currently, to update a role you will need to re-create, run `eksctl delete iamserviceaccount` followed by `eksctl create iamserviceaccount` to achieve that.

### Usage with config files

To manage `iamserviceaccounts` using config file, you will be looking to set `iam.withOIDC: true` and list account you want under `iam.serviceAccount`.

All of the commands support `--config-file`, you can manage _iamserviceaccounts_ the same way as _nodegroups_.
The `eksctl create iamserviceaccount` command supports `--include` and `--exclude` flags.
And the `eksctl delete iamserviceaccount` command supports `--only-missing` as well, so you can perform deletions the same way as nodegroups.

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
    attachPolicy: # inline policy can be defined along with `attachPolicyARNs`
      Version: "2012-10-17"
      Statement:
      - Effect: Allow
        Action:
        - "autoscaling:DescribeAutoScalingGroups"
        - "autoscaling:DescribeAutoScalingInstances"
        - "autoscaling:DescribeLaunchConfigurations"
        - "autoscaling:DescribeTags"
        - "autoscaling:SetDesiredCapacity"
        - "autoscaling:TerminateInstanceInAutoScalingGroup"
        Resource: '*'

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
