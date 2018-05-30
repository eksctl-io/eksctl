COMMIT:=$(shell git describe --dirty --always)

build: update-bindata
	go build -ldflags "-X main.commit=$(COMMIT)" ./cmd/eksctl

update-bindata:
	go generate ./pkg/eks

install-bindata:
	go get -u github.com/jteeuwen/go-bindata/...
