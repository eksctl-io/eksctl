# Addons

EKS Add-Ons is a new feature that lets you enable and manage Kubernetes operational
software for your AWS EKS clusters. At launch, EKS add-ons supports controlling the launch and version of the AWS VPC
CNI plugin through the EKS API

## Creating addons

You can specify what addons you want and what policies (if required) to attach to them in your config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: exmaple-cluster
  region: us-west-2
  version: "1.18"

iam:
  withOIDC: true

addons:
- name: vpc-cni
  version: 1.7.5 # optional
  attachPolicyARNs: #optional
  - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
  serviceAccountRoleARN: arn:aws:iam::aws:policy/AmazonEKSCNIAccess # optional
  tags: # optional
    team: eks
  attachPolicy: # optional
    Statement:
    - Effect: Allow
      Action:
      - ec2:AssignPrivateIpAddresses
      - ec2:AttachNetworkInterface
      - ec2:CreateNetworkInterface
      - ec2:DeleteNetworkInterface
      - ec2:DescribeInstances
      - ec2:DescribeTags
      - ec2:DescribeNetworkInterfaces
      - ec2:DescribeInstanceTypes
      - ec2:DetachNetworkInterface
      - ec2:ModifyNetworkInterfaceAttribute
      - ec2:UnassignPrivateIpAddresses
      Resource: '*'
```

You can specify at most one of `attachPolicy`, `attachPolicyARNs` and `serviceAccountRoleARN`.

!!!note
    In order to attach policies to addons your cluster must have `OIDC` enabled. If it's not enabled we ignore any policies
    attached.


You can then either have these addons created during the cluster creation process:
```console
eksctl create cluster -f config.yaml
```

Or you can create after cluster creation using the config file or CLI flags:

```console
eksctl create addon -f config.yaml
```

```console
eksctl create addon --name vpc-cni --version 1.7.5 --service-account-role-arn=<role-arn>
```

## Listing enabled addons

You can see what addons are enabled in your cluster by running:
```console
eksctl get addons --cluster <cluster-name>
```

## Setting the addon's version

Setting the version of the addon is optional. If the `version` field is empty in the request sent by `eksctl`, the EKS API will set it to the default version for that specific addon. More information about which version is the default version for specific addons can be found in the AWS documentation about EKS. Note that the default version might not necessarily be the latest version available. 

The addon version can be set to `latest`. Alternatively, the version can be set with the EKS build tag specified, such as `v1.7.5-eksbuild.1` or `v1.7.5-eksbuild.2`. It can also be set to the release version of the addon, such as `v1.7.5` or `1.7.5`, and the `eksbuild` suffix tag will be discovered and set for you.

See the section below on how to discover available addons and their versions.

## Discovering addons
You can discover what addons are available to install on your cluster by running:
```console
eksctl utils describe-addon-versions --cluster <cluster-name>
```

This will discover your cluster's kubernetes version and filter on that. Alternatively if you want to see what
addons are available for a particular kubernetes version you can run:
```console
eksctl utils describe-addon-versions --kubernetes-version <version>
```

## Updating addons
You can update your addons to newer versions and change what policies are attached by running:
```console
eksctl update addon -f config.yaml
```

```console
eksctl update addon --name vpc-cni --version 1.8.0 --service-account-role-arn=<new-role>
```

## Deleting addons
You can delete an addon by running:
```console
eksctl delete addon --cluster <cluster-name> --name <addon-name
```
This will delete the addon and any IAM roles associated to it.

When you delete your cluster all IAM roles associated to addons are also deleted.
