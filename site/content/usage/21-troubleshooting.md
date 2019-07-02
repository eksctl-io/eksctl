---
title: Troubleshooting
weight: 210
url: usage/troubleshooting
---

## Troubleshooting

### subnet ID "subnet-11111111" is not the same as "subnet-22222222"

Given a config file specifying subnets for a VPC like the following:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test
  region: us-east-1

vpc:
  subnets:
    public:
      us-east-1a: {id: subnet-11111111}
      us-east-1b: {id: subnet-22222222}
    private:
      us-east-1a: {id: subnet-33333333}
      us-east-1b: {id: subnet-44444444}

nodeGroups: []
```

An error `subnet ID "subnet-11111111" is not the same as "subnet-22222222"` means that the subnets specified are not 
placed in the right Availability zone. Check in the AWS console which is the right subnet ID for each Availability Zone.

In this example, the correct configuration for the VPC would be:

```yaml
vpc:
  subnets:
    public:
      us-east-1a: {id: subnet-22222222}
      us-east-1b: {id: subnet-11111111}
    private:
      us-east-1a: {id: subnet-33333333}
      us-east-1b: {id: subnet-44444444}
```
