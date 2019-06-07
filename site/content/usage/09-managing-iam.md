---
title: "IAM"
weight: 90
---

## IAM Add-on Policies

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
        appMesh: true        
        ebs: true        
        fsx: true        
        efs: true        
        albIngress: true        
        xRay: true        
        cloudWatch: true        
```

[comment]: <> (TODO: One section per addon and brief explanation of what it is )

## Managing IAM Users and Roles

EKS clusters use IAM users and roles to control access to the cluster. The rules are implemented in a config map 
called `aws-auth`. `eksctl` provides commands to read and edit this config map.


Get all identity mappings:

```bash
eksctl get iamidentitymapping --name my-cluster-1
```

Get all identity mappings matching a role:
```bash
eksctl get iamidentitymapping --name my-cluster-1 --role arn:aws:iam::123456:role/testing-role
``` 

Create an identity mapping:
```bash
 eksctl create iamidentitymapping --name  my-cluster-1 --role arn:aws:iam::123456:role/testing --group system:masters --username admin
```

Delete a mapping:

```bash
eksctl delete iamidentitymapping --name  my-cluster-1 --role arn:aws:iam::123456:role/testing
```

*Note*: this deletes a single mapping FIFO unless `--all`is given in which case it removes all matching. Will warn if 
more mappings matching this role are found.
