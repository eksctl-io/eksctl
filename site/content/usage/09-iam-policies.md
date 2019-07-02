---
title: "IAM policies"
weight: 90
url: usage/iam-policies
---

## Supported IAM add-on policies

Example of all supported add-on policies:

```yaml
nodeGroups:
  - name: ng-1
    instanceType: m5.xlarge
    desiredCapacity: 1
    iam:
      withAddonPolicies:
        imageBuilder: true
        autoScaler: true
        externalDNS: true
        certManager: true
        appMesh: true
        ebs: true
        fsx: true
        efs: true
        albIngress: true
        xRay: true
        cloudWatch: true
```

### Image Builder Policy

The `imageBuilder` policy allows for full ECR (Elastic Container Registry) access. This is useful for building, for
example, a CI server that needs to push images to ECR.

### EBS Policy

The `ebs` policy enables the new EBS CSI (Elastic Block Store Container Storage Interface) driver.

[comment]: <> (TODO: One section per addon and brief explanation of what it is )

## Adding a custom instance role

This example creates a nodegroup that reuses an existing IAM Instance Role from another cluster:

```yaml
apiVersion: eksctl.io/v1alpha4
kind: ClusterConfig
metadata:
  name: test-cluster-c-1
  region: eu-north-1

nodeGroups:
  - name: ng2-private
    instanceType: m5.large
    desiredCapacity: 1
    iam:
      instanceProfileARN: "arn:aws:iam::123:instance-profile/eksctl-test-cluster-a-3-nodegroup-ng2-private-NodeInstanceProfile-Y4YKHLNINMXC"
      instanceRoleARN: "arn:aws:iam::123:role/eksctl-test-cluster-a-3-nodegroup-NodeInstanceRole-DNGMQTQHQHBJ"
```

## Attaching policies by ARN

```yaml
nodeGroups:
  - name: my-special-nodegroup
    iam:
      attachPolicyARNs:
        - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy
        - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
        - arn:aws:iam::aws:policy/ElasticLoadBalancingFullAccess
        - arn:aws:iam::1111111111:policy/kube2iam
      withAddonPolicies:
        autoScaler: true
        imageBuilder: true
```

_Important_: if a nodegroup includes the `attachPolicyARNs` it must also include the default node policies, like in
this example (`AmazonEKSWorkerNodePolicy` and `AmazonEKS_CNI_Policy`).

[comment]: <> (TODO find better example and explain more)
