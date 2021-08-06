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
    labels: { role: web }
    instanceType: m5.xlarge
    desiredCapacity: 10
    privateNetworking: true
    amiFamily: AmazonLinux2
    containerRuntime: containerd
  - name: ng-2-api
    labels: { role: api }
    instanceType: m5.2xlarge
    desiredCapacity: 2
    privateNetworking: true
```

This value is set to `dockerd` by default to preserve backwards compatibility, but will soon be
deprecated.

The container runtime value is then passed to the bootstrap.sh script by `eksctl` as follows:

```bash
/etc/eks/bootstrap.sh "${CLUSTER_NAME}" \
  ...
  --container-runtime "${CONTAINER_RUNTIME}"
```

At the time of this writing the following container runtime values are allowed:

- containerd
- dockerd

## Managed Nodes

For managed nodes we don't explicitly provide a bootstrap script, and thus it's up to the user
to define a different runtime if they wish, using `overrideBootstrapCommand`:

```yaml
nodeGroups:
  - name: ng1
    instanceType: p2.xlarge
    ami: auto
  - name: ng2
    instanceType: m5.large
    ami: ami-custom1234
managedNodeGroups:
  - name: m-ng-2
    ami: ami-custom1234
    instanceType: m5.large
    overrideBootstrapCommand: |
      #!/bin/bash
      /etc/eks/bootstrap.sh <cluster-name> --container-runtime containerd
```