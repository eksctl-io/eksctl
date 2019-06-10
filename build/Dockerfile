FROM golang:1.12-alpine3.9

RUN apk add --no-cache \
      curl \
      git \
      make \
      gcc \
      musl-dev \
      && true

ENV CGO_ENABLED=0

ENV EKSCTL_BUILD $GOPATH/src/github.com/weaveworks/eksctl/build
COPY . $EKSCTL_BUILD

WORKDIR $EKSCTL_BUILD
RUN ./install.sh

WORKDIR $EKSCTL_BUILD/vendor/github.com/goreleaser/goreleaser
RUN go build -tags netgo && go install -tags netgo

WORKDIR $GOPATH
