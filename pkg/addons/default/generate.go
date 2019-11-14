package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/79fcb8675a3106424d2df2f8f3fedf41fc38bd4b/config/v1.5/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml

//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
