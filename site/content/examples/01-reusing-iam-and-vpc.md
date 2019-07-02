---
title: Example with custom IAM and VPC config
weight: 10
url: examples/reusing-iam-and-vpc
---
# Example with custom IAM and VPC config

This example shows how to create a cluster reusing pre-existing IAM and VPC resources:

```yaml
apiVersion: eksctl.io/v1alpha4
kind: ClusterConfig
metadata:
  name: my-test
  region: us-east-1
vpc:
  id: "vpc-11111"
  cidr: "152.28.0.0/16"
  subnets:
    private:
      us-east-1d:
          id: "subnet-1111"
          cidr: "152.28.152.0/21"
      us-east-1c:
          id: "subnet-11112"
          cidr: "152.28.144.0/21"
      us-east-1a:
          id: "subnet-11113"
          cidr: "152.28.136.0/21"
iam: 
  serviceRoleARN: "arn:aws:iam::11111:role/eks-base-service-role"

nodeGroups:
  - name: my-test-m5-private
    labels: {pool: my-test-m5-private}
    instanceType: m5.large
    desiredCapacity: 3
    minSize: 1
    maxSize: 15
    volumeSize: 50
    volumeType: gp2
    iam:
      instanceProfileARN: "arn:aws:iam::11111:instance-profile/eks-nodes-base-role"
      instanceRoleARN: "arn:aws:iam::1111:role/eks-nodes-base-role"
    privateNetworking: true
    securityGroups:
      withShared: true
      withLocal: true
      attachIDs: ['sg-11111', 'sg-11112']
    allowSSH: true
    sshPublicKeyName: 'my-instance-key'
    tags:
      'environment:basedomain': 'example.org'
```

[comment]: <> (TODO explain in more detail)
[comment]: <> (TODO mention why withLocal and withShared are needed)
