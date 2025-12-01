# EKS Capabilities

EKS Capabilities enable you to install and manage cloud-native tools on your EKS clusters through the EKS API. Capabilities provide a streamlined way to deploy and configure popular tools like AWS Controllers for Kubernetes (ACK), Kubernetes Resource Optimizer (KRO), and ArgoCD fully managed by Amazon EKS.

## Supported Capability Types

eksctl supports the following capability types:

- **ACK** - AWS Controllers for Kubernetes, enabling management of AWS services from Kubernetes
- **KRO** - Kubernetes Resource Optimizer for resource optimization and management
- **ARGOCD** - ArgoCD GitOps for continuous deployment and application management

## Creating Capabilities

### Using Configuration File

You can define capabilities in your cluster configuration file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-cluster
  region: us-west-2

capabilities:
  # AWS Controllers for Kubernetes (ACK)
  - name: ack-capability
    type: ACK
    deletePropagationPolicy: RETAIN
    attachPolicyARNs:
      - arn:aws:iam::aws:policy/AdministratorAccess
    tags:
      Environment: dev
      Team: platform

  # Kubernetes Resource Optimizer (KRO)
  - name: kro-capability
    type: KRO
    deletePropagationPolicy: RETAIN
    accessPolicies:
      - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy
        accessScope:
          type: cluster
    tags:
      Environment: production
      Team: platform
```

Create capabilities during cluster creation:
```console
eksctl create cluster -f config.yaml
```

Or create capabilities after cluster creation:
```console
eksctl create capability -f config.yaml
```

### Using CLI Flags

Create a capability using command-line flags:

```console
eksctl create capability \
  --cluster my-cluster \
  --name ack-capability \
  --type ACK \
  --attach-policy-arns arn:aws:iam::aws:policy/AdministratorAccess \
  --tags Environment=dev,Team=platform
```

## ArgoCD Configuration

ArgoCD capabilities require additional configuration for AWS Identity Center (IDC) integration:

```yaml
capabilities:
  - name: argocd-gitops
    type: ARGOCD
    deletePropagationPolicy: RETAIN
    attachPolicyARNs:
      - arn:aws:iam::aws:policy/AdministratorAccess
    accessPolicies:
      - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy
        accessScope:
          type: cluster
    configuration:
      argocd:
        namespace: argocd
        awsIdc:
          idcInstanceArn: arn:aws:sso:::instance/ssoins-1234567890abcdef
          idcRegion: us-west-2
        networkAccess:
          vpceIds:
            - vpce-1234567890abcdef0
        rbacRoleMappings:
          - role: ADMIN
            identities:
              - id: 38414300-1041-708a-01af-abc123defg
                type: SSO_USER
              - id: 08017340-8041-7093-3ffa-abc123defg
                type: SSO_GROUP
    tags:
      Environment: production
      Team: platform
```

### ArgoCD Configuration Options

- **namespace**: Kubernetes namespace for ArgoCD installation (optional, defaults to `argocd`)
- **awsIdc**: AWS Identity Center configuration (required for ARGOCD type)
  - **idcInstanceArn**: ARN of the AWS IDC instance
  - **idcRegion**: Region of the IDC instance (optional)
- **networkAccess**: Network access configuration (optional)
  - **vpceIds**: List of VPC endpoint IDs for private access
- **rbacRoleMappings**: RBAC role mappings for ArgoCD (optional)
  - **role**: ArgoCD role (ADMIN, EDITOR, VIEWER)
  - **identities**: List of SSO identities to map to the role

## IAM Configuration

### Using Existing IAM Role

Specify an existing IAM role ARN:

```yaml
capabilities:
  - name: my-capability
    type: ACK
    roleArn: arn:aws:iam::123456789012:role/MyCapabilityRole
```

### Auto-creating IAM Role with Policies

Let eksctl create the IAM role with attached policies:

```yaml
capabilities:
  - name: my-capability
    type: ACK
    attachPolicyARNs:
      - arn:aws:iam::aws:policy/AdministratorAccess
      - arn:aws:iam::123456789012:policy/MyCustomPolicy
```

### Using Inline Policy

Attach an inline policy document:

```yaml
capabilities:
  - name: my-capability
    type: ACK
    attachPolicy:
      Statement:
        - Effect: Allow
          Action:
            - s3:GetObject
            - s3:PutObject
          Resource: "arn:aws:s3:::my-bucket/*"
