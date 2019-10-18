---
title: "Custom AMI support"
weight: 50
url: usage/custom-ami-support
---

## Latest & Custom AMI Support

With the 0.1.2 release we have introduced the `--node-ami` flag for use when creating a cluster. This enables a number of advanced use cases such as using a custom AMI or querying AWS in realtime to determine which AMI to use (non-GPU and GPU instances).

The `--node-ami` can take the AMI image id for an image to explicitly use. It also can take the following 'special' keywords:

| Keyword   | Description                                                                                                         |
| --------- | ------------------------------------------------------------------------------------------------------------------- |
| static    | Indicates that the AMI images ids embedded into `eksctl` should be used. This relates to the static resolvers.      |
| auto      | Indicates that the AMI to use for the nodes should be found by querying AWS EC2. This relates to the auto resolver. |
| auto-ssm  | Indicates that the AMI to use for the nodes should be found by querying AWS SSM Parameter Store.                    |

If, for example, AWS release a new version of the EKS node AMIs and a new version of `eksctl` hasn't been released you can use the latest AMI by doing the following:

```
eksctl create cluster --node-ami=auto
```

With the 0.1.9 release we have introduced the `--node-ami-family` flag for use when creating the cluster. This makes it possible to choose between different officially supported EKS AMI families.

The `--node-ami-family` can take following keywords:

| Keyword                        |                                          Description                                         |
|--------------------------------|:--------------------------------------------------------------------------------------------:|
| AmazonLinux2                   | Indicates that the EKS AMI image based on Amazon Linux 2 should be used (default).           |
| Ubuntu1804                     | Indicates that the EKS AMI image based on Ubuntu 18.04 should be used.                       |
| WindowsServer2019FullContainer | Indicates that the EKS AMI image based on Windows Server 2019 Full Container should be used. |
| WindowsServer2019CoreContainer | Indicates that the EKS AMI image based on Windows Server 2019 Core Container should be used. |

<!-- TODO for 0.3.0
To use more advanced configuration options, [Cluster API](https://github.com/kubernetes-sigs/cluster-api):

```
eksctl apply --cluster-config advanced-cluster.yaml
```
-->
