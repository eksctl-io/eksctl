package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/61c8d18c0e097c0b7e8477e1afc61d6f2601295d/config/v1.6/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml

//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
