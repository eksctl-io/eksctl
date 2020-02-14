package defaultaddons

//go:generate /bin/sh -c "if [ '${EKSCTL_DOWNLOAD_ASSETS}' = 'true' ]; then curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/61c8d18c0e097c0b7e8477e1afc61d6f2601295d/config/v1.6/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml; fi"
//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
