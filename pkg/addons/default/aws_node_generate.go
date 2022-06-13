package defaultaddons

// Please refer to https://docs.aws.amazon.com/eks/latest/userguide/cni-upgrades.html
//go:generate curl --silent --location https://raw.githubusercontent.com/aws/amazon-vpc-cni-k8s/v1.11.2/config/master/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml
