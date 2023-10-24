# EKS Fully-Private Cluster

eksctl supports creation of fully-private clusters that have no outbound internet access and have only private subnets.
VPC endpoints are used to enable private access to AWS services.

This guide describes how to create a private cluster without outbound internet access.


## Creating a fully-private cluster

The only required field to create a fully-private cluster is `privateCluster.enabled`:

```yaml
privateCluster:
  enabled: true
```

???+ note
    Post cluster creation, eksctl commands that need access to the Kubernetes API server will have to be run from within the cluster's VPC, a peered VPC or using some other means like AWS Direct Connect. eksctl commands that need access to the EKS APIs will not work if they're being run from within the cluster's VPC. To fix this, [create an interface endpoint for Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/vpc-interface-endpoints.html) to privately access the Amazon Elastic Kubernetes Service (Amazon EKS) management APIs from your Amazon Virtual Private Cloud (VPC). In a future release, eksctl will add support to create this endpoint so it does not need to be manually created.
    Commands that need access to the OpenID Connect provider URL will need to be run from outside of your cluster's VPC once you've enabled AWS PrivateLink for Amazon EKS.
    Creating managed nodegroups will continue to work, and creating self-managed nodegroups will work as it needs access to the API server via the EKS [interface endpoint](https://aws.amazon.com/about-aws/whats-new/2022/12/amazon-eks-supports-aws-privatelink/) if the command is run from within the cluster's VPC, a peered VPC or using some other means like AWS Direct Connect.

???+ info
    VPC endpoints are charged by the hour and based on usage. More details about pricing can be found at
    [AWS PrivateLink pricing](https://aws.amazon.com/privatelink/pricing/)
    
!!! warning 
    Fully-private clusters are not supported in `eu-south-1`.

## Configuring private access to additional AWS services

To enable worker nodes to access AWS services privately, eksctl creates VPC endpoints for the following services:

- Interface endpoints for ECR (both `ecr.api` and `ecr.dkr`) to pull container images (AWS CNI plugin etc)
- A gateway endpoint for S3 to pull the actual image layers
- An interface endpoint for EC2 required by the `aws-cloud-provider` integration
- An interface endpoint for STS to support Fargate and IAM Roles for Services Accounts (IRSA)
- An interface endpoint for CloudWatch logging (`logs`) if CloudWatch logging is enabled

These VPC endpoints are essential for a functional private cluster, and as such, eksctl does not support configuring or
disabling them. However, a cluster might need private access to other AWS services (e.g., Autoscaling required by the Cluster Autoscaler).
These services can be specified in `privateCluster.additionalEndpointServices`, which instructs eksctl to create a VPC endpoint
for each of them.


For example, to allow private access to Autoscaling and CloudWatch logging:

```yaml
privateCluster:
  enabled: true
  additionalEndpointServices:
  # For Cluster Autoscaler
  - "autoscaling"
  # CloudWatch logging
  - "logs"
```

The endpoints supported in `additionalEndpointServices` are `autoscaling`, `cloudformation` and `logs`.

### Skipping endpoint creations

If a VPC has already been created with the necessary AWS endpoints set up and linked to the subnets described in the EKS documentation,
`eksctl` can skip creating them by providing the option `skipEndpointCreation` like this:

```yaml
privateCluster:
  enabled: true
  skipEndpointCreation: true
```

_Note_: this setting cannot be used together with `additionalEndpointServices`. It will skip all endpoint creation. Also, this setting is
only recommended if the endpoint <-> subnet topology is correctly set up. I.e.: subnet ids are correct, `vpce` routing is set up with prefix addresses,
all the necessary EKS endpoints are created and linked to the provided VPC. `eksctl` will not alter any of these resources.

## Nodegroups
Only private nodegroups (both managed and self-managed) are supported in a fully-private cluster because the cluster's VPC is created without
any public subnets. The `privateNetworking` field (`nodeGroup[*].privateNetworking` and `managedNodeGroup[*].privateNetworking`) must be
explicitly set. It is an error to leave `privateNetworking` unset in a fully-private cluster.


```yaml
nodeGroups:
- name: ng1
  instanceType: m5.large
  desiredCapacity: 2
  # privateNetworking must be explicitly set for a fully-private cluster
  # Rather than defaulting this field to `true`,
  # we require users to explicitly set it to make the behaviour
  # explicit and avoid confusion.
  privateNetworking: true

managedNodeGroups:
- name: m1
  instanceType: m5.large
  desiredCapacity: 2
  privateNetworking: true
```

## Cluster Endpoint Access
A fully-private cluster does not support modifying `clusterEndpointAccess` during cluster creation.
It is an error to set either `clusterEndpoints.publicAccess` or `clusterEndpoints.privateAccess`, as a fully-private cluster
can have private access only, and allowing modification of these fields can break the cluster.


## User-supplied VPC and subnets
eksctl supports creation of fully-private clusters using a pre-existing VPC and subnets. Only private subnets can be
specified and it's an error to specify subnets under `vpc.subnets.public`.

eksctl creates VPC endpoints in the supplied VPC and modifies route tables for the supplied subnets. Each subnet should
have an explicit route table associated with it because eksctl does not modify the main route table.

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: private-cluster
  region: us-west-2

privateCluster:
  enabled: true
  additionalEndpointServices:
  - "autoscaling"

vpc:
  subnets:
    private:
      us-west-2b:
        id: subnet-0818beec303f8419b
      us-west-2c:
        id: subnet-0d42ef09490805e2a
      us-west-2d:
        id: subnet-0da7418077077c5f9


nodeGroups:
- name: ng1
  instanceType: m5.large
  desiredCapacity: 2
  # privateNetworking must be explicitly set for a fully-private cluster
  # Rather than defaulting this field to true for a fully-private cluster, we require users to explicitly set it
  # to make the behaviour explicit and avoid confusion.
  privateNetworking: true

managedNodeGroups:
- name: m1
  instanceType: m5.large
  desiredCapacity: 2
  privateNetworking: true
```

## Managing a fully-private cluster

For all commands to work post cluster creation, eksctl will need private access to the EKS API server endpoint, and outbound
internet access (for `EKS:DescribeCluster`). Commands that do not need access to the API server will be supported if eksctl has
outbound internet access.

## Force-delete a fully-private cluster

Errors are likely to occur when deleting a fully-private cluster through eksctl since eksctl does not automatically have access to all of the cluster's resources. `--force` exists to solve this: it will force delete the cluster and continue when errors occur.

## Limitations
A limitation of the current implementation is that eksctl initially creates the cluster with both public and private endpoint
access enabled, and disables public endpoint access after all operations have completed.
This is required because eksctl needs access to the Kubernetes API server to allow self-managed nodes to join the cluster and
to support GitOps and Fargate. After these operations have completed, eksctl switches the cluster endpoint access to private-only.
This additional update does mean that creation of a fully-private cluster will take longer than for a standard cluster.
In the future, eksctl may switch to a VPC-enabled Lambda function to perform these API operations.


## Outbound access via HTTP proxy servers
eksctl is able to talk to the AWS APIs via a configured HTTP(S) proxy server,
however you will need to ensure you set your proxy exclusion list correctly.

Generally, you will need to ensure that requests for the VPC endpoint for your
cluster are not routed via your proxies by setting an appropriate `no_proxy`
environment variable including the value `.eks.amazonaws.com`.

If your proxy server performs "SSL interception" and you are using IAM Roles
for Service Accounts (IRSA), you will need to ensure that you explicitly bypass
SSL Man-in-the-Middle for the domain `oidc.<region>.amazonaws.com`. Failure to
do so will result in eksctl obtaining the incorrect root certificate thumbprint
for the OIDC provider, and the AWS VPC CNI plugin will fail to start due to
being unable to obtain IAM credentials, rendering your cluster inoperative.


## Further information

- [EKS Private Clusters][eks-private-clusters]

[eks-private-clusters]: https://docs.aws.amazon.com/eks/latest/userguide/private-clusters.html

