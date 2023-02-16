# Launch Template support for Managed Nodegroups

eksctl supports launching managed nodegroups using a provided [EC2 Launch Template](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-launch-templates.html).
This enables multiple customization options for nodegroups including providing custom AMIs and security groups, and passing user data for node bootstrapping.


## Creating managed nodegroups using a provided launch template

```yaml
# managed-cluster.yaml
# A cluster with two managed nodegroups
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: managed-cluster
  region: us-west-2

managedNodeGroups:
  - name: managed-ng-1
    launchTemplate:
      id: lt-12345
      version: "2" # optional (uses the default launch template version if unspecified)

  - name: managed-ng-2
    minSize: 2
    desiredCapacity: 2
    maxSize: 4
    labels:
      role: worker
    tags:
      nodegroup-name: managed-ng-2
    privateNetworking: true
    launchTemplate:
      id: lt-12345

```


## Upgrading a managed nodegroup to use a different launch template version

```shell
eksctl upgrade nodegroup --name=managed-ng-1 --cluster=managed-cluster --launch-template-version=3
```

???+ note
    If a launch template is using a custom AMI, then the new version should also use a custom AMI or the upgrade operation will fail


If a launch template is not using a custom AMI, the Kubernetes version to upgrade to can also be specified:

```shell
eksctl upgrade nodegroup --name=managed-ng-1 --cluster=managed-cluster --launch-template-version=3 --kubernetes-version=1.17
```


## Notes on custom AMI and launch template support
- When a launch template is provided, the following fields are not supported: `instanceType`, `ami`, `ssh.allow`, `ssh.sourceSecurityGroupIds`, `securityGroups`,
 `instancePrefix`, `instanceName`, `ebsOptimized`, `volumeEncrypted`, `volumeKmsKeyID`, `volumeIOPS`, `maxPodsPerNode`, `preBootstrapCommands`, `overrideBootstrapCommand` and `disableIMDSv1`.
- When using a custom AMI (`ami`), `overrideBootstrapCommand` must also be set to perform the bootstrapping.
- `overrideBootstrapCommand` can only be set when using a custom AMI.
- When a launch template is provided, tags specified in the nodegroup config apply to the EKS Nodegroup resource only and are not propagated to EC2 instances.
