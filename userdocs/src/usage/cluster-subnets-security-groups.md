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