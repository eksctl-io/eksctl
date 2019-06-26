package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/dd631108e61a977809f9b1a1c40232637e734184/config/v1.5/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml

//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
