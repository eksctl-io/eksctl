build:
	go build ./cmd/eksctl

update-bindata:
	go generate ./pkg/cfn
