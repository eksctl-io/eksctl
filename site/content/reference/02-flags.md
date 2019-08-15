---
title: "Flags reference"
weight: 20
---

# Flags reference

## alb-ingress-access

`--alb-ingress-access`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## appmesh-access

`--appmesh-access`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## asg-access

`--appmesh-access`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## authenticator-role-arn

`--authenticator-role-arn`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## auto-kubeconfig

`--auto-kubeconfig`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## cfn-role-arn

`--cfn-role-arn`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## cluster

`--cluster`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create nodegroup](../01-commands#create-nodegroup)

## color

`-C, --color string`

toggle colorized logs

**Valid options are**: true, false, fabulous (default "true")

### supported commands

All eksctl commands.

## config-file

`-f, --config-file string`

load configuration from a file (or stdin if set to '-')

### supported commands

- [create cluster](../01-commands#create-cluster)
- [create iamidentitymapping](../01-commands#create-iamidentitymapping)
- [create nodegroup](../01-commands#create-nodegroup)

## enable-types

`--enable-types stringArray`

List of CloudWatch logging types to enable

### supported commands

- [utils update-cluster-logging](../01-commands#utils-update-cluster-logging)

### config yaml

- [clusterLogging](../03-config-yaml#clusterLogging)
- [enableTypes](../03-config-yaml#enableTypes)

## exclude

`--exclude`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create nodegroup](../01-commands#create-nodegroup)

## external-dns-access

`--external-dns-access`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## full-ecr-access

`--full-ecr-access`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## group

`--group stringArray`

Kubernetes group to which `eksctl` will map an IAM role.

### supported commands

- [create iamidentitymapping](../01-commands#create-iamidentitymapping)

## help

`-h, --help`

help for this command

### supported commands

All eksctl commands.

## include

`--include`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create nodegroup](../01-commands#create-nodegroup)

## kubeconfig

`--kubeconfig`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## max-pods-per-node

`--max-pods-per-node`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## name

`-n, --name string`

The resource name to which the command will apply.

***Note***: If this value is not specified for `create` commands, the resource will be assigned an auto-generated name (e.g. "unique-mushroom-1565299533").

### supported commands

- [create cluster](../01-commands#create-cluster)
- [create iamidentitymapping](../01-commands#create-iamidentitymapping)
  - ***Note***: the string value of the name flag denotes the name of the EKS cluster for which the `create identitymapping` command will create the identity mapping and **NOT** the name of the identity mapping resource.
- [create nodegroup](../01-commands#create-nodegroup)
  - ***Note***: the auto-generated name for a nodegroup will start with "ng" (e.g. "ng-f06b88af")

## node-ami

`--node-ami`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-ami-family

`--node-ami-family`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-labels

`--node-labels`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-private-networking

`-P, --node-private-networking`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-type

`-t, --node-type`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-security-groups

`--node-security-groups`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-volume-size

`--node-volume-size`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-volume-type

`--node-volume-type`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## node-zones

`--node-zones strings`

AWS availability zones for an EKS cluster nodegroup.

Inherited from the cluster if unspecified.

### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)
  - Inherited from the cluster if unspecified.

### config yaml

- [availabilityZones](../03-config-yaml#availabilityZones)

## nodes

`-N, --nodes`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## nodes-max

`-M, --nodes-max`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## nodes-min

`-m, --nodes-min`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## profile

`-p, --profile string`

AWS credentials profile to use (overrides the AWS_PROFILE environment variable)

### supported commands

- [create cluster](../01-commands#create-cluster)
- [create iamidentitymapping](../01-commands#create-iamidentitymapping)
- [create nodegroup](../01-commands#create-nodegroup)

## region

`-r, --region string`

AWS region

### supported commands

- [create cluster](../01-commands#create-cluster)
- [create iamidentitymapping](../01-commands#create-iamidentitymapping)
- [create nodegroup](../01-commands#create-nodegroup)

## role

`--role string`

ARN of the IAM role to create.

### supported commands

- [create iamidentitymapping](../01-commands#create-iamidentitymapping)

## set-kubeconfig-context

`--set-kubeconfig-context`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## ssh-access

`--ssh-access`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## ssh-public-key

`--ssh-public-key`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)

## tags

`--tags stringToString`

A comma-separated list of KV pairs used to tag the AWS resources.  By default, no tags will be applied.

If your tag values include spaces, be sure to wrap the value string for this flag in quotes.

### example

```
--tags "Owner=John Doe,Team=Some Team"
```

### supported commands

- [create cluster](../01-commands#create-cluster)

## timeout

`--timeout duration`

Max wait time in any polling operations (default 25m0s)

### supported commands

- [create cluster](../01-commands#create-cluster)
- [create iamidentitymapping](../01-commands#create-iamidentitymapping)
- [create nodegroup](../01-commands#create-nodegroup)

## update-auth-configmap

`--update-auth-configmap`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create nodegroup](../01-commands#create-nodegroup)

## username

`--username string`

User name within Kubernetes to map to IAM role.

### supported commands

- [create iamidentitymapping](../01-commands#create-iamidentitymapping)

## verbose

`-v, --verbose int`

set log level, use 0 to silence, 4 for debugging and 5 for debugging with AWS debug logging (default 3)

### supported commands

All eksctl commands.

## version

`--version string`

Kubernetes version

### valid values

1.11, 1.12, 1.13 (default 1.13)

### supported commands

- [create cluster](../01-commands#create-cluster)
- [create nodegroup](../01-commands#create-nodegroup)
  - ***Note***: Additional valid values for nodegroups
    - **auto**: automatically inherit version from the control plane
    - **latest**: use latest version
    - defaults to **auto** for nodegroups

## vpc-cidr

`--vpc-cidr`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## vpc-from-kops-cluster

`--vpc-from-kops-cluster`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## vpc-nat-mode

`--vpc-nat-mode`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## vpc-private-subnets

`--vpc-private-subnets`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## vpc-public-subnets

`--vpc-public-subnets`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## without-nodegroup

`--without-nodegroup`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## write-kubeconfig

`--write-kubeconfig`
<!-- FIXME(bianca, sebastian): add description -->
### supported commands

- [create cluster](../01-commands#create-cluster)

## zones

`--zones strings`

AWS zones associated with the EKS cluster.  By default, eksctl will auto-select appropriate zones for the AWS region.

***Note***: In the `us-east-1` region you are likely to get `UnsupportedAvailabilityZoneException`. If you do, copy the suggested zones and pass them in as values for the `--zones` flag. This may occur in other regions, but less likely. You shouldn’t need to use --zone flag otherwise.

### example

```
eksctl create cluster --region=us-east-1 --zones=us-east-1a,us-east-1b,us-east-1d
```

### supported commands

- [create cluster](../01-commands#create-cluster)

### config yaml

- [availabilityZones](../03-config-yaml#availabilityZones)

# Scratch space
## create cluster flags

```
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
      --cfn-role-arn string   IAM role used by CloudFormation to call AWS API on your behalf
      --kubeconfig string               path to write kubeconfig (incompatible with --auto-kubeconfig) (default "/Users/sebastianbernheim/.kube/config")
      --authenticator-role-arn string   AWS IAM role to assume for authenticator
      --set-kubeconfig-context          if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten (default true)
      --auto-kubeconfig                 save kubeconfig file by cluster name, e.g. "/Users/sebastianbernheim/.kube/eksctl/clusters/unique-mushroom-1565299533"
      --write-kubeconfig                toggle writing of kubeconfig (default true)
```

## create identitymapping flags

```
done
```

## create nodegroup flags

```
      --cluster string          name of the EKS cluster to add the nodegroup to
      --include strings         nodegroups to include (list of globs), e.g.: 'ng-team-?,prod-*'
      --exclude strings         nodegroups to exclude (list of globs), e.g.: 'ng-team-?,prod-*'
      --update-auth-configmap   Remove nodegroup IAM role from aws-auth configmap (default true)

New nodegroup flags:
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
      --cfn-role-arn string   IAM role used by CloudFormation to call AWS API on your behalf
```
