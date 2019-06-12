package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/957f835f437fd76953f37e6b409ba6426bb4ce37/config/v1.4/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml
//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/957f835f437fd76953f37e6b409ba6426bb4ce37/config/v1.4/aws-k8s-cni-1.10.yaml?raw=1 --output assets/aws-node-1.10.yaml

//go:generate ${GOPATH}/bin/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
