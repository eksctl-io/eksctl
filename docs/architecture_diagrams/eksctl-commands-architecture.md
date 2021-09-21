# eksctl create cluster
![image](https://user-images.githubusercontent.com/19669095/133767149-d92f0669-18b6-4829-8712-996976781914.png)
Created using https://www.websequencediagrams.com/ and saved to my personal account (attached the raw-contents below)
raw value:
```
title eksctl create cluster
note right of eksctl: All AWS API calls use the aws-sdk-go library
eksctl -> STS API: Authenticates and establishes a session
eksctl -> CF API: creates CloudFormation stack for EKS cluster, VPC and Roles
note right of eksctl: Uses aws-iam-authenticator library to create an in-memory kubeconfig using our session
eksctl -> EKS Cluster: Creates/Modifies Kubernetes Resources
eksctl -> CF API: creates CloudFormation stack for Nodegroups
eksctl -> EKS API: Create EKS Addons
eksctl -> User Machine: Write Kubeconfig 
```

# eksctl create nodegroup
![eksctl create nodegroup (1)](https://user-images.githubusercontent.com/19669095/134142254-7ecf44be-d41c-4552-8048-4aa99050c7c9.png)

raw value:
```
title eksctl create nodegroup
note right of eksctl: All AWS API calls use the aws-sdk-go library
eksctl -> STS API: Authenticates and establishes a session
eksctl -> EKS API: Checks cluster status
eksctl -> CF API: Fetches VPC output from cluster stack
eksctl -> SSM API: Discover recommended AMI ID
eksctl -> CF API: Create stack for Nodegroup
note right of eksctl: Uses aws-iam-authenticator library to create an in-memory kubeconfig using our session
eksctl -> EKS Cluster: Polls Kubernetes until nodegroups report ready
```

# eksctl create iamserviceaccount
![eksctl create iamserviceaccount](https://user-images.githubusercontent.com/19669095/134138728-afd1e31d-1b14-4d26-acff-21ae632094d1.png)

raw value:
```
title eksctl create iamserviceaccount
note right of eksctl: All AWS API calls use the aws-sdk-go library
eksctl -> STS API: Authenticates and establishes a session
eksctl -> EKS API: Checks cluster status
eksctl -> IAM API: Checks OIDC is enabled
eksctl -> CF API: Create stack for IAM Role
note right of eksctl: Uses aws-iam-authenticator library to create an in-memory kubeconfig using our session
eksctl -> EKS Cluster: Creates and/or annotates Service Account
```

# eksctl delete cluster (including nodegroup  & IAMServiceAccount)
![eksctl delete cluster (including nodegroups   IAMServiceAccounts)](https://user-images.githubusercontent.com/19669095/134141603-a69874a6-b82d-43f4-ac3a-a99bdb826433.png)

raw value:
```
title eksctl delete cluster (including nodegroups & IAMServiceAccounts)

note right of eksctl: All AWS API calls use the aws-sdk-go library
eksctl -> STS API: Authenticates and establishes a session
eksctl -> EKS API: Checks cluster status
eksctl -> EKS API: Delete fargate profiles
eksctl -> CF API: Delete fargate IAM role stacks
eksctl -> EC2 API: Delete SSH key pairs
eksctl -> Local Machine: Delete kubeconfig entry for cluster
eksctl -> EC2 API: Delete any orphaned loadbalancer SGs
eksctl -> ELB API: Delete any orphaned loadbalancer
eksctl -> IAM API: Delete OIDC provider
eksctl -> CF API: Delete IAMServiceaccout, Nodegroup & Cluster stacks
eksctl -> EC2 API: Delete orphaned ENIs
```


