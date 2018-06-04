FROM golang:1.10.2-alpine3.7

RUN apk add --update \
      curl \
      git \
      make \
      && true

ENV DEP_VERSION v0.4.1

# Install dep to avoid vendoring all the things
RUN curl --silent --location "https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64" --output /usr/local/bin/dep \
    && chmod +x /usr/local/bin/dep

ENV EKSCTL_BUILD $GOPATH/src/github.com/weaveworks/eksctl/build
RUN mkdir -p "$(dirname ${EKSCTL_BUILD})"
COPY . $EKSCTL_BUILD

WORKDIR $EKSCTL_BUILD

RUN dep ensure

RUN go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata
RUN go install ./vendor/golang.org/x/tools/cmd/stringer

WORKDIR $EKSCTL_BUILD/vendor/github.com/goreleaser/goreleaser
RUN make build && go install

WORKDIR $GOPATH
