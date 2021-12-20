# Latest & Custom AMI support

!!! warning
    In a future as yet undecided release **unmanaged** nodegroups created with **custom AmazonLinux2** or **custom Ubuntu** images
    will need to have the `overrideBootstrapCommand` configuration option set. This is to ensure that nodes are able
    to join the cluster. For more information and to track this change please see [this issue](https://github.com/weaveworks/eksctl/issues/3563).

    Users setting the `ami` field on **unmanaged** nodegroups to `ami-XXXX` (i.e. setting a custom AMI) will start to see
    a warning message that they are being sent down a legacy code path. There is no action to take at this time, but in a future
    release the users seeing this warning will be affected by the above change.

## Setting the node AMI ID

The `--node-ami` flag enables a number of advanced use cases such as using a custom AMI or querying AWS in realtime to determine which AMI to use.
The flag can be used for both non-GPU and GPU images.

The flag can take the AMI image id for an image to explicitly use. It also can take the following 'special' keywords:

| Keyword   | Description                                                                                                         |
| --------- | ------------------------------------------------------------------------------------------------------------------- |
| auto      | Indicates that the AMI to use for the nodes should be found by querying AWS EC2. This relates to the auto resolver. |
| auto-ssm  | Indicates that the AMI to use for the nodes should be found by querying AWS SSM Parameter Store.                    |


!!! note
    When setting `--node-ami` to an ID string, `eksctl` will assume that a custom AMI has been requested.
    For managed nodes this will mean that `overrideBootstrapCommand` is required. For unmanaged nodes
    `overrideBootstrapCommand` is recommended for AmazonLinux2 and Ubuntu custom images.

CLI flag examples:
```sh
eksctl create cluster --node-ami=auto

# with a custom ami id
eksctl create cluster --node-ami=ami-custom1234
```

Config file example:
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
      /etc/eks/bootstrap.sh <cluster-name>
```

The `--node-ami` flag can also be used with `eksctl create nodegroup`.

## Setting the node AMI Family

The `--node-ami-family` can take following keywords:

| Keyword                        |                                          Description                                         |
|--------------------------------|:--------------------------------------------------------------------------------------------:|
| AmazonLinux2                   | Indicates that the EKS AMI image based on Amazon Linux 2 should be used (default).           |
| Ubuntu2004                     | Indicates that the EKS AMI image based on Ubuntu 20.04 LTS (Focal) should be used.           |
| Ubuntu1804                     | Indicates that the EKS AMI image based on Ubuntu 18.04 LTS (Bionic) should be used.          |
| Bottlerocket                   | Indicates that the EKS AMI image based on Bottlerocket should be used.                       |
| WindowsServer2019FullContainer | Indicates that the EKS AMI image based on Windows Server 2019 Full Container should be used. |
| WindowsServer2019CoreContainer | Indicates that the EKS AMI image based on Windows Server 2019 Core Container should be used. |
| WindowsServer2004CoreContainer | Indicates that the EKS AMI image based on Windows Server 2004 Core Container should be used. |
| WindowsServer20H2CoreContainer | Indicates that the EKS AMI image based on Windows Server 20H2 Core Container should be used. |

CLI flag example:
```sh
eksctl create cluster --node-ami-family=AmazonLinux2
```

Config file example:
```yaml
nodeGroups:
  - name: ng1
    instanceType: m5.large
    amiFamily: AmazonLinux2
managedNodeGroups:
  - name: m-ng-2
    instanceType: m5.large
    amiFamily: Ubuntu2004
```

The `--node-ami-family` flag can also be used with `eksctl create nodegroup`.
