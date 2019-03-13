package defaultaddons

//go:generate curl --silent --location https://github.com/aws/amazon-vpc-cni-k8s/blob/12136199717ddb05a75827951d52896b9652c323/config/v1.3/aws-k8s-cni.yaml?raw=1 --output assets/aws-node.yaml
//go:generate curl --silent --location https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/dns.yaml --output assets/coredns.yaml

//go:generate ${GOPATH}/bin/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
