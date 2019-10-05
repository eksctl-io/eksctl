package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/release-1.5.4/config/v1.5/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml

//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
