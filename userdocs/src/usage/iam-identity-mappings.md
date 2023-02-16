# Manage IAM users and roles

EKS clusters use IAM users and roles to control access to the cluster. The rules are implemented in a config map
called `aws-auth`. `eksctl` provides commands to read and edit this config map.

Get all identity mappings:

```bash
eksctl get iamidentitymapping --cluster <clusterName> --region=<region>
```

Get all identity mappings matching an arn:

```bash
eksctl get iamidentitymapping --cluster <clusterName> --region=<region> --arn arn:aws:iam::123456:role/testing-role
```

Create an identity mapping:

```bash
 eksctl create iamidentitymapping --cluster  <clusterName> --region=<region> --arn arn:aws:iam::123456:role/testing --group system:masters --username admin
```

The identity mappings can also be specified in ClusterConfig:

```yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-with-iamidentitymappings
  region: us-east-1

iamIdentityMappings:
  - arn: arn:aws:iam::000000000000:role/myAdminRole
    groups:
      - system:masters
    username: admin
    noDuplicateARNs: true # prevents shadowing of ARNs

  - arn: arn:aws:iam::000000000000:user/myUser
    username: myUser
    noDuplicateARNs: true # prevents shadowing of ARNs

  - serviceName: emr-containers
    namespace: emr # serviceName requires namespace

  - account: "000000000000" # account must be configured with no other options

nodeGroups:
  - name: ng-1
    instanceType: m5.large
    desiredCapacity: 1
```

```bash
 eksctl create iamidentitymapping -f cluster-with-iamidentitymappings.yaml
```

Delete an identity mapping:

```bash
eksctl delete iamidentitymapping --cluster  <clusterName> --region=<region> --arn arn:aws:iam::123456:role/testing
```

???+ note
    Above command deletes a single mapping FIFO unless `--all` is given in which case it removes all matching. Will warn if
    more mappings matching this role are found.

Create an account mapping:

```bash
 eksctl create iamidentitymapping --cluster  <clusterName> --region=<region> --account user-account
```

Delete an account mapping:

```bash
 eksctl delete iamidentitymapping --cluster  <clusterName> --region=<region> --account user-account
```
