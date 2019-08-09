---
title: "Command-Line Flags Reference"
weight: 10
---

# Command-Line Flags Reference

## color <a name="color"></a>

`-C, --color string`

toggle colorized logs

**Valid options are**: true, false, fabulous (default "true")

### supported commands
All eksctl commands.

## config-file <a name="config-file"></a>

`-f, --config-file string`

load configuration from a file (or stdin if set to '-')

### supported commands
- [create cluster](01-commands.md#create-common-flags)
- [create iamidentitymapping](01-commands.md#create-common-flags)
- [create nodegroup](01-commands.md#create-common-flags)

## help <a name="help"></a>
`-h, --help`

help for this command

### supported commands
All eksctl commands.

## name
`-n, --name string`

The resource name to which the command will apply.

***Note***: If this value is not specified for `create` commands, the resource will be assigned an auto-generated name (e.g. "unique-mushroom-1565299533").

### supported commands
- [create cluster](01-commands.md#create-cluster)
- [create iamidentitymapping](01-commands.md#create-iamidentitymapping)
  - ***Note***: the string value of the name flag denotes the name of the EKS cluster for which the `create identitymapping` command will create the identity mapping and **NOT** the name of the identity mapping resource.
- [create nodegroup](01-commands.md#create-nodegroup)

## tags
`--tags stringToString`

A list of KV pairs used to tag the AWS resources (default [])

### Example
`--tags "Owner=John Doe,Team=Some Team"`

## verbose <a name="verbose"></a>

`-v, --verbose int`

set log level, use 0 to silence, 4 for debugging and 5 for debugging with AWS debug logging (default 3)

### supported commands
All eksctl commands.

## version <a name="version"></a>

`--version string`

Kubernetes version

### valid values
1.11, 1.12, 1.13 (default 1.13)

### supported commands

- [create cluster](01-commands.md#create-cluster)
- [create nodegroup](01-commands.md#create-nodegroup)
  - ***Note***: Additional valid values
    - **auto**: automatically inherit version from the control plane
    - **latest**: use latest version
    - uses a default value of **auto** if left unspecified


# Scratch space
## create cluster flags

```
      --tags stringToString   A list of KV pairs used to tag the AWS resources (e.g. "Owner=John Doe,Team=Some Team") (default [])
  -r, --region string         AWS region
      --zones strings         (auto-select if unspecified)
      --version string        Kubernetes version (valid options: 1.11, 1.12, 1.13) (default "1.13")
  -f, --config-file string    load configuration from a file (or stdin if set to '-')
      --nodegroup-name string          name of the nodegroup (generated if unspecified, e.g. "ng-946e68f1")
      --without-nodegroup              if set, initial nodegroup will not be created
  -t, --node-type string               node instance type (default "m5.large")
  -N, --nodes int                      total number of nodes (for a static ASG) (default 2)
  -m, --nodes-min int                  minimum nodes in ASG (default 2)
  -M, --nodes-max int                  maximum nodes in ASG (default 2)
      --node-volume-size int           node volume size in GB
      --node-volume-type string        node volume type (valid options: gp2, io1, sc1, st1) (default "gp2")
      --max-pods-per-node int          maximum number of pods per node (set automatically if unspecified)
      --ssh-access                     control SSH access for nodes. Uses ~/.ssh/id_rsa.pub as default key path if enabled
      --ssh-public-key string          SSH public key to use for nodes (import from local path, or use existing EC2 key pair)
      --node-ami string                Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on version/region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care. (default "static")
      --node-ami-family string         Advanced use cases only. If 'AmazonLinux2' is supplied (default), then eksctl will use the official AWS EKS AMIs (Amazon Linux 2); if 'Ubuntu1804' is supplied, then eksctl will use the official Canonical EKS AMIs (Ubuntu 18.04). (default "AmazonLinux2")
  -P, --node-private-networking        whether to make nodegroup networking private
      --node-security-groups strings   Attach additional security groups to nodes, so that it can be used to allow extra ingress/egress access from/to pods
      --node-labels stringToString     Extra labels to add when registering the nodes in the nodegroup, e.g. "partition=backend,nodeclass=hugememory" (default [])
      --node-zones strings             (inherited from the cluster if unspecified)
      --asg-access            enable IAM policy for cluster-autoscaler
      --external-dns-access   enable IAM policy for external-dns
      --full-ecr-access       enable full access to ECR
      --appmesh-access        enable full access to AppMesh
      --alb-ingress-access    enable full access for alb-ingress-controller
      --vpc-cidr ipNet                 global CIDR to use for VPC (default 192.168.0.0/16)
      --vpc-private-subnets strings    re-use private subnets of an existing VPC
      --vpc-public-subnets strings     re-use public subnets of an existing VPC
      --vpc-from-kops-cluster string   re-use VPC from a given kops cluster
      --vpc-nat-mode string            VPC NAT mode, valid options: HighlyAvailable, Single, Disable (default "Single")
  -p, --profile string        AWS credentials profile to use (overrides the AWS_PROFILE environment variable)
      --timeout duration      max wait time in any polling operations (default 25m0s)
      --cfn-role-arn string   IAM role used by CloudFormation to call AWS API on your behalf
      --kubeconfig string               path to write kubeconfig (incompatible with --auto-kubeconfig) (default "/Users/sebastianbernheim/.kube/config")
      --authenticator-role-arn string   AWS IAM role to assume for authenticator
      --set-kubeconfig-context          if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten (default true)
      --auto-kubeconfig                 save kubeconfig file by cluster name, e.g. "/Users/sebastianbernheim/.kube/eksctl/clusters/unique-mushroom-1565299533"
      --write-kubeconfig                toggle writing of kubeconfig (default true)
```

## create identitymapping flags

```
      --role string          ARN of the IAM role to create
      --username string      User name within Kubernetes to map to IAM role
      --group stringArray    Group within Kubernetes to which IAM role is mapped
  -n, --name string          EKS cluster name
  -r, --region string        AWS region
  -f, --config-file string   load configuration from a file (or stdin if set to '-')
  -p, --profile string     AWS credentials profile to use (overrides the AWS_PROFILE environment variable)
      --timeout duration   max wait time in any polling operations (default 25m0s)

```

## create nodegroup flags

```
      --cluster string          name of the EKS cluster to add the nodegroup to
  -r, --region string           AWS region
      --version string          Kubernetes version (valid options: 1.11, 1.12, 1.13) [for nodegroups "auto" and "latest" can be used to automatically inherit version from the control plane or force latest] (default "auto")
  -f, --config-file string      load configuration from a file (or stdin if set to '-')
      --include strings         nodegroups to include (list of globs), e.g.: 'ng-team-?,prod-*'
      --exclude strings         nodegroups to exclude (list of globs), e.g.: 'ng-team-?,prod-*'
      --update-auth-configmap   Remove nodegroup IAM role from aws-auth configmap (default true)

New nodegroup flags:
  -n, --name string                    name of the new nodegroup (generated if unspecified, e.g. "ng-f06b88af")
  -t, --node-type string               node instance type (default "m5.large")
  -N, --nodes int                      total number of nodes (for a static ASG) (default 2)
  -m, --nodes-min int                  minimum nodes in ASG (default 2)
  -M, --nodes-max int                  maximum nodes in ASG (default 2)
      --node-volume-size int           node volume size in GB
      --node-volume-type string        node volume type (valid options: gp2, io1, sc1, st1) (default "gp2")
      --max-pods-per-node int          maximum number of pods per node (set automatically if unspecified)
      --ssh-access                     control SSH access for nodes. Uses ~/.ssh/id_rsa.pub as default key path if enabled
      --ssh-public-key string          SSH public key to use for nodes (import from local path, or use existing EC2 key pair)
      --node-ami string                Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on version/region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care. (default "static")
      --node-ami-family string         Advanced use cases only. If 'AmazonLinux2' is supplied (default), then eksctl will use the official AWS EKS AMIs (Amazon Linux 2); if 'Ubuntu1804' is supplied, then eksctl will use the official Canonical EKS AMIs (Ubuntu 18.04). (default "AmazonLinux2")
  -P, --node-private-networking        whether to make nodegroup networking private
      --node-security-groups strings   Attach additional security groups to nodes, so that it can be used to allow extra ingress/egress access from/to pods
      --node-labels stringToString     Extra labels to add when registering the nodes in the nodegroup, e.g. "partition=backend,nodeclass=hugememory" (default [])
      --node-zones strings             (inherited from the cluster if unspecified)

IAM addons flags:
      --asg-access            enable IAM policy for cluster-autoscaler
      --external-dns-access   enable IAM policy for external-dns
      --full-ecr-access       enable full access to ECR
      --appmesh-access        enable full access to AppMesh
      --alb-ingress-access    enable full access for alb-ingress-controller

AWS client flags:
  -p, --profile string        AWS credentials profile to use (overrides the AWS_PROFILE environment variable)
      --timeout duration      max wait time in any polling operations (default 25m0s)
      --cfn-role-arn string   IAM role used by CloudFormation to call AWS API on your behalf

```