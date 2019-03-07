package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/d40649ec7ea927d6a17c1d862f08d0bd1772897f/config/v1.3/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml
//go:generate curl --silent --location https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/dns.yaml --output assets/coredns.yaml

//go:generate ${GOPATH}/bin/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