```

Or via CLI:
```console
eksctl create capability \
  --cluster my-cluster \
  --name my-capability \
  --type ACK \
  --attach-policy '{"Statement":[{"Effect":"Allow","Action":["s3:GetObject"],"Resource":"*"}]}'
```

### Access Policies

Associate EKS access policies with the capability:

```yaml
capabilities:
  - name: my-capability
    type: KRO
    accessPolicies:
      - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy
        accessScope:
          type: cluster
      - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSAdminPolicy
        accessScope:
          type: namespace
          namespaces: ["default", "kube-system"]
```

## Listing Capabilities

View all capabilities in your cluster:

```console
eksctl get capabilities --cluster my-cluster
```

Get a specific capability:

```console
eksctl get capability --cluster my-cluster --name my-capability
```

Output in different formats:

```console
eksctl get capabilities --cluster my-cluster --output yaml
eksctl get capabilities --cluster my-cluster --output json
```

## Updating Capabilities

Update capability configuration using a config file:

```console
eksctl update capability -f config.yaml
```

## Deleting Capabilities

Delete a specific capability:

```console
eksctl delete capability --cluster my-cluster --name my-capability
```

Delete capabilities defined in a config file:

```console
eksctl delete capability -f config.yaml
```

## Configuration Options

### Common Options

All capability types support these common configuration options:

- **name**: Unique name for the capability (required)
- **type**: Capability type - ACK, KRO, or ARGOCD (required)
- **roleArn**: IAM role ARN (optional if IAM policies are provided)
- **deletePropagationPolicy**: Delete propagation policy (RETAIN or DELETE, defaults to RETAIN)
- **tags**: Key-value pairs to tag AWS resources
- **accessPolicies**: EKS access policies to associate with the capability
- **attachPolicyARNs**: List of IAM policy ARNs to attach
- **attachPolicy**: Inline IAM policy document
- **permissionsBoundary**: ARN of the permissions boundary policy

### Delete Propagation Policy

Controls what happens to capability resources when the capability is deleted:

- **RETAIN**: Keep capability resources after deletion (default)
- **DELETE**: Remove all capability resources when deleted

## Prerequisites

- EKS cluster must be in ACTIVE state
- For ArgoCD capabilities: AWS Identity Center (IDC) instance must be configured
- For network access configuration: VPC endpoints must be pre-provisioned
- Appropriate IAM permissions for capability management

## Examples

### Complete ACK Capability Example

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: ack-cluster
  region: us-west-2

autoModeConfig:
  enabled: true

capabilities:
  - name: s3-controller
    type: ACK
    deletePropagationPolicy: RETAIN
    attachPolicyARNs:
      - arn:aws:iam::aws:policy/AmazonS3FullAccess
    accessPolicies:
      - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy
        accessScope:
          type: cluster
    tags:
      Service: s3-controller
      Environment: production
      ManagedBy: eksctl
```

### Complete ArgoCD Capability Example

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: gitops-cluster
  region: us-west-2

autoModeConfig:
  enabled: true

capabilities:
  - name: argocd-gitops
    type: ARGOCD
    deletePropagationPolicy: RETAIN
    attachPolicyARNs:
      - arn:aws:iam::aws:policy/ReadOnlyAccess
    accessPolicies:
      - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy
        accessScope:
          type: cluster
    configuration:
      argocd:
        namespace: argocd-system
        awsIdc:
          idcInstanceArn: arn:aws:sso:::instance/ssoins-1234567890abcdef
          idcRegion: us-west-2
        rbacRoleMappings:
          - role: ADMIN
            identities:
              - id: admin-user-id
                type: SSO_USER
          - role: VIEWER
            identities:
              - id: readonly-group-id
                type: SSO_GROUP
    tags:
      Application: argocd
      Environment: production
```

## Troubleshooting

### Common Issues

1. **Cluster not ready**: Ensure your EKS cluster is in ACTIVE state before creating capabilities
2. **IAM permissions**: Verify you have sufficient IAM permissions to create roles and policies
3. **ArgoCD configuration**: Ensure AWS IDC instance ARN is correct and accessible
4. **Access policy association**: Access policies are associated after the capability is created by EKS

### Checking Capability Status

Monitor capability creation and status:

```console
eksctl get capabilities --cluster my-cluster
```

Check CloudFormation stacks for detailed information:

```console
aws cloudformation describe-stacks --stack-name eksctl-my-cluster-capability-<hash>
```