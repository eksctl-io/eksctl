---
title: Example with different IAM options
weight: 30
---

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-west-2

nodeGroups:
  - name: ng1
    instanceType: t3.small
    desiredCapacity: 1
    iam:
      withAddonPolicies:
        imageBuilder: true        
        autoScaler: true        
        externalDNS: true        
        appMesh: true        
        ebs: true        
        fsx: true        
        efs: true        
        albIngress: true        
        xRay: true        
        cloudWatch: true        
```
