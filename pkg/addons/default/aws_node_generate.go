package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/fcec46c23d0f8d85245977316cc5894b9eca746e/config/v1.6/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml
