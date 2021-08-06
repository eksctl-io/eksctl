# Define Container Runtime

For AL2 ( AmazonLinux2 ) AMIs it's possible to set container runtime to `containerd`.

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
- dockerd

## Managed Nodes

For managed nodes we don't explicitly provide a bootstrap script, and thus it's up to the user
to define a different runtime if they wish, using `overrideBootstrapCommand`:

```yaml
managedNodeGroups:
  - name: m-ng-1
    instanceType: m5.large
    overrideBootstrapCommand: |
      #!/bin/bash
      /etc/eks/bootstrap.sh <cluster-name> <other flags> --container-runtime containerd
```
