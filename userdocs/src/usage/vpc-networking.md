# Networking

By default `eksctl create cluster` will create a dedicated VPC for the cluster.
This is done in order to avoid interference with existing resources for a
variety of reasons, including security, but also because it is challenging to detect all settings in an existing VPC.

The default VPC CIDR used by `eksctl` is `192.168.0.0/16`. It is divided into 8 (`/19`) subnets (3 private, 3 public & 2 reserved).
The initial nodegroup is created in public subnets, with SSH access disabled unless `--allow-ssh` is specified.
The nodegroup by default allows inbound traffic from the control plane security group on ports 1025 - 65535.

!!! note
    In `us-east-1` eksctl only creates 2 public and 2 private subnets by default.

!!! important
    From `eksctl` version `0.17.0` and onwards public subnets will have the property `MapPublicIpOnLaunch` enabled, and
    the property `AssociatePublicIpAddress` disabled in the Auto Scaling Group for the nodegroups. This means that when
    creating a **new nodegroup** on a **cluster made with an earlier version** of `eksctl`, the nodegroup must **either** be private
    **or** have `MapPublicIpOnLaunch` enabled in its public subnets. Without one of these, the new nodes won't have access to
    the internet and won't be able to download the basic add-ons (CNI plugin, kube-proxy, etc.). To help set up
    subnets correctly for old clusters you can use the new command `eksctl utils update-legacy-subnet-settings`.

If the default functionality doesn't suit you, the following sections explain how to customize VPC configuration further:

- [VPC Configuration](vpc-configuration.md)
- [Subnet Settings](vpc-subnet-settings.md)
- [Cluster Access](vpc-cluster-access.md)
- [IP Family](vpc-ip-family.md)
