build: update-bindata
	go build ./cmd/eksctl

update-bindata:
	go generate ./pkg/eks
