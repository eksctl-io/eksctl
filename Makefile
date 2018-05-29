build: update-bindata
	go build ./cmd/eksctl

update-bindata:
	go generate ./pkg/eks

install-bindata:
	go get -u github.com/jteeuwen/go-bindata/...
