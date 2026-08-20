# Updating control plane subnets and security groups

## Updating control plane subnets
When a cluster is created with eksctl, a set of public and private subnets are created and passed to the EKS API.
EKS creates 2 to 4 cross-account elastic network interfaces (ENIs) in those subnets to enable communication between the EKS
managed Kubernetes control plane and your VPC.

To update the subnets used by the EKS control plane, run:

```console
eksctl utils update-cluster-vpc-config --cluster=<cluster> --control-plane-subnet-ids=subnet-1234,subnet-5678
```

To update the setting using a config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: cluster
  region: us-west-2

vpc:
  controlPlaneSubnetIDs: [subnet-1234, subnet-5678]
```

```console
eksctl utils update-cluster-vpc-config -f config.yaml
```

Without the `--approve` flag, eksctl only logs the proposed changes. Once you are satisfied with the proposed changes, rerun the command with
the  `--approve` flag.

## Restricting the control plane to private subnets at creation time

`eksctl utils update-cluster-vpc-config` changes the control plane subnets of an existing cluster, and it does so by
calling the EKS API directly rather than by updating the cluster's CloudFormation stack. That leaves the stack out of sync
with the actual configuration.

To place the control plane on private subnets from the outset, set `vpc.controlPlaneOnPrivateSubnets` when creating the
cluster. Only the private subnets are then passed to the EKS API, so the cross-account ENIs are never created in public
subnets, and the CloudFormation stack reflects the intended configuration from the start:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: cluster
  region: us-west-2

availabilityZones:
  - us-west-2a
  - us-west-2b

vpc:
  controlPlaneOnPrivateSubnets: true
```

```console
eksctl create cluster -f config.yaml
```

Public subnets are still created and are still used for NAT gateways, internet-facing load balancers, and any nodegroup
that is not private, so this setting does not turn the cluster into a fully-private one. To create a cluster with no
public subnets at all, see [EKS Private Cluster without Outbound Internet Access](/usage/eks-private-cluster/).

Note the following when using this setting:

- At least two private subnets across at least two availability zones are required, because EKS places its cross-account
  ENIs in a minimum of two zones. eksctl validates this before creating anything, except when pre-existing subnets are
  specified only by ID: their availability zone is not known until eksctl looks them up in EC2, which happens after
  validation, so insufficient AZ coverage is instead rejected by the EKS API.
- The private subnets must have a NAT gateway or the relevant VPC endpoints, otherwise nodes cannot reach the API server.
  When eksctl creates the VPC this is handled by the `vpc.nat` configuration.
- It cannot be combined with `vpc.controlPlaneSubnetIDs`; specify one or the other.
- It applies both when eksctl creates the VPC and when a pre-existing VPC is supplied via `vpc.id`.

## Updating control plane security groups
To manage traffic between the control plane and worker nodes, EKS supports passing additional security groups that are applied to the cross-account network interfaces
provisioned by EKS. To update the security groups for the EKS control plane, run:

```console
eksctl utils update-cluster-vpc-config --cluster=<cluster> --control-plane-security-group-ids=sg-1234,sg-5678 --approve
```

To update the setting using a config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: cluster
  region: us-west-2

vpc:
  controlPlaneSecurityGroupIDs: [sg-1234, sg-5678]
```

```console
eksctl utils update-cluster-vpc-config -f config.yaml
```

To update both control plane subnets and security groups for a cluster, run:

```console
eksctl utils update-cluster-vpc-config --cluster=<cluster> --control-plane-subnet-ids=<> --control-plane-security-group-ids=<> --approve
```

To update both fields using a config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: cluster
  region: us-west-2

vpc:
  controlPlaneSubnetIDs: [subnet-1234, subnet-5678]
  controlPlaneSecurityGroupIDs: [sg-1234, sg-5678]
```

```console
eksctl utils update-cluster-vpc-config -f config.yaml
```

For a complete example, refer to [cluster-subnets-sgs.yaml](https://github.com/eksctl-io/eksctl/blob/main/examples/38-cluster-subnets-sgs.yaml).