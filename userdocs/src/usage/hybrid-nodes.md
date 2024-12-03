# EKS Hybrid Nodes

## Introduction

AWS EKS introduces Hybrid Nodes, a new feature that enables you to run on-premises and edge applications on customer-managed infrastructure with the same AWS EKS clusters, features, and tools you use in the AWS Cloud. AWS EKS Hybird Nodes brings an AWS-managed Kubernetes experience to on-premises environments for customers to simplify and standardize how you run applications across on-premises, edge and cloud environments. Read more at [EKS Hybrid Nodes][eks-hybrid-nodes].

To facilitate support for this feature, eksctl introduces a new top-level field called `remoteNetworkConfig`. Any Hybrid Nodes related configuration shall be set up via this field, as part of the config file; there are no CLI flags counterparts. Additionally, at launch, any remote network config can only be set up during cluster creation and cannot be updated afterwards. This means, you won't be able to update existing clusters to use Hybrid Nodes. 

The `remoteNetworkConfig` section of the config file allows you to setup the two core areas when it comes to joining remote nodes to you EKS clusters: **networking** and **credentials**.  

## Networking

EKS Hybrid Nodes is ﬂexible to your preferred method of connecting your on-premises network(s) to an AWS VPC. There are several [documented options](https://docs.aws.amazon.com/whitepapers/latest/aws-vpc-connectivity-options/network-to-amazon-vpc-connectivity-options.html) available, including AWS Site-to-Site VPN and AWS Direct Connect, and you can choose the method that best fits your use case. In most of the methods you might choose, your VPC will be attached to either a virtual private gateway (VGW) or a transit gateway (TGW). If you rely on eksctl to create a VPC for you, eksctl will also configure, **within the scope of your VPC**, any networking related pre-requisites in order to facilitate communication between your EKS control plane and the remote nodes i.e.

- ingress/egress SG rules
- routes in the private subnets' route tables
- the VPC gateway attachment to the given TGW or VGW

Example config file:
 
```yaml
remoteNetworkConfig:
  vpcGatewayID: tgw-xxxx # either VGW or TGW to be attached to your VPC
  remoteNodeNetworks:
    # eksctl will create, behind the scenes, SG rules, routes, and a VPC gateway attachment,
    # to facilitate communication between remote network(s) and EKS control plane, via the attached gateway
    - cidrs: ["10.80.146.0/24"]
  remotePodNetworks:
    - cidrs: ["10.86.30.0/23"]
```

If your connectivity method of choice does not involve using a TGW or VGW, you must not rely on eksctl to create the VPC for you, and instead provide a pre-existing one. On a related note, if you are using a pre-existing VPC, eksctl won't make any amendments to it, and ensuring all networking requirements are in place falls under your responsibility.

???+ note
    eksctl does not setup any networking infrastructure outside your AWS VPC (i.e. any infrastructure from VGW/TGW to the remote networks)

## Credentials

EKS Hybrid Nodes use the AWS IAM Authenticator and temporary IAM credentials provisioned by either **AWS SSM** or **AWS IAM Roles Anywhere**
to authenticate with the EKS cluster. Similar to the self-managed nodegroups, if not otherwise provided, eksctl will create for you a Hybrid Nodes IAM Role to be assumed by the remote nodes. Additioanlly, when using IAM Roles Anywhere as your credentials provider, eksctl will setup a profile, and trust anchor based on a given certificate authority bundle (`iam.caBundleCert`) e.g.

```yaml
remoteNetworkConfig:
  iam:
    # the provider for temporary IAM credentials. Default is SSM.
    provider: IRA
    # the certificate authority bundle that serves as the root of trust,
    # used to validate the X.509 certificates provided by your nodes.
    # can only be set when provider is IAMRolesAnywhere.
    caBundleCert: xxxx
```

The ARN of the Hybrid Nodes Role created by eksctl is needed later in the process of joining your remote nodes to the cluster, to setup `NodeConfig` for `nodeadm`, and to create activations (if using SSM). To fetch it, use:

```bash
aws cloudformation describe-stacks \        
  --stack-name eksctl-<CLUSTER_NAME>-cluster \
  --query 'Stacks[].Outputs[?OutputKey==`RemoteNodesRoleARN`].[OutputValue]' \
  --output text
```

Similarly, if using IAM Roles Anywhere, you can fetch the ARN of the trust anchor and of the anywhere profile created by eksctl, amending the previous command by replacing `RemoteNodesRoleARN` with `RemoteNodesTrustAnchorARN` or `RemoteNodesAnywhereProfileARN`, respectively.

If you have a pre-existing IAM Roles Anywhere configuration in place, or you are using SSM, you can provide a IAM Role for Hybrid nodes via `remoteNetworkConfig.iam.roleARN`. Bear in mind that in this scenario, eksctl won't create the trust anchor and anywhere profile for you. e.g. 

```yaml
remoteNetworkConfig:
  iam:
    roleARN: arn:aws:iam::000011112222:role/HybridNodesRole
```

To map the role to a Kubernetes identity and authorise the remote nodes to join the EKS cluster, eksctl creates an access entry with Hybrid Nodes IAM Role as principal ARN and of type `HYBRID_LINUX`. i.e.

```bash
eksctl get accessentry --cluster my-cluster --principal-arn arn:aws:iam::000011112222:role/eksctl-my-cluster-clust-HybridNodesSSMRole-XiIAg0d29PkO --output json
[
    {
        "principalARN": "arn:aws:iam::000011112222:role/eksctl-my-cluster-clust-HybridNodesSSMRole-XiIAg0d29PkO",
        "kubernetesGroups": [
            "system:nodes"
        ]
    }
]
```

## Add-ons support

Container Networking Interface (CNI): The AWS VPC CNI can’t be used with hybrid nodes. The core capabilities of Cilium and Calico are supported for use with hybrid nodes. You can manage your CNI with your choice of tooling such as Helm. For more information, see [Configure a CNI for hybrid nodes](https://docs.aws.amazon.com/eks/latest/userguide/hybrid-nodes-cni.html).

???+ note
    If you install VPC CNI in your cluster for your self-managed or EKS-managed nodegroups, you have to use `v1.19.0-eksbuild.1` or later, as this includes an udpate to the add-on's daemonset to exclude it from being installed on Hybrid Nodes.

## Further references

- [EKS Hybrid Nodes UserDocs][eks-hybrid-nodes]
- [Launch Announcement][launch-announcement]

[eks-hybrid-nodes]: https://docs.aws.amazon.com/eks/latest/userguide/hybrid-nodes-overview.html
[launch-announcement]: https://aws.amazon.com/about-aws/whats-new/2024/12/amazon-eks-hybrid-nodes