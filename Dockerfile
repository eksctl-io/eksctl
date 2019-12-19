# Make sure to run the following commands after changes to this file are made:
# `make -f Makefile.docker update-build-image-manifest && make -f Makefile.docker push-build-image`

# This digest corresponds to golang:1.13.4-alpine3.10
FROM golang@sha256:679fe3791d2710d53efe26b05ba1c7083178d6375318b0c669b6bcd98f25c448 AS base

# Build-time dependencies
RUN apk add --no-cache \
    bash \
    curl \
    docker-cli \
    g++ \
    gcc \
    git \
    libsass-dev \
    make \
    musl-dev \
    && true

# Runtime dependencies. Build the root filesystem of the eksctl image at /out
RUN mkdir -p /out/etc/apk && cp -r /etc/apk/* /out/etc/apk/
RUN apk add --no-cache --initdb --root /out \
    alpine-baselayout \
    busybox \
    ca-certificates \
    coreutils \
    git \
    libc6-compat \
    openssh \
    && true

ENV KUBECTL_VERSION v1.11.5
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output /out/usr/local/bin/kubectl \
    && chmod +x /out/usr/local/bin/kubectl

# Remaining dependencies are controlled by go.mod
WORKDIR /src
ENV CGO_ENABLED=0 GOPROXY=https://proxy.golang.org,direct

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

COPY .requirements install-build-deps.sh go.mod go.sum /src/

# Install all build tools dependencies
RUN ./install-build-deps.sh

# The authenticator is a runtime dependency, so it needs to be in /out
RUN go install github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator \
    && mv $GOPATH/bin/aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator

# Add kubectl and aws-iam-authenticator to the PATH
ENV PATH="${PATH}:/out/usr/bin:/out/usr/local/bin"
