# Define Container Runtime

!!! warning
    Starting with Kubernetes version `1.24`, dockershim support has been deprecated. Therefore, if you create a cluster using `eksctl` on version `1.24` or higher, the information below no longer applies, and the only supported container runtime is `containerd`. Trying to set it otherwise will return a validation error. Additionally, AL2023 AMIs only support `containerd` regadless of K8s version.

    At some point, we will completely remove the option to set `containerRuntime` in config file, together with the support for older Kubernetes versions support (i.e. `1.22` or `1.23`).

For AL2 ( AmazonLinux2 ) and Windows AMIs, it's possible to set container runtime to `containerd`.

## Un-managed Nodes

For un-managed nodes, simply provide the following configuration when creating a new node:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: container-runtime-test
  region: us-west-2

nodeGroups:
  - name: ng-1
    instanceType: m5.xlarge
    desiredCapacity: 2
    amiFamily: AmazonLinux2
    containerRuntime: containerd
```

This value is set to `dockerd` by default to preserve backwards compatibility, but will soon be
deprecated.

_Note that there is no equivalent flag for setting the container runtime, this can only be done via a config file._

At the time of this writing the following container runtime values are allowed:

- containerd
- dockerd (docker for Windows)

## Managed Nodes

For managed nodes we don't explicitly provide a bootstrap script, and thus it's up to the user
to define a different runtime if they wish, using `overrideBootstrapCommand`.
The `overrideBootstrapCommand` option requires that you specify an AMI for the managed node group.

```yaml
managedNodeGroups:
  - name: m-ng-1
    ami: ami-XXXXXXXXXXXXXX
    instanceType: m5.large
    overrideBootstrapCommand: |
      #!/bin/bash
      /etc/eks/bootstrap.sh <cluster-name> <other flags> --container-runtime containerd
```
For Windows managed nodes, you will need to use a custom launch template with an ami-id and pass in the required bootstrap arguments in the userdata.
Read more on [creating a launch template](https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html), and [using a launch template with eksctl](launch-template-support.md).
